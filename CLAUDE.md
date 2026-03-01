# ig-cli Development Guide

## Build & Test Commands

```bash
make build      # Build binary → ./ig-cli
make test       # Run tests with race detector + enforce 80% coverage
make lint       # Run golangci-lint
make clean      # Remove build artifacts
```

## Project Structure

```
cmd/            # Cobra CLI commands (one file per command)
internal/
  api/          # Instagram Graph API client (Client interface + GraphClient)
  auth/         # Keychain storage, OAuth flow, token management
  config/       # Viper-backed config (~/.ig-cli/config.yaml)
  models/       # Data models (Media, Comment, Insight, etc.)
```

## Key Architecture

- **Interface-driven**: `api.Client` interface enables mock injection for tests
- **Injectable factories**: `clientFactory` and `keychainFactory` in `cmd/helpers.go` allow test overrides
- **OAuth flow**: `auth.OAuthFlow` has injectable `OpenBrowser` field (set to no-op in tests)
- **Token auto-refresh**: `auth.TokenManager.GetValidToken()` refreshes tokens within 24h of expiry
- **Config**: App ID in `~/.ig-cli/config.yaml`, App Secret + tokens in OS keychain via go-keyring

## Testing

- Tests use `httptest.Server` for real HTTP integration tests of `GraphClient`
- `MockClient` and `MockKeychain` for unit tests
- OAuth tests use `noopBrowser` to avoid opening real browsers
- Config tests use `t.TempDir()` + `t.Setenv("HOME", ...)` for isolation
- Coverage threshold: 80% total, 75% per-package, 70% per-file

## Git Hooks

Enable pre-push hook:
```bash
git config core.hooksPath .githooks
```

## Coding Conventions

- Run `gofmt` before committing (enforced by golangci-lint)
- Combine parameters of the same type: `func(a, b string)` not `func(a string, b string)`
- Use `http.NoBody` instead of `nil` for empty request bodies
- All output to stdout as JSON, errors/logs to stderr
