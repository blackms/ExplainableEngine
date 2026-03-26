# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.1.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.1.0
