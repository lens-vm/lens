// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package module

// Module is the representation of loaded lens module. This will often be sourced from a WASM binary
// but it does not have to be.
type Module struct {
	// Alloc allocates the given number of bytes in memory and returns the start index to the allocated block.
	Alloc func(size MemSize) (MemSize, error)

	// Transform transforms the data stored at the given start index, returning the start index of the result.
	Transform func(startIndex MemSize) (MemSize, error)

	// GetData returns the current state of the linear memory that this module uses.

	// Values written to the return slice will be made available to this module, however changes made by the
	// module after this function has been called are not guaranteed to be visible to the previously returned slice.
	GetData func() []byte

	// ownedBy hosts a reference to any object(s) that may be required to live in memory for the lifetime of this Module.
	//
	// This is very important when working with some libraries (such as wasmer-go), as without this dependencies of other members
	// of this Module may be garbage collected prematurely.
	OwnedBy any
}
