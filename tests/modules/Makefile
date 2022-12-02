build:
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_simple/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_simple2/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_rename/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_counter/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_filter/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_normalize/Cargo.toml"
	(cd "./as_wasm32_simple/" && npm install && npm run asbuild:debug)
