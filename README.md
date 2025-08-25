# Echoes — AI-Powered Emotional Journal & Memory Assistant

Echoes is a personal journaling and “memory companion” that helps you capture moments, reflect on emotions, and surface meaningful memories over time.  
It’s designed as a clean, extensible backend (Go) with a simple REST API you can use from any client (mobile, web, CLI). Optional AI features enable semantic search, summaries, and mood insights.

> Goal: be a portfolio-quality project that’s small enough to ship, but real enough to show professional practices (auth, persistence, testing, Docker, CI).

---

## Features

- **Entries CRUD:** create/read/update/delete journal entries with `title`, `body`, `mood`, `tags`, `timestamp`.
- **Auth & Sessions:** JWT-based auth (register/login/refresh); ready for frontend consumption.
- **Search:** keyword + (optional) semantic search over your entries.
- **Insights (optional AI):**
  - auto-summary and tone/mood extraction
  - “memory recall” suggestions (e.g., *this week last year*)
- **Attachments (optional):** link images/audio via URL; store metadata alongside entries.
- **Portable DB:** SQLite by default for zero-setup; Postgres ready for production.
- **Clean structure:** idiomatic Go modules, handlers/routes/models split, env-driven config.
- **Docker-ready:** one command to run everything locally.

---

## Tech Stack

- **Backend:** Go (Gin or net/http; idiomatic handlers, services)
- **Auth:** JWT (HS256)
- **DB:** SQLite (dev) / Postgres (prod) through GORM or database/sql
- **Embeddings/AI (optional):** pluggable provider (e.g., OpenAI or local), behind a simple interface

> **Note:** The repo is intentionally backend-first. Any client (React/Flutter/CLI) can talk to it over HTTP.

---

## Quickstart (Local)

### 1) Clone & setup
```bash
git clone https://github.com/CBYeuler/echoes.git
cd echoes
go mod tidy
```
# Server
APP_ENV=dev
APP_PORT=8080

# Auth
JWT_SECRET=replace-with-a-long-random-string
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h  # 7 days

# Database (choose ONE)
DB_DRIVER=sqlite
DB_DSN=./echoes.db           # used when DB_DRIVER=sqlite
# DB_DRIVER=postgres
# DB_DSN=host=localhost user=postgres password=postgres dbname=echoes port=5432 sslmode=disable TimeZone=UTC

# AI / Embeddings (optional)
AI_PROVIDER=none            # none|openai|local
OPENAI_API_KEY=             # if AI_PROVIDER=openai

# Feature flags
FF_EMBEDDINGS=false
FF_SUMMARIES=false
FF_SUGGESTIONS=false


go run ./cmd/server
# Server listening on :8080 (by default)

Optional Docker:
# Build
docker build -t echoes:dev .

# Run (SQLite volume persisted to ./data)
docker run --env-file ./.env -p 8080:8080 -v $(pwd)/data:/app/data echoes:dev

API Overview:
POST   /api/v1/auth/register     { email, password, name }
POST   /api/v1/auth/login        { email, password } -> { access_token, refresh_token }
POST   /api/v1/auth/refresh      { refresh_token } -> { access_token }
GET    /api/v1/auth/me           (Bearer) -> current user

Entries:
GET    /api/v1/entries           ?q=keyword&tag=...&limit=...&offset=...
POST   /api/v1/entries           { title, body, mood?, tags?[], created_at? }
GET    /api/v1/entries/:id
PUT    /api/v1/entries/:id
DELETE /api/v1/entries/:id

Entry model(baseline):
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "string",
  "body": "string",
  "mood": "very_low|low|neutral|high|very_high",
  "tags": ["life","study","relationship"],
  "created_at": "RFC3339",
  "updated_at": "RFC3339",
  "ai": {
    "summary": "optional short abstract",
    "keywords": ["optional","keywords"],
    "embedding_id": "optional-reference"
  }
}

Insights (optional AI):
POST   /api/v1/entries/:id/summarize     -> { summary }
POST   /api/v1/entries/:id/embed         -> { embedding_id }
GET    /api/v1/entries/similar?id=...    -> k-NN like results (if embeddings enabled)
GET    /api/v1/suggestions/recall        -> “On this day”, “Last week”, “Last year”

Endpoints may be toggled by feature flags.

Minimal Client Example (cURL):

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret123","name":"Batuhan"}'

# Login (get tokens)
TOKENS=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret123"}')

ACCESS=$(echo $TOKENS | jq -r '.access_token')

# Create an entry
curl -X POST http://localhost:8080/api/v1/entries \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Morning walk & sunlight",
    "body":"Felt grounded. Wrote for an hour under the camellia.",
    "mood":"high",
    "tags":["routine","gratitude"]
  }'
Project Structure:
echoes/
  cmd/
    server/
      main.go            # wire up routes, config, logger
  internal/
    config/              # env, flags, loader
    http/
      middleware/        # auth, logging, recover, cors
      handlers/          # entry, auth, health
      routes/            # register routes
    core/
      models/            # domain structs
      repo/              # db access (sqlite/postgres)
      services/          # business logic, AI adapters
      search/            # keyword/semantic search
    pkg/
      jwt/               # token utils
      logger/            # zap/logging
      respond/           # response helpers
  migrations/            # sql migrations (if using)
  .env.example
  Dockerfile
  docker-compose.yml     # optional (api + postgres)
  README.md

Development Tasks

 Auth: register/login/refresh + middleware

 Entries: CRUD + filtering (q, tags, date range)

 Search: LIKE or FTS for SQLite; ILIKE/TS for Postgres

 AI adapters: interface + provider(s) (OpenAI/local)

 Semantic search: store embeddings (SQLite table or vector extension / Postgres + pgvector)

 Suggestions: “On this day”, spaced-repetition recall

 Tests: unit tests for handlers/services; integration tests with SQLite

 Observability: request logs, error tracking, health/readiness

 CI: go vet, golangci-lint, tests on PR

 Docker: small image, production en

Environment Variables:
| Key               | Example                       | Notes                            |                                          |             |
| ----------------- | ----------------------------- | -------------------------------- | ---------------------------------------- | ----------- |
| `APP_PORT`        | `8080`                        | API port                         |                                          |             |
| `JWT_SECRET`      | `super-long-random-string`    | keep secret                      |                                          |             |
| `JWT_ACCESS_TTL`  | `15m`                         | access token lifetime            |                                          |             |
| `JWT_REFRESH_TTL` | `168h`                        | refresh token lifetime           |                                          |             |
| `DB_DRIVER`       | `sqlite` or `postgres`        | choose your DB                   |                                          |             |
| `DB_DSN`          | `./echoes.db` or Postgres DSN | connection string                |                                          |             |
| `AI_PROVIDER`     | \`none                        | openai                           | local\`                                  | optional AI |
| `OPENAI_API_KEY`  | `sk-...`                      | required if `AI_PROVIDER=openai` |                                          |             |
| `FF_*`            | \`true                        | false\`                          | feature flags (embeddings, summaries...) |             |


Disclaimer

This README describes the intended shape of Echoes as a clean, demonstrable backend with optional AI. If your local code differs, adjust endpoints/flags accordingly. Ship a thin vertical slice first (Auth + Entries CRUD), then turn on AI features behind flags.
License
MIT © 2025 Cem Batuhan Yaman
