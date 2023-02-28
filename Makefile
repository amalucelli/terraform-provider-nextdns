# This is used mostly for development.
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
VERSION := 0.1.0

# This is used mostly for development as this is how Terraform
# reads the plugin from the local plugins folder.
REGISTRY := github.com
ORG := amalucelli
NAME := nextdns

# Terraform requieres a specific format.
BINARY := terraform-provider-${NAME}

.PHONY: build
build:
	@go build -o ${BINARY}

.PHONY: install
install: build
	@mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${ORG}/${NAME}/${VERSION}/${GOOS}_${GOARCH}
	@mv ${BINARY} ~/.terraform.d/plugins/${REGISTRY}/${ORG}/${NAME}/${VERSION}/${GOOS}_${GOARCH}

.PHONY: clean
clean:
	@rm -rf examples/.terraform* examples/terraform.*

.PHONY: test
test:
	@go test ./...
