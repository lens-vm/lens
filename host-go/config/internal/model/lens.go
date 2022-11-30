/*
The model package contains a lens-file-format agnostic set of types that describe the contents
of a given lens file.
*/
package model

type Lens struct {
	// The LensModules that should be applied to the source data, declared in the order
	// in which they should be executed.
	Lenses []LensModule
}

type LensModule struct {
	// The path to the wasm binary containing the lens transform that you wish to be applied.
	Path string

	// If true, the module will be inversed.
	//
	// This may result in an error if the module does not provide an inverse function.
	Inverse bool

	// Any additional parameters that you wish to be passed to the lens transform.
	//
	// The lens module must expose a `set_param` function if values are provided here.
	Arguments map[string]any
}
