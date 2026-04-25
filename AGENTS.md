# AGENTS.md — CLIProxyAPIPlus Fork

## What This Is

Our maintained fork of `router-for-me/CLIProxyAPIPlus` — the upstream repo was **deleted** (returns 404).
This fork preserves **GitHub Copilot proxy support** that was never merged into the mainline `CLIProxyAPI`.

- **Fork**: `ravindrabarthwal/CLIProxyAPIPlus`
- **Upstream (dead)**: `router-for-me/CLIProxyAPIPlus`
- **Mainline (alive, no Copilot)**: `router-for-me/CLIProxyAPI` — v6.9.37+, 28k+ stars

## Why We Maintain This

The mainline `CLIProxyAPI` does **not** support GitHub Copilot as a provider.
Copilot support lived exclusively in the Plus repo, which was deleted.
We need Copilot proxy to power our AI infra at `ai.p8n.ai`.

## Our Custom Changes (on top of last upstream sync)

### 1. Copilot Responses API defaults (`5349cfea`)
- Changed `copilotOpenAIIntent` from `"conversation-panel"` → `"conversation-edits"`
- Added `applyGitHubCopilotResponsesDefaults()`: sets `store: false`, injects
  `reasoning.encrypted_content` into `include`, auto-sets `reasoning.summary: "auto"`
- **Why**: Aligns with how `pi-ai` sends requests to Copilot. Without this, reasoning
  features break and requests get persisted server-side unnecessarily.
- **Reference**: `badlogic/pi-mono` `packages/ai/src/providers/openai-responses.ts` and
  `packages/ai/src/providers/github-copilot-headers.ts`

### 2. Model definitions
We keep Copilot model definitions current by referencing `badlogic/pi-mono`
(`packages/ai/src/models.generated.ts`) as the source of truth for:
- Which models are available on GitHub Copilot
- Context window sizes, max tokens, supported endpoints
- Whether a model uses `openai-responses` or `anthropic-messages` API

Current Copilot models we support:
- GPT: gpt-4o variants, gpt-4.1, gpt-5, gpt-5-mini, gpt-5.1/5.2/5.3 + codex variants, **gpt-5.4**, **gpt-5.4-mini**
- Claude: haiku-4.5, opus-4.1/4.5/4.6/**4.7**, sonnet-4/4.5/4.6
- Gemini: 2.5-pro, 3-pro-preview, 3.1-pro-preview, 3-flash-preview
- Other: grok-code-fast-1, oswe-vscode-prime

## How pi.dev Does It (Reference)

Pi talks to Copilot **directly** — no proxy. All the headers, defaults, and model
definitions live client-side in `badlogic/pi-mono`:

| What | Where in pi-mono |
|------|------------------|
| `Openai-Intent: conversation-edits` | `packages/ai/src/providers/github-copilot-headers.ts:29` |
| `store: false` | `packages/ai/src/providers/openai-responses.ts:218` |
| `include: ["reasoning.encrypted_content"]` | `packages/ai/src/providers/openai-responses.ts:243` |
| `reasoning.summary: "auto"` | `packages/ai/src/providers/openai-responses.ts:241` |
| Model definitions | `packages/ai/src/models.generated.ts` (search `github-copilot`) |
| Copilot baseUrl | `https://api.individual.githubcopilot.com` |
| Claude models use | `anthropic-messages` API (not `openai-responses`) |
| GPT models use | `openai-responses` API |

## Adding New Models

When a new model appears on Copilot:

1. Check `badlogic/pi-mono` `packages/ai/src/models.generated.ts` for the `github-copilot` provider section
2. Note: `contextWindow`, `maxTokens`, `api` type (`openai-responses` vs `anthropic-messages`)
3. Add to `internal/registry/model_definitions.go` in `GetGitHubCopilotModels()`
4. GPT models → `SupportedEndpoints: []string{"/chat/completions", "/responses"}`
5. Claude models → `SupportedEndpoints: []string{"/chat/completions"}`
6. If model supports reasoning → add `Thinking` field

## Key Files

| File | Purpose |
|------|---------|
| `internal/registry/model_definitions.go` | Copilot model definitions (`GetGitHubCopilotModels()`) |
| `internal/runtime/executor/github_copilot_executor.go` | Copilot request execution, headers, responses defaults |
| `internal/runtime/executor/github_copilot_executor_test.go` | Tests for Copilot executor |
| `internal/translator/` | Request/response translation between API formats |
| `cmd/server/` | Main entry point |
| `Dockerfile` | Official multi-arch Docker build |
| `.goreleaser.yml` | Release binary build config |

## Deployment

Deployed as Docker container on Azure VM (`vm-openclaw` at `20.219.50.12`).
See `sentry-azure-k8s/wow-agents/deploy/` for:
- `Dockerfile.cli-proxy` — multi-stage build from this fork
- `docker-compose.cli-proxy-public.yml` — service definition
- `configs/cli-proxy-config.yaml` — runtime configuration
- `configs/Caddyfile.cli-proxy` — TLS termination (`ai.p8n.ai` → proxy:8317)
- `scripts/upgrade-cli-proxy-api-public.sh` — upgrade with auto-rollback

## Build & Test

```bash
# Run tests
go test ./...

# Build locally
go build -o cli-proxy-api-plus ./cmd/server/

# Docker build (from sentry-azure-k8s/wow-agents/deploy/)
docker compose -f docker-compose.cli-proxy-public.yml build

# Pin to specific commit
CLI_PROXY_REF=abc123 docker compose -f docker-compose.cli-proxy-public.yml build
```
