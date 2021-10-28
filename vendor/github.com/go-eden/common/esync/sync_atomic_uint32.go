package esync

import (
	"fmt"
	"sync/atomic"
)

// AtomicUint32 support atomic operation
type AtomicUint32 struct {
	value uint32
}

func (t *AtomicUint32) Inc() uint32 {
	return t.Add(1)
}

func (t *AtomicUint32) Add(v uint32) uint32 {
	return atomic.AddUint32(&t.value, v)
}

func (t *AtomicUint32) Set(v uint32) {
	atomic.StoreUint32(&t.value, v)
}

func (t *AtomicUint32) Get() uint32 {
	return atomic.LoadUint32(&t.value)
}

func (t *AtomicUint32) Swap(v uint32) uint32 {
	return atomic.SwapUint32(&t.value, v)
}

func (t *AtomicUint32) String() string {
	return fmt.Sprint(t.Get())
}
