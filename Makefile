.PHONY: format
format:
	 goimports -w *.go

.PHONY: lint
lint: format
	golangci-lint run --config ./golangci-lint.yml ./...
