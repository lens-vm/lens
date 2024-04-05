// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package module

import "io"

// Instance is the representation of loaded lens module. This will often be sourced from a WASM binary
// but it does not have to be.
type Instance struct {
	// Alloc allocates the given number of bytes in memory and returns the start index to the allocated block.
	Alloc func(size MemSize) (MemSize, error)

	// Transform transforms the data stored at the given start index, returning the start index of the result.
	//
	// The next function provided should return a wasm memory pointer to the next source item to be transformed.
	Transform func(next func() MemSize) (MemSize, error)

	// Memory returns an io.ReadWriter that can be used to read or write to the
	// linear memory that this module uses starting at the given offset.
	//
	// Values written to memory will be made available to this module, however changes made by the
	// module after this function has been called are not guaranteed to be visible to the returned io.Reader.
	Memory func(offset int32) io.ReadWriter

	// OwnedBy hosts a reference to any object(s) that may be required to live in memory for the lifetime of this Module.
	//
	// This is very important when working with some libraries (such as wasmer-go), as without this, dependencies of other members
	// of this Module may be garbage collected prematurely.
	OwnedBy any
}

// SliceReadWriter converts a byte slice into an io.ReadWriter.
type SliceReadWriter struct {
	data   []byte
	offset int32
}

// NewSliceReadWriter returns a new SliceReadWriter that reads and writes
// from the given byte slice starting at the given offset.
func NewSliceReadWriter(data []byte, offset int32) *SliceReadWriter {
	return &SliceReadWriter{data: data, offset: offset}
}

// Read implements the io.Reader interface.
func (s *SliceReadWriter) Read(dst []byte) (int, error) {
	n := copy(dst, s.data[s.offset:])
	s.offset = s.offset + int32(n)
	return n, nil
}

// Write implements the io.Writer interface.
func (s *SliceReadWriter) Write(src []byte) (int, error) {
	n := copy(s.data[s.offset:], src)
	s.offset = s.offset + int32(n)
	return n, nil
}
