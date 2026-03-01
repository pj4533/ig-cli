# Authentication Guide

## Overview

ig-cli uses the Instagram Graph API OAuth flow to authenticate. The flow produces a **long-lived access token** (valid for 60 days) that is stored securely in the OS keychain.

## Token Lifecycle

1. **Authorization Code** — User authorizes in browser, Instagram redirects with a code
2. **Short-lived Token** — Code is exchanged for a 1-hour token via `POST /oauth/access_token`
3. **Long-lived Token** — Short token is exchanged for a 60-day token via `GET /oauth/access_token?grant_type=ig_exchange_token`
4. **Auto-refresh** — When a token is within 24 hours of expiry, `ig` automatically refreshes it

## Setup Flow

### 1. Create Meta Developer App

See [README prerequisites](../README.md#prerequisites).

### 2. Configure Credentials

```bash
ig auth setup
```

This stores:
- **App ID** in `~/.ig-cli/config.yaml` (not sensitive — appears in OAuth URLs)
- **App Secret** in OS keychain (sensitive — used for token exchange)

### 3. Connect Account

```bash
ig auth add
```

This:
1. Starts a local HTTP server on port 8080
2. Opens the Instagram authorization page in your browser
3. After you authorize, Instagram redirects to `http://localhost:8080/callback`
4. The CLI exchanges the code for a long-lived token
5. Token and user profile are stored in the keychain

### 4. Verify

```bash
ig auth list
```

Shows connected accounts with token expiry dates.

## Required Scopes

- `instagram_basic` — Read profile and media
- `instagram_manage_insights` — Read insights and audience data
- `pages_show_list` — Required for business discovery
- `pages_read_engagement` — Required for engagement metrics

## Troubleshooting

### "Invalid redirect_uri"

Make sure `http://localhost:8080/callback` is listed in your Meta App's "Valid OAuth Redirect URIs".

### "Port 8080 already in use"

Another process is using port 8080. Stop it before running `ig auth add`.

### Token expired

Tokens auto-refresh within 24 hours of expiry. If a token has fully expired (>60 days), run `ig auth add` again.

## Multi-account Support

You can connect multiple accounts:

```bash
ig auth add  # Connect first account (becomes default)
ig auth add  # Connect second account

ig media list                    # Uses default account
ig media list --account user2    # Uses specific account
```
