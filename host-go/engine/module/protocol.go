package module

import (
	"encoding/binary"
)

// LenType is the type used to represent the byte length of an item transmitted to/from a lens module.
type LenType uint32

// LenSize is the size in bytes of `LenType`
var LenSize int32 = 4

// TypeIdType is the type used to represent the type of an item transmitted to/from a lens module.
//
// Positive values represent valid values, negative values represent errors, 0 is undefined.
type TypeIdType int8

// TypeIdSize is the size in bytes of `TypeIdType`
var TypeIdSize int32 = 1

// LenByteOrder is the byte order in which the byte length is written to memory.
//
// This is independent to the system byte order, and does not have to match it.
var TypeIdByteOrder = binary.LittleEndian

// LenByteOrder is the byte order in which the byte length is written to memory.
//
// This is independent to the system byte order, and does not have to match it.
var LenByteOrder = binary.LittleEndian

// MemSize is the memory size of the module runtime.
//
// This is independent to the host system memory size and does not have to match it.
// This type must only be alias, otherwise calling the wasm funcs will fail.
type MemSize = int32

const (
	JSONTypeID TypeIdType = 1
)

// IsError returns true if the given typeId is an error type.
//
// Otherwise returns false.
func IsError(typeId TypeIdType) bool {
	return typeId < 0
}
