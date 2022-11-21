test:
	cargo build --target wasm32-unknown-unknown --manifest-path "../tests/modules/rust_wasm32_simple/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "../tests/modules/rust_wasm32_simple2/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "../tests/modules/rust_wasm32_rename/Cargo.toml"
	(cd "../tests/modules/as_wasm32_simple/" && npm run asbuild:debug)
	go clean -testcache
	go test ./...
