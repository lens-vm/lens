// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package module

import (
	"encoding/binary"
	"math"
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
	ErrTypeID  TypeIdType = -1
	NilTypeID  TypeIdType = 0
	JSONTypeID TypeIdType = 1

	// A type id that denotes the end of stream.
	//
	// If recieved it signals that the end of the stream has been reached and that the source will no longer yield
	// new values.
	EOSTypeID TypeIdType = math.MaxInt8
)

// IsError returns true if the given typeId is an error type.
//
// Otherwise returns false.
func (typeId TypeIdType) IsError() bool {
	return typeId < 0
}

// IsEOS returns true if the given typeId declares that the end of stream has been reached.
//
// Otherwise returns false.
func (typeId TypeIdType) IsEOS() bool {
	return typeId == EOSTypeID
}
