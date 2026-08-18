# git-lore CLI (optional extension)

.PHONY: build test clean run

VERSION ?= $(shell cat ../../version.txt 2>/dev/null || echo 0.1.0)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/git-lore ./cmd/git-lore

test:
	go test ./...

clean:
	rm -rf bin/ dist/

run: build
	./bin/git-lore serve --repo ../.. --open
