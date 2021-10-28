package esync

import (
	"fmt"
	"math"
	"sync/atomic"
)

// AtomicFloat64 support atomic operation
type AtomicFloat64 struct {
	value uint64
}

func (t *AtomicFloat64) Set(v float64) {
	atomic.StoreUint64(&t.value, math.Float64bits(v))
}

func (t *AtomicFloat64) Get() float64 {
	return math.Float64frombits(atomic.LoadUint64(&t.value))
}

func (t *AtomicFloat64) Swap(v float64) float64 {
	return math.Float64frombits(atomic.SwapUint64(&t.value, math.Float64bits(v)))
}

func (t *AtomicFloat64) String() string {
	return fmt.Sprint(t.Get())
}
