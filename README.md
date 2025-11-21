# Echoes — AI-Powered Emotional Journal & Memory Assistant



![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Build](https://img.shields.io/badge/build-passing-brightgreen?style=flat)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Docker](https://img.shields.io/badge/docker-ready-2496ED?style=flat&logo=docker)
![Issues](https://img.shields.io/github/issues/CBYeuler/echoes)
![Stars](https://img.shields.io/github/stars/CBYeuler/echoes)
![Last Commit](https://img.shields.io/github/last-commit/CBYeuler/echoes)
![Code Style](https://img.shields.io/badge/code%20style-gofmt-79D4FD.svg)

Echoes is a personal journaling and “memory companion” backend that helps you capture moments, reflect emotionally, and revisit meaningful memories over time.  
It is built with **Go**, features a clean architecture, and exposes a **REST API** that any client (mobile, web, CLI) can use.  
Optional AI abilities enable **semantic search**, **auto-summaries**, **tone analysis**, and **memory recall**.

> **Project Goal:** A portfolio-quality backend showcasing real-world engineering: authentication, persistence, search, feature flags, Docker, and clean modular code.

---

## Features

- **Journal Entries CRUD**  
  Create, read, update, delete entries with: `title`, `body`, `mood`, `tags`, timestamps.

- **Authentication & Sessions**  
  JWT-based (access + refresh), clean middleware, secure routing.

- **Search System**  
  - Keyword search  
  - Optional semantic embeddings  
  - Optional vector-based similarity search

- **AI Insights (Optional)**  
  - Auto-summary  
  - Mood/tone extraction  
  - “This week last year” / “on this day” recall suggestions

- **Attachments (Optional)**  
  Store metadata for image/audio URLs.

- **DB Layer**  
  - SQLite (default, zero-config)  
  - PostgreSQL support

- **Clean Architecture**  
  Handlers → Services → Repository → Models  
  Fully maintainable & test-friendly.

- **Docker Ready**  
  One command to run the entire stack.

---

## Tech Stack

- **Language:** Go
- **API Layer:** Gin or net/http
- **Auth:** JWT (HS256)
- **DB:** SQLite (dev), PostgreSQL (prod)
- **ORM / DB Layer:** GORM or database/sql
- **AI:** Pluggable provider (OpenAI or Local)
- **Container:** Docker & docker-compose

---

## Quickstart

### 1) Clone the repository
```bash
git clone https://github.com/CBYeuler/echoes.git
cd echoes
go mod tidy
```

## Create `.env`:
```bash
APP_ENV=dev
APP_PORT=8080

# Auth
JWT_SECRET=replace-with-a-long-random-string
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# Database
DB_DRIVER=sqlite
DB_DSN=./echoes.db
# For Postgres:
# DB_DRIVER=postgres
# DB_DSN=host=localhost user=postgres password=postgres dbname=echoes port=5432 sslmode=disable TimeZone=UTC

# AI (optional)
AI_PROVIDER=none
OPENAI_API_KEY=

# Feature Flags
FF_EMBEDDINGS=false
FF_SUMMARIES=false
FF_SUGGESTIONS=false
```

## Run the server:
```bash
go run ./cmd/server
```

### Default adress:
`http://localhost:8080`

## Run With Docker:
```bash
docker build -t echoes:dev .

docker run --env-fil ./.env -p 8080:8080 -v $(pwd)/data:/app/data echoes:dev
```

## API Overview

- Auth:
```bash
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
GET    /api/v1/auth/me
```
- Entires:
```bash
GET    /api/v1/entries
POST   /api/v1/entries
GET    /api/v1/entries/:id
PUT    /api/v1/entries/:id
DELETE /api/v1/entries/:id
```
- AI(optional):
```bash
POST   /api/v1/entries/:id/summarize
POST   /api/v1/entries/:id/embed
GET    /api/v1/entries/similar
GET    /api/v1/suggestions/recall
```


# Feature Flags
FF_EMBEDDINGS=false
FF_SUMMARIES=false
FF_SUGGESTIONS=false


### Entry Model Example:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "string",
  "body": "string",
  "mood": "very_low | low | neutral | high | very_high",
  "tags": ["life","study","relationship"],
  "created_at": "RFC3339",
  "updated_at": "RFC3339",
  "ai": {
    "summary": "optional short abstract",
    "keywords": ["optional","keywords"],
    "embedding_id": "optional"
  }
}
```

## Minimal Client Example (cURL):
- Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret123","name":"Batuhan"}'
```
- Login (retrieve access token)
```bash
TOKENS=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret123"}')

ACCESS=$(echo $TOKENS | jq -r '.access_token')
```
 - Create Entry
```bash
curl -X POST http://localhost:8080/api/v1/entries \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Morning walk & sunlight",
    "body":"Felt grounded. Wrote for an hour under the camellia.",
    "mood":"high",
    "tags":["routine","gratitude"]
  }'
```

## Environment Variables

| Key              | Example                        | Notes                                      |
|------------------|--------------------------------|--------------------------------------------|
| **APP_PORT**     | `8080`                         | API port                                   |
| **JWT_SECRET**   | `super-long-random-string`     | Keep this secret                           |
| **JWT_ACCESS_TTL** | `15m`                        | Access token lifetime                      |
| **JWT_REFRESH_TTL** | `168h`                      | Refresh token lifetime                     |
| **DB_DRIVER**    | `sqlite` or `postgres`         | Choose your DB                             |
| **DB_DSN**       | `./echoes.db` or Postgres DSN  | Connection string                          |
| **AI_PROVIDER**  | `none`, `openai`, `local`      | Optional AI provider                       |
| **OPENAI_API_KEY** | `sk-...`                     | Required if `AI_PROVIDER=openai`           |
| **FF_*`**        | `true` / `false`               | Feature flags (embeddings, summaries, etc) |

---

MIT © 2025 **Cem Batuhan Yaman**

