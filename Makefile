lint:
	@golangci-lint run -E gofumpt --fix

.PHONY: lint