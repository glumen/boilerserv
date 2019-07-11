export GOBIN = $(PWD)/bin

GOLANGCI_LINT = $(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION = v1.17.1

.PHONY: CI
CI:lint test

.PHONY: lint
lint:$(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

.PHONY: test
test:
	go test ./...

# Tools
$(GOLANGCI_LINT):
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
