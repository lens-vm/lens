package module

import (
	"encoding/binary"
	"reflect"
)

// LenType is the type used to represent the byte length of an item transmitted to/from a lens module.
type LenType uint32

// LenSize is the size in bytes of `LenType`
var LenSize = int32(reflect.TypeOf(LenType(0)).Size())

// LenByteOrder is the byte order in which the byte length is written to memory.
//
// This is independent to the system byte order, and does not have to match it.
var LenByteOrder = binary.LittleEndian

// MemSize is the memory size of the module runtime.
//
// This is independent to the host system memory size and does not have to match it.
// This type must only be alias, otherwise calling the wasm funcs will fail.
type MemSize = int32
