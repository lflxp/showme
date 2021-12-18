package esync

import (
	"fmt"
	"sync/atomic"
)

type AtomicInt struct {
	value int64
}

func (t *AtomicInt) Inc() int {
	return t.Add(1)
}

func (t *AtomicInt) Add(v int) int {
	return int(atomic.AddInt64(&t.value, int64(v)))
}

func (t *AtomicInt) Set(v int) {
	atomic.StoreInt64(&t.value, int64(v))
}

func (t *AtomicInt) Get() int {
	return int(atomic.LoadInt64(&t.value))
}

func (t *AtomicInt) Swap(v int) int {
	return int(atomic.SwapInt64(&t.value, int64(v)))
}

func (t *AtomicInt) String() string {
	return fmt.Sprint(t.Get())
}
