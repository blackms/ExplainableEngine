# Explainable Engine Dashboard -- Architecture Document

**Version:** 1.0.0
**Date:** 2026-03-26
**Status:** Approved for Sprint 8
**Author:** System Architecture

> This document provides the complete technical architecture for the Explainable Engine
> web dashboard. It is intended to give a developer everything needed to start coding.

---

## Table of Contents

1. [Overview](#1-overview)
2. [Tech Stack Decision](#2-tech-stack-decision)
3. [Project Structure](#3-project-structure)
4. [Architecture Decision Records](#4-architecture-decision-records)
5. [API Client and Type Definitions](#5-api-client-and-type-definitions)
6. [Component Architecture](#6-component-architecture)
7. [Data Flow and Component Interaction](#7-data-flow-and-component-interaction)
8. [Performance Considerations](#8-performance-considerations)
9. [Security](#9-security)
10. [Deployment Strategy](#10-deployment-strategy)
11. [Testing Strategy](#11-testing-strategy)
12. [Sprint 8 Implementation Plan](#12-sprint-8-implementation-plan)

---

## 1. Overview

The Explainable Engine Dashboard is a web frontend that provides a visual interface for
interacting with the ExplainableEngine Go API. The API is deployed on GCP Cloud Run at:

```
https://explainable-engine-516741092583.europe-west1.run.app
```

### API Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/v1/explain` | Create explanation |
| GET | `/api/v1/explain/{id}` | Retrieve explanation |
| GET | `/api/v1/explain/{id}/graph?format=json\|dot\|mermaid` | Graph export |
| GET | `/api/v1/explain/{id}/narrative?level=basic\|advanced&lang=en\|it` | Narrative |
| POST | `/api/v1/explain/{id}/what-if` | Sensitivity analysis |
| GET | `/health` | Health check |

### Design Principles

- The dashboard never calls the Go API directly from the browser.
- All API access goes through Next.js API routes (BFF pattern).
- Server Components handle initial data fetching; Client Components handle interactivity.
- TypeScript types mirror the Go models exactly -- no ad-hoc frontend types.

---

## 2. Tech Stack Decision

### 2.1 Framework: Next.js 14+ (App Router)

**Decision:** Next.js 14 with the App Router.

**Rationale:**
- Server-Side Rendering for SEO on shareable audit and explanation pages.
- React Server Components reduce JavaScript shipped to the client.
- Built-in API routes serve as a Backend-for-Frontend (BFF) layer.
- File-based routing maps cleanly to our page hierarchy.
- First-class Vercel deployment with zero-config, with Cloud Run as a fallback.

**Alternative considered:** Vite + React SPA. Simpler to set up, but no SSR, no built-in
BFF layer, and shareable explanation URLs would render blank until JS loads.

### 2.2 Language: TypeScript (strict mode)

**Decision:** TypeScript with `strict: true` in `tsconfig.json`.

**Rationale:**
- The domain models are complex (recursive components, nested breakdowns, graph
  structures). Type safety prevents an entire class of bugs at compile time.
- Strict mode catches nullable access, implicit `any`, and missing return types.
- The Go API has well-defined JSON contracts. TypeScript interfaces enforce that the
  frontend stays in sync.

### 2.3 Styling: Tailwind CSS + shadcn/ui

**Decision:** Tailwind CSS v3 for utility-first styling, with shadcn/ui as the component library.

**Rationale:**
- Tailwind eliminates context-switching between CSS files and components.
- shadcn/ui provides accessible, unstyled primitives (Dialog, Dropdown, Table, Tabs) that
  we own as source code -- not an opaque dependency.
- Consistent design system out of the box with Tailwind's design tokens.
- shadcn/ui components use Radix UI primitives, which handle keyboard navigation, focus
  trapping, and ARIA attributes correctly.

**Alternative considered:** Chakra UI. Heavier runtime, harder to customize at the
token level, and adds a CSS-in-JS layer we do not need.

### 2.4 State Management: TanStack Query v5 (React Query)

**Decision:** TanStack Query for all server state. No Redux, no Zustand.

**Rationale:**
- The dashboard is read-heavy. Almost all state comes from the API.
- TanStack Query provides: caching, background refetching, stale-while-revalidate,
  optimistic updates, and request deduplication.
- The what-if simulator is the only feature with meaningful client-side state, and
  `useState` / `useReducer` handles it without a global store.
- Redux would add boilerplate for zero benefit in this domain.

### 2.5 Graph Visualization: React Flow (@xyflow/react)

**Decision:** React Flow v12 for rendering the causal DAG.

**Rationale:**
- Purpose-built for node-based graphs in React.
- Supports custom node and edge components, which we need for styling target, component,
  and sub_component nodes differently.
- Built-in zoom, pan, minimap, and controls.
- Handles DAG layout via dagre or elkjs integration.
- Active maintenance and large community.

**Alternative considered:** D3.js. More powerful for arbitrary visualizations, but
significantly more code to maintain. D3 manages its own DOM, which conflicts with React's
reconciliation. React Flow is the right abstraction level for our use case.

### 2.6 Charts: Recharts

**Decision:** Recharts for bar charts, treemaps, and gauges.

**Rationale:**
- Declarative React API that composes naturally with our component tree.
- Handles the breakdown bar chart, contribution treemap, and confidence gauge.
- Lightweight compared to Tremor (which bundles Tailwind and its own design system).

### 2.7 Auth: NextAuth.js v5 with Google Provider

**Decision:** NextAuth.js (Auth.js) for authentication.

**Rationale:**
- First-class Next.js integration (middleware, server components, API routes).
- Google OAuth covers the team's identity provider.
- Session management via encrypted JWTs stored in HTTP-only cookies.
- Easy to add additional providers later (GitHub, Azure AD).

### 2.8 Testing: Vitest + React Testing Library + Playwright

**Decision:** Three-layer testing strategy.

| Layer | Tool | Scope |
|-------|------|-------|
| Unit | Vitest | Utility functions, hooks, type guards |
| Component | React Testing Library + Vitest | Individual components, user interactions |
| E2E | Playwright | Full page flows, API integration, visual regression |

**Rationale:**
- Vitest is faster than Jest, natively supports ESM, and shares Vite's transform pipeline.
- React Testing Library enforces testing from the user's perspective.
- Playwright provides cross-browser E2E coverage and is the current industry standard.

---

## 3. Project Structure

```
dashboard/
├── app/                              # Next.js App Router
│   ├── layout.tsx                    # Root layout: sidebar, header, providers
│   ├── page.tsx                      # Home: recent explanations list
│   ├── loading.tsx                   # Root loading skeleton
│   ├── error.tsx                     # Root error boundary
│   ├── explain/
│   │   └── [id]/
│   │       ├── page.tsx              # Explanation detail (summary + breakdown)
│   │       ├── loading.tsx           # Detail skeleton
│   │       ├── graph/
│   │       │   └── page.tsx          # Full-screen interactive graph
│   │       └── whatif/
│   │           └── page.tsx          # What-if simulator
│   ├── audit/
│   │   └── page.tsx                  # Paginated audit trail
│   ├── monitor/
│   │   └── page.tsx                  # Live monitoring feed
│   ├── playground/
│   │   └── page.tsx                  # API playground (submit raw JSON)
│   ├── settings/
│   │   └── page.tsx                  # User preferences (language, theme)
│   └── api/                          # BFF API routes
│       ├── explain/
│       │   ├── route.ts              # POST: proxy create explanation
│       │   └── [id]/
│       │       ├── route.ts          # GET: proxy retrieve explanation
│       │       ├── graph/route.ts    # GET: proxy graph export
│       │       ├── narrative/route.ts # GET: proxy narrative
│       │       └── whatif/route.ts   # POST: proxy what-if
│       └── health/
│           └── route.ts              # GET: proxy health check
├── components/
│   ├── explanation/
│   │   ├── SummaryCard.tsx           # Target value, confidence, hash
│   │   ├── BreakdownChart.tsx        # Horizontal bar chart of contributions
│   │   ├── BreakdownTreemap.tsx      # Treemap alternative view
│   │   ├── ConfidenceGauge.tsx       # Radial gauge for confidence score
│   │   ├── NarrativeViewer.tsx       # Rendered narrative text (basic/advanced)
│   │   ├── DriverRanking.tsx         # Ordered list of top drivers
│   │   └── MetadataPanel.tsx         # Version, hash, created_at, computation_type
│   ├── graph/
│   │   ├── ExplanationGraph.tsx      # React Flow canvas with DAG layout
│   │   ├── TargetNode.tsx            # Custom node: OUTPUT type (highlighted)
│   │   ├── ComponentNode.tsx         # Custom node: component type
│   │   ├── SubComponentNode.tsx      # Custom node: sub_component type
│   │   ├── WeightedEdge.tsx          # Custom edge: shows weight label
│   │   └── GraphControls.tsx         # Zoom, fit, export PNG/SVG, format toggle
│   ├── whatif/
│   │   ├── WhatIfSimulator.tsx       # Main container: sliders + results
│   │   ├── ComponentSlider.tsx       # Single component value slider
│   │   ├── ComparisonView.tsx        # Side-by-side original vs modified
│   │   ├── DiffTable.tsx             # Component diffs table
│   │   └── SensitivityRanking.tsx    # Ranked sensitivity results
│   ├── audit/
│   │   ├── ExplanationTable.tsx      # Virtualized, sortable, filterable table
│   │   ├── FilterPanel.tsx           # Search by target, date range, confidence
│   │   └── ExportButton.tsx          # Export as CSV or PDF
│   ├── monitor/
│   │   ├── LiveFeed.tsx              # Auto-polling explanation feed
│   │   ├── AlertBadge.tsx            # Low-confidence or high-missing-impact alert
│   │   └── StatsCards.tsx            # Aggregate metrics (count, avg confidence)
│   ├── layout/
│   │   ├── Sidebar.tsx               # Navigation sidebar
│   │   ├── Header.tsx                # Top bar with breadcrumbs, user menu
│   │   └── ThemeToggle.tsx           # Light/dark mode switch
│   └── ui/                           # shadcn/ui primitives (generated, not hand-written)
│       ├── button.tsx
│       ├── card.tsx
│       ├── dialog.tsx
│       ├── dropdown-menu.tsx
│       ├── input.tsx
│       ├── label.tsx
│       ├── select.tsx
│       ├── skeleton.tsx
│       ├── slider.tsx
│       ├── table.tsx
│       ├── tabs.tsx
│       └── tooltip.tsx
├── lib/
│   ├── api/
│   │   ├── client.ts                 # Typed fetch wrapper for BFF routes
│   │   ├── types.ts                  # TypeScript types mirroring Go models
│   │   └── hooks.ts                  # TanStack Query hooks for each endpoint
│   ├── graph/
│   │   └── layout.ts                 # dagre layout: API graph -> React Flow nodes/edges
│   ├── utils.ts                      # Formatting, color scales, number helpers
│   └── constants.ts                  # API base URL, polling intervals, defaults
├── public/
│   └── favicon.ico
├── tests/
│   ├── unit/
│   │   ├── lib/
│   │   │   ├── api-client.test.ts
│   │   │   └── graph-layout.test.ts
│   │   └── components/
│   │       ├── SummaryCard.test.tsx
│   │       └── BreakdownChart.test.tsx
│   └── e2e/
│       ├── explanation-detail.spec.ts
│       ├── whatif-simulator.spec.ts
│       └── audit-trail.spec.ts
├── .env.local.example                # NEXT_PUBLIC_* and server-side env vars
├── .eslintrc.json
├── .prettierrc
├── next.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── vitest.config.ts
├── playwright.config.ts
└── package.json
```

### Key Structural Decisions

- **`app/api/` routes mirror the Go API paths.** Each BFF route is a thin proxy that
  adds auth headers, rewrites the URL to the Go backend, and optionally transforms the
  response shape.
- **`components/` is organized by feature, not by type.** A developer working on the
  what-if simulator finds everything in `components/whatif/`.
- **`lib/graph/layout.ts` is the bridge between the API response and React Flow.** It
  converts `GraphNodeResponse[]` and `GraphEdgeResponse[]` into React Flow's `Node[]` and
  `Edge[]` with positions computed by dagre.
- **`ui/` contains only shadcn/ui generated files.** These are not hand-written. Run
  `npx shadcn-ui@latest add button` to add new primitives.

---

## 4. Architecture Decision Records

### ADR-001: BFF Pattern via Next.js API Routes

**Context:** The browser needs to communicate with the Go API deployed on Cloud Run.

**Decision:** The browser never calls the Go API directly. All requests go through
Next.js API routes (`app/api/`), which proxy to the Go backend.

**Rationale:**
- **Security:** The Go API URL is never exposed to the browser. API keys or service
  account tokens stay server-side.
- **Caching:** The BFF layer can add response caching (e.g., `Cache-Control` headers,
  in-memory LRU cache for hot explanations).
- **Data transformation:** The BFF can reshape responses for the frontend's needs without
  changing the Go API contract.
- **Auth enforcement:** NextAuth.js middleware runs before BFF routes, rejecting
  unauthenticated requests before they reach the Go API.
- **CORS elimination:** Browser-to-BFF is same-origin. No CORS configuration needed.

**Trade-off:** Adds one network hop (browser -> Next.js -> Go API). Latency increase is
roughly 5-15ms per request on the same GCP region, which is acceptable.

**Implementation:**

```typescript
// app/api/explain/[id]/route.ts
import { getServerSession } from "next-auth";
import { authOptions } from "@/lib/auth";

const API_BASE = process.env.EXPLAINABLE_API_URL;

export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const session = await getServerSession(authOptions);
  if (!session) {
    return Response.json({ error: "Unauthorized" }, { status: 401 });
  }

  const res = await fetch(`${API_BASE}/api/v1/explain/${params.id}`, {
    headers: { "X-Request-Id": crypto.randomUUID() },
    next: { revalidate: 3600 }, // Cache for 1 hour (explanations are immutable)
  });

  const data = await res.json();
  return Response.json(data, { status: res.status });
}
```

---

### ADR-002: Server Components for Data Fetching

**Context:** Explanation detail pages load structured data (breakdown, graph, narrative)
that does not require client-side interactivity on first render.

**Decision:** Use React Server Components for initial data fetching on detail pages.
Client Components are used only for interactive elements: graph pan/zoom, what-if sliders,
audit table sorting/filtering.

**Rationale:**
- Server Components fetch data on the server, eliminating a client-side loading waterfall.
- Less JavaScript shipped to the browser. The breakdown chart data, narrative text, and
  metadata panel are rendered as HTML on the server.
- Client Components are loaded only when needed (graph canvas, slider inputs).

**Boundary definition:**

| Component | Type | Why |
|-----------|------|-----|
| `SummaryCard` | Server | Static display of values |
| `MetadataPanel` | Server | Static display |
| `NarrativeViewer` | Server | Static text rendering |
| `BreakdownChart` | Client (`"use client"`) | Recharts requires browser APIs |
| `ConfidenceGauge` | Client | Recharts |
| `ExplanationGraph` | Client | React Flow requires browser APIs |
| `WhatIfSimulator` | Client | Interactive sliders, real-time updates |
| `ExplanationTable` | Client | Sorting, filtering, virtualization |
| `LiveFeed` | Client | Polling, auto-update |

---

### ADR-003: Graph Rendering Strategy

**Context:** The Go API supports three graph export formats: JSON, DOT, and Mermaid.

**Decision:** Fetch the graph as JSON from the API and render it with React Flow. DOT and
Mermaid are available as export formats only (download buttons).

**Rationale:**
- JSON gives us structured data (`GraphNodeResponse[]`, `GraphEdgeResponse[]`) that maps
  directly to React Flow's node/edge model.
- React Flow provides interactivity (click nodes, hover for details, zoom, pan) that
  static DOT/Mermaid renders cannot match.
- Custom node components let us style target nodes differently from component nodes, show
  confidence as a color gradient, and highlight missing data nodes in red.
- Export to PNG/SVG is handled by React Flow's built-in `toImage()` utility. DOT and
  Mermaid exports use the API's `?format=dot` and `?format=mermaid` endpoints.

**Layout algorithm:** dagre (via `@dagrejs/dagre`) computes node positions from the DAG
structure. Layout direction is bottom-to-top (leaves at the bottom, target at the top) to
match the causal flow: inputs cause the output.

**Layout pipeline:**

```
API GraphResponse
    |
    v
lib/graph/layout.ts  (transform to React Flow format, compute dagre positions)
    |
    v
React Flow <ReactFlow nodes={nodes} edges={edges} />
```

---

### ADR-004: Real-time Updates via Polling

**Context:** The monitoring page shows a live feed of recent explanations.

**Decision:** Use polling at a 10-second interval via TanStack Query's `refetchInterval`.
Do not use WebSockets or Server-Sent Events.

**Rationale:**
- Cloud Run has a 15-minute request timeout. Long-lived WebSocket connections are not
  natively supported and require workarounds (keep-alive pings, reconnection logic).
- SSE works on Cloud Run but adds complexity to the Go API, which currently has no
  streaming endpoints.
- Polling at 10s is sufficient for a monitoring dashboard. The Go API's explanation
  creation is not high-frequency (tens per minute, not thousands).
- TanStack Query handles polling natively: `useQuery({ refetchInterval: 10_000 })`.
  It also pauses polling when the browser tab is inactive.

**Migration path:** If sub-second latency is needed in the future, add an SSE endpoint to
the Go API (`GET /api/v1/explain/stream`) and consume it with `EventSource` in a Client
Component. The BFF route would proxy the SSE stream.

---

### ADR-005: Deployment Strategy

**Context:** The Go API runs on Cloud Run in `europe-west1`. The dashboard needs a
deployment target.

**Decision:** Deploy to Vercel initially. Migrate to Cloud Run if needed.

**Rationale for starting with Vercel:**
- Zero-config deployment for Next.js (Vercel is the creator of Next.js).
- Automatic preview deployments for every PR.
- Edge Functions for API routes (low latency across regions).
- Built-in analytics and performance monitoring.
- Free tier covers early development.

**Rationale for Cloud Run fallback:**
- If we need the dashboard and API in the same GCP project for network policy or billing
  reasons, we containerize the Next.js app with `next start` in a Docker image.
- Cloud Run supports Next.js standalone output mode (`output: "standalone"` in
  `next.config.ts`), which produces a minimal Node.js server.

**Vercel environment variables:**

```
EXPLAINABLE_API_URL=https://explainable-engine-516741092583.europe-west1.run.app
NEXTAUTH_URL=https://dashboard.explainableengine.io
NEXTAUTH_SECRET=<generated>
GOOGLE_CLIENT_ID=<from GCP console>
GOOGLE_CLIENT_SECRET=<from GCP console>
```

---

### ADR-006: No Global State Management Library

**Context:** Several state management options exist for React (Redux, Zustand, Jotai,
MobX).

**Decision:** Do not add a global state management library. Use TanStack Query for server
state and React's built-in `useState`/`useReducer` for local UI state.

**Rationale:**
- Over 90% of the dashboard's state is server state (explanations, graphs, narratives).
  TanStack Query handles this with caching, deduplication, and background refetching.
- The only meaningful client state is in the what-if simulator (slider values, comparison
  mode toggle). This is local to the `WhatIfSimulator` component tree and does not need
  to be global.
- Adding Redux or Zustand would create a parallel data layer that duplicates what TanStack
  Query already does, adding indirection without benefit.

---

## 5. API Client and Type Definitions

### 5.1 TypeScript Types

These types mirror the Go models in `internal/models/`. The field names match the JSON
tags exactly.

```typescript
// lib/api/types.ts

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

export interface Component {
  id?: string;
  name: string;
  value: number;
  weight: number;
  confidence: number;
  missing?: boolean;
  components?: Component[];
}

export interface ExplainOptions {
  include_graph: boolean;
  include_drivers: boolean;
  max_drivers: number;
  max_depth: number;
  missing_threshold: number;
}

export interface ExplainRequest {
  target: string;
  value: number;
  components: Component[];
  options?: ExplainOptions;
  metadata?: Record<string, string>;
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

export interface BreakdownItem {
  node_id: string;
  label: string;
  value: number;
  weight: number;
  absolute_contribution: number;
  percentage: number;
  confidence: number;
  children?: BreakdownItem[];
}

export interface DriverItem {
  name: string;
  impact: number;
  rank: number;
}

export interface GraphNodeResponse {
  id: string;
  label: string;
  value: number;
  confidence: number;
  node_type: "input" | "computed" | "output" | "missing";
}

export interface GraphEdgeResponse {
  source: string;
  target: string;
  weight: number;
  transformation_type: "weighted_sum" | "normalization" | "threshold" | "custom";
}

export interface GraphResponse {
  nodes: GraphNodeResponse[];
  edges: GraphEdgeResponse[];
}

export interface DependencyNode {
  node_id: string;
  label: string;
  depth: number;
  relation?: string;
  children?: DependencyNode[];
}

export interface DependencyTree {
  root: DependencyNode;
  depth: number;
  total_nodes: number;
}

export interface ConfidenceDetail {
  overall: number;
  per_node: Record<string, number>;
}

export interface ExplainMetadata {
  version: string;
  created_at: string;              // ISO-8601
  deterministic_hash: string;      // "sha256:..."
  computation_type: "additive" | "multiplicative" | "custom";
}

export interface ExplainResponse {
  id: string;
  target: string;
  final_value: number;
  confidence: number;
  breakdown: BreakdownItem[];
  top_drivers: DriverItem[];
  missing_impact: number;
  graph?: GraphResponse;
  dependency_tree?: DependencyTree;
  confidence_detail?: ConfidenceDetail;
  metadata: ExplainMetadata;
  original_request?: ExplainRequest;
}

// ---------------------------------------------------------------------------
// Narrative types
// ---------------------------------------------------------------------------

export interface NarrativeResult {
  explanation_id: string;
  level: "basic" | "advanced";
  language: "en" | "it";
  narrative: string;
  confidence_level: string;
  has_missing_data: boolean;
}

// ---------------------------------------------------------------------------
// What-if types
// ---------------------------------------------------------------------------

export interface Modification {
  component: string;
  new_value: number;
}

export interface WhatIfRequest {
  modifications: Modification[];
}

export interface ComponentDiff {
  name: string;
  original_value: number;
  modified_value: number;
  delta_value: number;
  delta_percentage: number;
  original_contribution: number;
  modified_contribution: number;
}

export interface SensitivityRanking {
  name: string;
  impact: number;
  rank: number;
}

export interface SensitivityResult {
  original_value: number;
  modified_value: number;
  delta_value: number;
  delta_percentage: number;
  component_diffs: ComponentDiff[];
  sensitivity_ranking: SensitivityRanking[];
}

// ---------------------------------------------------------------------------
// Health types
// ---------------------------------------------------------------------------

export interface HealthResponse {
  status: "ok" | "degraded";
  version: string;
  uptime_seconds: number;
  timestamp: string;
}

// ---------------------------------------------------------------------------
// Error types
// ---------------------------------------------------------------------------

export interface ValidationDetail {
  field: string;
  issue: string;
}

export interface ApiError {
  error: {
    code: "BAD_REQUEST" | "NOT_FOUND" | "VALIDATION_ERROR" | "INTERNAL_ERROR";
    message: string;
    details?: ValidationDetail[];
    request_id: string;
  };
}
```

### 5.2 API Client

```typescript
// lib/api/client.ts

import type {
  ExplainRequest,
  ExplainResponse,
  GraphResponse,
  NarrativeResult,
  WhatIfRequest,
  SensitivityResult,
  HealthResponse,
} from "./types";

const BASE = "/api";

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new ApiClientError(res.status, body);
  }

  return res.json() as Promise<T>;
}

export class ApiClientError extends Error {
  constructor(
    public status: number,
    public body: unknown,
  ) {
    super(`API error: ${status}`);
    this.name = "ApiClientError";
  }
}

export const api = {
  explain: {
    create(data: ExplainRequest): Promise<ExplainResponse> {
      return request(`${BASE}/explain`, {
        method: "POST",
        body: JSON.stringify(data),
      });
    },

    get(id: string): Promise<ExplainResponse> {
      return request(`${BASE}/explain/${id}`);
    },

    graph(id: string, format: "json" | "dot" | "mermaid" = "json"): Promise<GraphResponse> {
      return request(`${BASE}/explain/${id}/graph?format=${format}`);
    },

    narrative(
      id: string,
      level: "basic" | "advanced" = "basic",
      lang: "en" | "it" = "en",
    ): Promise<NarrativeResult> {
      return request(`${BASE}/explain/${id}/narrative?level=${level}&lang=${lang}`);
    },

    whatIf(id: string, data: WhatIfRequest): Promise<SensitivityResult> {
      return request(`${BASE}/explain/${id}/whatif`, {
        method: "POST",
        body: JSON.stringify(data),
      });
    },
  },

  health(): Promise<HealthResponse> {
    return request(`${BASE}/health`);
  },
};
```

### 5.3 TanStack Query Hooks

```typescript
// lib/api/hooks.ts

import { useQuery, useMutation } from "@tanstack/react-query";
import { api } from "./client";
import type { ExplainRequest, WhatIfRequest } from "./types";

export const queryKeys = {
  explanation: (id: string) => ["explanation", id] as const,
  graph: (id: string) => ["explanation", id, "graph"] as const,
  narrative: (id: string, level: string, lang: string) =>
    ["explanation", id, "narrative", level, lang] as const,
  health: () => ["health"] as const,
};

export function useExplanation(id: string) {
  return useQuery({
    queryKey: queryKeys.explanation(id),
    queryFn: () => api.explain.get(id),
    staleTime: Infinity,           // Explanations are immutable
  });
}

export function useGraph(id: string) {
  return useQuery({
    queryKey: queryKeys.graph(id),
    queryFn: () => api.explain.graph(id, "json"),
    staleTime: Infinity,
  });
}

export function useNarrative(
  id: string,
  level: "basic" | "advanced" = "basic",
  lang: "en" | "it" = "en",
) {
  return useQuery({
    queryKey: queryKeys.narrative(id, level, lang),
    queryFn: () => api.explain.narrative(id, level, lang),
    staleTime: Infinity,
  });
}

export function useCreateExplanation() {
  return useMutation({
    mutationFn: (data: ExplainRequest) => api.explain.create(data),
  });
}

export function useWhatIf(id: string) {
  return useMutation({
    mutationFn: (data: WhatIfRequest) => api.explain.whatIf(id, data),
  });
}

export function useHealth() {
  return useQuery({
    queryKey: queryKeys.health(),
    queryFn: () => api.health(),
    refetchInterval: 30_000,       // Poll health every 30s
  });
}
```

---

## 6. Component Architecture

### 6.1 Explanation Detail Page (`app/explain/[id]/page.tsx`)

This is the primary page. It displays a complete explanation with all sub-views.

```
+------------------------------------------------------------------+
| Header: breadcrumb (Home > Explanations > expl_abc123)           |
+------------------------------------------------------------------+
| SummaryCard                          | ConfidenceGauge            |
| - Target: market_regime_score        | - 82.8% (radial gauge)    |
| - Final value: 0.72                  | - Label: "high"           |
| - Missing impact: 0.0               |                            |
+--------------------------------------+----------------------------+
| Tabs: [Breakdown] [Graph] [Narrative] [Metadata]                 |
+------------------------------------------------------------------+
| Tab: Breakdown                                                    |
| +--------------------------------------------------------------+ |
| | BreakdownChart (horizontal bar)                               | |
| | - trend_strength: 44.4% ████████████████                     | |
| | - momentum:       25.0% █████████                             | |
| | - volatility:     20.8% ████████                              | |
| +--------------------------------------------------------------+ |
| | DriverRanking (list)                                          | |
| | 1. trend_strength  impact: 0.44                               | |
| | 2. momentum        impact: 0.25                               | |
| | 3. volatility      impact: 0.208                              | |
| +--------------------------------------------------------------+ |
+------------------------------------------------------------------+
```

### 6.2 Graph Page (`app/explain/[id]/graph/page.tsx`)

Full-screen React Flow canvas with the causal DAG.

```
+------------------------------------------------------------------+
| GraphControls: [Zoom In] [Zoom Out] [Fit] [Export PNG] [Export DOT]|
+------------------------------------------------------------------+
|                                                                    |
|          [trend_strength]                                          |
|              0.8 | w=0.4                                          |
|                  \                                                 |
|     [volatility] ------ [market_regime_score]                      |
|         0.5 | w=0.3          0.72                                  |
|                  /                                                 |
|          [momentum]                                                |
|              0.6 | w=0.3                                          |
|                                                                    |
+------------------------------------------------------------------+
```

Node styling by type:

| Node Type | Background | Border | Icon |
|-----------|------------|--------|------|
| `output` | Blue-600 | 2px solid Blue-800 | Target icon |
| `input` | Green-100 | 1px solid Green-400 | Circle |
| `computed` | Gray-100 | 1px solid Gray-400 | Square |
| `missing` | Red-100 | 1px dashed Red-400 | Warning triangle |

Edge styling: stroke width proportional to weight. Label shows `w=0.4` on hover.

### 6.3 What-If Simulator (`app/explain/[id]/whatif/page.tsx`)

```
+------------------------------------------------------------------+
| Original: market_regime_score = 0.72                              |
+------------------------------------------------------------------+
| Component Sliders:                                                |
| trend_strength:  [====|=======] 0.80  (range: 0.0 - 1.0)        |
| volatility:      [==|=========] 0.50                              |
| momentum:        [===|========] 0.60                              |
|                                                [Run Analysis]     |
+------------------------------------------------------------------+
| ComparisonView:                                                   |
| +---------------------------+---------------------------+         |
| | Original                  | Modified                  |         |
| | Value: 0.72               | Value: 0.78               |         |
| | Confidence: 0.828         | Confidence: 0.828         |         |
| +---------------------------+---------------------------+         |
| Delta: +0.06 (+8.3%)                                             |
+------------------------------------------------------------------+
| DiffTable:                                                        |
| Component       | Original | Modified | Delta   | Delta %        |
| trend_strength  | 0.80     | 0.90     | +0.10   | +12.5%        |
| volatility      | 0.50     | 0.50     | 0.00    | 0.0%          |
| momentum        | 0.60     | 0.60     | 0.00    | 0.0%          |
+------------------------------------------------------------------+
```

**Interaction model:**
1. User adjusts sliders (client state via `useState`).
2. On "Run Analysis" click (or after 300ms debounce), call `useWhatIf` mutation.
3. Display `SensitivityResult` in the comparison view.
4. No navigation -- everything happens in-page.

---

## 7. Data Flow and Component Interaction

### 7.1 Request Flow

```
Browser                 Next.js (BFF)              Go API (Cloud Run)
   |                        |                            |
   |-- fetch /api/explain/X |                            |
   |                        |-- GET /api/v1/explain/X -->|
   |                        |<-- 200 ExplainResponse ----|
   |                        |                            |
   |                        | (optional: transform,      |
   |                        |  cache, add headers)       |
   |                        |                            |
   |<-- 200 JSON -----------|                            |
   |                        |                            |
   | TanStack Query caches  |                            |
   | the response. Next     |                            |
   | request serves from    |                            |
   | cache (staleTime:      |                            |
   | Infinity for immutable |                            |
   | explanations).         |                            |
```

### 7.2 Graph Rendering Pipeline

```
1. Server Component fetches ExplainResponse (includes graph if include_graph was true)
       |
       v
2. OR: Client Component calls useGraph(id) to fetch graph separately
       |
       v
3. lib/graph/layout.ts transforms API types to React Flow types:
       |
       |   GraphNodeResponse[] --> ReactFlow.Node[]
       |   - Maps node_type to custom component (TargetNode, ComponentNode, etc.)
       |   - Runs dagre layout to compute x, y positions
       |   - Sets node dimensions based on label length
       |
       |   GraphEdgeResponse[] --> ReactFlow.Edge[]
       |   - Maps source/target IDs
       |   - Sets edge label to weight value
       |   - Sets animated: true for high-weight edges
       |
       v
4. <ReactFlow nodes={nodes} edges={edges} nodeTypes={nodeTypes} />
       |
       v
5. React Flow renders interactive SVG canvas
```

### 7.3 What-If Interaction Flow

```
1. Page loads: useExplanation(id) fetches original ExplainResponse
       |
       v
2. WhatIfSimulator extracts components from original_request
   and initializes slider state: { [name]: value }
       |
       v
3. User adjusts slider for "trend_strength" from 0.8 to 0.9
       |
       v
4. Debounce (300ms) triggers useWhatIf mutation:
   POST /api/explain/{id}/whatif
   Body: { modifications: [{ component: "trend_strength", new_value: 0.9 }] }
       |
       v
5. BFF proxies to Go API:
   POST /api/v1/explain/{id}/what-if
       |
       v
6. Go API returns SensitivityResult
       |
       v
7. ComparisonView renders original_value vs modified_value
   DiffTable renders component_diffs
   SensitivityRanking renders sensitivity_ranking
```

---

## 8. Performance Considerations

### 8.1 Code Splitting and Lazy Loading

The graph visualization library (`@xyflow/react` + `@dagrejs/dagre`) adds approximately
150KB gzipped to the bundle. It must be lazy-loaded.

```typescript
// components/graph/ExplanationGraph.tsx
"use client";

import dynamic from "next/dynamic";

const ReactFlow = dynamic(
  () => import("@xyflow/react").then((mod) => mod.ReactFlow),
  { ssr: false, loading: () => <GraphSkeleton /> }
);
```

Recharts is similarly lazy-loaded for chart components.

### 8.2 Caching Strategy

| Data Type | staleTime | gcTime | Rationale |
|-----------|-----------|--------|-----------|
| Explanation (by ID) | `Infinity` | 30 min | Immutable. Never refetch. |
| Graph (by ID) | `Infinity` | 30 min | Immutable. Same as explanation. |
| Narrative (by ID+level+lang) | `Infinity` | 30 min | Deterministic output. |
| What-if result | 0 (no cache) | 0 | Depends on slider state. |
| Health check | 30s | 5 min | Should reflect current state. |
| Audit list | 60s | 5 min | New explanations appear. |
| Live feed | 10s | 1 min | Near-real-time. |

### 8.3 Debouncing

The what-if simulator debounces API calls by 300ms after the last slider change. This
prevents flooding the API during rapid slider adjustments.

```typescript
import { useDebouncedCallback } from "use-debounce";

const debouncedAnalyze = useDebouncedCallback(
  (modifications: Modification[]) => {
    whatIfMutation.mutate({ modifications });
  },
  300,
);
```

### 8.4 Table Virtualization

The audit trail table uses `@tanstack/react-virtual` to render only visible rows. For a
dataset of 10,000 explanations, only ~20 rows are in the DOM at any time.

### 8.5 Image and Asset Optimization

- Next.js `<Image>` component for any static images (logo, icons).
- Graph PNG exports are generated client-side via `html-to-image` (React Flow integration),
  not server-rendered.
- SVG exports use React Flow's built-in SVG serialization.

---

## 9. Security

### 9.1 Authentication Flow

```
1. User visits dashboard
       |
       v
2. NextAuth.js middleware checks session cookie
       |
       +-- No session --> redirect to /api/auth/signin (Google OAuth)
       |
       +-- Valid session --> proceed to page
       |
       v
3. BFF API routes check session server-side before proxying
```

### 9.2 Environment Variables

| Variable | Location | Exposure |
|----------|----------|----------|
| `EXPLAINABLE_API_URL` | Server only | Never sent to browser |
| `NEXTAUTH_SECRET` | Server only | Signs session tokens |
| `GOOGLE_CLIENT_ID` | Server only | OAuth flow |
| `GOOGLE_CLIENT_SECRET` | Server only | OAuth flow |
| `NEXTAUTH_URL` | Server only | Callback URL |

No `NEXT_PUBLIC_` prefixed variables expose API URLs or secrets.

### 9.3 Request Headers

BFF routes add:
- `X-Request-Id`: generated UUID for tracing across frontend and backend.
- `Authorization`: service account token if Go API requires authentication in the future.

### 9.4 Input Validation

- The playground page validates JSON payloads client-side before submission using Zod
  schemas derived from the TypeScript types.
- The BFF layer validates path parameters (explanation ID format: `^expl_[a-zA-Z0-9]{6,64}$`).
- All user input in sliders is bounded by min/max values from the original component data.

---

## 10. Deployment Strategy

### 10.1 Vercel (Primary)

```
GitHub repo (dashboard/ directory)
    |
    v
Vercel Project
    - Framework: Next.js (auto-detected)
    - Root Directory: dashboard/
    - Build Command: next build
    - Output Directory: .next
    - Node.js version: 20
    - Region: europe-west (ew1) -- same as Cloud Run
    |
    v
Production: https://dashboard.explainableengine.io
Preview: https://<branch>.dashboard.explainableengine.io
```

### 10.2 Cloud Run (Fallback)

```dockerfile
# dashboard/Dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

`next.config.ts` must include `output: "standalone"` for the Docker build.

---

## 11. Testing Strategy

### 11.1 Unit Tests (Vitest)

Target: utility functions and data transformation logic.

```
tests/unit/lib/graph-layout.test.ts   -- dagre layout produces valid positions
tests/unit/lib/api-client.test.ts     -- error handling, request formatting
tests/unit/lib/utils.test.ts          -- number formatting, color scales
```

### 11.2 Component Tests (Vitest + React Testing Library)

Target: individual components render correctly given props.

```
tests/unit/components/SummaryCard.test.tsx
  - renders target name and final value
  - renders missing impact warning when > 0

tests/unit/components/BreakdownChart.test.tsx
  - renders one bar per breakdown item
  - bars are sorted by percentage descending

tests/unit/components/DriverRanking.test.tsx
  - renders drivers in rank order
  - shows impact value formatted to 3 decimal places
```

### 11.3 E2E Tests (Playwright)

Target: full page flows against a mocked or real API.

```
tests/e2e/explanation-detail.spec.ts
  - navigate to explanation, verify summary card content
  - switch tabs (breakdown, graph, narrative, metadata)
  - graph page renders nodes matching API response

tests/e2e/whatif-simulator.spec.ts
  - adjust slider, verify debounced API call
  - comparison view shows delta values
  - reset sliders restores original values

tests/e2e/audit-trail.spec.ts
  - table loads with paginated data
  - filter by target name narrows results
  - export CSV downloads file
```

### 11.4 Coverage Targets

| Layer | Target |
|-------|--------|
| Unit (lib/) | 90% |
| Component | 80% |
| E2E | Critical paths covered (not measured by line coverage) |

---

## 12. Sprint 8 Implementation Plan

### Sprint Goal

Deliver a functional dashboard with explanation detail view, interactive graph, and
what-if simulator. Authentication and audit trail are deferred to Sprint 9.

### Task Breakdown

| # | Task | Files | Priority | Estimate |
|---|------|-------|----------|----------|
| T-001 | Project scaffold: `create-next-app`, Tailwind, shadcn/ui, TanStack Query | `dashboard/*` | P0 | 2h |
| T-002 | TypeScript types from Go models | `lib/api/types.ts` | P0 | 1h |
| T-003 | API client and BFF routes | `lib/api/client.ts`, `app/api/**` | P0 | 3h |
| T-004 | TanStack Query hooks | `lib/api/hooks.ts` | P0 | 1h |
| T-005 | Root layout (sidebar, header) | `app/layout.tsx`, `components/layout/*` | P0 | 2h |
| T-006 | Home page: recent explanations list | `app/page.tsx` | P1 | 2h |
| T-007 | Explanation detail page: SummaryCard, ConfidenceGauge | `app/explain/[id]/page.tsx`, `components/explanation/*` | P0 | 3h |
| T-008 | BreakdownChart and DriverRanking | `components/explanation/BreakdownChart.tsx`, `DriverRanking.tsx` | P0 | 3h |
| T-009 | NarrativeViewer with level/language toggle | `components/explanation/NarrativeViewer.tsx` | P1 | 2h |
| T-010 | Graph layout engine (dagre) | `lib/graph/layout.ts` | P0 | 3h |
| T-011 | React Flow graph page with custom nodes | `app/explain/[id]/graph/page.tsx`, `components/graph/*` | P0 | 4h |
| T-012 | Graph export (PNG, SVG, DOT, Mermaid) | `components/graph/GraphControls.tsx` | P1 | 2h |
| T-013 | What-if simulator: sliders and mutation | `app/explain/[id]/whatif/page.tsx`, `components/whatif/*` | P0 | 4h |
| T-014 | ComparisonView and DiffTable | `components/whatif/ComparisonView.tsx`, `DiffTable.tsx` | P0 | 2h |
| T-015 | API playground page | `app/playground/page.tsx` | P2 | 3h |
| T-016 | Unit tests for lib/ | `tests/unit/lib/*` | P0 | 2h |
| T-017 | Component tests | `tests/unit/components/*` | P1 | 3h |
| T-018 | E2E tests for critical paths | `tests/e2e/*` | P1 | 3h |
| T-019 | Vercel deployment configuration | `vercel.json`, env vars | P0 | 1h |

**Total estimate:** ~45 hours (roughly 6 working days for one developer).

### Sprint 8 Definition of Done

- A user can paste a JSON payload in the playground and see the full explanation rendered.
- The explanation detail page shows summary, breakdown chart, confidence gauge, narrative,
  and metadata.
- The graph page renders an interactive DAG with custom-styled nodes.
- The what-if simulator allows adjusting component values and displays the sensitivity
  result.
- BFF routes proxy all requests to the Go API without exposing the backend URL.
- Unit and component tests pass with 80%+ coverage on `lib/`.
- The dashboard is deployed to Vercel and accessible via a preview URL.

### Deferred to Sprint 9

- Authentication (NextAuth.js + Google OAuth)
- Audit trail page (paginated table, filters, CSV/PDF export)
- Live monitoring page (polling feed, alert badges)
- Settings page (language, theme persistence)
- Dark mode

---

## Appendix A: Package Dependencies

```json
{
  "dependencies": {
    "next": "^14.2",
    "react": "^18.3",
    "react-dom": "^18.3",
    "@xyflow/react": "^12.0",
    "@dagrejs/dagre": "^1.1",
    "@tanstack/react-query": "^5.50",
    "recharts": "^2.12",
    "next-auth": "^5.0.0-beta",
    "zod": "^3.23",
    "use-debounce": "^10.0",
    "clsx": "^2.1",
    "tailwind-merge": "^2.3"
  },
  "devDependencies": {
    "typescript": "^5.5",
    "@types/react": "^18.3",
    "@types/node": "^20",
    "tailwindcss": "^3.4",
    "postcss": "^8.4",
    "autoprefixer": "^10.4",
    "vitest": "^1.6",
    "@testing-library/react": "^16.0",
    "@testing-library/jest-dom": "^6.4",
    "@playwright/test": "^1.45",
    "eslint": "^8.57",
    "eslint-config-next": "^14.2",
    "prettier": "^3.3",
    "prettier-plugin-tailwindcss": "^0.6",
    "@tanstack/react-virtual": "^3.8"
  }
}
```

## Appendix B: Environment Setup

```bash
# Create the project
npx create-next-app@latest dashboard \
  --typescript --tailwind --eslint --app --src-dir=false \
  --import-alias="@/*"

cd dashboard

# Install dependencies
npm install @xyflow/react @dagrejs/dagre @tanstack/react-query \
  recharts next-auth@beta zod use-debounce clsx tailwind-merge

npm install -D vitest @testing-library/react @testing-library/jest-dom \
  @playwright/test @tanstack/react-virtual prettier prettier-plugin-tailwindcss

# Add shadcn/ui
npx shadcn-ui@latest init
npx shadcn-ui@latest add button card dialog dropdown-menu input label \
  select skeleton slider table tabs tooltip

# Create env file
cp .env.local.example .env.local
# Edit .env.local with actual values
```

## Appendix C: Key Configuration Files

### next.config.ts

```typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",          // Required for Cloud Run Docker deployment
  reactStrictMode: true,
  experimental: {
    typedRoutes: true,           // Type-safe <Link href="..."> at compile time
  },
};

export default nextConfig;
```

### tsconfig.json (key overrides)

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "paths": {
      "@/*": ["./*"]
    }
  }
}
```

### tailwind.config.ts

```typescript
import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class",
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Semantic colors for node types
        "node-output": { DEFAULT: "#2563eb", light: "#dbeafe" },
        "node-input": { DEFAULT: "#16a34a", light: "#dcfce7" },
        "node-computed": { DEFAULT: "#6b7280", light: "#f3f4f6" },
        "node-missing": { DEFAULT: "#dc2626", light: "#fef2f2" },
      },
    },
  },
  plugins: [],
};

export default config;
```
