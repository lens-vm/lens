// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pipes

import "github.com/sourcenetwork/immutable/enumerable"

// Pipe extends the Enumerable interface to allow for more efficient communication between
// lens modules.
type Pipe[T any] interface {
	enumerable.Enumerable[T]

	// Bytes returns the current value of the enumerable as a byte array.  This includes
	// the length specifier.
	Bytes() ([]byte, error)
}
