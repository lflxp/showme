package esync

import (
	"fmt"
	"sync/atomic"
)

type AtomicInt64 struct {
	value int64
}

func (t *AtomicInt64) Inc() int64 {
	return t.Add(1)
}

func (t *AtomicInt64) Add(v int64) int64 {
	return atomic.AddInt64(&t.value, v)
}

func (t *AtomicInt64) Set(v int64) {
	atomic.StoreInt64(&t.value, v)
}

func (t *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&t.value)
}

func (t *AtomicInt64) Swap(v int64) int64 {
	return atomic.SwapInt64(&t.value, v)
}

func (t *AtomicInt64) String() string {
	return fmt.Sprint(t.Get())
}
