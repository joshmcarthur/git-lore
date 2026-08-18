# git-lore CLI (optional extension)

.PHONY: web ensure-web build test clean run dist

VERSION ?= $(shell cat ../../version.txt 2>/dev/null || echo 0.1.0)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

web:
	cd web && npm ci && npm run build

# go:embed requires webdist at compile time; build it if missing.
ensure-web:
	@test -f internal/server/webdist/index.html || $(MAKE) web

build: web
	go build $(LDFLAGS) -o bin/git-lore ./cmd/git-lore

test: ensure-web
	go test ./...

clean:
	rm -rf bin/ dist/ web/node_modules/ internal/server/webdist

run: build
	./bin/git-lore serve --repo ../.. --open

dev-api: ensure-web
	go run $(LDFLAGS) ./cmd/git-lore serve --repo ../.. --addr 127.0.0.1:9473

dev-web:
	cd web && npm run dev

# Cross-compile release archives into dist/ (requires webdist already built).
# Usage: make dist VERSION=0.2.0
dist: ensure-web
	@mkdir -p dist
	@set -e; \
	for pair in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do \
	  goos=$${pair%/*}; goarch=$${pair#*/}; \
	  ext=""; [ "$$goos" = windows ] && ext=".exe"; \
	  out="git-lore_$${goos}_$${goarch}$$ext"; \
	  echo "building $$out"; \
	  GOOS=$$goos GOARCH=$$goarch CGO_ENABLED=0 \
	    go build $(LDFLAGS) -o "dist/$$out" ./cmd/git-lore; \
	  archive="git-lore_$(VERSION)_$${goos}_$${goarch}"; \
	  if [ "$$goos" = windows ]; then \
	    (cd dist && zip -q "$$archive.zip" "$$out" && rm "$$out"); \
	  else \
	    tar -C dist -czf "dist/$$archive.tar.gz" "$$out" && rm "dist/$$out"; \
	  fi; \
	done
	@ls -la dist/
