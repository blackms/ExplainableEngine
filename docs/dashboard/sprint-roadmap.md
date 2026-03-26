# ExplainableEngine Dashboard — Sprint Roadmap

> Sprints 8-12 | 2 weeks each | React/Next.js + TypeScript + Tailwind
>
> Backend API: `https://explainable-engine-516741092583.europe-west1.run.app`

---

## Sprint 8 — Dashboard Foundation + Explanation Explorer Core

**Dates:** 2026-04-13 to 2026-04-26

**Goal:** Stand up the dashboard project with core layout and deliver the first usable explanation views so users can inspect AI decisions immediately.

### User Stories

| ID     | Title                        | Priority |
|--------|------------------------------|----------|
| US-101 | Explanation Summary View     | High     |
| US-102 | Breakdown Chart              | High     |

### Key Deliverables

- Project scaffolding: Next.js 14 (App Router) + TypeScript + Tailwind CSS
- Shared layout: sidebar navigation, top header, main content area
- API integration layer (typed fetch wrapper, error handling, auth headers)
- Basic routing (`/`, `/explanations`, `/explanations/:id`)
- Explanation summary card showing score, model, timestamp, and top factors
- Breakdown bar/donut chart visualizing factor contributions

### Dependencies

- Backend API endpoints: `GET /explain`, `GET /explain/:id`
- Design tokens / style guide (can be minimal; Tailwind defaults are acceptable for MVP)

### Risks

| Risk | Mitigation |
|------|-----------|
| API schema changes during sprint | Pin to current OpenAPI contract; add adapter layer |
| Next.js App Router learning curve | Timebox exploration to day 1; fall back to Pages Router if needed |

---

## Sprint 9 — Graph Visualization + What-if Simulator + Narrative

**Dates:** 2026-04-27 to 2026-05-10

**Goal:** Enable users to explore dependency graphs, read natural-language narratives, and run what-if simulations to understand how input changes affect outcomes.

### User Stories

| ID     | Title                              | Priority |
|--------|------------------------------------|----------|
| US-103 | Dependency Graph Visualization     | High     |
| US-104 | Confidence Panel                   | Medium   |
| US-105 | Narrative Viewer                   | Medium   |
| US-201 | What-if Simulator with Sliders     | High     |
| US-202 | Side-by-side Comparison            | Medium   |

### Key Deliverables

- Interactive dependency graph (react-flow or cytoscape.js) with zoom, pan, node details
- Confidence panel showing reliability score, interval, data-quality indicators
- Narrative viewer rendering markdown/plain-text explanations from `/narrative` endpoint
- What-if simulator: slider controls per input feature, live re-computation
- Side-by-side diff view comparing original vs. simulated explanation

### Dependencies

- Backend endpoints: `GET /graph`, `GET /confidence`, `GET /narrative`, `POST /whatif`
- Graph library selection finalized in Sprint 8 spike

### Risks

| Risk | Mitigation |
|------|-----------|
| Graph rendering performance on large models (>50 nodes) | Implement virtualization; limit default depth; add expand-on-click |
| What-if latency from backend | Show loading skeleton; consider debounce on slider changes |

---

## Sprint 10 — Audit Trail + Authentication

**Dates:** 2026-05-11 to 2026-05-24

**Goal:** Provide a searchable audit log of all explanations and secure the dashboard with Google OAuth so only authorized users can access it.

### User Stories

| ID     | Title                          | Priority |
|--------|--------------------------------|----------|
| US-301 | Explanation List with Pagination | High   |
| US-302 | Search and Filter              | High     |
| US-303 | Detail View                    | Medium   |
| US-601 | User Login with Google OAuth   | High     |
| US-203 | Sensitivity Ranking            | Medium   |

### Key Deliverables

- Paginated explanation history table with sorting
- Search bar + filter chips (date range, model, score threshold, tags)
- Full detail view for a single historical explanation (reuses Sprint 8 components)
- Google OAuth integration (NextAuth.js) with session management
- Protected routes middleware
- Sensitivity ranking table showing which features have the highest impact

### Dependencies

- Backend: paginated list endpoint, search/filter query parameters
- Google Cloud OAuth client ID configured
- Sprint 8/9 explanation components (reused in detail view)

### Risks

| Risk | Mitigation |
|------|-----------|
| OAuth redirect issues on deployed environment | Test with production URL early; configure allowed origins |
| Large audit datasets slow down pagination | Server-side pagination with cursor-based approach |

---

## Sprint 11 — Export + Live Monitoring

**Dates:** 2026-05-25 to 2026-06-07

**Goal:** Let users export explanations for compliance and monitor live model behavior with anomaly detection alerts.

### User Stories

| ID     | Title                    | Priority |
|--------|--------------------------|----------|
| US-304 | PDF/CSV Export           | Medium   |
| US-401 | Live Monitoring Feed     | High     |
| US-402 | Anomaly Alerts           | High     |
| US-403 | Summary Statistics       | Medium   |
| US-204 | Save Scenarios           | Low      |

### Key Deliverables

- Client-side PDF generation (e.g., react-pdf) and CSV download for explanations
- Real-time monitoring feed using SSE or WebSocket (live explanation stream)
- Anomaly alert cards with severity, timestamp, and affected features
- Summary statistics dashboard: total explanations, avg confidence, anomaly rate
- Save/load what-if scenarios to local storage or backend

### Dependencies

- Backend: SSE/WebSocket endpoint for live feed, anomaly detection logic
- Sprint 9 what-if simulator (for save scenarios feature)
- Sprint 10 auth (monitoring requires authenticated access)

### Risks

| Risk | Mitigation |
|------|-----------|
| WebSocket connection stability on Cloud Run | Fall back to SSE; implement reconnection logic |
| PDF rendering differences across browsers | Use server-side PDF generation as fallback |

---

## Sprint 12 — API Playground + Polish

**Dates:** 2026-06-08 to 2026-06-21

**Goal:** Ship the API playground for developers, finalize settings, and polish the entire dashboard for production readiness.

### User Stories

| ID     | Title                    | Priority |
|--------|--------------------------|----------|
| US-501 | API Playground           | High     |
| US-502 | Code Generator           | Medium   |
| US-602 | API Key Management       | High     |
| US-603 | User Preferences         | Low      |

### Key Deliverables

- Interactive API playground: endpoint selector, JSON body editor, response viewer
- Code snippet generator (cURL, Python, Go, JavaScript)
- API key CRUD interface with scopes and expiry
- User preferences panel (theme, default filters, notification settings)
- Performance optimization: code splitting, lazy loading, image optimization
- Mobile responsive layout adjustments
- Final QA pass: cross-browser testing, accessibility audit (WCAG 2.1 AA)

### Dependencies

- All previous sprints complete
- Backend: API key management endpoints
- OpenAPI spec for playground auto-discovery

### Risks

| Risk | Mitigation |
|------|-----------|
| Scope creep from polish tasks | Strict timebox: 3 days max for optimization, 2 days for QA |
| API key security concerns | Server-side key hashing; display key only once on creation |

---

## Summary Timeline

```
Sprint 8  [Apr 13 - Apr 26]  Foundation + Explorer Core
Sprint 9  [Apr 27 - May 10]  Graph + What-if + Narrative
Sprint 10 [May 11 - May 24]  Audit Trail + Auth
Sprint 11 [May 25 - Jun 07]  Export + Monitoring
Sprint 12 [Jun 08 - Jun 21]  Playground + Polish
```

## Epic-to-Sprint Mapping

| Epic | Sprint(s) |
|------|-----------|
| A: Explanation Explorer (US-101..105) | 8, 9 |
| B: What-if Simulator (US-201..204) | 9, 10, 11 |
| C: Audit Trail (US-301..304) | 10, 11 |
| D: Live Monitoring (US-401..403) | 11 |
| E: API Playground (US-501..502) | 12 |
| F: Settings & Auth (US-601..603) | 10, 12 |
