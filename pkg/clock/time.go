package clock

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/antlabs/stl/list"
)

const (
	haveStop = uint32(1)
)

// CancelFunc aborts a scheduled task.
type CancelFunc func()

// First use sync.Mutex to implement the function
// Use cas optimization later
type timepoint struct {
	timeNode
	sync.Mutex

	// |---16bit---|---16bit---|------32bit-----|
	// |---level---|---index---|-------seq------|
	// level is 1 in the near plate, and 2 in the T2ToTt[0] plate
	// index is the index value of the respective plate
	// seq increment id
	version uint64
}

func newTimeHead(level uint64, index uint64) *timepoint {
	head := &timepoint{}
	head.version = genVersionHeight(level, index)
	head.Init()
	return head
}

func genVersionHeight(level uint64, index uint64) uint64 {
	return level<<(32+16) | index<<32
}

func (t *timepoint) lockPushBack(node *timeNode, level uint64, index uint64) {
	t.Lock()
	defer t.Unlock()
	if atomic.LoadUint32(&node.stop) == haveStop {
		return
	}

	t.AddTail(&node.Head)
	atomic.StorePointer(&node.list, unsafe.Pointer(t))
	// Update the version information of the node
	atomic.StoreUint64(&node.version, atomic.LoadUint64(&t.version))
}

type timeNode struct {
	expire     uint64
	userExpire time.Duration
	callback   func()
	stop       uint32
	list       unsafe.Pointer // Store header information
	version    uint64         // Save node version information
	isSchedule bool

	list.Head
}

// A timeNode node has 4 states
// 1. Exist in the initialization linked list
// 2. Moved to tmp linked list
// 3.1 and 3.2 are the status of if else
// 3.1 is moved to the new list
// 3.2 Direct execution
// 1 and 3.1 status is no problem
// 2 and 3.2 states will be operations without lock protection, and there will be contention
func (t *timeNode) Stop() {

	atomic.StoreUint32(&t.stop, haveStop)

	// Use the version number algorithm to let timeNode know if it has been moved
	// The version of timeNode is the same as the version of the header, indicating that it has not been moved and can be deleted directly
	// If it is not the same, it may be in the 2nd or 3.2 state, use lazy deletion
	cpyList := (*timepoint)(atomic.LoadPointer(&t.list))
	cpyList.Lock()
	defer cpyList.Unlock()
	if atomic.LoadUint64(&t.version) != atomic.LoadUint64(&cpyList.version) {
		return
	}

	cpyList.Del(&t.Head)
}
