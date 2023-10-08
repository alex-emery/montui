GOPATH := $(shell go env GOPATH)
GOBIN := $(if $(GOPATH),$(GOPATH)/bin,$(HOME)/go/bin)
GOLINT_VERSION := v1.52.2

.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLINT_VERSION) 
	$(GOBIN)/golangci-lint run --timeout=10m ./pkg/... ./cmd/... 

.PHONY: test 
test:
	go test -race ./... 

.PHONY: build
build:
	go build ./cmd/main.go

.PHONY: generate
generate:
	oapi-codegen -package nordigen ./pkg/nordigen/swagger.json > ./pkg/nordigen/nordigen.gen.go