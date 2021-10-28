package esync

import (
	"fmt"
	"math"
	"sync/atomic"
)

// AtomicFloat32 support atomic operation
type AtomicFloat32 struct {
	value uint32
}

func (t *AtomicFloat32) Set(v float32) {
	atomic.StoreUint32(&t.value, math.Float32bits(v))
}

func (t *AtomicFloat32) Get() float32 {
	return math.Float32frombits(atomic.LoadUint32(&t.value))
}

func (t *AtomicFloat32) Swap(v float32) float32 {
	return math.Float32frombits(atomic.SwapUint32(&t.value, math.Float32bits(v)))
}

func (t *AtomicFloat32) String() string {
	return fmt.Sprint(t.Get())
}
