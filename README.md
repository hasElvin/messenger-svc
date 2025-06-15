# Messenger Service 📨

A simple message-sending microservice built with Go, Gin, PostgreSQL, Redis, and Docker.  
It supports sending SMS-style messages to a webhook, stores metadata in Redis, and saves messages in PostgreSQL.

---

## 🚀 Features

- Auto-sender that checks DB and sends unsent messages on interval
- Redis caching of message delivery metadata
- REST API to control auto-sender
- PostgreSQL-backed message persistence
- Automatic seeding of sample data for testing

---

## 📦 Architecture

This project follows a **hexagonal (clean) architecture** style:
- `core/ports`: Interfaces (use cases)
- `core/services`: Business logic
- `internal/adapters`: External dependencies (DB, Redis, HTTP webhook)

---

## 🧪 Sample Test Data

On startup, the app inserts **10 sample messages** into the database so you can test endpoints right away.

---

## 📚 API Endpoints

| Method | Endpoint     | Description                  |
|--------|--------------|------------------------------|
| GET    | `/ping`      | Health check                 |
| POST   | `/start`     | Start auto-sender            |
| POST   | `/stop`      | Stop auto-sender             |
| GET    | `/sent`      | List all sent messages       |

---

## 🐳 Run with Docker Compose

> Make sure you have Docker + Docker Compose installed

---

## 🪜 Step by step

1. Clone the repository

```bash
git clone https://github.com/hasElvin/messenger-svc.git
cd messenger-svc
```

2. If you're using custom credentials or ports, edit your .env file or config/config.go.
   By default, the app connects to:
- `localhost:5432`: PostgreSQL
- `localhost:6379`: Redis

No API key is required for webhook usage unless you configure one.

3. Build and run the containers
```bash
docker compose up --build
```
On first run, this will:
- Build the Go application
- Start PostgreSQL and Redis
- Auto-migrate the DB schema
- Insert 10 sample messages
- Launch the API server on `localhost:8080`
- Start the message sender automatically

4. Redis cache check
```bash
docker compose exec redis redis-cli
127.0.0.1:6379> KEYS *
127.0.0.1:6379> GET msg:1
```
After messages are sent, their delivery metadata is stored in Redis.
Each key looks like `msg:<id>` with a corresponding value
