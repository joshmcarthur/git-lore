# git-lore — shell plugin + serve binary

.PHONY: web ensure-web build serve-bin test clean run install dist

VERSION ?= $(shell cat ../../version.txt 2>/dev/null || echo 0.1.0)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
PREFIX ?= $(HOME)/.local

web:
	cd web && npm ci && npm run build

# go:embed requires webdist at compile time; build it if missing.
ensure-web:
	@test -f internal/server/webdist/index.html || $(MAKE) web

serve-bin: ensure-web
	@mkdir -p lib/git-lore
	go build $(LDFLAGS) -o lib/git-lore/serve ./cmd/git-lore-serve
	@chmod +x bin/git-lore lib/git-lore/*.sh

build: web serve-bin

test: ensure-web
	go test ./internal/...
	@if command -v bats >/dev/null 2>&1; then \
	  bats test/; \
	else \
	  echo "bats not installed — skipping shell tests (install bats-core)"; \
	fi

clean:
	rm -rf dist/ web/node_modules/ internal/server/webdist lib/git-lore/serve lib/git-lore/serve.exe

run: build
	./bin/git-lore serve --repo ../.. --open

install: build
	install -d "$(PREFIX)/bin" "$(PREFIX)/lib/git-lore"
	install -m 755 bin/git-lore "$(PREFIX)/bin/git-lore"
	install -m 755 lib/git-lore/*.sh "$(PREFIX)/lib/git-lore/"
	install -m 755 lib/git-lore/serve "$(PREFIX)/lib/git-lore/serve"
	@echo "Installed to $(PREFIX)/bin/git-lore (ensure $(PREFIX)/bin is on PATH)"

# Cross-compile release archives: bin/git-lore + lib/git-lore/{*.sh,serve}
# Usage: make dist VERSION=0.2.0
dist: ensure-web
	@mkdir -p dist
	@set -e; \
	for pair in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do \
	  goos=$${pair%/*}; goarch=$${pair#*/}; \
	  ext=""; [ "$$goos" = windows ] && ext=".exe"; \
	  staging="dist/staging_$${goos}_$${goarch}"; \
	  rm -rf "$$staging"; \
	  mkdir -p "$$staging/bin" "$$staging/lib/git-lore"; \
	  cp bin/git-lore "$$staging/bin/git-lore"; \
	  cp lib/git-lore/*.sh "$$staging/lib/git-lore/"; \
	  echo "building serve for $$goos/$$goarch"; \
	  GOOS=$$goos GOARCH=$$goarch CGO_ENABLED=0 \
	    go build $(LDFLAGS) -o "$$staging/lib/git-lore/serve$$ext" ./cmd/git-lore-serve; \
	  archive="git-lore_$(VERSION)_$${goos}_$${goarch}"; \
	  if [ "$$goos" = windows ]; then \
	    (cd "$$staging" && zip -qr "../$$archive.zip" bin lib) && rm -rf "$$staging"; \
	  else \
	    tar -C "$$staging" -czf "dist/$$archive.tar.gz" bin lib && rm -rf "$$staging"; \
	  fi; \
	done
	@ls -la dist/
