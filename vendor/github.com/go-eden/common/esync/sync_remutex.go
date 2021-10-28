package esync

import (
	"fmt"
	"github.com/go-eden/common/goid"
	"sync"
	"sync/atomic"
)

// ReMutex Reentrant Mutex
type ReMutex struct {
	sync.Mutex
	holder  int64
	counter int32
}

func (t *ReMutex) Lock() {
	var gid = goid.Gid()

	if atomic.LoadInt64(&t.holder) != gid {
		t.Mutex.Lock()
		atomic.StoreInt64(&t.holder, gid)
		t.counter = 1
		return
	}

	// relock
	t.counter++
}

func (t *ReMutex) Unlock() {
	var gid = goid.Gid()

	if atomic.LoadInt64(&t.holder) != gid {
		panic(fmt.Sprintf("Lock conflict, held by g[%v]", t.holder))
	}
	if t.counter > 0 {
		t.counter--
	}
	if t.counter == 0 {
		atomic.StoreInt64(&t.holder, 0)
		t.Mutex.Unlock()
	}
}
