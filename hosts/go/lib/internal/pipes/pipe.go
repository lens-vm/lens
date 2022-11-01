package pipes

import "lens-host/lib/enumerable"

// Pipe extends the Enumerable interface to allow for more efficient communication between
// lens modules.
type Pipe[T any] interface {
	enumerable.Enumerable[T]

	// Result returns the current value of the enumerable as a byte array.  This includes
	// the length specifier.
	Result() []byte
}
