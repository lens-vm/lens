package modules

import (
	"path"
	"runtime"
)

// WasmPath1 contains a simple wasm32 rust lens that takes an item of `type1` and transforms it to `type2`.
var WasmPath1 string = getPathRelativeToProjectRoot(
	"/tests/modules/rust_wasm32_simple/target/wasm32-unknown-unknown/debug/rust_wasm32_simple.wasm",
)

// WasmPath2 contains a simple wasm32 rust lens that takes an item of `type2` and adds 1 to its age.
var WasmPath2 string = getPathRelativeToProjectRoot(
	"tests/modules/rust_wasm32_simple2/target/wasm32-unknown-unknown/debug/rust_wasm32_simple2.wasm",
)

// WasmPath3 contains a simple wasm32 assemblyScript lens that takes an item of `type2` and adds 10 to its age.
var WasmPath3 string = getPathRelativeToProjectRoot(
	"tests/modules/as_wasm32_simple/build/debug/as_wasm32_simple.wasm",
)

// WasmPath4 contains a wasm32 rust lens that takes two additional properties and a map and renames one of the properties.
var WasmPath4 string = getPathRelativeToProjectRoot(
	"/tests/modules/rust_wasm32_rename/target/wasm32-unknown-unknown/debug/rust_wasm32_rename.wasm",
)

func getPathRelativeToProjectRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(filename)))
	return path.Join(root, relativePath)
}
