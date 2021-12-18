package esync

import (
	"fmt"
	"sync/atomic"
)

// AtomicBool support atomic operation
type AtomicBool struct {
	value int32
}

func (t *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreInt32(&t.value, 1)
	} else {
		atomic.StoreInt32(&t.value, 0)
	}
}

func (t *AtomicBool) Get() bool {
	return atomic.LoadInt32(&t.value) != 0
}

func (t *AtomicBool) Swap(v bool) bool {
	var _v int32
	if v {
		_v = 1
	} else {
		_v = 0
	}
	return atomic.SwapInt32(&t.value, _v) != 0
}

func (t *AtomicBool) String() string {
	return fmt.Sprint(t.Get())
}
