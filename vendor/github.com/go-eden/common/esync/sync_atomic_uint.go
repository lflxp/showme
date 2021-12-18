package esync

import (
	"fmt"
	"sync/atomic"
)

type AtomicUint struct {
	value uint64
}

func (t *AtomicUint) Inc() uint {
	return t.Add(1)
}

func (t *AtomicUint) Add(v uint) uint {
	return uint(atomic.AddUint64(&t.value, uint64(v)))
}

func (t *AtomicUint) Set(v uint) {
	atomic.StoreUint64(&t.value, uint64(v))
}

func (t *AtomicUint) Get() uint {
	return uint(atomic.LoadUint64(&t.value))
}

func (t *AtomicUint) Swap(v uint) uint {
	return uint(atomic.SwapUint64(&t.value, uint64(v)))
}

func (t *AtomicUint) String() string {
	return fmt.Sprint(t.Get())
}
