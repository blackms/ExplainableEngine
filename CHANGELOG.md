# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.2.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.2.0
[0.1.0]: https://github.com/blackms/ExplainableEngine/releases/tag/v0.1.0
