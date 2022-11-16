package tests

import (
	"path"
	"runtime"
)

// wasmPath1 contains a simple wasm32 rust lens that takes an item of `type1` and transforms it to `type2`.
var wasmPath1 string = getPathRelativeToProjectRoot(
	"/tests/modules/rust_wasm32_simple/target/wasm32-unknown-unknown/debug/rust_wasm32_simple.wasm",
)

// wasmPath2 contains a simple wasm32 rust lens that takes an item of `type2` and adds 1 to its age.
var wasmPath2 string = getPathRelativeToProjectRoot(
	"tests/modules/rust_wasm32_simple2/target/wasm32-unknown-unknown/debug/rust_wasm32_simple2.wasm",
)

// wasmPath3 contains a simple wasm32 assemblyScript lens that takes an item of `type2` and adds 10 to its age.
var wasmPath3 string = getPathRelativeToProjectRoot(
	"tests/modules/as_wasm32_simple/build/debug/as_wasm32_simple.wasm",
)

type type1 struct {
	Name string
	Age  int
}

type type2 struct {
	FullName string
	Age      int
}

func getPathRelativeToProjectRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(path.Dir(filename))))
	return path.Join(root, relativePath)
}
