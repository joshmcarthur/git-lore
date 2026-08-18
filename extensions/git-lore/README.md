# git-lore CLI

Optional extension for [git-lore](../..): a Go CLI that can abstract Lore
operations for people who prefer commands over agent skills.

This PR ships the **backend**: CLI entrypoint, lore API over `git`, and
`serve` with a placeholder HTML page. The Vue UI lands in a stacked follow-up.

Skills under `skills/` remain the primary surface.

## Build and run

```bash
make build
./bin/git-lore serve --repo ../.. --open
./bin/git-lore help
```

Flags for `serve`: `--repo <path>`, `--addr host:port` (default
`127.0.0.1:9473`), `--open`.

## Screenshots

Placeholder page until the Vue UI PR lands:

![Placeholder serve page](docs/screenshots/serve-placeholder.png)

## Commands

| Command | Status | Purpose |
|---------|--------|---------|
| `serve` | API + placeholder UI | Lore HTTP API; browser UI in next PR |
| `list` / `show` / `create` / … | planned | CLI counterparts to skill workflows |
