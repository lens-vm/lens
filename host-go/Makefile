.PHONY: build
build:
	(cd cli && go build -o ../build/host-go)

test:
	$(MAKE) -C ../tests/modules build
	go clean -testcache
	go test ./...
