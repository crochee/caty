OUT_DIR = ./out
PROJECT = obs
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
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		if [ $(GO_VERSION) -gt 16 ]; then \
            go install  github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1; \
        else \
            echo "please update go to 1.16+"; \
        fi \
	fi
	@bash ./scripts/ci-lint.sh;