# CLIProxyAPIPlus Fork

Maintained fork that tracks `router-for-me/CLIProxyAPI` and keeps GitHub Copilot provider support from the deleted Plus project.

## Source of Truth

| Remote | Purpose |
|--------|---------|
| `official` | Living upstream: `https://github.com/router-for-me/CLIProxyAPI` |
| `origin` | Maintained fork: `ravindrabarthwal/CLIProxyAPIPlus` |
| `upstream` | Historical/dead Plus upstream; do not sync from it |

Always treat official `CLIProxyAPI` as the base and keep the Plus delta limited to GitHub Copilot support.

## Plus Delta

Preserve these Copilot areas when merging official changes:

| Area | Files |
|------|-------|
| Copilot auth | `internal/auth/copilot/`, `sdk/auth/github_copilot.go`, `internal/cmd/github_copilot_login.go` |
| Copilot execution | `internal/runtime/executor/github_copilot_executor.go` and tests |
| Copilot models | `internal/registry/model_definitions.go`, `internal/registry/github_copilot_models_test.go` |
| Copilot config/aliases | `config.example.yaml`, `internal/config/oauth_model_alias_defaults.go` |
| Copilot management/login wiring | `cmd/server/main.go`, `internal/api/server.go`, `internal/api/handlers/management/auth_files.go` |

## Sync Policy

When merging `official/main`:

1. Prefer official for shared runtime, Amp, translators, config loading, storage, and management infrastructure.
2. Re-apply only the Copilot provider/auth/executor/model wiring listed above.
3. Keep import paths aligned with the module version in `go.mod`.
4. Do not reintroduce deleted Plus-only providers unless explicitly requested.

## Commands

```bash
# Sync official upstream
git fetch official --prune
git merge official/main

# Verify after conflict resolution
gofmt -w $(git diff --name-only -- '*.go')
go test ./...
go build -o cli-proxy-api-plus ./cmd/server/
```
