GO = go
GOSRC = $(shell find . -type f -name '*.go')
TEST_TIMEOUT=10s

.PHONY: all
all: format deps test full-lint doc

.PHONY: help
help:
	@echo "\nTARGETS:\n"
	@sed -n 's/^\.PHONY: \(\S*\)/\1/p' < Makefile

.PHONY: deps
deps:
	 $(GO) mod tidy -v && $(GO) mod verify

.PHONY: lint
lint:
	$(info #Running lint...)
	golangci-lint run --new-from-rev=origin/$(shell (grep -s DEFAULT_BRANCH: .gitlab-ci.yml | cut -d: -f2 | sed 's/^\s*$$/master/' | tr -d ' \t'; echo master) | head -n1) ./...

.PHONY: full-lint
full-lint:
	$(info #Running full lint...)
	golangci-lint run ./...

.PHONY: format
format:
	find . -name '*.go' | xargs $(GO) tool goimports -local github.com/my-mail-ru/ -l -w

.PHONY: test
test:
	$(GO) test $(GOFLAGS) -v     \
	    -timeout $(TEST_TIMEOUT) \
	    -covermode count         \
	    -coverprofile cover.out  \
	    -coverpkg ./...          \
	    ./...
	$(GO) tool cover -html cover.out -o cover.html

.PHONY: generate
generate: 
	$(GO) generate

.PHONY: doc
doc: README.md # doc/index.html

README.md: $(GOSRC)
	gomarkdoc --output $@ ./...

doc/index.html: $(GOSRC)
	godoc-static -destination=doc .
