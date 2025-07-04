# Messenger Service 📨

A simple message-sending microservice built with Go, Gin, PostgreSQL, Redis, and Docker.  
It supports sending SMS-style messages to a webhook, stores metadata in Redis, and saves messages in PostgreSQL.

---

## 🚀 Features

- Auto-sender that checks DB and sends unsent messages on interval
- Redis caching of message delivery metadata
- REST API to control auto-sender
- Backoff strategy after pre-defined retry limit
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

On startup, the app inserts **10 sample messages** into the database so you can test endpoints right away (please note that message sender will also auto-start handling new / pending messages)

---

## 📚 API Endpoints

### 🔗 Interactive Documentation
Test all endpoints using interactive Swagger UI: **[API Documentation](https://messenger-svc-gfsy.onrender.com/docs/index.html)**

Note: no security measures or middlewares have been implemented to make the development and testing processes easier (can be added in the next release upon request)

### Available Endpoints
To send curl or postman requests, you can use base link `https://messenger-svc-gfsy.onrender.com` followed by:

| Method | Endpoint     | Description                  |
|--------|--------------|------------------------------|
| POST   | `/start`     | Start auto-sender            |
| POST   | `/stop`      | Stop auto-sender             |
| GET    | `/sent`      | List all sent messages       |

For testing purposes only, you can use the following utility endpoints:

| Method | Endpoint | Description                  |
|--------|----------|------------------------------|
| GET    | `/ping`  | Health check                 |
| POST   | `/seed`  | Seeds 10 sample data into db |
| DELETE | `/clear` | Clears database              |
---

## 🐳 Run with Docker Compose

> Make sure you have Docker + Docker Compose installed

---

## 🧪 Unit Tests

> A complete list of unit tests is available in the `internal/core/tests` package.

---

## 🪜 Step by step

1. Clone the repository

```bash
git clone https://github.com/hasElvin/messenger-svc.git
cd messenger-svc
```

2. If you're using custom credentials or ports, edit your .env file or config/config.go.
- For Redis, you can add connection url into config.yaml file or relevant environment variable.
- No API key is required for webhook usage unless you configure one.
- If run locally, by default the app connects to: `localhost:5432` for PostgreSQL



3. Build and run the containers
```bash
docker compose up --build
```
On first run, this will:
- Build the Go application
- Start PostgreSQL and Redis
- Auto-migrate the DB schema
- Insert 10 sample messages
- Launch the API server on `localhost:8080` (unless overwritten in .env or by hosting provider)
- Start the message sender automatically

4. Redis cache check
```bash
docker compose exec redis redis-cli
127.0.0.1:6379> KEYS *
127.0.0.1:6379> GET msg:1
```
After messages are sent, their delivery metadata is stored in Redis.
Each key looks like `msg:<id>` with a corresponding value
---

## 📝 Notes
- Webhook url has been constructed in a way that only returns static msgId just because the dynamic values in custom actions are only supported in their paid plan.
- Webhook url might get expired from time to time. I will monitor myself, but in case of expiration, feel free to generate your own and add it to config.yaml or relevant environment variable.
- The send interval between the messages, the message character limit, and maximum retry allowance limit in case of failed webhook calls are in the config.yaml for the purpose of simplicity. If needed, they can easily be incorporated into the endpoint params. 