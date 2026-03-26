# Explainable Engine — Product Specification

## Visione

Trasformare qualsiasi output numerico o decisionale in una catena causale esplicita, interrogabile e verificabile.

## Design Principle

> "Non spiegare il risultato. Rappresenta il processo che lo genera."

---

## Personas

### 1. Quant / Trader (AIP)
- "Perche questo segnale e AGGRESSIVE?"
- "Quale indicatore pesa di piu?"

### 2. Analyst / Sales
- "Spiegami questo numero in modo semplice"
- "Posso fidarmi?"

### 3. System / Backend
- Bisogno di API -> explanation machine-readable

---

## Core Capabilities

1. Breakdown numerico
2. Causal graph (dipendenze)
3. Confidence propagation
4. What-if explanation
5. Natural explanation (LLM opzionale)

---

## Epics & User Stories

### EPIC 1 — Basic Explainability

**US-001 — Breakdown di un valore** (MVP)
- As a quant, I want vedere come e costruito un valore, so that posso fidarmi del risultato
- AC: input valore + componenti → output lista componenti, peso (%), contributo numerico

**US-002 — Tracciare le dipendenze** (MVP)
- As a developer, I want vedere da dove nasce un valore, so that posso fare debug
- AC: output include parent nodes, depth traversal, full dependency tree

**US-003 — Visualizzazione grafo** (post-MVP)
- As a user, I want vedere il grafo decisionale, so that capisco visivamente
- AC: nodes + edges, direction chiara, livelli (layers)

### EPIC 2 — Confidence & Reliability

**US-004 — Confidence score** (MVP)
- As a user, I want sapere quanto e affidabile un risultato
- AC: ogni nodo ha confidence score, aggregazione finale weighted propagation

**US-005 — Missing data impact** (post-MVP)
- As a user, I want sapere quanto incide un dato mancante
- AC: identificazione nodi mancanti, impatto su output (%), warning esplicito

### EPIC 3 — What-if Explainability (post-MVP)

**US-006 — Sensitivity analysis**
**US-007 — Critical drivers**

### EPIC 4 — API-first System

**US-008 — Submit explanation request** (MVP)
- As a system (AIP), I want inviare un payload, so that ottengo explanation
- AC: POST /explain, input JSON generico, response strutturata

**US-009 — Retrieve explanation** (MVP)
- As a system, I want recuperare explanation gia calcolate
- AC: GET /explain/{id}, caching support

### EPIC 5 — Natural Language Layer (v1.5, post-MVP)

**US-010 — Human explanation**

---

## MVP Scope

- US-001 (breakdown)
- US-002 (dependencies)
- US-004 (confidence base)
- US-008 (POST explain)
- US-009 (GET explain)

## Explicitly OUT of MVP

- LLM explanation
- UI complessa
- Simulator integration
- Fancy graph UI

---

## Non-Functional Requirements

- Performance: < 200ms simple, < 1s complex graph
- Scalability: stateless API, K8s-ready
- Auditability: ogni explanation persistita e versionata
- Determinismo: stessa input → stessa output

---

## Input Model

```json
{
  "target": "market_regime_score",
  "value": 0.72,
  "components": [
    {
      "name": "trend_strength",
      "value": 0.8,
      "weight": 0.4,
      "confidence": 0.9
    }
  ]
}
```

## Output Model

```json
{
  "final_value": 0.72,
  "confidence": 0.81,
  "breakdown": [...],
  "top_drivers": [...],
  "missing_impact": 0.12,
  "graph": {...}
}
```

## Explainability Types

- additive (weighted sum) — MVP
- rule-based — post-MVP
- normalized scoring — post-MVP
- black-box wrapper — post-MVP
