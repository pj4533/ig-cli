# Contributing

## Setup

```bash
git clone https://github.com/pj4533/ig-cli.git
cd ig-cli
go mod download
```

### Enable Git hooks

```bash
git config core.hooksPath .githooks
```

This runs `make lint && make test` before every push.

## Development Workflow

1. Create a branch for your change
2. Make your changes
3. Run lint: `make lint`
4. Run tests: `make test`
5. Open a pull request

## Code Style

- Run `gofmt` on all Go files (enforced by linter)
- Combine parameters of the same type: `func(a, b string)`
- Use `http.NoBody` instead of `nil` for empty HTTP request bodies
- All public functions need doc comments

## Testing

### Running tests

```bash
make test
```

This runs tests with the race detector and enforces coverage thresholds:
- 80% total
- 75% per package
- 70% per file

### Writing tests

- Use `httptest.Server` for testing real HTTP interactions
- Use `api.MockClient` for testing command logic
- Use `auth.MockKeychain` instead of the OS keychain
- Set `OAuthFlow.OpenBrowser = noopBrowser` in OAuth tests
- Use `t.TempDir()` and `t.Setenv("HOME", ...)` for config isolation

### Test file locations

Test files live alongside the code they test:
- `internal/api/client_test.go`
- `internal/auth/keychain_test.go`
- `cmd/data_commands_test.go`

## CI

GitHub Actions runs on every push and pull request:
1. **Lint** — golangci-lint
2. **Test** — tests + coverage enforcement
3. **Build** — compilation check

## Pull Request Guidelines

- Keep PRs focused on a single change
- Include tests for new functionality
- Ensure `make lint` and `make test` pass
- Update documentation if behavior changes
