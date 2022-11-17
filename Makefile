test:
	$(MAKE) -C ./host-go test
	$(MAKE) -C ./tests/integration test
