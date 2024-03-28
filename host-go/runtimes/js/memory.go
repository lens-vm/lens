// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build js

package js

import (
	"syscall/js"
)

type memory struct {
	array  js.Value
	offset int32
}

func newMemory(buffer js.Value, offset int32) *memory {
	array := js.Global().Get("Uint8Array").New(buffer)
	return &memory{
		array:  array,
		offset: offset,
	}
}

func (m *memory) Read(dst []byte) (int, error) {
	src := m.array.Call("subarray", m.offset)
	n := js.CopyBytesToGo(dst, src)
	m.offset = m.offset + int32(n)
	return n, nil
}

func (m *memory) Write(src []byte) (int, error) {
	dst := m.array.Call("subarray", m.offset)
	n := js.CopyBytesToJS(dst, src)
	m.offset = m.offset + int32(n)
	return n, nil
}
