package clock

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/antlabs/stl/list"
	syncutil "github.com/lthibault/util/sync"
)

const (
	DefaultAccuracy = time.Millisecond * 10

	nearShift  = 8
	nearSize   = 1 << nearShift
	levelShift = 6
	levelSize  = 1 << levelShift
	nearMask   = nearSize - 1
	levelMask  = levelSize - 1
)

type Clock struct {
	// Must come first.  See:  https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	tick syncutil.Ctr64 // monotonically increasing

	t1     [nearSize]*timepoint     // 256 slots
	t2Tot5 [4][levelSize]*timepoint // 4x64 time-scales

	accuracy     time.Duration
	ticks        func(time.Time) time.Duration
	curTimePoint time.Duration
}

func New(opt ...Option) *Clock {
	var c Clock
	for _, option := range withDefault(opt) {
		option(&c)
	}

	for i := 0; i < nearSize; i++ {
		c.t1[i] = newTimeHead(1, uint64(i))

	}

	for i := 0; i < 4; i++ {
		for j := 0; j < levelSize; j++ {
			c.t2Tot5[i][j] = newTimeHead(uint64(i+2), uint64(j))
		}
	}

	return &c
}

func maxVal() uint64 {
	return (1 << (nearShift + 4*levelShift)) - 1
}

func levelMax(index int) uint64 {
	return 1 << (nearShift + index*levelShift)
}

func (c *Clock) Accuracy() time.Duration { return c.accuracy }

func (c *Clock) index(n int) uint64 {
	return (uint64(c.tick) >> (nearShift + levelShift*n)) & levelMask
}

func (c *Clock) add(node *timeNode, tick uint64) *timeNode {
	var head *timepoint
	expire := node.expire
	idx := expire - tick

	level, index := uint64(1), uint64(0)

	if idx < nearSize {

		index = uint64(expire) & nearMask
		head = c.t1[index]

	} else {

		max := maxVal()
		for i := 0; i <= 3; i++ {

			if idx > max {
				idx = max
				expire = idx + tick
			}

			if uint64(idx) < levelMax(i+1) {
				index = uint64(expire >> (nearShift + i*levelShift) & levelMask)
				head = c.t2Tot5[i][index]
				level = uint64(i) + 2
				break
			}
		}
	}

	if head == nil {
		panic("not found head")
	}

	head.lockPushBack(node, level, index)

	return node
}

func (c *Clock) After(d time.Duration, callback func()) (cancel func()) {
	tick := c.tick.Load()

	node := &timeNode{
		expire:   uint64(d/c.accuracy + time.Duration(tick)),
		callback: callback,
	}

	return c.add(node, tick).Stop
}

func (c *Clock) getExpire(expire time.Duration, tick uint64) time.Duration {
	return expire/c.accuracy + time.Duration(tick)
}

func (c *Clock) Ticker(d time.Duration, callback func()) (cancel func()) {
	tick := c.tick.Load()

	node := &timeNode{
		userExpire: d,
		expire:     uint64(c.getExpire(d, tick)),
		callback:   callback,
		isSchedule: true,
	}

	return c.add(node, tick).Stop
}

// Move linked list
func (c *Clock) cascade(levelIndex int, index int) {

	tmp := newTimeHead(0, 0)

	l := c.t2Tot5[levelIndex][index]
	l.Lock()
	if l.Len() == 0 {
		l.Unlock()
		return
	}

	l.ReplaceInit(&tmp.Head)

	// Every time an element of the linked list is moved away, the version is modified
	atomic.AddUint64(&l.version, 1)
	l.Unlock()

	offset := unsafe.Offsetof(tmp.Head)
	tmp.ForEachSafe(func(pos *list.Head) {
		node := (*timeNode)(pos.Entry(offset))
		c.add(node, c.tick.Load())
	})

}

// moveAndExec function function
//1. Move to the near list first
//2. When the near list node is empty, move some nodes from the upper layer to the next layer
//3. Execute again
func (c *Clock) moveAndExec() {

	// // time overflow?
	// if uint32(c.tick) == 0 {
	// 	// TODO
	// 	// return
	// }

	// If the plate on this layer does not have a timer,
	// move some from the plate on the upper layer at this time
	index := c.tick & nearMask
	if index == 0 {
		for i := 0; i <= 3; i++ {
			index2 := c.index(i)
			c.cascade(i, int(index2))
			if index2 != 0 {
				break
			}
		}
	}

	c.tick.Incr()

	c.t1[index].Lock()
	if c.t1[index].Len() == 0 {
		c.t1[index].Unlock()
		return
	}

	head := newTimeHead(0, 0)
	t1 := c.t1[index]
	t1.ReplaceInit(&head.Head)
	atomic.AddUint64(&t1.version, 1)
	c.t1[index].Unlock()

	// Execute, the timer in the linked list
	offset := unsafe.Offsetof(head.Head)

	head.ForEachSafe(func(pos *list.Head) {
		val := (*timeNode)(pos.Entry(offset))
		head.Del(pos)

		if atomic.LoadUint32(&val.stop) == haveStop {
			return
		}

		go val.callback()

		if val.isSchedule {
			tick := uint64(c.tick)
			// The tick must be subtracted by 1.
			// The current callback is called and already contains a time slice.
			// If you don’t subtract this time slice, every time there is one more time slice,
			// it becomes an accumulator, and the periodic timer will progressively become inaccurate.
			val.expire = uint64(c.getExpire(val.userExpire, tick-1))
			c.add(val, tick)
		}
	})

}

func (c *Clock) Advance(t time.Time) {
	if c.curTimePoint == 0 {
		c.curTimePoint = c.ticks(t)
	}

	ts := c.ticks(t)
	if ts < c.curTimePoint {
		c.curTimePoint = ts
		return
	}

	diff := ts - c.curTimePoint
	c.curTimePoint = ts

	for i := time.Duration(0); i < diff; i++ {
		c.moveAndExec()
	}
}
