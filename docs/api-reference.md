# API Reference

Every CLI command maps to one or more Instagram Graph API endpoints.

## Media

### `ig media list`

Lists media posts for the authenticated user.

- **Endpoint:** `GET /{user-id}/media`
- **Fields:** `id,caption,media_type,media_url,permalink,timestamp,like_count,comments_count`
- **Flags:** `--limit` (cap results)
- **Pagination:** Auto-follows `paging.next` URLs

### `ig media insights <media-id>`

Gets insight metrics for a specific post.

- **Endpoint:** `GET /{media-id}/insights`
- **Metrics:** `impressions,reach,engagement,saved,video_views,likes,comments,shares`

## Comments

### `ig comments list <media-id>`

Lists comments on a media post.

- **Endpoint:** `GET /{media-id}/comments`
- **Fields:** `id,text,username,timestamp,like_count`
- **Flags:** `--limit`

### `ig comments replies <comment-id>`

Lists replies to a specific comment.

- **Endpoint:** `GET /{comment-id}/replies`
- **Fields:** `id,text,username,timestamp,like_count`
- **Flags:** `--limit`

## Insights

### `ig insights account`

Gets account-level insight metrics.

- **Endpoint:** `GET /{user-id}/insights`
- **Metrics:** `impressions,reach,profile_views,website_clicks,follower_count,email_contacts,phone_call_clicks,text_message_clicks,get_directions_clicks`
- **Flags:** `--period` (day, week, days_28)

### `ig insights audience`

Gets audience demographic data.

- **Endpoint:** `GET /{user-id}/insights`
- **Metrics:** `audience_city,audience_country,audience_gender_age,audience_locale`
- **Period:** `lifetime` (fixed)

## Discovery

### `ig discover <username>`

Looks up a public Business/Creator account.

- **Endpoint:** `GET /{user-id}?fields=business_discovery.fields(...){username}`
- **Fields:** `id,username,name,biography,followers_count,media_count,profile_picture_url,website`

## Authentication

### `ig auth setup`

Stores App ID in config and App Secret in keychain. No API calls.

### `ig auth add`

Runs the OAuth flow:

1. `POST https://api.instagram.com/oauth/access_token` — Exchange code for short-lived token
2. `GET https://graph.facebook.com/v22.0/oauth/access_token?grant_type=ig_exchange_token` — Get long-lived token
3. `GET https://graph.instagram.com/v22.0/me?fields=id,username,name` — Get user profile

### `ig auth list`

Reads from local config and keychain. No API calls.

### `ig auth remove <username>`

Removes from config and keychain. No API calls.
