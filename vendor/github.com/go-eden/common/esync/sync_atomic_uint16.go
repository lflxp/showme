package esync

import (
	"fmt"
	"sync/atomic"
)

// AtomicUint16 support atomic operation
type AtomicUint16 struct {
	value int32
}

func (t *AtomicUint16) Inc() uint16 {
	return t.Add(1)
}

func (t *AtomicUint16) Add(v int) uint16 {
	return uint16(atomic.AddInt32(&t.value, int32(v)))
}

func (t *AtomicUint16) Set(v uint16) {
	atomic.StoreInt32(&t.value, int32(v))
}

func (t *AtomicUint16) Get() uint16 {
	return uint16(atomic.LoadInt32(&t.value))
}

func (t *AtomicUint16) Swap(v uint16) uint16 {
	return uint16(atomic.SwapInt32(&t.value, int32(v)))
}

func (t *AtomicUint16) String() string {
	return fmt.Sprint(t.Get())
}
