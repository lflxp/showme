package etime

import (
	"fmt"
	"github.com/go-eden/common/esync"
	"runtime/debug"
	"time"
)

type Interval struct {
	closed   esync.AtomicBool
	f        func()
	timer    *time.Timer
	duration time.Duration
}

func NewInterval(d time.Duration, f func()) *Interval {
	t := &Interval{
		duration: d,
		f:        f,
	}
	t.timer = time.AfterFunc(d, t.exec)
	return t
}

func (t *Interval) exec() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("interval timer-func panic: \n" + string(debug.Stack()))
		}
		if !t.closed.Get() {
			t.timer.Reset(t.duration)
		}
	}()
	if !t.closed.Get() {
		t.f()
	}
}

func (t *Interval) Close() {
	if t.closed.Swap(true) == false {
		t.timer.Stop()
	}
}
