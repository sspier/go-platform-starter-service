# go-platform-starter-service

A minimal Go HTTP service intended as a starting point for platform and infrastructure-focused work.  
The service exposes basic endpoints that are commonly used by load balancers, health checks, and platform tooling.


## Features

- HTTP server written in Go
- Configurable port via `PORT` environment variable (default: `8080`)
- `GET /health` endpoint for health checks
- `GET /version` endpoint exposing service and version metadata

## Getting Started

Install dependencies and tidy the module:

```bash
go mod tidy
```

Run the service:

```bash
make run
# or
go run ./cmd/server
```

By default, the service listens on http://localhost:8080.

## Endpoints

### GET /health

```bash
curl http://localhost:8080/health
```

Example response:

```json
{
  "status": "ok",
  "service": "go-platform-starter-service"
}
```

### GET /version

```bash
curl http://localhost:8080/version
```

Example response:

```json
{
  "service": "go-platform-starter-service",
  "version": "0.1.0"
}
```

### Configuration

- PORT (optional): port for the HTTP server to listen on.

Example:

```bash
PORT=9090 go run ./cmd/server
# service listens on http://localhost:9090
```

### Build

``` bash
make build
# binary will be in ./bin/go-platform-starter-service
```

### Test

```bash
make test
# currently runs go test ./...
```