package esync

import (
	"fmt"
	"sync/atomic"
)

type AtomicUint64 struct {
	value uint64
}

func (t *AtomicUint64) Inc() uint64 {
	return t.Add(1)
}

func (t *AtomicUint64) Add(v uint64) uint64 {
	return atomic.AddUint64(&t.value, v)
}

func (t *AtomicUint64) Set(v uint64) {
	atomic.StoreUint64(&t.value, v)
}

func (t *AtomicUint64) Get() uint64 {
	return atomic.LoadUint64(&t.value)
}

func (t *AtomicUint64) Swap(v uint64) uint64 {
	return atomic.SwapUint64(&t.value, v)
}

func (t *AtomicUint64) String() string {
	return fmt.Sprint(t.Get())
}
