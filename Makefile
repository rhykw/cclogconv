INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
TEST ?= $(shell go list ./... | grep -v vendor)
VERSION = $(shell cat VERSION)
REVISION = $(shell git describe --always)
GO ?= CGO_ENABLED=0 GO111MODULE=on go

default: build
ci: depsdev test vet

depsdev: ## Installing dependencies for development
	$(GO) get -u github.com/tcnksm/ghr

build: ## Build as linux binary
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building$(RESET)"
	$(GO) build -o cclogconv cmd/cclogconv/main.go

test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb: ## test data
	mkdir -p test/tmp
	tar -C test/tmp -zxf test/GeoLite2-Country_20181113.tar.gz

test/tmp/test.log:
	mkdir -p test/tmp
	tar -C test/tmp -zxf test/test_log.tar.gz

test: test/tmp/GeoLite2-Country_20181113/GeoLite2-Country.mmdb test/tmp/test.log ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -cover -v $(TEST) -timeout=30s -parallel=4
	$(GO) test $(TEST)

vet: ## Exec go vet
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Vetting$(RESET)"
	$(GO) vet $(TEST)

clean:
	-rm -f cclogconv
	-rm -rf test/tmp

.PHONY: default test

