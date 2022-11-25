test:
	$(MAKE) -C ../modules build
	$(MAKE) -C ../../host-go build
	go clean -testcache
	go test ./...
