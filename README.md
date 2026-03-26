# Explainable Engine

Transform any numerical or decisional output into an explicit, queryable, verifiable causal chain. The engine decomposes a target value into weighted components, builds a dependency graph, computes confidence scores, detects missing data, ranks key drivers, and generates human-readable narrative explanations.

## Quick Start

### Local (Go)

```bash
go run ./cmd/server          # starts on :8000 with in-memory store
```

Or use Make:

```bash
make run                     # same as above
make build                   # compile to bin/server
make test                    # run all tests
make lint                    # go vet + golangci-lint
```

### Docker

```bash
docker compose up -d         # builds and starts with SQLite persistence
```

Or build manually:

```bash
docker build -t explainable-engine .
docker run -p 8000:8000 explainable-engine
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check with uptime and version |
| `POST` | `/api/v1/explain` | Create a new explanation |
| `GET` | `/api/v1/explain/{id}` | Retrieve an explanation by ID |
| `GET` | `/api/v1/explain/{id}/graph?format=` | Export dependency graph (json, dot, mermaid) |
| `GET` | `/api/v1/explain/{id}/narrative?level=&lang=` | Generate narrative (level: basic/advanced, lang: en/it) |
| `POST` | `/api/v1/explain/{id}/what-if` | What-if sensitivity analysis |

## Example Requests

### Create an explanation

```bash
curl -X POST http://localhost:8000/api/v1/explain \
  -H "Content-Type: application/json" \
  -d '{
    "target": "monthly_cost",
    "components": [
      {"name": "compute", "value": 500, "weight": 0.6},
      {"name": "storage", "value": 200, "weight": 0.3},
      {"name": "network", "value": 100, "weight": 0.1}
    ]
  }'
```

### Retrieve an explanation

```bash
curl http://localhost:8000/api/v1/explain/{id}
```

### Export dependency graph

```bash
curl "http://localhost:8000/api/v1/explain/{id}/graph?format=mermaid"
```

Supported formats: `json` (default), `dot`, `mermaid`.

### Generate narrative

```bash
curl "http://localhost:8000/api/v1/explain/{id}/narrative?level=basic&lang=en"
```

Parameters: `level` (basic, advanced), `lang` (en, it).

### What-if analysis

```bash
curl -X POST http://localhost:8000/api/v1/explain/{id}/what-if \
  -H "Content-Type: application/json" \
  -d '{
    "modifications": [
      {"component": "compute", "new_value": 300}
    ]
  }'
```

### Health check

```bash
curl http://localhost:8000/health
```

Response:

```json
{"status":"healthy","uptime":"2m30s","version":"0.4.0"}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8000` | HTTP listen port |
| `STORAGE_BACKEND` | `memory` | Storage backend: `memory` or `sqlite` |
| `SQLITE_PATH` | `explanations.db` | Path to SQLite database file |
| `LOG_LEVEL` | `INFO` | Logging level: `DEBUG`, `INFO`, `WARNING` |
| `CORS_ORIGINS` | `*` | Comma-separated allowed origins, or `*` for all |

## Architecture

The engine follows a layered architecture:

- **cmd/server** -- HTTP entry point, wires dependencies and middleware.
- **internal/api** -- HTTP handlers, routing, request/response helpers.
- **internal/middleware** -- CORS and structured JSON logging middleware.
- **internal/engine** -- Core computation: DAG construction, breakdown, confidence propagation, driver analysis, sensitivity (what-if), graph serialization, and narrative generation.
- **internal/models** -- Shared data types (request, response, graph nodes/edges).
- **internal/storage** -- Persistence interface with in-memory (LRU) and SQLite implementations.

Request flow: HTTP request enters through CORS and logging middleware, then the router applies recovery, request-ID, and timing middleware before dispatching to the appropriate handler. The handler delegates to the orchestrator, which builds a DAG from the request, runs all analysis passes, and returns a structured response.

## License

See LICENSE file for details.
