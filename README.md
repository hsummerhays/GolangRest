# Golang API Boilerplate 🚀

A robust, production-ready REST API written in Go, focusing on modern backend engineering practices, concurrency, and observability.

This repository is designed not just as a standard CRUD API, but as an **engineering showcase** of how to structure and scale a real-world Go service.

---

## 🏗 Architecture & Design Decisions

### **1. Clean Architecture & Dependency Injection**
The application adheres strictly to separation of concerns:
- **Handlers (Transport Layer):** Manage HTTP translation, request decoding, and JSON responses.
- **Services (Business Logic):** Handle validations and orchestration.
- **Repository (Data Layer):** Interfaces with SQLite using WAL mode for safe concurrent I/O.

All components use **Dependency Injection**, making them completely unit-testable. The database repository is abstracted via interfaces, meaning we can swap SQLite out for PostgreSQL simply by injecting a different implementation—no business logic changes required.

### **2. Concurrency: The Worker Pool Pattern**
Go routines are powerful, but unbounded concurrency can crash a database.
This API features a bespoke **Worker Pool** (`pkg/worker/pool.go`) to handle background processing for bulk data imports:
- **The Problem:** The `POST /products/batch` endpoint needs to import thousands of records. If done synchronously, the HTTP request hangs. If we spawn a goroutine per insert, we exhaust database connections.
- **The Solution:** The batch endpoint decodes the payload, pushes jobs onto a buffered channel, and immediately returns `202 Accepted`. A fixed pool of workers pulls jobs from the channel and executes them, ensuring predictable database load without dropping requests.

### **3. Context Propagation & Resilience**
- **Timeout Middleware:** Every HTTP request runs under a `60-second context timeout` middleware.
- **Context Passing:** Contexts flow deeply from the HTTP layer `r.Context()` through the service layer, all the way into the `sql.DB` repository via `ExecContext` and `QueryContext`. If a client hangs up early, the database query aborts, saving CPU and DB resources.
- **SQLite WAL Mode & Busy Timeouts:** By enabling Write-Ahead Logging (`journal_mode=WAL`) and setting a busy timeout (`busy_timeout=5000`), the database survives high-concurrency writes coming from the worker pool without throwing "database is locked" exceptions.

### **4. Production Readiness**
- **Graceful Shutdown:** The server listens for `SIGTERM` and `SIGINT`. Upon termination, it stops accepting new requests, drains inflight connections, and waits for the worker pool to finish its buffered jobs before exiting. Zero data loss.
- **Observability (Health & Readiness):**
  - `/health` ensures the binary is running (Liveness probe).
  - `/ready` actively pings the database via `PingContext()` to ensure dependencies are alive (Readiness probe).
- **Rate Limiting:** A custom `RateLimiter` middleware protects the application from DDoS and brute-force attacks on a per-IP basis using a token bucket algorithm.
- **API Documentation:** The service leverages Swagger (`swaggo`) to automatically generate OpenAPI specs and hosts a live Swagger UI at `/swagger/index.html`.

---

## 🛠️ Tech Stack
- **Go 1.23+**
- **Router:** `go-chi/chi` (Standard library compatible)
- **Database:** `modernc.org/sqlite` (CGO-free pure Go SQLite driver)
- **Validation:** `go-playground/validator`
- **Docs:** `swaggo/http-swagger`

---

## 🚦 Getting Started

### **1. Run Tests**
```bash
go test ./... -v
```

### **2. Run Locally**
```bash
# Generate Swagger docs
swag init -g cmd/server/main.go

# Start server
go run cmd/server/main.go
```

The application will be available at `http://localhost:8080`.

### **3. Swagger UI**
Interactive documentation is available at:
👉 `http://localhost:8080/swagger/index.html`

### **4. Try the Batch Importer**
You can test the concurrency mechanism yourself:
```bash
curl -X POST -H "Content-Type: application/json" -d '[
  {"name": "Cortado", "price": 3.75},
  {"name": "Macchiato", "price": 3.25},
  {"name": "Flat White", "price": 4.25}
]' http://localhost:8080/products/batch
```
You will instantly receive a `202 Accepted` response while the workers persist the products in the background!
