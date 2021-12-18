package esync

import (
	"fmt"
	"sync/atomic"
)

type AtomicInt32 struct {
	value int32
}

func (t *AtomicInt32) Inc() int32 {
	return t.Add(1)
}

func (t *AtomicInt32) Add(v int32) int32 {
	return atomic.AddInt32(&t.value, v)
}

func (t *AtomicInt32) Set(v int32) {
	atomic.StoreInt32(&t.value, v)
}

func (t *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&t.value)
}

func (t *AtomicInt32) Swap(v int32) int32 {
	return atomic.SwapInt32(&t.value, v)
}

func (t *AtomicInt32) String() string {
	return fmt.Sprint(t.Get())
}
