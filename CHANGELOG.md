# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-26

### Added — Dashboard (Sprints 8-12)

- **Explanation Explorer** (Sprint 8): Summary card with confidence gauge, interactive breakdown chart with drill-down, driver ranking, detail page ([#95](https://github.com/blackms/ExplainableEngine/pull/95))
- **Dependency Graph** (Sprint 9): Interactive DAG with React Flow + dagre layout, custom colored nodes, zoom/pan/minimap, full-screen view ([#96](https://github.com/blackms/ExplainableEngine/pull/96))
- **What-if Simulator** (Sprint 9): Per-component sliders with debounced API calls, side-by-side comparison, diff table, sensitivity ranking ([#96](https://github.com/blackms/ExplainableEngine/pull/96))
- **Narrative Viewer** (Sprint 9): Basic/advanced toggle, EN/IT language support, copy-to-clipboard ([#96](https://github.com/blackms/ExplainableEngine/pull/96))
- **Confidence Panel** (Sprint 9): Per-node confidence bars, missing data impact warnings ([#96](https://github.com/blackms/ExplainableEngine/pull/96))
- **Audit Trail** (Sprint 10): Paginated table, search/filter (target, confidence range, dates), cursor-based pagination ([#97](https://github.com/blackms/ExplainableEngine/pull/97))
- **Authentication** (Sprint 10): NextAuth.js v5 with Google OAuth, login page, session management, middleware ([#97](https://github.com/blackms/ExplainableEngine/pull/97))
- **List API Endpoint** (Sprint 10): `GET /api/v1/explain` with pagination, filters; `GET /api/v1/stats` ([#97](https://github.com/blackms/ExplainableEngine/pull/97))
- **PDF/CSV Export** (Sprint 11): Single explanation PDF report, filtered audit list CSV download ([#98](https://github.com/blackms/ExplainableEngine/pull/98))
- **Live Monitoring** (Sprint 11): 10s polling feed, anomaly alerts (low confidence, high missing), summary statistics dashboard ([#98](https://github.com/blackms/ExplainableEngine/pull/98))
- **Scenario Manager** (Sprint 11): Save what-if scenarios to localStorage, load into simulator, compare up to 3 side-by-side ([#98](https://github.com/blackms/ExplainableEngine/pull/98))
- **API Playground** (Sprint 12): Interactive explorer for all 6 endpoints, JSON editor, response viewer with timing ([#99](https://github.com/blackms/ExplainableEngine/pull/99))
- **Code Generator** (Sprint 12): curl/Python/Go/JavaScript snippets, auto-update, copy-to-clipboard ([#99](https://github.com/blackms/ExplainableEngine/pull/99))
- **API Key Management** (Sprint 12): Generate/revoke keys, one-time display, masked listing ([#99](https://github.com/blackms/ExplainableEngine/pull/99))
- **User Preferences** (Sprint 12): Theme (light/dark/auto), default language, narrative level ([#99](https://github.com/blackms/ExplainableEngine/pull/99))

### All 21 Dashboard User Stories Complete

US-101 through US-603 across 6 epics: Explorer, What-if, Audit, Monitoring, Playground, Settings.

## [0.5.0] - 2026-03-26

### Added

- **GCP Deployment**: Full production deployment on Google Cloud Platform ([#66](https://github.com/blackms/ExplainableEngine/pull/66), [#67](https://github.com/blackms/ExplainableEngine/pull/67))
  - Cloud Run v2 service with autoscaling (0-10 instances)
  - Cloud SQL PostgreSQL 16 (private IP, automated backups, PITR)
  - VPC with private networking and serverless VPC connector
  - Secret Manager for database credentials
  - Artifact Registry for container images
- **PostgreSQL Storage Backend**: Production-grade persistence with pgx driver, connection pooling, and auto-migration ([#66](https://github.com/blackms/ExplainableEngine/pull/66))
- **Database Migration System**: Embedded SQL migrations with schema versioning, auto-run on startup, rollback support ([#67](https://github.com/blackms/ExplainableEngine/pull/67))
- **Terraform IaC**: Complete infrastructure as code — VPC, Cloud SQL, Cloud Run, Secret Manager, IAM ([#66](https://github.com/blackms/ExplainableEngine/pull/66))
- **Cloud Build Pipeline**: `cloudbuild.yaml` — test, build, push, deploy in one pipeline ([#67](https://github.com/blackms/ExplainableEngine/pull/67))
- **GitHub Actions CI/CD**: PR testing (`ci.yml`) and automated deployment (`deploy.yml`) with Workload Identity Federation ([#67](https://github.com/blackms/ExplainableEngine/pull/67))
- **Production Smoke Tests**: `scripts/smoke-test.sh` — 6-endpoint validation script ([#67](https://github.com/blackms/ExplainableEngine/pull/67))

### Production URL

`https://explainable-engine-516741092583.europe-west1.run.app`

## [0.4.0] - 2026-03-26

### Added

- **Narrative Engine**: Template-based human-readable explanations with basic (1 sentence) and advanced (paragraph) levels, EN/IT language support, confidence classification (high/moderate/low), and missing data warnings ([#50](https://github.com/blackms/ExplainableEngine/pull/50))
- **Narrative Endpoint**: `GET /api/v1/explain/{id}/narrative?level=basic|advanced&lang=en|it` ([#50](https://github.com/blackms/ExplainableEngine/pull/50))
- **Structured JSON Logging**: Request logging middleware with timestamp, method, path, status, duration, request_id; configurable via `LOG_LEVEL` env var ([#50](https://github.com/blackms/ExplainableEngine/pull/50))
- **CORS Middleware**: Configurable cross-origin support via `CORS_ORIGINS` env var ([#50](https://github.com/blackms/ExplainableEngine/pull/50))
- **Docker Support**: Multi-stage Dockerfile + docker-compose.yml with SQLite persistence volume ([#50](https://github.com/blackms/ExplainableEngine/pull/50))
- **Enhanced Health**: `/health` now returns uptime and version 0.4.0 ([#50](https://github.com/blackms/ExplainableEngine/pull/50))

### Changed

- **README**: Complete rewrite with all 6 API endpoints, curl examples, environment variables, architecture overview ([#50](https://github.com/blackms/ExplainableEngine/pull/50))

### User Stories Completed

- US-010: Human explanation (Natural Language Layer)

## [0.3.0] - 2026-03-26

### Added

- **Graph Serializer**: Multi-format export — JSON, DOT (Graphviz), Mermaid — with strategy pattern, deterministic output, and node coloring by type (INPUT=green, OUTPUT=pink, COMPUTED=blue, MISSING=grey) ([#49](https://github.com/blackms/ExplainableEngine/pull/49))
- **Graph Export Endpoint**: `GET /api/v1/explain/{id}/graph?format=json|dot|mermaid` with content-type negotiation ([#49](https://github.com/blackms/ExplainableEngine/pull/49))
- **Sensitivity Analyzer**: What-if analysis engine that clones requests, applies modifications, recomputes via orchestrator, and produces per-component diffs with sensitivity ranking ([#49](https://github.com/blackms/ExplainableEngine/pull/49))
- **What-if Endpoint**: `POST /api/v1/explain/{id}/what-if` for scenario exploration with transient results (not persisted) ([#49](https://github.com/blackms/ExplainableEngine/pull/49))
- **OriginalRequest preservation**: ExplainResponse now stores the original request for what-if recomputation ([#49](https://github.com/blackms/ExplainableEngine/pull/49))

### User Stories Completed

- US-003: Visualizzazione grafo
- US-006: Sensitivity analysis

## [0.2.0] - 2026-03-26

### Added

- **Missing Data Analyzer**: Detect MISSING nodes in the causal graph and compute their impact as ratio of missing weights to total weights, with configurable threshold warnings ([#48](https://github.com/blackms/ExplainableEngine/pull/48))
- **Driver Analyzer**: Dedicated module for ranking top drivers by normalized impact score (|contribution| x confidence), with multi-level flattening and deterministic tie-breaking ([#48](https://github.com/blackms/ExplainableEngine/pull/48))
- **Component.Missing flag**: Optional `missing` boolean on request components to mark data as unavailable ([#48](https://github.com/blackms/ExplainableEngine/pull/48))
- **ExplainOptions.MissingThreshold**: Configurable threshold (default 0.2) for missing data warnings ([#48](https://github.com/blackms/ExplainableEngine/pull/48))
- **Sprint 2 Integration Tests**: 7 new tests covering missing data, driver normalization, and backward compatibility ([#48](https://github.com/blackms/ExplainableEngine/pull/48))

### Changed

- **TopDrivers normalization**: Driver impacts are now normalized to [0, 1] range where the top driver has impact = 1.0 (previously raw absolute values) ([#48](https://github.com/blackms/ExplainableEngine/pull/48))
- **Orchestrator v2**: Pipeline extended with missing data analysis and dedicated driver analysis steps; removed inline driver logic ([#48](https://github.com/blackms/ExplainableEngine/pull/48))

### User Stories Completed

- US-005: Missing data impact
- US-007: Critical drivers

## [0.1.0] - 2026-03-26

### Added

- **Core Engine**: Custom DAG graph engine with adjacency lists, Kahn's topological sort, cycle detection, and deterministic iteration ([#46](https://github.com/blackms/ExplainableEngine/pull/46))
- **Breakdown Engine**: Recursive contribution analysis computing absolute and percentage contributions for each component ([#46](https://github.com/blackms/ExplainableEngine/pull/46))
- **Dependency Resolver**: Build dependency trees from causal graphs with DFS and BFS traversal ([#46](https://github.com/blackms/ExplainableEngine/pull/46))
- **Confidence Propagation**: Bottom-up weighted average confidence computation via topological order with full propagation path tracking ([#46](https://github.com/blackms/ExplainableEngine/pull/46))
- **Explanation Orchestrator**: Full pipeline (graph build, breakdown, dependency, confidence, top drivers, deterministic SHA-256 hash, UUID v4 generation) ([#46](https://github.com/blackms/ExplainableEngine/pull/46))
- **REST API**: `POST /api/v1/explain` and `GET /api/v1/explain/{id}` endpoints with JSON request/response ([#47](https://github.com/blackms/ExplainableEngine/pull/47))
- **Health Endpoint**: `GET /health` returning service status and version ([#47](https://github.com/blackms/ExplainableEngine/pull/47))
- **Middleware Stack**: Recovery (panic to 500), request ID (UUID in X-Request-Id), timing (X-Processing-Time-Ms) ([#47](https://github.com/blackms/ExplainableEngine/pull/47))
- **Storage Layer**: `ExplanationStore` interface with InMemoryStore (LRU, thread-safe) and SQLiteStore (pure Go via modernc.org/sqlite) ([#47](https://github.com/blackms/ExplainableEngine/pull/47))
- **Data Models**: Node, Edge, ExplanationGraph, ExplainRequest/Response, Contribution, ConfidenceResult with full JSON serialization
- **Test Suite**: 61 tests across engine (33), API (10), storage (13), models (5) with 83.7% coverage
- **Determinism Verification**: 100-iteration determinism tests ensuring identical output for same input

### User Stories Completed (MVP)

- US-001: Breakdown di un valore
- US-002: Tracciare le dipendenze
- US-004: Confidence score
- US-008: Submit explanation request (POST /explain)
- US-009: Retrieve explanation (GET /explain/{id})

[1.0.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v1.0.0
[0.5.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.5.0
[0.4.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.4.0
[0.3.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.3.0
[0.2.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.2.0
[0.1.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.1.0
