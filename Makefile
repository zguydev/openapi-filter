LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(PATH):$(LOCAL_BIN)

PKG:=github.com/zguydev/openapi-filter
FILTER_ENTRYPOINT:=.
FILTER_BIN:=$(LOCAL_BIN)/openapi-filter

GOLANGCI_LINT_VERSION:=v2.1.6

ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: .app-deps
.app-deps:
	go mod tidy

.PHONY: .bin-deps
.bin-deps:
	$(info "Installing bin deps...")
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: install
install: .bin-deps .app-deps

.PHONY: build
build:
	@mkdir -p $(LOCAL_BIN)
	@go build -o $(FILTER_BIN) $(FILTER_ENTRYPOINT)

.PHONY: clean
clean:
	@rm -i $(LOCAL_BIN)/*

.PHONY: lint
lint:
	@$(LOCAL_BIN)/golangci-lint run ./...

.PHONY: test
test:
	@go test -v -coverprofile=./coverage.out $(PKG)/...

.PHONY: cover
cover:
	@go test -cover -coverprofile ./coverage.out $(PKG)/...

.PHONY: show_cover
show_cover:
	@go tool cover -func=coverage.out

.PHONY: show_cover_html
show_cover_html:
	@go tool cover -html=coverage.out
