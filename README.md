# ig-cli

[![CI](https://github.com/pj4533/ig-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/pj4533/ig-cli/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line tool for Instagram analytics using the Instagram Graph API (v22.0). Pull post data, comments, engagement metrics, and audience demographics directly from your terminal.

## Installation

```bash
go install github.com/pj4533/ig-cli@latest
```

Or build from source:

```bash
git clone https://github.com/pj4533/ig-cli.git
cd ig-cli
make build
```

## Prerequisites

### 1. Switch to a Business or Creator Account

Your Instagram account must be a **Business** or **Creator** account (free to switch):

1. Open Instagram app > Settings > Account
2. Tap "Switch to Professional Account"
3. Choose **Business** or **Creator**
4. Complete the setup

### 2. Create a Meta Developer App

1. Go to [developers.facebook.com](https://developers.facebook.com)
2. Click "My Apps" > "Create App"
3. Choose "Business" type
4. Name your app and click "Create"

### 3. Add Instagram Graph API

1. In your app dashboard, click "Add Products"
2. Find "Instagram Graph API" and click "Set Up"

### 4. Configure OAuth Redirect

1. Go to your app's Instagram Graph API settings
2. Under "Valid OAuth Redirect URIs", add:
   ```
   http://localhost:8080/callback
   ```
3. Save changes

### 5. Development Mode

Your app runs in **Development Mode** by default — no App Review needed for your own accounts. Add yourself as a test user:

1. Go to App Dashboard > Roles > Test Users
2. Add your Instagram account

## Quick Start

```bash
# 1. Configure your Meta App credentials
ig auth setup

# 2. Connect your Instagram account (opens browser for OAuth)
ig auth add

# 3. List your recent posts
ig media list

# 4. View insights for a specific post
ig media insights <media-id>

# 5. Check account-level analytics
ig insights account
```

## Commands

| Command | Description |
|---------|-------------|
| `ig auth setup` | Configure Meta App ID and Secret |
| `ig auth add` | Connect an Instagram account via OAuth |
| `ig auth list` | List connected accounts with token expiry |
| `ig auth remove <username>` | Disconnect an account |
| `ig media list` | List posts with like/comment counts |
| `ig media insights <media-id>` | Detailed metrics for a post |
| `ig comments list <media-id>` | All comments on a post |
| `ig comments replies <comment-id>` | Replies to a comment |
| `ig insights account` | Account-level reach, views, follower growth |
| `ig insights audience` | Demographics (age, gender, city, country) |
| `ig discover <username>` | Look up a public Business/Creator account |
| `ig completion [bash\|zsh\|fish]` | Generate shell completions |

### Global Flags

| Flag | Description |
|------|-------------|
| `--account, -a` | Select which connected account to use |
| `--verbose, -v` | Enable debug logging to stderr |

### Per-Command Flags

| Flag | Commands | Description |
|------|----------|-------------|
| `--limit` | `media list`, `comments list`, `comments replies` | Cap number of results |
| `--period` | `insights account` | Time period: `day`, `week`, or `days_28` |

## Output

All data is output as pretty-printed JSON to stdout. Errors are output as JSON to stderr. This makes it easy to pipe into `jq` for further processing:

```bash
# Get the top 5 posts by likes
ig media list | jq 'sort_by(.like_count) | reverse | .[0:5]'

# Get all comment usernames on a post
ig comments list <media-id> | jq '.[].username'
```

## Configuration

Configuration is stored at `~/.ig-cli/config.yaml`. Sensitive credentials (App Secret, tokens) are stored in the OS keychain.

## Rate Limits

The Instagram Graph API has rate limits. The CLI monitors usage and warns at 80% capacity. See [Meta's rate limiting docs](https://developers.facebook.com/docs/graph-api/overview/rate-limiting/) for details.

## Documentation

- [Authentication Guide](docs/authentication.md)
- [API Reference](docs/api-reference.md)
- [Architecture](docs/architecture.md)
- [Contributing](docs/contributing.md)

## License

MIT - see [LICENSE](LICENSE)
