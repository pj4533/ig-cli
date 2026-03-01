# Architecture

## Package Diagram

```
main.go
  └─ cmd/
       ├─ root.go          (Cobra root, global flags)
       ├─ helpers.go        (getClient, outputJSON, factories)
       ├─ auth*.go          (auth commands)
       ├─ media*.go         (media commands)
       ├─ comments*.go      (comments commands)
       ├─ insights*.go      (insights commands)
       ├─ discover.go       (discover command)
       └─ completion.go     (shell completions)
  └─ internal/
       ├─ api/
       │    ├─ client.go     (Client interface + GraphClient)
       │    ├─ auth.go       (token exchange methods)
       │    ├─ media.go      (ListMedia, GetMediaInsights)
       │    ├─ comments.go   (ListComments, ListReplies)
       │    ├─ insights.go   (GetAccountInsights, GetAudienceDemographics)
       │    ├─ discovery.go  (DiscoverUser)
       │    ├─ errors.go     (APIError type)
       │    └─ mock_client.go (MockClient for tests)
       ├─ auth/
       │    ├─ keychain.go   (KeychainStore interface, OSKeychain, MockKeychain)
       │    ├─ oauth.go      (OAuthFlow)
       │    └─ token.go      (TokenManager)
       ├─ config/
       │    └─ config.go     (Config load/save, account management)
       └─ models/
            └─ models.go     (all data models)
```

## Design Decisions

### Interface-driven API client

`api.Client` is the central abstraction. All commands depend on the interface, not `GraphClient` directly. This enables:

- **MockClient** for unit tests without HTTP
- **httptest.Server** for integration tests
- Potential future implementations (caching, logging wrappers)

### Injectable factories

`cmd/helpers.go` exposes `clientFactory` and `keychainFactory` variables. Tests override these to inject mocks without touching the OS keychain or real API.

### Token transparency

Commands call `getClient()` which internally:
1. Loads config to find the active account
2. Gets the token via `TokenManager.GetValidToken()` (auto-refreshes if near expiry)
3. Returns a configured `api.Client`

No command ever deals with auth directly.

### Credential storage

- **App ID** → `~/.ig-cli/config.yaml` (not sensitive, appears in OAuth URLs)
- **App Secret** → OS keychain via `go-keyring`
- **Access tokens** → OS keychain
- **Token expiry** → OS keychain (as Unix timestamp)

### Auto-pagination

`autoPaginate[T]()` is a generic function that follows `paging.next` URLs. The `--limit` flag caps results by stopping pagination early.

### JSON-only output

- stdout: pretty-printed JSON data
- stderr: JSON error objects, debug logs via `slog`

This makes output composable with `jq` and other CLI tools.
