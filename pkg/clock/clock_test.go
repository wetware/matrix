package clock

import (
	"context"
	"math"
	"sync/atomic"
	"testing"
	"time"

	syncutil "github.com/lthibault/util/sync"
	"github.com/stretchr/testify/assert"
)

var t0 = time.Date(2021, 04, 9, 8, 0, 0, 0, time.UTC)

func TestTicker(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("Accuracy=10ms", func(t *testing.T) {
		t.Parallel()

		c := New()

		var ctr syncutil.Ctr
		c.Ticker(time.Millisecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0) // init the clock
		c.Advance(t0.Add(time.Millisecond * 1010))

		assert.Eventually(t, func() bool {
			return ctr.Num() == 10
		}, time.Millisecond*100, time.Millisecond*10)
	})

	t.Run("Accuracy=10µs", func(t *testing.T) {
		t.Parallel()

		c := New(WithAccuracy(time.Microsecond * 10))

		var ctr syncutil.Ctr
		c.Ticker(time.Microsecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0) // init the clock
		c.Advance(t0.Add(time.Microsecond * 1010))

		assert.Eventually(t, func() bool {
			return ctr.Num() == 10
		}, time.Millisecond*100, time.Millisecond*10)
	})
}

func Test_maxVal(t *testing.T) {
	t.Parallel()

	assert.Equal(t, maxVal(), uint64(math.MaxUint32))
}

func Test_LevelMax(t *testing.T) {
	t.Parallel()

	assert.Equal(t, levelMax(1), uint64(1<<(nearShift+levelShift)))
	assert.Equal(t, levelMax(2), uint64(1<<(nearShift+2*levelShift)))
	assert.Equal(t, levelMax(3), uint64(1<<(nearShift+3*levelShift)))
	assert.Equal(t, levelMax(4), uint64(1<<(nearShift+4*levelShift)))
}

func Test_GenVersion(t *testing.T) {
	t.Parallel()

	assert.Equal(t, genVersionHeight(1, 0xf), uint64(0x0001000f00000000))
	assert.Equal(t, genVersionHeight(1, 64), uint64(0x0001004000000000))
}

func Test_hour(t *testing.T) {
	t.Parallel()

	c := New(incr)

	testHour := new(bool)
	done := make(chan struct{}, 1)
	c.After(time.Hour, func() {
		*testHour = true
		done <- struct{}{}
	})

	expire := c.getExpire(time.Hour, 0)
	for i := 0; i < int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-done:
	case <-time.After(time.Second / 100):
	}
	assert.True(t, *testHour)
}

func Test_ScheduleFunc_5s(t *testing.T) {
	t.Parallel()

	c := New(incr)

	var first5 int32
	ctx, cancel := context.WithCancel(context.Background())

	const total = int32(1000)

	testTime := time.Second * 5

	c.Ticker(testTime, func() {
		atomic.AddInt32(&first5, 1)
		if atomic.LoadInt32(&first5) == total {
			cancel()
		}

	})

	expire := c.getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5)
}

func Test_ScheduleFunc_hour(t *testing.T) {
	t.Parallel()

	c := New(incr)

	var first5 int32
	ctx, cancel := context.WithCancel(context.Background())

	const total = int32(100)
	testTime := time.Hour

	c.Ticker(testTime, func() {
		atomic.AddInt32(&first5, 1)
		if atomic.LoadInt32(&first5) == total {
			cancel()
		}

	})

	expire := c.getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5)
}

func Test_ScheduleFunc_day(t *testing.T) {
	t.Parallel()

	c := New(incr)

	var first5 int32
	ctx, cancel := context.WithCancel(context.Background())

	const total = int32(10)
	testTime := time.Hour * 24

	c.Ticker(testTime, func() {
		atomic.AddInt32(&first5, 1)
		if atomic.LoadInt32(&first5) == total {
			cancel()
		}

	})

	expire := c.getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		c.Advance(time.Time{})
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5)
}

func TestAfter(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("Accuracy=10ms", func(t *testing.T) {
		t.Parallel()

		c := New()

		var ctr syncutil.Ctr
		c.After(time.Millisecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0)
		c.Advance(t0.Add(time.Second))

		assert.Eventually(t, func() bool {
			return ctr.Num() == 1
		}, time.Millisecond*100, time.Millisecond*10)
	})

	t.Run("Accuracy=µs", func(t *testing.T) {
		t.Parallel()

		c := New(WithAccuracy(time.Microsecond * 10))

		var ctr syncutil.Ctr
		c.After(time.Microsecond*100, func() {
			ctr.Incr()
		})

		c.Advance(t0)
		c.Advance(t0.Add(time.Second))

		assert.Eventually(t, func() bool {
			return ctr.Num() == 1
		}, time.Millisecond*100, time.Millisecond*10)
	})
}

func Test_Node_Stop_1(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := New()

	count := uint32(0)
	stop := c.After(time.Millisecond*10, func() {
		atomic.AddUint32(&count, 1)
	})

	go func() {
		time.Sleep(time.Millisecond * 30)
		stop()
		cancel()
	}()

	run(ctx, c, time.Millisecond*10)
	assert.NotEqual(t, count, 1)
}

func Test_Node_Stop(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := New()

	count := uint32(0)
	stop := c.After(time.Millisecond*100, func() {
		atomic.AddUint32(&count, 1)
	})
	stop()
	go func() {
		time.Sleep(time.Millisecond * 200)
		cancel()
	}()

	run(ctx, c, time.Millisecond*10)
	assert.NotEqual(t, count, 1)
}

func run(ctx context.Context, c *Clock, d time.Duration) {
	var t = t0
	for ctx.Err() == nil {
		c.Advance(t)
		t = t0.Add(d)
	}

	// tk := time.NewTicker(d)
	// defer tk.Stop()

	// for {
	// 	select {
	// 	case t := <-tk.C:
	// 		c.Advance(t)
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }
}

func incr(c *Clock) {
	WithAccuracy(-1)(c)
	c.ticks = func(time.Time) time.Duration {
		return c.curTimePoint + 1
	}
}