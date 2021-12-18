// +build go1.16

package routine

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	gDead = 6
)

//go:linkname runtimeG runtime.g
type runtimeG struct {
}

//go:linkname runtimeAtomicAllG runtime.atomicAllG
func runtimeAtomicAllG() (**runtimeG, uintptr)

//go:linkname runtimeReadgstatus runtime.readgstatus
func runtimeReadgstatus(g *runtimeG) uint32

//go:linkname runtimeIsSystemGoroutine runtime.isSystemGoroutine
func runtimeIsSystemGoroutine(gp *runtimeG, fixed bool) bool

// getAllGoidByNative retrieve all goid through runtime.atomicAllG
func getAllGoidByNative() (goids []int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("get all goid failed: %v", e))
			goids = nil
		}
	}()
	root, n := runtimeAtomicAllG()
	goids = make([]int64, 0, n)
	for i := uintptr(0); i < n; i++ {
		gp := *(**runtimeG)(unsafe.Pointer(uintptr(unsafe.Pointer(root)) + i*ptrSize))
		if runtimeReadgstatus(gp) == gDead || runtimeIsSystemGoroutine(gp, false) {
			continue
		}
		gid := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(gp)) + goidOffset))
		goids = append(goids, *gid)
	}
	return
}
