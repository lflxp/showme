package esync

import (
	"fmt"
	"sync/atomic"
)

type AtomicInt16 struct {
	value int32
}

func (t *AtomicInt16) Inc() int16 {
	return t.Add(1)
}

func (t *AtomicInt16) Add(v int16) int16 {
	return int16(atomic.AddInt32(&t.value, int32(v)))
}

func (t *AtomicInt16) Set(v int16) {
	atomic.StoreInt32(&t.value, int32(v))
}

func (t *AtomicInt16) Get() int16 {
	return int16(atomic.LoadInt32(&t.value))
}

func (t *AtomicInt16) Swap(v int16) int16 {
	return int16(atomic.SwapInt32(&t.value, int32(v)))
}

func (t *AtomicInt16) String() string {
	return fmt.Sprint(t.Get())
}
