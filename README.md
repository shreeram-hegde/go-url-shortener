Go URL Shortener

A production-style URL shortening service built in Go using net/http with clean layered architecture, pluggable storage backends, background expiration cleanup, and environment-based configuration.

Features

REST API using Go net/http

Clean architecture: handler / service / store

Pluggable storage backends:

In-memory store

SQLite persistence

URL expiry support with background cleanup worker

Configurable via environment variables

/health endpoint for health checks

Concurrency-safe design

Architecture

HTTP (handlers) → Service (business logic) → Store (persistence)

Store is interface-based, so switching from in-memory to SQLite does not affect business logic.

Project structure:

cmd/server/main.go
internal/
  handler/
  service/
  store/
  model/

API Endpoints
1. Create Short URL

POST /shorten

Request:

{
  "url": "https://google.com",
  "expiry_minutes": 60
}


Response:

{
  "short_url": "http://localhost:8080/Ab3XkQ"
}

2. Redirect

GET /{code}

Example:

GET /Ab3XkQ → redirects to https://google.com

3. Health Check

GET /health

Response:

200 OK

Configuration

Configuration is done using environment variables (or a .env file in development).

Example .env:

PORT=8080
BASE_URL=http://localhost:8080
STORE=sqlite
SQLITE_PATH=data.db

Environment Variables
Variable	Description	Default
PORT	HTTP server port	8080
BASE_URL	Base URL used in short links	http://localhost:PORT

STORE	Storage backend: memory or sqlite	memory
SQLITE_PATH	Path to SQLite DB file	data.db
Running Locally
1. Clone repository
git clone https://github.com/your-username/go-url-shortener.git
cd go-url-shortener

2. Create .env
PORT=8080
BASE_URL=http://localhost:8080
STORE=sqlite
SQLITE_PATH=data.db

3. Run
go run ./cmd/server

Testing with curl

Create short URL:

curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://google.com","expiry_minutes":1}'


Open the returned short URL in your browser.

Design Decisions

Uses Go interfaces for storage to allow multiple backends.

In-memory store uses map + RWMutex for concurrency safety.

SQLite store uses database/sql with schema migration on startup.

Expired URLs are cleaned up periodically by a background goroutine.

Business logic is isolated from HTTP layer for testability and clean design.

Possible Improvements

Add Redis cache

Add rate limiting middleware

Add authentication

Add Docker support

Add tests and CI