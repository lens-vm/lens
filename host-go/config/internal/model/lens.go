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

	// Any additional parameters that you wish to be passed to the lens transform.
	//
	// The lens module must expose a `set_param` function if values are provided here.
	Arguments map[string]any
}
