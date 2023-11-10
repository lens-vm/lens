.PHONY: deps\:build
deps\:build:
	rustup target add wasm32-unknown-unknown

.PHONY: build
build:
	$(MAKE) deps:build
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_simple/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_simple2/Cargo.toml"
	cargo build --target wasm32-unknown-unknown --manifest-path "./rust_wasm32_rename/Cargo.toml"
	(cd "./as_wasm32_simple/" && npm install && npm run asbuild:debug)
