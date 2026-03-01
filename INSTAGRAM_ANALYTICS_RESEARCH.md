# Instagram Analytics Research

> Research compiled February 28, 2026

## Table of Contents

- [Executive Summary](#executive-summary)
- [The API Landscape](#the-api-landscape)
- [Instagram Graph API Deep Dive](#instagram-graph-api-deep-dive)
- [Authentication & Setup](#authentication--setup)
- [What You CAN Do (Official API)](#what-you-can-do-official-api)
- [What You CANNOT Do (Official API)](#what-you-cannot-do-official-api)
- [Existing CLI Tools & Libraries](#existing-cli-tools--libraries)
- [Third-Party Commercial Platforms](#third-party-commercial-platforms)
- [Unofficial / Private API Tools](#unofficial--private-api-tools)
- [Recommendation: Build vs Buy](#recommendation-build-vs-buy)
- [Sources](#sources)

---

## Executive Summary

**Bottom line: We should build a CLI.** Here's why:

1. The **Instagram Graph API** (official) gives you solid access to your own accounts' analytics — post metrics, comments, commenter usernames, reach, impressions, audience demographics — all without app review, using Development Mode.

2. The one big gap: **you cannot get a list of individual users who liked your posts** through the official API. You only get aggregate like counts. This has been a deliberate privacy restriction since 2018.

3. **No good CLI tool exists today.** The space is littered with abandoned projects. The few that use the official API have 1-12 GitHub stars and are unmaintained. The popular tools (Instaloader, Instagrapi, Osintgram) all use unofficial scraping/private APIs and risk account bans.

4. The **Development Mode loophole** is key: you can add your own accounts as test users on your Meta App and access all analytics permissions without going through Meta's notoriously difficult App Review process. Perfect for managing your own set of accounts.

5. Rate limit: 200 API calls per hour per Instagram account. For personal analytics use, this is more than sufficient.

---

## The API Landscape

### What Exists Today (February 2026)

| API | Status | Account Type | Use Case |
|-----|--------|-------------|----------|
| Instagram Graph API | **Active** | Business/Creator only | Full analytics, comments, publishing |
| Instagram API with Instagram Login | **Active** (since July 2024) | Business/Creator only | Simplified OAuth without Facebook Page |
| Instagram Basic Display API | **Dead** (December 4, 2024) | Was for Personal accounts | N/A — shut down |
| Facebook Graph API | **Active** | Via linked Facebook Page | Backend for Instagram Graph API |

### Key Timeline of Changes

- **July 2024**: Instagram API with Instagram Login launched — no Facebook Page required for auth
- **December 4, 2024**: Basic Display API permanently shut down — personal accounts lost all API access
- **January 2025**: Several profile metrics deprecated (profile_views, website_clicks, etc.)
- **2025**: Rate limits slashed from 5,000 to 200 calls/hour/account (no warning)
- **April 2025**: Graph API v22 — `views` replaces `impressions`/`plays` as unified metric

### The Fundamental Rule

Everything goes through `graph.facebook.com`. Instagram does not have a separate, independent API. All requests hit:
```
GET https://graph.facebook.com/v22.0/{node-id}?fields={fields}&access_token={token}
```

You need a **Meta Developer App** regardless of which auth path you choose.

---

## Instagram Graph API Deep Dive

### What Data Is Available Per Post

```
GET /{ig-media-id}?fields=id,caption,media_type,media_url,permalink,timestamp,like_count,comments_count
```

Fields:
- `id` — media object ID
- `caption` — post caption text
- `media_type` — IMAGE, VIDEO, CAROUSEL_ALBUM
- `media_url` — direct URL to media
- `permalink` — Instagram URL
- `timestamp` — ISO 8601 timestamp
- `like_count` — aggregate number (NOT individual likers)
- `comments_count` — aggregate number

### Post-Level Insights (Performance Metrics)

```
GET /{ig-media-id}/insights?metric=reach,impressions,engagement,saves,shares,views
```

Available metrics per post:
- `reach` — unique accounts that saw it
- `views` — total views (replaces `impressions` for content after July 2, 2024)
- `engagement` — total interactions
- `likes` — like count
- `comments` — comment count
- `saves` — save count
- `shares` — share count
- `video_views` — (Reels only, deprecated in v22 in favor of `views`)

### Comments & Commenter Data

```
GET /{ig-media-id}/comments?fields=id,text,username,timestamp,like_count,replies
```

This is **fully available** for your own posts:
- `text` — the actual comment content
- `username` — who wrote it
- `timestamp` — when it was posted
- `like_count` — likes on the comment
- `replies` — nested reply objects with same fields

You can also reply to, hide, and delete comments via the API.

**Limitation:** You get the commenter's `username` but NOT their full profile (bio, follower count, etc.). To look up a commenter's public profile, you'd use the `business_discovery` endpoint separately.

### Account-Level Insights

```
GET /{ig-user-id}/insights?metric=reach,views,follower_count&period=day
```

- `reach` — unique accounts reached (day/week/28-day periods)
- `views` — total content views
- `follower_count` — net follower change per day
- `follows_count` — total followers (lifetime metric)

### Audience Demographics

Available through account insights:
- Age ranges (13-17, 18-24, 25-34, 35-44, 45-54, 55-64, 65+)
- Gender distribution
- Top cities (up to 45)
- Top countries (up to 45)

**Requires 100+ followers** to return demographic data.

### Looking Up Other Public Accounts

```
GET /{your-ig-user-id}?fields=business_discovery.fields(username,followers_count,media_count,media{like_count,comments_count,caption,timestamp})&username={target-username}
```

For any public Business/Creator account you can see:
- Follower/following counts
- Post count
- Their recent posts' like counts, comment counts, captions, timestamps

You CANNOT see their insights (reach, impressions) — only public-facing metrics.

### Rate Limits

- **200 API calls per hour per Instagram account**
- All requests count (including failures and pagination)
- Hashtag searches: 30 unique hashtags per 7-day rolling window
- Content publishing: 25 posts/day max
- Monitor via `X-Business-Use-Case-Usage` response header

---

## Authentication & Setup

### Prerequisites

1. Instagram account(s) must be **Business or Creator** type (free to switch in app settings)
2. Create a **Meta Developer Account** at developers.facebook.com
3. Create a **Meta App** (type: Business)
4. Add **Instagram Graph API** as a product

### Two Auth Paths

**Path 1: Instagram Login (Simpler, New)**
- Direct Instagram OAuth — no Facebook Page needed
- Scopes: `instagram_business_basic`, `instagram_business_manage_comments`, `instagram_manage_insights`
- Short-lived token → exchange for 60-day long-lived token

**Path 2: Facebook Login (Legacy, Better for Multi-Account)**
- Authenticate via Facebook Page linked to IG account
- More complex but supports Business Manager for agency-style multi-account
- Scopes: `pages_show_list`, `instagram_basic`, `instagram_manage_insights`, `pages_read_engagement`

### The Development Mode Shortcut (Key for Your Use Case)

**You do NOT need App Review for your own accounts.**

In Development Mode:
- Add your Instagram accounts as test users/developers in your Meta App
- All permissions work immediately — insights, comments, post data, demographics
- No App Review submission required
- No business verification required
- Only limitation: app cannot be used by accounts you haven't manually added

This is the practical path for managing your own accounts. App Review is only needed if you want to distribute your tool to other users.

### App Review (If You Ever Need It)

Meta's App Review process is notoriously difficult:
- Requires screen recordings demonstrating each permission's use
- Requires published privacy policy on a live website
- 2-7 business day review cycle per submission
- Rejections are common; some developers report 14+ submission attempts
- Business verification is a separate additional hurdle
- One developer community described it as "so damn hard to get approved"

**For personal use: skip this entirely by staying in Development Mode.**

---

## What You CAN Do (Official API)

For your own Business/Creator accounts in Development Mode:

| Capability | Endpoint | Notes |
|-----------|----------|-------|
| List all your posts | `GET /{user-id}/media` | Paginated, all post types |
| Post metrics (reach, views, engagement) | `GET /{media-id}/insights` | Per-post performance |
| Like count per post | `GET /{media-id}?fields=like_count` | Aggregate number only |
| Comment count per post | `GET /{media-id}?fields=comments_count` | Aggregate number |
| All comment text + commenter usernames | `GET /{media-id}/comments` | Full text, username, timestamp |
| Reply to comments | `POST /{media-id}/comments` | Programmatic replies |
| Account reach/views over time | `GET /{user-id}/insights` | Day/week/28-day periods |
| Follower growth tracking | `GET /{user-id}/insights?metric=follower_count` | Net change per day |
| Audience demographics | `GET /{user-id}/insights` | Age, gender, city, country |
| Story metrics | `GET /{media-id}/insights` | Views, reach, replies (while live) |
| Competitor public metrics | `business_discovery` endpoint | Followers, post counts, like/comment counts |
| Manage multiple accounts | Via test users in Dev Mode | Each account = separate auth |

---

## What You CANNOT Do (Official API)

| Capability | Status | Alternative |
|-----------|--------|-------------|
| **List individual likers** (who liked a post) | Not available | Private API / scraping (risky) |
| **List followers/following** (enumerable) | Not available | Private API / scraping (risky) |
| Competitor's reach/impressions | Not available | Only public metrics via `business_discovery` |
| Historical data beyond API windows | Not stored by API | Must build your own DB to accumulate |
| Personal account access | Impossible since Dec 2024 | Must convert to Business/Creator |
| Story viewer list | Not available | Only aggregate metrics |
| DM content (without Business Messaging) | Not available | Private API (risky) |

---

## Existing CLI Tools & Libraries

### Tools Using the Official API (Recommended Approach)

These are safe but all have very small communities:

| Tool | Language | Stars | Last Updated | Notes |
|------|----------|-------|-------------|-------|
| [gramoco-cli](https://github.com/alexmarqs/gramoco-cli) | TypeScript | 1 | May 2025 | Extracts posts/comments to Excel |
| [instagram-insights](https://github.com/PardhuMadipalli/instagram-insights) | Python | 12 | Mar 2021 | Best time to post, hashtag analytics |
| [ig-mcp](https://github.com/jlbadano/ig-mcp) | — | 79 | Recent | MCP server for AI apps, official API |
| [instagram-analytics-mcp](https://github.com/BilalTariq01/instagram-analytics-mcp) | — | 2 | Feb 2026 | MCP server for insights |
| [SocialMediaAnalytics](https://github.com/ckoutavas/SocialMediaAnalytics) | Python | 4 | Mar 2023 | Returns Pandas DataFrames |

**Verdict:** Nothing production-ready exists. The official-API CLI space is wide open.

### Popular Unofficial Tools (Account Ban Risk)

| Tool | Language | Stars | Method | Status | Risk |
|------|----------|-------|--------|--------|------|
| [Osintgram](https://github.com/Datalux/Osintgram) | Python | 12,400 | Scraping | Unmaintained (2021) | High |
| [Instaloader](https://github.com/instaloader/instaloader) | Python | 11,700 | Scraping | Active (v4.15) | Moderate-High |
| [dilame/instagram-private-api](https://github.com/dilame/instagram-private-api) | Node.js | 6,400 | Private API | Partly commercialized | High |
| [Instagrapi](https://github.com/subzeroid/instagrapi) | Python | 5,900 | Private API | Active (v2.3.0) | High |
| [instagram-cli](https://github.com/supreme-gg-gg/instagram-cli) | TS/Python | 1,600 | Private API | Active | High — TUI client, not analytics |
| [instagram_monitor](https://github.com/misiektoja/instagram_monitor) | Python | 787 | Scraping | Dec 2024 | Moderate-High |

**Ban evidence is well-documented:**
- Instaloader GitHub Issue #2555: "IG threatening to ban account because of instaloader activity"
- Instaloader GitHub Issue #1937: "I have now lost my second account on Instagram for using this script"

### Compiled Language Tools (Rust/Go/Swift)

All are either tiny, archived, or abandoned:
- `instagram-scraper-rs` (Rust, 18 stars, archived)
- `insta-tools` (Go, 2 stars, cookie-based)
- `SwiftInstagram` (Swift, 576 stars, archived 2018 — used dead API)

---

## Third-Party Commercial Platforms

If you don't want to build, these are the paid alternatives:

| Platform | Price | Strengths |
|----------|-------|-----------|
| **Iconosquare** | $33-120/mo | Most Instagram-specific; historical data, benchmarking |
| **Sprout Social** | $249+/mo/user | Enterprise-grade; unified paid+organic reporting |
| **Hootsuite** | $99+/mo | Multi-platform; large team management |
| **Later** | $18-80/mo | Scheduling-focused with analytics |
| **Buffer** | Free tier available | Indie-friendly; basic analytics |
| **AgencyAnalytics** | $59-179/mo | White-label client reports; 80+ integrations |

**What they provide that the raw API doesn't:**
- Historical data persistence (they poll the API continuously and store results)
- Competitor benchmarking
- Trend visualization over arbitrary time ranges
- Best-time-to-post analysis
- Automated reporting / PDF exports

**What they still can't do:**
- Show individual likers (same API limitation applies)
- Access personal accounts

---

## Unofficial / Private API Tools

### Instagrapi (Most Feature-Rich)

The [instagrapi](https://github.com/subzeroid/instagrapi) Python library reverse-engineers Instagram's mobile app API. It can access things the official API cannot:
- Individual likers list
- Follower/following lists
- Story viewer lists
- Insights (even via private API)
- Direct messages

**Risks:**
- Account bans are common without proxies and careful rate limiting
- Repo warns: "more suits for testing or research than a working business"
- Instagram regularly changes its private API, breaking the library
- Clear ToS violation

### Commercial Private API Services

**HikerAPI** (hikerapi.com) — built on top of instagrapi:
- Industrialized private API access
- 4-5 million requests/day
- Pricing: $0.02/request (start) to $0.0006/request (ultra)
- They absorb the ban risk on their infrastructure

### Scraping Legal Status

- **hiQ v. LinkedIn** (Ninth Circuit): Scraping publicly available data does not violate CFAA
- **BUT**: Meta has pursued breach of contract (ToS) claims against logged-in scrapers
- **Practical risk**: Account termination, not prosecution
- **Key distinction**: If you need to log in to access the data, scraping is legally riskier

---

## Recommendation: Build vs Buy

### Build a CLI — Here's Why

1. **The official API gives you 80% of what you want** — post metrics, comment text + commenters, reach, engagement, audience demographics, multi-account support.

2. **The missing 20% (individual likers) is a hard gap** that even commercial platforms can't fill through legitimate means.

3. **No good CLI exists.** The space has literally zero well-maintained official-API CLI tools. Everything is either dead, tiny, or uses risky unofficial methods.

4. **Development Mode makes setup easy.** No app review needed for your own accounts.

5. **You can augment over time.** Start with official API, optionally layer in `business_discovery` for competitor public data, and build your own historical database.

### Suggested Architecture

```
ig-cli/
├── cmd/                    # CLI commands (Swift or Go)
│   ├── posts.swift         # List posts with metrics
│   ├── comments.swift      # Fetch/analyze comments
│   ├── insights.swift      # Account-level insights
│   ├── audience.swift      # Demographics
│   └── competitors.swift   # Public competitor data
├── api/
│   ├── graph_api.swift     # Facebook Graph API client
│   ├── auth.swift          # Token management (long-lived tokens)
│   └── models.swift        # Data models
├── storage/
│   └── database.swift      # Local SQLite for historical data accumulation
└── config/
    └── accounts.swift      # Multi-account config management
```

### Setup Steps

1. Convert Instagram accounts to Business or Creator (free, in app settings)
2. Create Meta Developer App at developers.facebook.com
3. Add Instagram Graph API product
4. Add your accounts as test users (Development Mode)
5. Generate long-lived access tokens (60-day, auto-refreshable)
6. Start pulling data

### What the CLI Could Do (Day 1)

- `ig posts list` — show recent posts with like/comment/reach counts
- `ig posts insights <post-id>` — detailed metrics for a post
- `ig comments list <post-id>` — all comments with usernames and text
- `ig comments analyze <post-id>` — sentiment/engagement analysis on comments
- `ig insights account` — account-level reach, views, follower changes
- `ig audience demographics` — age, gender, location breakdowns
- `ig competitors check <username>` — public metrics for competitor accounts

### Historical Data Strategy

The API only returns rolling windows. To build trend analysis:
- Poll daily and store results in local SQLite
- Track follower count deltas over time
- Accumulate per-post metrics at regular intervals
- This is exactly what commercial platforms charge $33-249/month for

---

## Sources

### Official Documentation
- [Instagram Graph API Documentation (via Meta)](https://developers.facebook.com/docs/instagram-api/)
- [Instagram Platform API with Instagram Login Guide](https://gist.github.com/PrenSJ2/0213e60e834e66b7e09f7f93999163fc)

### API Guides & Analysis
- [Instagram Graph API: Complete Developer Guide for 2026 — Elfsight](https://elfsight.com/blog/instagram-graph-api-complete-developer-guide-for-2026/)
- [How to Use the Instagram Graph API for Audience Insights — Phyllo](https://www.getphyllo.com/post/how-to-use-the-instagram-graph-api-for-audience-insight-iv)
- [Instagram Graph API: Who is it good for? — Phyllo](https://www.getphyllo.com/post/instagram-graph-api-who-is-it-good-for)
- [Instagram Analytics API: Comprehensive Overview — Data365](https://data365.co/blog/instagram-analytics-api)
- [Getting and Replying to Comments on Instagram with Graph API — Justin Stolpe](https://justinstolpe.com/blog_code/instagram_graph_api/comments_and_replies.php)

### API Changes & Deprecations
- [Instagram Basic Display API Deprecation — WPZOOM](https://www.wpzoom.com/documentation/instagram-widget/basic-display-api-deprecation/)
- [Instagram Insights Metrics Deprecation (January 2025) — Emplifi](https://docs.emplifi.io/platform/latest/home/instagram-media-and-profile-insights-metrics-depre)
- [Instagram Insights Metrics Deprecation (April 2025) — Emplifi](https://docs.emplifi.io/platform/latest/home/instagram-insights-metrics-deprecation-april-2025)
- [Instagram Rate Limits Deep Dive — Marketing Scoop](https://www.marketingscoop.com/marketing/instagrams-api-rate-limits-a-deep-dive-for-developers-and-marketers-in-2024/)
- [Meta Graph API v22.0 Updates — Swipe Insight](https://web.swipeinsight.app/posts/facebook-launches-graph-api-v22-0-and-marketing-api-v22-0-for-developers-14179)

### App Review & Developer Experience
- [Why it is so damn hard to get approved for Meta Graph API — GraphAPI Substack](https://graphapi.substack.com/p/why-it-is-so-damn-hard-to-get-approved)
- [Has anyone managed to pass Meta's Access Verification? — Hacker News](https://news.ycombinator.com/item?id=43895515)
- [Meta App Approval Guide — Saurabh Dhar](https://www.saurabhdhar.com/blog/meta-app-approval-guide)
- [Instagram App Review — Chatwoot Developer Docs](https://developers.chatwoot.com/self-hosted/instagram-app-review)

### GitHub Tools (Official API)
- [gramoco-cli](https://github.com/alexmarqs/gramoco-cli) — TypeScript CLI for Graph API
- [ig-mcp](https://github.com/jlbadano/ig-mcp) — MCP server for Instagram Graph API
- [instagram-analytics-mcp](https://github.com/BilalTariq01/instagram-analytics-mcp) — MCP analytics server
- [instagram-insights](https://github.com/PardhuMadipalli/instagram-insights) — Python analytics CLI

### GitHub Tools (Unofficial / Scraping)
- [Instaloader](https://github.com/instaloader/instaloader) — Python scraper (11.7k stars)
- [Instagrapi](https://github.com/subzeroid/instagrapi) — Python private API (5.9k stars)
- [Osintgram](https://github.com/Datalux/Osintgram) — OSINT tool (12.4k stars, unmaintained)
- [instagram-private-api (Node.js)](https://github.com/dilame/instagram-private-api) — Node.js private API (6.4k stars)
- [instagram_monitor](https://github.com/misiektoja/instagram_monitor) — Real-time monitoring

### Scraping Risks & Legal
- [Is Instagram Scraping Legal? 2025 Guide — SociaVault](https://sociavault.com/blog/instagram-scraping-legal-2025)
- [hiQ v. LinkedIn — Apify Blog](https://blog.apify.com/hiq-v-linkedin/)
- [Instaloader Ban Reports — GitHub Issues #2555, #1937](https://github.com/instaloader/instaloader/issues/2555)

### Commercial Platforms
- [Top 14 Instagram Analytics Tools — Sprout Social](https://sproutsocial.com/insights/instagram-analytics-tools/)
- [21 Social Media Analytics Tools — Hootsuite](https://blog.hootsuite.com/social-media-analytics-tools/)
- [7 Best Social Media Management Tools — Iconosquare](https://www.iconosquare.com/blog/best-social-media-management-tools)

### CrowdTangle (Historical Context)
- [Meta Shuts Down CrowdTangle — Social Media Today](https://www.socialmediatoday.com/news/meta-announces-shut-down-crowdtangle-monitoring-app/710358/)
- [Meta killed CrowdTangle — Engadget](https://www.engadget.com/big-tech/meta-killed-crowdtangle-an-invaluable-research-tool-because-what-it-showed-was-inconvenient-121700584.html)
