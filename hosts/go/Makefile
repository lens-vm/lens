test:
	cargo build --target wasm32-unknown-unknown --manifest-path "../tests/modules/rust_wasm32_simple/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "../tests/modules/rust_wasm32_simple2/Cargo.toml"
	go clean -testcache
	go test ./...
