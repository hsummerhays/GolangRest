# Golang REST App

A simple Golang REST application featuring a basic health check endpoint.

## Prerequisites

- [Go](https://go.dev/doc/install) (1.16 or higher)

## Getting Started

1. **Clone the repository** (if applicable):
   ```bash
   git clone <repository-url>
   cd GolangRest
   ```

2. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

3. **Build the application** (optional):
   ```bash
   go build
   ./golangrest
   ```

The server will start and listen on port `8080`.

## Endpoints

### Health Check

- **URL:** `/health`
- **Method:** `GET`
- **Response:**
  ```json
  {
      "status": "UP"
  }
  ```

You can test this endpoint by navigating to [http://localhost:8080/health](http://localhost:8080/health) in your web browser or by using `curl`:

```bash
curl http://localhost:8080/health
```

### Get Products

- **URL:** `/products`
- **Method:** `GET`
- **Response:**
  ```json
  [
    {
      "id": 1,
      "name": "Espresso",
      "price": 3.5
    },
    ...
  ]
  ```

### Create Product

- **URL:** `/products`
- **Method:** `POST`
- **Body:**
  ```json
  {
      "name": "Mocha",
      "price": 5.00
  }
  ```
- **Response:**
  ```json
  {
      "id": 4,
      "name": "Mocha",
      "price": 5.00
  }
  ```
