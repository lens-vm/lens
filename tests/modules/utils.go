// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package modules

import (
	_ "embed"
)

var (
	// RustWasm32Simple contains a simple wasm32 rust lens that takes an item of `type1` and transforms it to `type2`.
	//go:embed rust_wasm32_simple/target/wasm32-unknown-unknown/debug/rust_wasm32_simple.wasm
	RustWasm32Simple []byte
	// RustWasm32Simple2 contains a simple wasm32 rust lens that takes an item of `type2` and adds 1 to its age.
	//
	// Module also supplies an inverse function.
	//go:embed rust_wasm32_simple2/target/wasm32-unknown-unknown/debug/rust_wasm32_simple2.wasm
	RustWasm32Simple2 []byte
	// AsWasm32Simple contains a wasm32 rust lens that takes two additional properties and a map and renames one of the properties.
	//go:embed as_wasm32_simple/build/debug/as_wasm32_simple.wasm
	AsWasm32Simple []byte
	// RustWasm32Rename contains a wasm32 rust lens that takes two additional properties and a map and renames one of the properties.
	//go:embed rust_wasm32_rename/target/wasm32-unknown-unknown/debug/rust_wasm32_rename.wasm
	RustWasm32Rename []byte
	// RustWasm32Counter contains a wasm32 rust lens that sets an id using a counter.
	//go:embed rust_wasm32_counter/target/wasm32-unknown-unknown/debug/rust_wasm32_counter.wasm
	RustWasm32Counter []byte
	// RustWasm32Filter contains a wasm32 rust lens that only returns values where `__type` == "pass".
	//go:embed rust_wasm32_filter/target/wasm32-unknown-unknown/debug/rust_wasm32_filter.wasm
	RustWasm32Filter []byte
	// RustWasm32Normalize contains a wasm32 rust lens that turns books into multiple pages.
	//go:embed rust_wasm32_normalize/target/wasm32-unknown-unknown/debug/rust_wasm32_normalize.wasm
	RustWasm32Normalize []byte
)
