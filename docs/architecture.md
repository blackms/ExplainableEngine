# Explainable Engine -- Architecture Document

**Versione:** 1.0.0
**Data:** 2026-03-26
**Stato:** Approvato per Sprint 1

> **Principio guida:** "Non spiegare il risultato. Rappresenta il processo che lo genera."

---

## Indice

1. [Visione e Obiettivi](#1-visione-e-obiettivi)
2. [Graph-First Data Model](#2-graph-first-data-model)
3. [Struttura del Progetto](#3-struttura-del-progetto)
4. [Moduli di Dominio](#4-moduli-di-dominio)
5. [Tipi di Spiegazione](#5-tipi-di-spiegazione)
6. [API Layer](#6-api-layer)
7. [Storage Strategy](#7-storage-strategy)
8. [Requisiti Non-Funzionali](#8-requisiti-non-funzionali)
9. [Tech Stack](#9-tech-stack)
10. [Architecture Decision Records](#10-architecture-decision-records)
11. [Diagrammi](#11-diagrammi)
12. [Piano di Implementazione Sprint 1](#12-piano-di-implementazione-sprint-1)

---

## 1. Visione e Obiettivi

Explainable Engine e' un sistema standalone che trasforma qualsiasi output numerico o decisionale in una **catena causale esplicita, interrogabile e verificabile**.

Il sistema non e' un layer di "post-hoc explanation". E' una struttura dati che **rappresenta il processo computazionale stesso**, rendendolo navigabile, persistente e deterministico.

### Obiettivi primari

- Rendere ogni calcolo numerico ispezionabile fino alle foglie (input grezzi)
- Garantire determinismo: stesso input produce sempre lo stesso grafo esplicativo
- Fornire API per navigare, interrogare e verificare qualsiasi risultato
- Supportare diversi tipi di computazione (additiva, rule-based, normalizzata, black-box)

---

## 2. Graph-First Data Model

Questo e' il **cuore** dell'intero sistema. Ogni spiegazione e' un grafo diretto aciclico (DAG) dove i nodi rappresentano valori e gli archi rappresentano trasformazioni.

### 2.1 Node

```python
class NodeType(str, Enum):
    INPUT = "INPUT"           # Dato grezzo, foglia del grafo
    COMPUTED = "COMPUTED"     # Valore calcolato da altri nodi
    OUTPUT = "OUTPUT"         # Risultato finale esposto all'utente
    MISSING = "MISSING"       # Dato atteso ma non disponibile

class ComputationType(str, Enum):
    ADDITIVE = "ADDITIVE"           # Somma pesata -- MVP
    RULE_BASED = "RULE_BASED"       # Regole condizionali -- post-MVP
    NORMALIZED = "NORMALIZED"       # Scoring normalizzato -- post-MVP
    BLACKBOX = "BLACKBOX"           # Wrapper per sistemi opachi -- post-MVP

class Node(BaseModel):
    id: str                              # UUID v4 o identificativo dominio
    name: str                            # Nome leggibile (es. "punteggio_credito")
    value: float                         # Valore numerico del nodo
    confidence: float = Field(ge=0, le=1)  # Confidenza [0, 1]
    node_type: NodeType
    computation_type: ComputationType | None = None  # None per INPUT
    metadata: dict[str, Any] = {}        # Dati dominio-specifici
```

**Semantica dei tipi di nodo:**

| NodeType | Significato | Ha figli? | confidence |
|----------|-------------|-----------|------------|
| INPUT | Dato grezzo proveniente dall'esterno | No | 1.0 (o definita dal chiamante) |
| COMPUTED | Valore derivato da altri nodi | Si | Propagata dai figli |
| OUTPUT | Risultato finale del grafo | Si | Propagata dai figli |
| MISSING | Input atteso ma non disponibile | No | 0.0 |

### 2.2 Edge

```python
class TransformationType(str, Enum):
    WEIGHTED_SUM = "WEIGHTED_SUM"       # target += source.value * weight
    NORMALIZATION = "NORMALIZATION"     # target = source.value / weight (dove weight = max)
    THRESHOLD = "THRESHOLD"             # target = 1 if source.value >= weight else 0
    CUSTOM = "CUSTOM"                   # Trasformazione descritta in metadata

class Edge(BaseModel):
    source_id: str
    target_id: str
    weight: float                        # Peso/parametro della trasformazione
    transformation_type: TransformationType
    metadata: dict[str, Any] = {}        # Es. {"formula": "x * 0.3 + bias"}
```

**Invariante:** un Edge `(A -> B)` significa "A contribuisce al calcolo di B". La direzione e' **dalla causa all'effetto**.

### 2.3 ExplanationGraph

```python
class ExplanationGraph(BaseModel):
    id: str                              # UUID v4
    target_node_id: str                  # Nodo OUTPUT che questo grafo spiega
    nodes: list[Node]
    edges: list[Edge]
    created_at: datetime
    version: str                         # Semantic versioning del grafo
    deterministic_hash: str              # SHA-256 di (nodes_sorted + edges_sorted)

    @computed_field
    def node_index(self) -> dict[str, Node]:
        """Indice per accesso O(1) ai nodi per id."""
        return {n.id: n for n in self.nodes}
```

### 2.4 Calcolo del deterministic_hash

Il deterministic hash garantisce che lo stesso input produca sempre lo stesso hash. L'algoritmo e':

```
1. Serializzare ogni Node come tupla ordinata dei suoi campi (escluso metadata non-deterministico)
2. Ordinare i nodi per id
3. Serializzare ogni Edge come tupla ordinata dei suoi campi
4. Ordinare gli archi per (source_id, target_id)
5. Concatenare le rappresentazioni
6. Calcolare SHA-256 della stringa risultante
```

Campi inclusi nel hash:
- Node: id, name, value, confidence, node_type, computation_type
- Edge: source_id, target_id, weight, transformation_type

Campi esclusi: metadata, created_at (non deterministici per natura).

---

## 3. Struttura del Progetto

```
ExplainableEngine/
|-- docs/
|   |-- architecture.md          # Questo documento
|   |-- adr/                     # Architecture Decision Records
|-- src/
|   |-- explainable_engine/
|   |   |-- __init__.py
|   |   |-- main.py              # Entrypoint FastAPI
|   |   |-- core/
|   |   |   |-- __init__.py
|   |   |   |-- graph.py         # Costruzione, traversal, validazione grafo
|   |   |   |-- breakdown.py     # Scomposizione numerica (US-001)
|   |   |   |-- confidence.py    # Propagazione confidenza (US-004)
|   |   |   |-- dependencies.py  # Albero dipendenze (US-002)
|   |   |   |-- hashing.py       # Deterministic hash computation
|   |   |-- models/
|   |   |   |-- __init__.py
|   |   |   |-- node.py          # Node, NodeType, ComputationType
|   |   |   |-- edge.py          # Edge, TransformationType
|   |   |   |-- graph.py         # ExplanationGraph
|   |   |   |-- requests.py      # Request/Response DTOs per API
|   |   |-- api/
|   |   |   |-- __init__.py
|   |   |   |-- routes.py        # FastAPI router principale
|   |   |   |-- dependencies.py  # FastAPI Depends (storage, services)
|   |   |-- storage/
|   |   |   |-- __init__.py
|   |   |   |-- base.py          # Abstract storage interface
|   |   |   |-- memory.py        # In-memory storage (MVP)
|   |   |   |-- sqlite.py        # SQLite storage (persistence MVP)
|   |   |-- services/
|   |   |   |-- __init__.py
|   |   |   |-- explanation.py   # Orchestrazione: costruisci grafo, calcola, persisti
|-- tests/
|   |-- __init__.py
|   |-- conftest.py              # Fixtures condivise
|   |-- core/
|   |   |-- test_graph.py
|   |   |-- test_breakdown.py
|   |   |-- test_confidence.py
|   |   |-- test_dependencies.py
|   |   |-- test_hashing.py
|   |-- api/
|   |   |-- test_routes.py
|   |-- storage/
|   |   |-- test_memory.py
|   |   |-- test_sqlite.py
|-- pyproject.toml
|-- README.md
```

---

## 4. Moduli di Dominio

### 4.1 `core/graph.py` -- Motore del Grafo

Responsabilita':
- Costruzione del grafo da input strutturato
- Validazione strutturale (DAG, no cicli, no nodi orfani)
- Traversal (BFS/DFS dal nodo target verso le foglie)
- Risoluzione delle dipendenze (ordinamento topologico)

```python
class GraphEngine:
    """
    Costruisce e manipola ExplanationGraph.
    Utilizza NetworkX internamente per operazioni su grafi.
    """

    def build(self, nodes: list[Node], edges: list[Edge],
              target_node_id: str) -> ExplanationGraph:
        """
        Costruisce un ExplanationGraph validato.
        Raises: CycleDetectedError, OrphanNodeError, MissingTargetError
        """
        ...

    def validate(self, graph: ExplanationGraph) -> list[ValidationError]:
        """Validazione strutturale completa del grafo."""
        ...

    def get_ancestors(self, graph: ExplanationGraph,
                      node_id: str) -> list[Node]:
        """Tutti i nodi che contribuiscono (direttamente o indirettamente) a node_id."""
        ...

    def get_leaves(self, graph: ExplanationGraph) -> list[Node]:
        """Tutti i nodi foglia (INPUT o MISSING)."""
        ...

    def topological_sort(self, graph: ExplanationGraph) -> list[str]:
        """Ordine di computazione: dalle foglie al target."""
        ...
```

**Dettaglio implementativo:** internamente si usa `networkx.DiGraph`. La conversione da/a `ExplanationGraph` avviene ai bordi del modulo. Il modello Pydantic rimane la rappresentazione canonica.

### 4.2 `core/breakdown.py` -- Scomposizione Numerica (US-001)

Dato un nodo COMPUTED o OUTPUT, produce la scomposizione del suo valore nei contributi dei nodi figli diretti.

```python
class BreakdownResult(BaseModel):
    target_node_id: str
    target_value: float
    contributions: list[Contribution]
    remainder: float  # Differenza non spiegata (dovrebbe essere ~0 per ADDITIVE)

class Contribution(BaseModel):
    source_node_id: str
    source_node_name: str
    source_value: float
    weight: float
    contribution: float          # source_value * weight
    percentage: float            # contribution / target_value * 100

class BreakdownEngine:
    def breakdown(self, graph: ExplanationGraph,
                  node_id: str) -> BreakdownResult:
        """
        Scompone il valore di node_id nei contributi dei predecessori diretti.
        Per ADDITIVE: contribution = source.value * edge.weight
        """
        ...

    def full_breakdown(self, graph: ExplanationGraph,
                       node_id: str) -> list[BreakdownResult]:
        """
        Scomposizione ricorsiva fino alle foglie.
        Restituisce un BreakdownResult per ogni nodo intermedio.
        """
        ...
```

**Regola di calcolo per ADDITIVE (MVP):**

```
target.value = SUM(source_i.value * edge_i.weight)
contribution_i = source_i.value * edge_i.weight
percentage_i = (contribution_i / target.value) * 100
remainder = target.value - SUM(contribution_i)  # Deve essere ~0
```

### 4.3 `core/confidence.py` -- Propagazione Confidenza (US-004)

La confidenza si propaga dalle foglie verso il target. L'algoritmo e' una **media pesata dei figli**, dove i pesi sono i pesi degli archi.

```python
class ConfidenceEngine:
    def propagate(self, graph: ExplanationGraph) -> ExplanationGraph:
        """
        Ricalcola la confidenza di tutti i nodi COMPUTED e OUTPUT
        a partire dai nodi INPUT/MISSING.
        Restituisce un nuovo ExplanationGraph con confidenze aggiornate.
        """
        ...

    def node_confidence(self, graph: ExplanationGraph,
                        node_id: str) -> float:
        """Confidenza di un singolo nodo."""
        ...
```

**Algoritmo di propagazione:**

```
Per ogni nodo in ordine topologico (foglie -> target):
    Se node_type == INPUT:
        confidence = confidence (definita dal chiamante, default 1.0)
    Se node_type == MISSING:
        confidence = 0.0
    Se node_type in (COMPUTED, OUTPUT):
        predecessors = [(source, edge) for source, edge in incoming_edges]
        total_weight = SUM(|edge.weight| for _, edge in predecessors)
        confidence = SUM(source.confidence * |edge.weight| / total_weight
                         for source, edge in predecessors)
```

**Proprieta' garantite:**
- Se tutti gli input hanno confidence = 1.0, il target avra' confidence = 1.0
- Se un input e' MISSING (confidence = 0.0), la confidence del target diminuisce proporzionalmente al peso di quell'input
- La confidence e' sempre in [0, 1]

### 4.4 `core/dependencies.py` -- Albero Dipendenze (US-002)

Costruisce una vista ad albero delle dipendenze di un nodo, navigabile livello per livello.

```python
class DependencyNode(BaseModel):
    node_id: str
    node_name: str
    node_type: NodeType
    value: float
    confidence: float
    depth: int
    children: list["DependencyNode"] = []

class DependencyEngine:
    def build_tree(self, graph: ExplanationGraph,
                   node_id: str, max_depth: int | None = None) -> DependencyNode:
        """
        Costruisce l'albero delle dipendenze radicato in node_id.
        max_depth limita la profondita' (None = completo).
        """
        ...

    def list_at_depth(self, graph: ExplanationGraph,
                      node_id: str, depth: int) -> list[Node]:
        """Tutti i nodi a una specifica profondita' nel sotto-grafo di node_id."""
        ...
```

**Nota:** il grafo e' un DAG, non un albero. Un nodo puo' apparire in piu' rami dell'albero delle dipendenze. Questo e' intenzionale: ogni percorso causale viene rappresentato esplicitamente. Il campo `DependencyNode.node_id` permette di riconoscere i duplicati.

### 4.5 `core/hashing.py` -- Deterministic Hash

```python
class HashEngine:
    def compute_hash(self, nodes: list[Node], edges: list[Edge]) -> str:
        """
        Calcola SHA-256 deterministico.
        Ordina nodi per id, archi per (source_id, target_id).
        Serializza solo campi deterministici.
        """
        ...
```

---

## 5. Tipi di Spiegazione

### 5.1 Additive (WEIGHTED_SUM) -- MVP Sprint 1

Il tipo piu' comune. Il valore di un nodo e' la somma pesata dei suoi predecessori.

```
score_finale = credito * 0.4 + reddito * 0.35 + storico * 0.25
```

Grafo risultante:
```
[credito: 750] --0.4--> [score_finale: 587.5]
[reddito: 500] --0.35-> [score_finale: 587.5]
[storico: 200] --0.25-> [score_finale: 587.5]
```

Breakdown:
```
score_finale = 587.5
  credito:  750 * 0.4  = 300.0  (51.1%)
  reddito:  500 * 0.35 = 175.0  (29.8%)
  storico:  200 * 0.25 =  50.0  (8.5%)
  ---
  Somma contributi: 525.0
  Remainder: 62.5  (se presente un bias o termine costante)
```

### 5.2 Rule-Based (THRESHOLD) -- Post-MVP

Decisioni basate su regole condizionali. Ogni regola diventa un nodo COMPUTED con edge di tipo THRESHOLD.

### 5.3 Normalized Scoring (NORMALIZATION) -- Post-MVP

Valori normalizzati rispetto a un massimo o a una distribuzione. Edge di tipo NORMALIZATION dove weight rappresenta il fattore di normalizzazione.

### 5.4 Black-Box Wrapper (BLACKBOX) -- Post-MVP

Per sistemi opachi (ML models, API esterne). Il nodo ha `computation_type = BLACKBOX` e la spiegazione si basa su feature importance o SHAP values passati nel metadata.

---

## 6. API Layer

### 6.1 Endpoint Design

Tutti gli endpoint sono sotto il prefisso `/api/v1/`.

```
POST   /api/v1/explanations              # Crea un nuovo ExplanationGraph
GET    /api/v1/explanations/{id}          # Recupera un grafo per id
GET    /api/v1/explanations/{id}/breakdown/{node_id}       # Breakdown numerico
GET    /api/v1/explanations/{id}/dependencies/{node_id}    # Albero dipendenze
GET    /api/v1/explanations/{id}/confidence/{node_id}      # Dettaglio confidenza
GET    /api/v1/explanations/{id}/verify                    # Verifica determinismo
DELETE /api/v1/explanations/{id}          # Rimuovi un grafo
```

### 6.2 Request/Response Models

**Creazione grafo:**

```python
class CreateExplanationRequest(BaseModel):
    nodes: list[NodeInput]
    edges: list[EdgeInput]
    target_node_id: str
    version: str = "1.0.0"

class NodeInput(BaseModel):
    id: str
    name: str
    value: float
    confidence: float = 1.0
    node_type: NodeType
    computation_type: ComputationType | None = None
    metadata: dict[str, Any] = {}

class EdgeInput(BaseModel):
    source_id: str
    target_id: str
    weight: float
    transformation_type: TransformationType = TransformationType.WEIGHTED_SUM
    metadata: dict[str, Any] = {}
```

**Response standard:**

```python
class ExplanationResponse(BaseModel):
    id: str
    target_node_id: str
    target_value: float
    target_confidence: float
    deterministic_hash: str
    created_at: datetime
    version: str
    node_count: int
    edge_count: int

class BreakdownResponse(BaseModel):
    target_node_id: str
    target_value: float
    contributions: list[ContributionResponse]
    remainder: float

class ContributionResponse(BaseModel):
    source_node_id: str
    source_node_name: str
    source_value: float
    weight: float
    contribution: float
    percentage: float
```

### 6.3 Error Handling

Errori strutturati con codice dominio:

```python
class ExplainableError(BaseModel):
    code: str          # Es. "CYCLE_DETECTED", "NODE_NOT_FOUND"
    message: str
    details: dict[str, Any] = {}
```

Mapping HTTP:
| Codice dominio | HTTP Status |
|----------------|-------------|
| GRAPH_NOT_FOUND | 404 |
| NODE_NOT_FOUND | 404 |
| CYCLE_DETECTED | 422 |
| ORPHAN_NODE | 422 |
| INVALID_GRAPH | 422 |
| HASH_MISMATCH | 409 |

---

## 7. Storage Strategy

### 7.1 Abstract Interface

```python
class ExplanationStorage(ABC):
    @abstractmethod
    async def save(self, graph: ExplanationGraph) -> str:
        """Persiste il grafo. Restituisce l'id."""
        ...

    @abstractmethod
    async def get(self, graph_id: str) -> ExplanationGraph | None:
        """Recupera il grafo per id. None se non esiste."""
        ...

    @abstractmethod
    async def delete(self, graph_id: str) -> bool:
        """Elimina il grafo. True se esisteva."""
        ...

    @abstractmethod
    async def exists(self, graph_id: str) -> bool:
        """Verifica esistenza."""
        ...
```

### 7.2 In-Memory Storage (MVP)

- `dict[str, ExplanationGraph]` in memoria
- Nessuna persistenza tra restart
- Utilizzato per testing e sviluppo rapido

### 7.3 SQLite Storage (Persistence MVP)

Schema:

```sql
CREATE TABLE explanation_graphs (
    id TEXT PRIMARY KEY,
    target_node_id TEXT NOT NULL,
    graph_json TEXT NOT NULL,          -- ExplanationGraph serializzato
    deterministic_hash TEXT NOT NULL,
    version TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE(deterministic_hash, version)
);

CREATE INDEX idx_graphs_hash ON explanation_graphs(deterministic_hash);
CREATE INDEX idx_graphs_created ON explanation_graphs(created_at);
```

**Scelta di design:** il grafo viene serializzato come JSON in una singola colonna. Questo e' sufficiente per il MVP e permette query per id e hash. Se in futuro servono query sui singoli nodi/archi, si migrera' a un modello relazionale completo o a un graph database.

### 7.4 Selezione dello Storage

Lo storage e' iniettato tramite FastAPI `Depends`. La configurazione avviene via variabile d'ambiente:

```
EXPLAINABLE_STORAGE=memory   # default per dev
EXPLAINABLE_STORAGE=sqlite   # per persistence
EXPLAINABLE_SQLITE_PATH=./data/explanations.db
```

---

## 8. Requisiti Non-Funzionali

### 8.1 Performance

| Scenario | Obiettivo | Metrica |
|----------|-----------|---------|
| Grafo semplice (< 20 nodi) | < 200ms | Tempo totale API (creazione + breakdown) |
| Grafo complesso (< 500 nodi) | < 1s | Tempo totale API (creazione + breakdown) |
| Lookup per id | < 50ms | Tempo di risposta GET |

### 8.2 Determinismo

**Garanzia assoluta:** dati gli stessi nodi e archi (stessi id, valori, pesi), il sistema DEVE produrre:
- Lo stesso `deterministic_hash`
- Lo stesso `BreakdownResult`
- Lo stesso `DependencyNode` tree
- La stessa `confidence` propagata

Questa garanzia e' verificabile chiamando `GET /api/v1/explanations/{id}/verify`, che ricalcola il hash e lo confronta con quello persistito.

### 8.3 Statelessness

- Nessuno stato in-process tra richieste (tranne lo storage iniettato)
- Ogni richiesta e' autocontenuta
- Compatibile con deployment multi-replica su Kubernetes
- Health check: `GET /health`
- Readiness check: `GET /ready`

### 8.4 Persistenza e Versioning

- Ogni ExplanationGraph creato viene persistito
- Il campo `version` permette di tracciare l'evoluzione dello schema
- Il `deterministic_hash` + `version` formano una chiave unica: stesso calcolo, stessa versione, un solo grafo

---

## 9. Tech Stack

| Componente | Tecnologia | Motivazione |
|------------|------------|-------------|
| Linguaggio | Python 3.11+ | Ecosistema data science, typing avanzato |
| Web framework | FastAPI | Async, auto-documentazione OpenAPI, Pydantic nativo |
| Validazione | Pydantic v2 | Performance, validazione rigorosa, serializzazione |
| Grafo | NetworkX | Maturo, API completa per DAG, topological sort, cycle detection |
| Persistence | SQLite | Zero-config, file-based, sufficiente per MVP |
| Testing | pytest | Standard de facto, fixture system, parametrize |
| Async DB | aiosqlite | Accesso SQLite async-compatible per FastAPI |
| Build | pyproject.toml | Standard PEP 621 |

### Dipendenze esatte (pyproject.toml)

```toml
[project]
name = "explainable-engine"
version = "0.1.0"
requires-python = ">=3.11"
dependencies = [
    "fastapi>=0.110",
    "uvicorn[standard]>=0.27",
    "pydantic>=2.6",
    "networkx>=3.2",
    "aiosqlite>=0.20",
]

[project.optional-dependencies]
dev = [
    "pytest>=8.0",
    "pytest-asyncio>=0.23",
    "httpx>=0.27",         # Per TestClient async
    "ruff>=0.3",           # Linting + formatting
]
```

---

## 10. Architecture Decision Records

### ADR-001: Graph-First, Non Tree

**Contesto:** serve una struttura dati per rappresentare catene causali.

**Decisione:** DAG (Directed Acyclic Graph), non albero.

**Motivazione:**
- Un valore puo' contribuire a piu' nodi (fan-out). Esempio: il "reddito" influenza sia "score_credito" sia "score_rischio".
- Un nodo puo' ricevere contributi da piu' sorgenti (fan-in). Questo e' il caso standard di una somma pesata.
- Un albero forzerebbe la duplicazione dei nodi condivisi, perdendo la semantica "stesso dato, usi multipli".
- NetworkX offre supporto nativo per DAG con cycle detection e topological sort.

**Conseguenze:**
- L'albero delle dipendenze (US-002) e' una **vista** derivata dal grafo, non la struttura primaria.
- Un nodo puo' apparire in piu' rami dell'albero (deduplicazione a carico del consumer).

### ADR-002: API Stateless

**Contesto:** il sistema deve essere deployabile su Kubernetes con scaling orizzontale.

**Decisione:** API completamente stateless. Nessuna sessione, nessun stato in-process.

**Motivazione:**
- Ogni replica puo' servire qualsiasi richiesta
- Lo storage e' l'unico stato condiviso (SQLite per MVP, migrabile a PostgreSQL)
- Nessun bisogno di sticky sessions o session affinity
- Semplifica il testing: ogni test e' indipendente

**Conseguenze:**
- Il grafo completo viene caricato da storage a ogni richiesta (accettabile per dimensioni MVP)
- Per grafi molto grandi, si considerera' un layer di caching (Redis) in futuro

### ADR-003: Determinismo Garantito

**Contesto:** per audit e compliance, lo stesso input deve produrre lo stesso output.

**Decisione:** deterministic hash basato su SHA-256 dei campi strutturali, ordinamento canonico.

**Motivazione:**
- L'ordinamento canonico (nodi per id, archi per source+target) elimina la dipendenza dall'ordine di inserimento
- SHA-256 e' collision-resistant e standard
- I campi non deterministici (metadata, created_at) sono esclusi dal hash

**Conseguenze:**
- Il metadata non influenza l'identita' del grafo
- Due grafi con stessi nodi/archi ma metadata diversi avranno lo stesso hash
- Il campo `version` e' parte della chiave unica (hash + version), non del hash stesso

### ADR-004: Confidenza come Media Pesata

**Contesto:** serve un algoritmo per propagare la confidenza dai nodi foglia al target.

**Decisione:** media pesata, dove i pesi sono i valori assoluti dei pesi degli archi.

**Formula:**
```
confidence(node) = SUM(confidence(child_i) * |weight_i|) / SUM(|weight_i|)
```

**Motivazione:**
- Intuitivo: un input con peso maggiore ha piu' influenza sulla confidenza
- Composabile: funziona ricorsivamente su grafi di qualsiasi profondita'
- Conservativo: un singolo MISSING node riduce la confidenza proporzionalmente
- Semplice da spiegare all'utente finale

**Conseguenze:**
- Un nodo MISSING (confidence=0) con peso alto abbassa significativamente la confidenza
- La confidenza del target e' sempre <= max(confidenza figli) (per costruzione della media)
- Non cattura correlazioni tra input (accettabile per MVP)

### ADR-005: Storage Ibrido Memory/SQLite

**Contesto:** serve persistenza per audit, ma anche velocita' per sviluppo.

**Decisione:** interfaccia astratta `ExplanationStorage` con due implementazioni: in-memory e SQLite.

**Motivazione:**
- In-memory per test unitari e sviluppo locale rapido
- SQLite per persistenza senza infrastruttura esterna
- L'interfaccia astratta permette di aggiungere PostgreSQL, Redis o altri backend senza modificare il dominio

**Conseguenze:**
- Il grafo viene serializzato come JSON per SQLite (semplice ma non queryabile sui singoli nodi)
- Per query avanzate sui nodi, si passera' a uno schema relazionale completo o graph DB

---

## 11. Diagrammi

### 11.1 Component Diagram (C4 - Level 2)

```
+--------------------------------------------------+
|              Explainable Engine                    |
|                                                    |
|  +-------------+    +-------------------------+    |
|  |   API       |    |   Services              |    |
|  |  (FastAPI)  |--->|  ExplanationService     |    |
|  |  routes.py  |    |  (orchestrazione)       |    |
|  +-------------+    +-------------------------+    |
|                          |          |              |
|                          v          v              |
|              +----------+  +----------------+      |
|              |  Core    |  |   Storage      |      |
|              |----------|  |----------------|      |
|              | graph    |  | base (ABC)     |      |
|              | breakdown|  | memory         |      |
|              | confidence| | sqlite         |      |
|              | deps     |  +----------------+      |
|              | hashing  |                          |
|              +----------+                          |
|                    |                               |
|                    v                               |
|              +----------+                          |
|              | Models   |                          |
|              | (Pydantic)|                         |
|              +----------+                          |
+--------------------------------------------------+
```

### 11.2 Data Flow -- Creazione di una Spiegazione

```
Client                API              Service          Core           Storage
  |                    |                  |               |               |
  |-- POST /explain -->|                  |               |               |
  |                    |-- create() ----->|               |               |
  |                    |                  |-- build() --->|               |
  |                    |                  |               |-- validate    |
  |                    |                  |               |-- topo_sort   |
  |                    |                  |               |-- propagate   |
  |                    |                  |               |   confidence  |
  |                    |                  |               |-- compute     |
  |                    |                  |               |   hash        |
  |                    |                  |<-- graph -----|               |
  |                    |                  |-- save() -----|-------------->|
  |                    |                  |<-- id --------|---------------|
  |                    |<-- response -----|               |               |
  |<-- 201 Created ----|                  |               |               |
```

### 11.3 Data Flow -- Breakdown Query

```
Client                API              Service          Core           Storage
  |                    |                  |               |               |
  |-- GET /breakdown ->|                  |               |               |
  |                    |-- breakdown() -->|               |               |
  |                    |                  |-- get() ------|-------------->|
  |                    |                  |<-- graph -----|---------------|
  |                    |                  |-- breakdown ->|               |
  |                    |                  |               |-- find node   |
  |                    |                  |               |-- get edges   |
  |                    |                  |               |-- compute     |
  |                    |                  |               |   contributions
  |                    |                  |<-- result ----|               |
  |                    |<-- response -----|               |               |
  |<-- 200 OK ---------|                  |               |               |
```

---

## 12. Piano di Implementazione Sprint 1

### Obiettivo Sprint 1

Implementare il flusso completo per computazioni **additive** (WEIGHTED_SUM): creazione grafo, breakdown numerico, albero dipendenze, propagazione confidenza.

### Task breakdown

| # | Task | File | Priorita' |
|---|------|------|-----------|
| 1 | Definire modelli Pydantic (Node, Edge, ExplanationGraph) | `models/*.py` | P0 |
| 2 | Implementare GraphEngine (build, validate, traversal) | `core/graph.py` | P0 |
| 3 | Implementare HashEngine | `core/hashing.py` | P0 |
| 4 | Implementare BreakdownEngine (solo ADDITIVE) | `core/breakdown.py` | P0 |
| 5 | Implementare ConfidenceEngine | `core/confidence.py` | P0 |
| 6 | Implementare DependencyEngine | `core/dependencies.py` | P1 |
| 7 | Implementare InMemoryStorage | `storage/memory.py` | P0 |
| 8 | Implementare ExplanationService | `services/explanation.py` | P0 |
| 9 | Implementare API routes | `api/routes.py` | P1 |
| 10 | Implementare SQLiteStorage | `storage/sqlite.py` | P1 |
| 11 | Test unitari per ogni modulo core | `tests/core/*.py` | P0 |
| 12 | Test integrazione API | `tests/api/*.py` | P1 |
| 13 | Setup progetto (pyproject.toml, CI) | root | P0 |

### Definition of Done -- Sprint 1

- [ ] Un client puo' inviare nodi e archi via POST e ricevere un ExplanationGraph persistito
- [ ] Il breakdown di un nodo OUTPUT restituisce i contributi percentuali dei suoi input
- [ ] L'albero delle dipendenze e' navigabile livello per livello
- [ ] La confidenza si propaga correttamente (verificato con test parametrizzati)
- [ ] Lo stesso input produce sempre lo stesso deterministic_hash
- [ ] Copertura test >= 90% sui moduli core
- [ ] Tempo di risposta < 200ms per grafi con < 20 nodi

---

## Appendice A: Glossario

| Termine | Definizione |
|---------|-------------|
| **DAG** | Directed Acyclic Graph -- grafo orientato senza cicli |
| **Nodo foglia** | Nodo senza predecessori (INPUT o MISSING) |
| **Nodo target** | Nodo OUTPUT la cui spiegazione e' richiesta |
| **Breakdown** | Scomposizione del valore di un nodo nei contributi dei suoi predecessori |
| **Confidenza** | Misura [0,1] di quanto un valore e' affidabile, basata sulla completezza degli input |
| **Deterministic hash** | Hash SHA-256 che garantisce: stessi nodi+archi = stesso hash |
| **Fan-in** | Numero di archi entranti in un nodo |
| **Fan-out** | Numero di archi uscenti da un nodo |
