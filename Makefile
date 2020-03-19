SHELL = /bin/bash

default: build

.PHONY: clean
.PHONY: deps
.PHONY: build
.PHONY: build-init-container
.PHONY: test-units
.PHONY: test-features
.PHONY: test-enemies
.PHONY: test
.PHONY: lint
.PHONY: fakes
.PHONY: run
.PHONY: test-kill
.PHONY: run-kill
.PHONY: env-variable-set
.PHONY: push

# #### CLEAN ####
clean:
	rm -rf build/*
	go clean --modcache

# #### DEPS ####
deps:
	go mod download

# #### BUILD ####
UAAGO_SRC = $(shell find . -name "*.go" | grep -v "_test\." )

bin/uaago: $(UAAGO_SRC) deps
	go build -o bin/uaago main.go

build: bin/uaago

# #### TEST ####
lint:
	git ls-files | grep '.go$$' | xargs goimports -l -w

test: lint deps
	ginkgo -r .
