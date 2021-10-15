OUT_DIR = ./out
PROJECT = cca
GO_VERSION=$(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)

.PHONY: all
all:
	@echo "wait for write"

.PHONY: run
run:
	@CGO_ENABLED=0 GOARCH=amd64 go run ./cmd/$(PROJECT)

.PHONY: swag
swag:
	@hash swagger > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
    	go install  github.com/go-swagger/go-swagger/cmd/swagger@0.27.0; \
	fi
	@swagger generate spec -o ./docs/swagger.json && swagger flatten --with-expand ./docs/swagger.json -o ./docs/swagger.json
.PHONY: lint
lint:
	@./scripts/ci.sh
.PHONY: layout
layout:
	@./scripts/layout.sh