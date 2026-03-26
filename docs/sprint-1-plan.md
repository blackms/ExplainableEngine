# Sprint 1 Plan -- Explainable Engine MVP

**Sprint Goal:** Costruire il nucleo funzionante del sistema di spiegabilita: dato un output numerico/decisionale, produrre una catena causale esplicita, interrogabile e verificabile tramite API REST.

**Durata:** 10 giorni lavorativi
**Tech Stack:** Python 3.11+, FastAPI, Pydantic v2, NetworkX, SQLite, pytest

---

## Indice

1. [Task Breakdown](#1-task-breakdown)
2. [Definition of Done](#2-definition-of-done)
3. [Rischi e Mitigazioni](#3-rischi-e-mitigazioni)
4. [Metriche di Successo](#4-metriche-di-successo)

---

## 1. Task Breakdown

### Fase 1: Fondamenta (Giorno 1--2)

#### T-001 -- Project Setup

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Creare la struttura del progetto Python con `pyproject.toml`, dipendenze, directory layout, configurazione linter/formatter, e Makefile con comandi di base (`make test`, `make run`, `make lint`). |
| **User Story** | Infrastruttura (nessuna US specifica -- prerequisito per tutte) |
| **Complessita** | S |
| **Dipendenze** | Nessuna |

**Acceptance Criteria:**
- [ ] `pyproject.toml` configurato con tutte le dipendenze (fastapi, uvicorn, pydantic>=2.0, networkx, pytest, httpx, ruff)
- [ ] Struttura directory creata: `src/explainable_engine/`, `tests/`, `docs/`
- [ ] `python -m pytest` eseguibile senza errori (anche con 0 test)
- [ ] `uvicorn src.explainable_engine.main:app` avvia il server senza errori
- [ ] Makefile funzionante con target `test`, `run`, `lint`

**File da creare:**
```
pyproject.toml
Makefile
src/explainable_engine/__init__.py
src/explainable_engine/main.py           # FastAPI app factory
src/explainable_engine/config.py         # Settings con pydantic-settings
tests/__init__.py
tests/conftest.py
```

---

#### T-002 -- Modelli dati core (Pydantic v2)

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Definire i modelli Pydantic v2 che rappresentano il grafo esplicativo: nodo, arco, grafo, richiesta di spiegazione e risposta. Questi modelli sono il contratto dati di tutto il sistema. |
| **User Story** | US-001, US-002, US-004 (modelli condivisi) |
| **Complessita** | M |
| **Dipendenze** | T-001 |

**Acceptance Criteria:**
- [ ] `Node` contiene: `id`, `label`, `value` (float), `unit` (opzionale), `confidence` (float, 0.0--1.0), `metadata` (dict opzionale)
- [ ] `Edge` contiene: `source`, `target`, `weight` (float, 0.0--1.0), `relation_type` (enum: "contributes_to", "depends_on", "derived_from")
- [ ] `ExplanationGraph` contiene: lista di `Node`, lista di `Edge`, `root_node_id`
- [ ] `ExplainRequest` contiene: `value` (float), `label` (str), `components` (lista di componenti con peso e valore), `metadata` (opzionale)
- [ ] `ExplainResponse` contiene: `id` (UUID), `timestamp`, `graph` (ExplanationGraph), `breakdown` (lista contributi), `confidence` (float aggregato), `dependency_tree` (struttura ad albero)
- [ ] Tutti i modelli hanno `model_config` con `json_schema_extra` contenente un esempio valido
- [ ] Validazione custom: la somma dei pesi dei componenti in `ExplainRequest` deve essere <= 1.0 (con warning, non errore, per flessibilita)
- [ ] Serializzazione/deserializzazione JSON round-trip verificata con test

**File da creare:**
```
src/explainable_engine/models/__init__.py
src/explainable_engine/models/node.py
src/explainable_engine/models/edge.py
src/explainable_engine/models/graph.py
src/explainable_engine/models/request.py
src/explainable_engine/models/response.py
tests/test_models.py
```

---

#### T-003 -- Graph Engine (wrapper NetworkX)

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Implementare il motore grafo che costruisce un `networkx.DiGraph` a partire dai modelli Pydantic. Fornisce operazioni di base: aggiunta nodi/archi, ricerca nodi, attraversamento BFS/DFS, estrazione sotto-grafi. |
| **User Story** | US-002 (prerequisito per la risoluzione delle dipendenze) |
| **Complessita** | M |
| **Dipendenze** | T-002 |

**Acceptance Criteria:**
- [ ] `GraphEngine.build(explanation_graph: ExplanationGraph) -> nx.DiGraph` costruisce il grafo correttamente
- [ ] `GraphEngine.get_ancestors(node_id) -> list[str]` restituisce tutti i nodi antenati
- [ ] `GraphEngine.get_descendants(node_id) -> list[str]` restituisce tutti i nodi discendenti
- [ ] `GraphEngine.get_subtree(node_id) -> ExplanationGraph` estrae il sotto-grafo radicato in un nodo
- [ ] `GraphEngine.topological_sort() -> list[str]` restituisce l'ordine topologico
- [ ] Gestione errore per grafi ciclici (deve lanciare `CyclicGraphError`)
- [ ] Gestione errore per nodi/archi non esistenti (deve lanciare `NodeNotFoundError`)

**File da creare:**
```
src/explainable_engine/engine/__init__.py
src/explainable_engine/engine/graph_engine.py
src/explainable_engine/engine/exceptions.py
tests/test_graph_engine.py
```

---

### Fase 2: Logica Core (Giorno 3--5)

#### T-004 -- Breakdown Engine (scomposizione di un valore)

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Dato un valore finale e i suoi componenti (con pesi), calcolare il contributo assoluto e percentuale di ciascun componente. Supportare breakdown multi-livello (ricorsivo sui sotto-nodi). |
| **User Story** | US-001 |
| **Complessita** | M |
| **Dipendenze** | T-002, T-003 |

**Acceptance Criteria:**
- [ ] `BreakdownEngine.compute(graph, root_node_id) -> list[Contribution]` restituisce la lista dei contributi
- [ ] Ogni `Contribution` contiene: `node_id`, `label`, `value`, `weight`, `absolute_contribution` (weight * value), `percentage` (contributo / totale * 100)
- [ ] La somma delle `percentage` di tutti i contributi diretti del root e pari a 100% (con tolleranza 0.01%)
- [ ] Supporto breakdown ricorsivo: ogni contributo puo avere sotto-contributi (campo `children: list[Contribution]`)
- [ ] Input con un solo componente: il contributo e il 100%
- [ ] Input con componenti a peso zero: inclusi nella lista con contributo 0%
- [ ] Risultato deterministico: stessi input producono sempre lo stesso output (stessa sequenza)

**File da creare:**
```
src/explainable_engine/engine/breakdown.py
src/explainable_engine/models/contribution.py
tests/test_breakdown.py
```

---

#### T-005 -- Dependency Resolver (albero delle dipendenze)

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Dato un nodo nel grafo, costruire l'albero delle dipendenze completo: quali nodi contribuiscono al valore finale, a che profondita, e con quale relazione. L'output e una struttura ad albero navigabile. |
| **User Story** | US-002 |
| **Complessita** | M |
| **Dipendenze** | T-003 |

**Acceptance Criteria:**
- [ ] `DependencyResolver.resolve(graph, node_id) -> DependencyTree` costruisce l'albero completo
- [ ] `DependencyTree` contiene: `root` (nodo radice), `depth` (profondita massima), `total_nodes` (conteggio)
- [ ] Ogni nodo dell'albero (`DependencyNode`) contiene: `node_id`, `label`, `depth`, `relation` (tipo di arco), `children` (sotto-nodi)
- [ ] Attraversamento in profondita (DFS): `resolver.traverse_dfs(tree) -> list[DependencyNode]`
- [ ] Attraversamento in ampiezza (BFS): `resolver.traverse_bfs(tree) -> list[DependencyNode]`
- [ ] Nodo foglia (senza dipendenze): albero con `depth=0` e `children=[]`
- [ ] Nodo con dipendenze circolari: errore gestito (`CyclicDependencyError`)

**File da creare:**
```
src/explainable_engine/engine/dependency.py
src/explainable_engine/models/dependency.py
tests/test_dependency.py
```

---

#### T-006 -- Confidence Propagation (propagazione della confidenza)

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Ogni nodo ha un punteggio di confidenza (0.0--1.0). La confidenza del nodo radice e calcolata come media pesata delle confidenze dei nodi figli, usando i pesi degli archi. Supportare propagazione ricorsiva bottom-up. |
| **User Story** | US-004 |
| **Complessita** | M |
| **Dipendenze** | T-003 |

**Acceptance Criteria:**
- [ ] `ConfidenceEngine.propagate(graph, root_node_id) -> ConfidenceResult` calcola la confidenza aggregata
- [ ] `ConfidenceResult` contiene: `overall_confidence` (float), `node_confidences` (dict node_id -> float), `propagation_path` (lista di step)
- [ ] Formula di propagazione: `confidence(parent) = sum(weight_i * confidence_i) / sum(weight_i)` per tutti i figli diretti
- [ ] Propagazione ricorsiva: prima si calcolano le foglie, poi si risale (bottom-up)
- [ ] Nodo foglia: la confidenza rimane quella dichiarata nel nodo
- [ ] Tutti i pesi a zero: la confidenza del parent e 0.0
- [ ] Confidenza sempre nel range [0.0, 1.0] -- clamped se necessario
- [ ] `propagation_path` traccia ogni step del calcolo per verificabilita

**File da creare:**
```
src/explainable_engine/engine/confidence.py
src/explainable_engine/models/confidence.py
tests/test_confidence.py
```

---

#### T-007 -- Explanation Orchestrator

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Componente che orchestra l'intera pipeline: riceve un `ExplainRequest`, costruisce il grafo, esegue breakdown, risolve le dipendenze, propaga la confidenza, e assembla la `ExplainResponse` finale. |
| **User Story** | US-001, US-002, US-004 (integrazione) |
| **Complessita** | L |
| **Dipendenze** | T-004, T-005, T-006 |

**Acceptance Criteria:**
- [ ] `Orchestrator.explain(request: ExplainRequest) -> ExplainResponse` produce una risposta completa
- [ ] La risposta contiene: `id` (UUID v4), `timestamp` (ISO 8601), `breakdown`, `dependency_tree`, `confidence`, `graph`
- [ ] L'orchestrator invoca in sequenza: graph engine -> breakdown -> dependency -> confidence
- [ ] Errori in qualsiasi fase producono un errore strutturato (non un crash generico)
- [ ] Tempo di esecuzione < 100ms per un grafo con 20 nodi (misurato con test)
- [ ] Output deterministico verificato: due chiamate con stesso input producono lo stesso JSON (esclusi `id` e `timestamp`)

**File da creare:**
```
src/explainable_engine/engine/orchestrator.py
tests/test_orchestrator.py
```

---

### Fase 3: API Layer (Giorno 6--7)

#### T-008 -- POST /api/v1/explain

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Endpoint REST che accetta un payload JSON (`ExplainRequest`), invoca l'orchestrator, salva il risultato nello storage, e restituisce la `ExplainResponse` completa. |
| **User Story** | US-008 |
| **Complessita** | M |
| **Dipendenze** | T-007, T-010 |

**Acceptance Criteria:**
- [ ] `POST /api/v1/explain` accetta un JSON body conforme a `ExplainRequest`
- [ ] Risposta 200 con `ExplainResponse` completa (breakdown, dipendenze, confidenza)
- [ ] Risposta 422 per payload non valido (errore di validazione Pydantic con dettagli)
- [ ] Risposta 500 per errori interni (messaggio generico, dettagli nei log)
- [ ] Header `X-Request-Id` nella risposta (stesso UUID della spiegazione)
- [ ] Header `X-Processing-Time-Ms` nella risposta
- [ ] Content-Type: `application/json`
- [ ] Il risultato viene persistito nello storage (verificabile con GET successiva)

**File da creare/modificare:**
```
src/explainable_engine/api/__init__.py
src/explainable_engine/api/v1/__init__.py
src/explainable_engine/api/v1/router.py
src/explainable_engine/api/v1/explain.py
src/explainable_engine/main.py              # registrare il router
```

---

#### T-009 -- GET /api/v1/explain/{id}

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Endpoint REST che recupera una spiegazione precedentemente calcolata e salvata, dato il suo UUID. |
| **User Story** | US-009 |
| **Complessita** | S |
| **Dipendenze** | T-008, T-010 |

**Acceptance Criteria:**
- [ ] `GET /api/v1/explain/{id}` restituisce la `ExplainResponse` salvata
- [ ] Risposta 200 con la spiegazione completa (identica a quella restituita dal POST)
- [ ] Risposta 404 con messaggio `{"detail": "Explanation not found"}` per UUID inesistente
- [ ] Risposta 422 per UUID malformato
- [ ] Il JSON restituito e byte-per-byte identico a quello del POST originale (esclusi header HTTP)

**File da creare/modificare:**
```
src/explainable_engine/api/v1/explain.py    # aggiungere GET handler
```

---

#### T-010 -- Storage Layer (in-memory + SQLite)

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Implementare un layer di persistenza con due backend: in-memory (per test e sviluppo) e SQLite (per persistenza). Usare un'interfaccia astratta per permettere di switchare tra i due. |
| **User Story** | US-009 (prerequisito per il caching) |
| **Complessita** | M |
| **Dipendenze** | T-002 |

**Acceptance Criteria:**
- [ ] Interfaccia `ExplanationStore` (Protocol) con metodi: `save(response: ExplainResponse) -> None`, `get(id: UUID) -> ExplainResponse | None`, `exists(id: UUID) -> bool`
- [ ] `InMemoryStore` implementa l'interfaccia con un dizionario
- [ ] `SQLiteStore` implementa l'interfaccia con un database SQLite (una tabella `explanations` con colonne `id TEXT PRIMARY KEY`, `data JSON`, `created_at TIMESTAMP`)
- [ ] `SQLiteStore` crea la tabella automaticamente al primo utilizzo (`CREATE TABLE IF NOT EXISTS`)
- [ ] Round-trip test: salvare e recuperare una `ExplainResponse`, verificare uguaglianza
- [ ] `InMemoryStore`: limite configurabile di entries (default 1000, LRU eviction)
- [ ] Configurazione via `config.py`: scegliere il backend tramite variabile d'ambiente `STORE_BACKEND=memory|sqlite`

**File da creare:**
```
src/explainable_engine/storage/__init__.py
src/explainable_engine/storage/base.py         # Protocol
src/explainable_engine/storage/memory.py
src/explainable_engine/storage/sqlite.py
tests/test_storage.py
```

---

#### T-011 -- Health Endpoint + Error Handling globale

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Aggiungere endpoint di health check e middleware globale per la gestione degli errori (exception handlers FastAPI). |
| **User Story** | Infrastruttura (supporto operativo) |
| **Complessita** | S |
| **Dipendenze** | T-001 |

**Acceptance Criteria:**
- [ ] `GET /health` restituisce `{"status": "healthy", "version": "0.1.0"}`
- [ ] Exception handler globale per `CyclicGraphError` -> 400 con messaggio esplicativo
- [ ] Exception handler globale per `NodeNotFoundError` -> 404
- [ ] Exception handler globale per `ValidationError` (Pydantic) -> 422 con dettagli
- [ ] Exception handler globale per eccezioni non gestite -> 500 con messaggio generico
- [ ] Logging strutturato (JSON) per ogni errore con traceback

**File da creare/modificare:**
```
src/explainable_engine/api/health.py
src/explainable_engine/api/error_handlers.py
src/explainable_engine/main.py               # registrare handlers
```

---

### Fase 4: Qualita (Giorno 8--10)

#### T-012 -- Unit test per Graph Engine

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Suite completa di unit test per `GraphEngine`: costruzione grafo, attraversamento, gestione errori. |
| **User Story** | US-002 (verifica) |
| **Complessita** | M |
| **Dipendenze** | T-003 |

**Acceptance Criteria:**
- [ ] Test grafo vuoto (0 nodi)
- [ ] Test grafo con singolo nodo (nessun arco)
- [ ] Test grafo lineare (A -> B -> C)
- [ ] Test grafo a diamante (A -> B, A -> C, B -> D, C -> D)
- [ ] Test rilevamento cicli
- [ ] Test `get_ancestors` e `get_descendants` su grafo complesso (>= 10 nodi)
- [ ] Test `topological_sort` verifica ordine valido
- [ ] Coverage >= 90% per `graph_engine.py`

**File da creare:**
```
tests/test_graph_engine.py                    # espandere i test di T-003
tests/fixtures/__init__.py
tests/fixtures/sample_graphs.py              # fixture riusabili
```

---

#### T-013 -- Unit test per Breakdown + Confidence

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Suite di unit test per `BreakdownEngine` e `ConfidenceEngine`, inclusi edge case numerici. |
| **User Story** | US-001, US-004 (verifica) |
| **Complessita** | M |
| **Dipendenze** | T-004, T-006 |

**Acceptance Criteria:**
- [ ] Breakdown: test con 2, 5, 10 componenti -- verifica somma percentuali = 100%
- [ ] Breakdown: test componente con peso 0
- [ ] Breakdown: test componente con valore negativo
- [ ] Breakdown: test ricorsivo multi-livello (3 livelli)
- [ ] Confidence: test propagazione su grafo lineare (3 nodi)
- [ ] Confidence: test propagazione su grafo a diamante
- [ ] Confidence: test con tutte le confidenze a 1.0 -> risultato 1.0
- [ ] Confidence: test con tutte le confidenze a 0.0 -> risultato 0.0
- [ ] Confidence: test con pesi non uniformi
- [ ] Coverage >= 90% per `breakdown.py` e `confidence.py`

**File da creare:**
```
tests/test_breakdown.py                      # espandere i test di T-004
tests/test_confidence.py                     # espandere i test di T-006
```

---

#### T-014 -- Integration test per API

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Test end-to-end dell'API usando `httpx.AsyncClient` con la TestClient di FastAPI. Verificare il flusso completo POST -> GET. |
| **User Story** | US-008, US-009 (verifica) |
| **Complessita** | M |
| **Dipendenze** | T-008, T-009 |

**Acceptance Criteria:**
- [ ] Test flusso completo: POST /api/v1/explain con payload valido -> 200, verifica struttura risposta
- [ ] Test flusso completo: POST -> GET con lo stesso ID -> risposte identiche
- [ ] Test POST con payload vuoto -> 422
- [ ] Test POST con payload malformato -> 422 con dettagli di validazione
- [ ] Test GET con UUID inesistente -> 404
- [ ] Test GET /health -> 200
- [ ] Test header `X-Request-Id` presente nella risposta POST
- [ ] Test header `X-Processing-Time-Ms` presente e valore numerico
- [ ] Tempo di risposta < 200ms per payload con 10 componenti (assertion nel test)

**File da creare:**
```
tests/integration/__init__.py
tests/integration/test_api.py
tests/integration/conftest.py                # fixture per TestClient
tests/fixtures/sample_requests.py            # payload di esempio
```

---

#### T-015 -- Test di determinismo

| Campo | Dettaglio |
|---|---|
| **Descrizione** | Verificare che il sistema sia completamente deterministico: dato lo stesso input, l'output e sempre identico (esclusi campi variabili come UUID e timestamp). |
| **User Story** | Trasversale (requisito di qualita) |
| **Complessita** | S |
| **Dipendenze** | T-007 |

**Acceptance Criteria:**
- [ ] Eseguire `Orchestrator.explain()` 100 volte con lo stesso input
- [ ] Verificare che `breakdown`, `dependency_tree`, `confidence` siano identici in tutte le 100 esecuzioni
- [ ] Verificare che l'ordine degli elementi nelle liste sia sempre lo stesso
- [ ] Test con 3 payload diversi (semplice, medio, complesso)

**File da creare:**
```
tests/test_determinism.py
```

---

## Diagramma delle Dipendenze

```
T-001 (setup)
  |
  +---> T-002 (modelli) ----+---> T-004 (breakdown) ------+
  |       |                  |                              |
  |       +---> T-003 (graph engine) --+--> T-005 (deps) --+--> T-007 (orchestrator)
  |       |                            |                    |       |
  |       +---> T-010 (storage)        +--> T-006 (conf) --+       |
  |                 |                                               |
  +---> T-011 (health)                                              |
                                                                    |
                    T-008 (POST) <------- T-007 + T-010 ------------+
                      |
                    T-009 (GET) <-------- T-008 + T-010
                      |
              T-012, T-013, T-014, T-015 (test)
```

---

## 2. Definition of Done

Lo sprint e considerato DONE quando TUTTE le seguenti condizioni sono soddisfatte:

| # | Criterio | Verifica |
|---|---|---|
| 1 | **US-001** (Breakdown): dato un valore e i suoi componenti, il sistema restituisce contributi assoluti e percentuali corretti | Test T-013 passano |
| 2 | **US-002** (Dipendenze): dato un nodo, il sistema restituisce l'albero delle dipendenze con profondita e relazioni | Test T-012 passano |
| 3 | **US-004** (Confidence): il sistema calcola e propaga la confidenza bottom-up con formula pesata | Test T-013 passano |
| 4 | **US-008** (POST /explain): l'endpoint accetta JSON, produce spiegazione completa, persiste il risultato | Test T-014 passano |
| 5 | **US-009** (GET /explain/{id}): l'endpoint restituisce la spiegazione cached identica all'originale | Test T-014 passano |
| 6 | **Test coverage** > 80% su tutto il codice sorgente | `pytest --cov` report |
| 7 | **Performance**: risposta API < 200ms per spiegazioni semplici (< 20 nodi) | Assertion nel test T-014 |
| 8 | **Determinismo**: stesso input produce sempre stesso output | Test T-015 passa (100 iterazioni) |
| 9 | **Documentazione**: API contract (OpenAPI auto-generato da FastAPI) accessibile su `/docs` | Verifica manuale |
| 10 | **Zero errori** su `ruff check` e `ruff format --check` | `make lint` passa |

---

## 3. Rischi e Mitigazioni

### Rischio 1: Complessita del modello dati

**Descrizione:** I modelli Pydantic per grafo, breakdown e dipendenze sono interconnessi. Errori nella progettazione iniziale si propagano a tutti i task successivi, causando rework.

**Probabilita:** Media | **Impatto:** Alto

**Mitigazione:**
- T-002 include test di round-trip serializzazione fin dal primo giorno
- Code review obbligatoria sui modelli prima di procedere alla Fase 2
- Usare `model_config` con esempi concreti per validare il design con gli stakeholder

---

### Rischio 2: Performance della propagazione di confidenza

**Descrizione:** Su grafi con molti nodi (>100), la propagazione ricorsiva bottom-up potrebbe essere lenta o consumare troppa memoria, specialmente con grafi densi.

**Probabilita:** Bassa | **Impatto:** Medio

**Mitigazione:**
- Usare l'ordinamento topologico di NetworkX (gia in T-003) per processare i nodi nell'ordine corretto senza ricorsione
- Implementare memoizzazione dei risultati intermedi in `ConfidenceEngine`
- Il requisito MVP e < 200ms per 20 nodi; l'ottimizzazione per grafi grandi e fuori scope per Sprint 1

---

### Rischio 3: Determinismo non garantito

**Descrizione:** Python `dict` mantiene l'ordine di inserimento dal 3.7+, ma `set` e iterazioni su NetworkX non garantiscono un ordine stabile. Questo puo causare output non deterministici.

**Probabilita:** Media | **Impatto:** Alto

**Mitigazione:**
- Mai usare `set` per risultati esposti all'utente; usare sempre `list` con ordinamento esplicito
- Ordinare esplicitamente i nodi per `node_id` in ogni punto in cui si itera sul grafo
- T-015 (test determinismo) esegue 100 iterazioni per catturare instabilita
- Usare `sorted()` su tutti gli output di NetworkX (es. `sorted(graph.predecessors(node))`)

---

## 4. Metriche di Successo

| Metrica | Target | Come si misura |
|---|---|---|
| User stories completate | 5/5 (US-001, US-002, US-004, US-008, US-009) | Acceptance criteria di ogni US |
| POST /explain funzionante | Risposta 200 con payload di esempio | Test integration T-014 |
| GET /explain/{id} funzionante | Risposta 200 identica al POST originale | Test integration T-014 |
| Confidenza propagata correttamente | Formula pesata verificata su 5+ scenari | Test unit T-013 |
| Tempo di risposta | < 200ms (P99 su 100 richieste, 20 nodi) | Assertion nel test T-014 |
| Test coverage | > 80% | `pytest --cov` |
| Determinismo | 100% su 100 iterazioni x 3 payload | Test T-015 |
| Build pulita | 0 warning ruff, 0 test failure | `make lint && make test` |

---

## Appendice: Payload di Esempio

Payload di riferimento per POST /api/v1/explain (da usare nei test e nella documentazione):

```json
{
  "value": 85.5,
  "label": "Punteggio finale candidato",
  "components": [
    {
      "id": "experience",
      "label": "Esperienza lavorativa",
      "value": 90.0,
      "weight": 0.35,
      "confidence": 0.95
    },
    {
      "id": "skills",
      "label": "Competenze tecniche",
      "value": 82.0,
      "weight": 0.30,
      "confidence": 0.88,
      "components": [
        {
          "id": "python",
          "label": "Python",
          "value": 95.0,
          "weight": 0.50,
          "confidence": 0.92
        },
        {
          "id": "system_design",
          "label": "System Design",
          "value": 70.0,
          "weight": 0.50,
          "confidence": 0.85
        }
      ]
    },
    {
      "id": "culture_fit",
      "label": "Fit culturale",
      "value": 78.0,
      "weight": 0.20,
      "confidence": 0.70
    },
    {
      "id": "education",
      "label": "Formazione",
      "value": 88.0,
      "weight": 0.15,
      "confidence": 0.99
    }
  ]
}
```

Risposta attesa (struttura semplificata):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2026-03-26T10:00:00Z",
  "breakdown": [
    {
      "node_id": "experience",
      "label": "Esperienza lavorativa",
      "value": 90.0,
      "weight": 0.35,
      "absolute_contribution": 31.5,
      "percentage": 36.84
    },
    {
      "node_id": "skills",
      "label": "Competenze tecniche",
      "value": 82.0,
      "weight": 0.30,
      "absolute_contribution": 24.6,
      "percentage": 28.77,
      "children": [
        {
          "node_id": "python",
          "label": "Python",
          "value": 95.0,
          "weight": 0.50,
          "absolute_contribution": 47.5,
          "percentage": 57.93
        },
        {
          "node_id": "system_design",
          "label": "System Design",
          "value": 70.0,
          "weight": 0.50,
          "absolute_contribution": 35.0,
          "percentage": 42.07
        }
      ]
    },
    {
      "node_id": "culture_fit",
      "label": "Fit culturale",
      "value": 78.0,
      "weight": 0.20,
      "absolute_contribution": 15.6,
      "percentage": 18.25
    },
    {
      "node_id": "education",
      "label": "Formazione",
      "value": 88.0,
      "weight": 0.15,
      "absolute_contribution": 13.2,
      "percentage": 15.44
    }
  ],
  "confidence": 0.879,
  "dependency_tree": {
    "root": "final_score",
    "depth": 2,
    "total_nodes": 6
  }
}
```
