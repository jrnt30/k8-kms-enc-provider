# Heavily inspired by https://github.com/heptio/authenticator/ Makefile/GoReleaser
default: build
GORELEASER := $(shell command -v goreleaser 2> /dev/null)
PROTOC := $(shell command -v protoc 2> /dev/null)

.PHONY: build test format check-formatting check-tools generate-protobufs

check: check-formatting check-tools

check-tools:
ifndef GORELEASER
	$(error "goreleaser not found (`go get -u -v github.com/goreleaser/goreleaser` to fix)")
endif
ifndef PROTOC
	$(error "PROTOC not found (`go get -u -v github.com/goreleaser/goreleaser` to fix)")
endif
	@true

check-formatting:
	@if [ ! `find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -l {} + | wc -l` -eq 0 ]; then \
		echo "Changes present in go files.  Run 'make format' to clean"; \
		exit 127; \
	fi

format:
	find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -w {} +;

release: check-tools
	$(GORELEASER) --rm-dist

generate-protobufs: check-tools
	$(PROTOC) --go_out plugins=grpc:v1beta1/ --proto_path proto/ proto/*.proto

build-local: check-tools
	$(GORELEASER) release --skip-publish

build: check-tools generate-protobufs format
	$(GORELEASER) --skip-publish --rm-dist --snapshot
