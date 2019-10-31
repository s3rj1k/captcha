GO_BIN ?= go
ENV_BIN ?= env

export PATH := $(PATH):/usr/local/go/bin

all: test lint

update:
	$(ENV_BIN) GOPROXY=direct GOPRIVATE=github.com/s3rj1k/* $(GO_BIN) get -u
	$(GO_BIN) get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GO_BIN) get -u github.com/mgechev/revive
	$(GO_BIN) mod tidy

test:
	$(GO_BIN) test -failfast ./...

lint:
	golangci-lint run ./...
	revive -config revive.toml -exclude ./vendor/... ./...