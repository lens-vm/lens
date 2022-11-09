// tiny-go wasm requires package name to be main... https://github.com/tinygo-org/tinygo/issues/2703
package main

import (
	"unsafe"
)

// tiny-go doesnt really support pure wasm, and forces the import of a whole bunch of extra junk that we dont want:
// https://github.com/tinygo-org/tinygo/issues/1383
// https://github.com/tinygo-org/tinygo/issues/3068

// tiny-go wasm requires a main func
func main() {}

type sliceHeader struct {
	addr unsafe.Pointer
	len  int
	cap  int
}

func alloc(size int) uintptr {
	buf := make([]byte, 0, size)
	//return reflect.ValueOf(buf).Pointer()
	return uintptr((*sliceHeader)(unsafe.Pointer(&buf)).addr)
}

func transform(ptr uintptr) uintptr {
	panic("")
}
