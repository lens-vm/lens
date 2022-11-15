package pipes

import "github.com/lens-vm/lens/host-go/engine/enumerable"

// Pipe extends the Enumerable interface to allow for more efficient communication between
// lens modules.
type Pipe[T any] interface {
	enumerable.Enumerable[T]

	// Bytes returns the current value of the enumerable as a byte array.  This includes
	// the length specifier.
	Bytes() []byte
}