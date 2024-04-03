.PHONY: deps\:test
deps\:test:
	$(MAKE) -C ../tests/modules build

.PHONY: build
build:
	(cd cli && go build -o ../build/host-go.exe)

.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test:
	@$(MAKE) deps:test
	@$(MAKE) clean
	go test ./...

.PHONY: test\:no-deps
test\:no-deps:
	@$(MAKE) clean
	go test ./...

.PHONY: test\:ci
test\:ci:
# We do not make the deps here, the ci does that seperately to avoid compiling stuff
# multiple times etc.
	gotestsum --format testname ./...

.PHONY: test\:js
test\:js:
	GOOS=js GOARCH=wasm gotestsum --format testname -- -exec wasmbrowsertest ./...