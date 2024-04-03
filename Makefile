.PHONY: deps\:test-ci
deps\:test-ci:
	go install gotest.tools/gotestsum@latest
	$(MAKE) -C ./host-go build
	$(MAKE) -C ./tests/modules build

.PHONY: deps\:test
deps\:test:
	$(MAKE) -C ./host-go build
	$(MAKE) -C ./tests/modules build

.PHONY: deps\:test-js
deps\:test-js:
	go install github.com/agnivade/wasmbrowsertest@latest
	$(MAKE) -C ./host-go build
	$(MAKE) -C ./tests/modules build

.PHONY: test
test:
	$(MAKE) deps:test
	$(MAKE) --no-print-directory -C ./host-go test:no-deps
	$(MAKE) --no-print-directory -C ./tests/integration test:no-deps

.PHONY: test\:ci
test\:ci:
	$(MAKE) --no-print-directory -C ./host-go test:ci
	$(MAKE) --no-print-directory -C ./tests/integration test:ci

.PHONY: test\:scripts
test\:scripts:
	@$(MAKE) -C ./tools/scripts/ test

.PHONY: test\:js
test\:js:
	$(MAKE) --no-print-directory -C ./host-go test:js