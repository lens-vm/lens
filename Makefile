.PHONY: test
test:
	$(MAKE) -C ./host-go test
	$(MAKE) -C ./tests/integration test

.PHONY: test\:scripts
test\:scripts:
	@$(MAKE) -C ./tools/scripts/ test
