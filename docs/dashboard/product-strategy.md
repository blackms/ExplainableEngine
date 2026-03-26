# Explainable Engine Dashboard -- Product Strategy

**Versione:** 1.0.0
**Data:** 2026-03-26
**Stato:** Draft
**Autore:** Product Management

---

## Indice

1. [Vision Statement](#1-vision-statement)
2. [Use Cases](#2-use-cases)
3. [Personas](#3-personas)
4. [Feature Prioritization Matrix](#4-feature-prioritization-matrix)
5. [Success Metrics](#5-success-metrics)
6. [Risks & Mitigations](#6-risks--mitigations)

---

## 1. Vision Statement

La Dashboard di Explainable Engine e' l'interfaccia visuale che rende accessibile a utenti non tecnici il potere dell'explainability engine. Oggi l'engine espone sei endpoint API (health, POST/GET explain, graph export, narrative, what-if) che producono catene causali esplicite, interrogabili e verificabili -- ma solo un developer con `curl` o Postman puo' realmente utilizzarli. La dashboard colma questo gap: trasforma dati strutturati in grafo, breakdown numerici e analisi di sensitivita' in esperienze visive intuitive. Per un Compliance Officer significa audit trail navigabile e esportabile in un click, senza dipendere dal team IT. Per un Quant significa debug interattivo dei modelli con drill-down nel grafo causale. Per un Sales significa un report chiaro e brandizzato da condividere con il cliente in cinque minuti. La dashboard non e' un "nice-to-have" cosmetico: e' il prodotto che trasforma un motore tecnico in valore di business misurabile -- compliance senza sanzioni, debugging senza attesa, client retention senza ambiguita'. Ogni decisione di design segue il principio fondante dell'engine: "Non spiegare il risultato. Rappresenta il processo che lo genera."

---

## 2. Use Cases

### UC-001: Audit Trail & Compliance Review

| Campo | Dettaglio |
|-------|-----------|
| **UC-ID** | UC-001 |
| **Name** | Audit Trail & Compliance Review |
| **Actor** | Elena -- Compliance Officer / Auditor esterno |
| **Business Context** | La regolamentazione europea (MiFID II art. 25, EU AI Act art. 13-14, GDPR art. 22) richiede che ogni decisione automatizzata sia spiegabile, tracciabile e riproducibile. Le sanzioni per non-compliance partono da centinaia di migliaia di euro e possono raggiungere il 10% del fatturato annuo globale per violazioni EU AI Act. Ogni audit non superato non e' solo una multa: e' un rischio reputazionale che impatta la capacita' di attrarre clienti istituzionali. Il costo medio di un'ispezione regolamentare mal preparata include settimane di lavoro del team legale e IT per ricostruire manualmente le catene causali. La dashboard elimina questo costo trasformando l'audit da processo manuale ed error-prone a workflow self-service. |
| **Trigger** | (1) Richiesta formale di audit da parte dell'autorita' di vigilanza (Consob, Banca d'Italia, ESMA). (2) Revisione periodica interna (trimestrale/annuale). (3) Reclamo di un cliente su una decisione automatizzata. (4) Incident post-mortem che richiede ricostruzione della catena decisionale. |
| **Preconditions** | L'utente ha credenziali di accesso con ruolo "compliance" o "auditor". Le explanation relative al periodo sotto esame sono state persistite nell'engine (storage SQLite/PostgreSQL). L'engine e' operativo e raggiungibile. |

**Main Flow:**

1. Elena accede alla dashboard e seleziona la sezione "Audit Trail".
2. Il sistema presenta un'interfaccia di ricerca avanzata con filtri: intervallo temporale (date picker), target metric name, intervallo di confidence score, computation type, explanation ID specifico.
3. Elena imposta i filtri per il trimestre sotto audit (es. Q4 2025) e per le explanation con confidence < 0.7 (focus su decisioni a bassa affidabilita').
4. Il sistema interroga `GET /api/v1/explain/{id}` per ogni risultato e presenta una lista paginata (cursor-based pagination, max 100 per pagina) con: ID, target, final_value, confidence, data creazione, deterministic_hash.
5. Elena seleziona una explanation dalla lista.
6. Il sistema carica la vista dettagliata: breakdown numerico completo (contribution, contribution_pct per ogni componente), grafo causale interattivo (via `GET /api/v1/explain/{id}/graph`), narrative in linguaggio naturale (via `GET /api/v1/explain/{id}/narrative`).
7. Elena verifica la riproducibilita': il sistema mostra il `deterministic_hash` e permette di ricalcolarlo on-demand per confermare l'integrita'.
8. Elena seleziona le explanation da includere nel report e clicca "Export".
9. Il sistema genera un documento PDF contenente: intestazione con data e identificativo audit, lista explanation con breakdown completo, grafo causale come immagine SVG, narrative testuale, hash deterministico per verifica, firma digitale del timestamp di esportazione.
10. Elena scarica il PDF e/o lo invia via email direttamente dalla dashboard.

**Alternative Flows:**

- **AF-001a: Nessun risultato.** Se i filtri non producono risultati, il sistema mostra un messaggio chiaro ("Nessuna explanation trovata per i criteri selezionati") e suggerisce di allargare l'intervallo temporale.
- **AF-001b: Explanation corrotta o non riproducibile.** Se il ricalcolo del deterministic_hash non corrisponde a quello persistito, il sistema mostra un warning rosso ("INTEGRITY CHECK FAILED") con dettagli sulla discrepanza. Elena puo' esportare anche questa evidenza.
- **AF-001c: Export massivo.** Se Elena seleziona piu' di 500 explanation, il sistema offre un'esportazione asincrona con notifica via email al completamento.
- **AF-001d: Auditor esterno senza account.** Un Compliance Officer puo' generare un link di accesso temporaneo (read-only, scadenza 72h) per un auditor esterno, limitato alle explanation selezionate.

**Postconditions:**

- Elena ha ottenuto un report PDF completo e verificabile con tutte le explanation del periodo sotto audit.
- Ogni explanation nel report include il deterministic_hash che permette la verifica indipendente.
- Il sistema ha registrato un audit log dell'accesso e dell'esportazione (chi, quando, cosa).

**Business Value:**

- Riduzione del 80% del tempo di preparazione audit (da settimane a ore).
- Eliminazione del rischio di sanzioni per mancata tracciabilita' (esposizione potenziale: 1M-50M EUR a seconda della normativa violata).
- Riduzione del costo legale per ispezioni (stima: 50-100K EUR per ispezione mal preparata).
- Abilitazione alla certificazione EU AI Act, prerequisito per servire clienti istituzionali dal 2025.

**Priority:** **Must-have** -- Senza audit trail, il prodotto non e' vendibile a clienti regolamentati (banche, assicurazioni, SGR). E' un prerequisito legale, non una feature.

---

### UC-002: Quantitative Model Debugging

| Campo | Dettaglio |
|-------|-----------|
| **UC-ID** | UC-002 |
| **Name** | Quantitative Model Debugging |
| **Actor** | Marco -- Senior Quant / Data Scientist |
| **Business Context** | In un hedge fund o trading desk, un segnale quantitativo errato che passa inosservato puo' generare perdite dirette nell'ordine di centinaia di migliaia di euro in una singola sessione di trading. Il tempo di debugging e' direttamente correlato al P&L: ogni ora in piu' per identificare un errore nel modello e' un'ora in cui il segnale sbagliato potrebbe essere ancora attivo. Oggi i quant fanno debugging manuale: scaricano i dati, ricostruiscono il calcolo in un notebook, confrontano i valori. Questo processo richiede in media 2-4 ore per investigation. La dashboard riduce questo tempo a minuti, rendendo la catena causale immediatamente navigabile. |
| **Trigger** | (1) Un segnale di trading produce un risultato inatteso (es. market_regime_score = 0.95 quando il mercato e' laterale). (2) Il confidence score di un output cala improvvisamente sotto soglia. (3) Un backtest mostra discrepanze rispetto ai risultati attesi. (4) Un nuovo componente viene aggiunto al modello e serve verificarne l'impatto. |
| **Preconditions** | Il quant ha accesso alla dashboard con ruolo "analyst". L'explanation relativa al calcolo sotto indagine esiste nell'engine. Il grafo causale e' stato generato con `include_graph: true` e `include_drivers: true`. |

**Main Flow:**

1. Marco accede alla dashboard e cerca l'explanation per ID o per target metric name (es. "market_regime_score").
2. Il sistema presenta la vista principale: il valore finale (final_value), il confidence score aggregato, e il breakdown dei top_drivers ordinati per impact.
3. Marco nota che `trend_strength` ha impact = 0.44 (rank 1) ma il mercato non sta mostrando alcun trend. Clicca su `trend_strength` nel breakdown.
4. Il sistema effettua un drill-down: carica il sotto-grafo di `trend_strength` mostrando i suoi sub_components (se presenti) con il relativo sub_breakdown. Se `trend_strength` non ha sub_components, mostra i suoi attributi dettagliati (value, weight, confidence, metadata).
5. Marco attiva la vista grafo interattiva (`GET /api/v1/explain/{id}/graph`): il DAG viene renderizzato con nodi cliccabili. I nodi sono colorati per confidence (verde > 0.8, giallo 0.5-0.8, rosso < 0.5). Gli edge mostrano i pesi.
6. Marco clicca su un nodo nel grafo per ispezionarne i dettagli: valore, peso, contributo assoluto e percentuale, confidence, tipo di computazione.
7. Marco identifica il problema: un componente con confidence = 0.3 (dato parzialmente mancante) sta alimentando `trend_strength` con un peso alto. La confidence propagata non e' sufficiente a declassare il risultato.
8. Marco utilizza la funzione "What-if" integrata (link a UC-004) per simulare cosa succederebbe con quel componente rimosso o con un valore diverso.
9. Marco annota le sue conclusioni direttamente nella dashboard (campo note associato all'explanation) per il team.

**Alternative Flows:**

- **AF-002a: Grafo molto profondo (> 5 livelli).** Il sistema offre un controllo di max_depth per limitare la profondita' di rendering. Un minimap mostra la posizione attuale nel grafo complessivo.
- **AF-002b: Grafo molto largo (> 50 nodi allo stesso livello).** Il sistema raggruppa automaticamente i nodi per contribution_pct, mostrando in dettaglio solo i top contributors e collassando gli altri in un nodo aggregato "Altri (N componenti, X% totale)".
- **AF-002c: Confronto tra explanation.** Marco puo' aprire due explanation side-by-side per confrontare i breakdown e identificare cosa e' cambiato tra due calcoli dello stesso target in momenti diversi.

**Postconditions:**

- Marco ha identificato la root cause del risultato inatteso navigando il grafo causale.
- Il tempo di investigation e' stato ridotto da ore a minuti.
- Le annotazioni di Marco sono persistite e visibili al team.

**Business Value:**

- Riduzione del tempo medio di debugging da 2-4 ore a 10-20 minuti (saving: 15-30 ore/mese per quant).
- Prevenzione di perdite da segnali errati (esposizione: variabile, ma un singolo segnale sbagliato su un portafoglio da 100M EUR puo' costare 50-500K EUR).
- Accelerazione del ciclo di sviluppo dei modelli: feedback loop piu' veloce significa modelli migliori in produzione prima.

**Priority:** **Must-have** -- E' il use case primario che ha originato l'intero progetto (AIP trading platform). Senza questo, l'engine non serve al suo pubblico principale.

---

### UC-003: Client Communication & Reporting

| Campo | Dettaglio |
|-------|-----------|
| **UC-ID** | UC-003 |
| **Name** | Client Communication & Reporting |
| **Actor** | Roberto -- Sales Director / Account Manager / Relationship Manager |
| **Business Context** | Nel settore dell'asset management e del wealth management, la fiducia del cliente e' direttamente correlata alla capacita' di spiegare i risultati. Il 34% dei clienti istituzionali cita "mancanza di trasparenza" come motivo principale di switch verso un altro provider (fonte tipica: surveys di settore). Su un AUM medio di 50M EUR per cliente istituzionale con fee dell'1%, ogni cliente perso rappresenta 500K EUR/anno di ricavi. La domanda piu' frequente in una client review e': "Perche' questo numero?" -- e oggi il Sales deve chiedere al Quant, che deve fare un'analisi ad-hoc, che richiede 1-2 giorni. La dashboard permette al Sales di rispondere in autonomia, in tempo reale, durante la call con il cliente. |
| **Trigger** | (1) Client review periodica (mensile/trimestrale). (2) Il cliente chiede spiegazioni su un risultato specifico (es. "perche' il risk score del mio portafoglio e' salito?"). (3) Onboarding di un nuovo cliente che vuole capire la metodologia. (4) Richiesta di report personalizzato da parte del cliente. |
| **Preconditions** | Il Sales ha accesso alla dashboard con ruolo "sales" (read-only sui dati, con permessi di export e condivisione). Le explanation rilevanti esistono nell'engine. Il template di report aziendale e' configurato (logo, colori, disclaimer legale). |

**Main Flow:**

1. Roberto accede alla dashboard e cerca l'explanation per il portafoglio del cliente (per target metric name o per tag/label associato).
2. Il sistema mostra la vista "Client View" -- una versione semplificata della explanation ottimizzata per la comunicazione: valore finale in evidenza, i 3-5 top drivers con barre percentuali colorate, un indicatore visivo di confidence (icona semaforo o gauge).
3. Roberto clicca "Genera Narrative" per ottenere la spiegazione in linguaggio naturale (via `GET /api/v1/explain/{id}/narrative`).
4. Il sistema presenta la narrative in italiano (o nella lingua configurata per il cliente). Roberto puo' switchare lingua (IT/EN/DE/FR) se il narrative engine lo supporta.
5. Roberto personalizza il report: seleziona quali sezioni includere (breakdown si/no, grafo si/no, narrative si/no), aggiunge una nota personale per il cliente.
6. Roberto clicca "Genera Report PDF".
7. Il sistema produce un PDF brandizzato con: logo e colori aziendali, intestazione con nome cliente e data, valore finale con breakdown visuale (bar chart), narrative in linguaggio naturale, disclaimer legale (configurabile), footer con deterministic_hash per verifica.
8. Roberto scarica il PDF o lo invia direttamente via email al cliente dalla dashboard (integrazione email opzionale).
9. Il sistema registra l'invio nel CRM log (o espone l'evento via webhook per integrazione con CRM esterni).

**Alternative Flows:**

- **AF-003a: Narrative non disponibile.** Se l'endpoint narrative restituisce un errore o il testo e' vuoto, il sistema mostra il breakdown numerico con un messaggio: "Narrative non disponibile. Breakdown numerico incluso nel report."
- **AF-003b: Multi-explanation report.** Roberto puo' selezionare piu' explanation (es. tutti i fattori di rischio di un portafoglio) e generare un report consolidato con indice e sezioni.
- **AF-003c: Client self-service portal.** In futuro, il cliente stesso potrebbe accedere a una vista read-only dedicata (link sicuro con scadenza) per consultare le explanation del proprio portafoglio senza intermediazione del Sales.

**Postconditions:**

- Il cliente ha ricevuto un report chiaro, brandizzato e verificabile.
- Roberto non ha avuto bisogno di coinvolgere il team Quant.
- L'invio e' tracciato per compliance.

**Business Value:**

- Riduzione del churn rate attribuibile a "mancanza di trasparenza" (target: -50%, saving potenziale: 250K EUR/anno per ogni cliente istituzionale salvato).
- Riduzione del tempo di preparazione client review da 1-2 giorni a 30 minuti.
- Riduzione della dipendenza Sales -> Quant per spiegazioni ad-hoc (liberando tempo dei quant per attivita' a piu' alto valore).
- Differenziazione competitiva: "noi vi mostriamo esattamente come arriviamo ai nostri numeri".

**Priority:** **Should-have** -- Altissimo impatto su retention e revenue, ma dipende da UC-001 (audit trail) e UC-002 (explorer) come base. Realizzabile in seconda fase senza bloccare il lancio.

---

### UC-004: What-if Scenario Analysis

| Campo | Dettaglio |
|-------|-----------|
| **UC-ID** | UC-004 |
| **Name** | What-if Scenario Analysis |
| **Actor** | Sofia -- Risk Manager / Portfolio Manager / Strategist |
| **Business Context** | Lo stress testing e' un requisito regolamentare (ICAAP/ILAAP, Solvency II, EBA guidelines) e una necessita' operativa. Ogni decisione di portafoglio si basa su assunzioni: "cosa succede al nostro risk score se la volatilita' raddoppia?" La risposta a questa domanda determina hedge strategies, allocazioni e limiti di rischio. Oggi lo stress testing richiede di ri-eseguire l'intero modello con input modificati -- un processo che puo' richiedere da minuti a ore a seconda della complessita'. L'endpoint `what-if` di ExplainableEngine permette di ottenere la risposta in millisecondi, e la dashboard lo rende accessibile via slider interattivi. Decisioni migliori, piu' veloci, meglio documentate. |
| **Trigger** | (1) Analisi pre-trade: "Se aumento l'esposizione su questo settore, come cambia il risk score?" (2) Stress testing regolamentare periodico. (3) Evento di mercato (flash crash, annuncio BCE, crisi geopolitica) che richiede ricalcolo immediato degli scenari. (4) Richiesta del CdA per scenario analysis su specifiche ipotesi. |
| **Preconditions** | L'explanation di base esiste nell'engine. L'utente ha ruolo "analyst" o "risk_manager". L'endpoint `POST /api/v1/explain/{id}/what-if` e' operativo. |

**Main Flow:**

1. Sofia accede a un'explanation esistente (es. "portfolio_risk_score" con final_value = 0.72).
2. Sofia clicca "What-if Analysis" e il sistema presenta la modalita' simulazione.
3. Per ogni componente del breakdown, il sistema mostra uno slider interattivo: il valore attuale e' pre-impostato, il range e' calcolato come [0, max(value * 3, 1.0)] o configurabile.
4. Sofia muove lo slider di `volatility` da 0.5 a 1.0 (scenario: "volatilita' raddoppia").
5. Il sistema chiama `POST /api/v1/explain/{id}/what-if` con gli override `{"volatility": 1.0}` e mostra in tempo reale (< 200ms): il nuovo final_value, il delta rispetto al baseline, il nuovo breakdown con le percentuali aggiornate, una heatmap che evidenzia quali componenti sono stati piu' impattati.
6. Sofia modifica un secondo slider (`momentum` da 0.6 a 0.2) per uno scenario combinato.
7. Il sistema ricalcola con entrambi gli override e aggiorna la vista.
8. Sofia clicca "Salva Scenario" per dare un nome a questa configurazione (es. "Stress: vol raddoppia + momentum crolla").
9. Sofia clicca "Confronta Scenari" e il sistema presenta una tabella side-by-side: Baseline vs Scenario 1 vs Scenario 2, con delta evidenziati.
10. Sofia esporta il confronto scenari come PDF o CSV per il report al CdA.

**Alternative Flows:**

- **AF-004a: Slider con sub-components.** Se un componente ha sub_components, lo slider modifica il valore aggregato. Un toggle permette di espandere e modificare i singoli sub_components.
- **AF-004b: Scenario predefiniti.** Il sistema offre template di scenario comuni: "Market crash (-30% equities)", "Rate hike (+200bps)", "Black swan (confidence drop to 0.2 on all inputs)". Questi template popolano automaticamente gli slider.
- **AF-004c: Sensitivity auto-discovery.** Un pulsante "Trova componente piu' sensibile" esegue una sensitivity analysis automatica (variando ogni input di +/-10%) e ordina i componenti per impatto marginale.
- **AF-004d: Vincoli di range.** Se un utente imposta un valore fuori dal range logico (es. probabilita' > 1.0), il sistema mostra un warning ma permette l'override con conferma.

**Postconditions:**

- Sofia ha quantificato l'impatto di scenari specifici sul risultato.
- Gli scenari sono salvati e nominati per riferimento futuro.
- Il confronto scenari e' esportabile per reporting.

**Business Value:**

- Decisioni di portafoglio piu' informate (impatto diretto sul rendimento aggiustato per il rischio).
- Compliance con requisiti di stress testing regolamentare (ICAAP, Solvency II).
- Riduzione del tempo di scenario analysis da ore a secondi.
- Documentazione automatica degli scenari analizzati (audit trail per il regolatore).

**Priority:** **Should-have** -- L'endpoint what-if esiste gia' nel backend. La dashboard lo rende utilizzabile da utenti non tecnici. Alto valore, effort moderato.

---

### UC-005: Real-time Decision Monitoring

| Campo | Dettaglio |
|-------|-----------|
| **UC-ID** | UC-005 |
| **Name** | Real-time Decision Monitoring |
| **Actor** | Operations Manager / Trading Desk Supervisor |
| **Business Context** | In un contesto di trading automatizzato, i segnali vengono generati continuamente. Un segnale con bassa confidence che passa inosservato puo' generare trade errati per ore prima che qualcuno se ne accorga. Il costo di un ritardo di rilevazione e' direttamente proporzionale al volume di trading: su un desk che esegue 1000 operazioni/giorno con un ticket medio di 100K EUR, anche un errore che impatta lo 0.1% delle operazioni per 4 ore costa 40-100K EUR. La dashboard di monitoraggio in tempo reale e' il "pannello di controllo" che permette al desk di reagire immediatamente ad anomalie nella catena esplicativa. |
| **Trigger** | Sempre attivo durante le ore di mercato. Alert specifici scattano quando: (1) una explanation ha confidence < soglia configurabile. (2) Il missing_impact supera una soglia (troppi dati mancanti). (3) Il numero di explanation con anomalie supera un threshold nel periodo. (4) Un deterministic_hash non corrisponde (possibile data corruption). |
| **Preconditions** | L'engine riceve un flusso continuo di explanation requests dal sistema di trading (AIP). La dashboard e' connessa e aggiornata (polling o websocket). Le soglie di alert sono configurate. |

**Main Flow:**

1. L'Operations Manager apre la dashboard di monitoraggio ("Control Room").
2. Il sistema presenta una vista real-time con: contatore di explanation processate (ultimi 5min/1h/24h), distribuzione della confidence (histogram), lista delle ultime N explanation in ordine cronologico inverso con status badge (verde/giallo/rosso), alert attivi in evidenza in cima.
3. Il sistema esegue polling periodico (ogni 10-30 secondi, configurabile) su `GET /api/v1/explain/{id}` per le explanation piu' recenti, oppure riceve eventi via webhook.
4. Una nuova explanation arriva con confidence = 0.45 (sotto soglia 0.6). Il sistema mostra un alert: badge rosso sulla riga, notifica toast in-app, incremento del contatore alert attivi.
5. L'Operations Manager clicca sull'alert e il sistema apre la vista dettagliata dell'explanation (collegamento a UC-002).
6. L'Operations Manager valuta la situazione e puo': (a) marcare l'alert come "acknowledged" (visto, in gestione), (b) marcare come "resolved" con nota, (c) escalare via notifica al team (email/Slack webhook).
7. La dashboard aggiorna continuamente le metriche aggregate: confidence media mobile, trend del missing_impact, numero di anomalie per fascia oraria.

**Alternative Flows:**

- **AF-005a: Picco di volume.** Se il numero di explanation/minuto supera una soglia (es. 10x la media), il sistema mostra un alert giallo "Volume anomaly" per segnalare possibili problemi upstream.
- **AF-005b: Engine non raggiungibile.** Se il polling fallisce per 3 cicli consecutivi, la dashboard mostra un banner "Engine unreachable" con timestamp dell'ultimo dato ricevuto.
- **AF-005c: Filtri per target.** L'Operations Manager puo' filtrare la vista per target metric (es. solo "market_regime_score") per concentrarsi su un singolo modello.

**Postconditions:**

- Tutte le anomalie sono state rilevate e gestite entro il ciclo di polling.
- Lo storico degli alert e' persistito per audit.
- Le metriche aggregate sono disponibili per reporting operativo.

**Business Value:**

- Riduzione del tempo di rilevazione anomalie da ore a secondi.
- Prevenzione di perdite da segnali errati non intercettati.
- Compliance con requisiti di supervisione operativa.
- Dati operativi per il miglioramento continuo dei modelli.

**Priority:** **Should-have** -- Alto valore operativo, ma richiede un flusso continuo di dati (presuppone che l'engine sia integrato con un sistema di trading attivo). Puo' essere lanciato dopo le funzionalita' core.

---

### UC-006: API Integration Playground

| Campo | Dettaglio |
|-------|-----------|
| **UC-ID** | UC-006 |
| **Name** | API Integration Playground |
| **Actor** | Alex -- Backend Developer / Integration Engineer |
| **Business Context** | Il time-to-value per un nuovo cliente dipende dalla velocita' di integrazione dell'API. Il ciclo tipico e': firma contratto -> il developer del cliente legge la documentazione -> prova le API -> integra nel sistema -> go-live. Oggi questo ciclo richiede in media 2-4 settimane. Ogni settimana in meno e' una settimana in piu' di fee ricorrenti. Inoltre, un'esperienza di integrazione frustrante genera support ticket costosi (30-60 minuti ciascuno per il team di engineering) e rischio di abbandono pre-go-live. Un playground interattivo con esempi pronti, code generation e test in-browser riduce il ciclo di integrazione e il volume di support ticket. |
| **Trigger** | (1) Un nuovo cliente deve integrare ExplainableEngine nel proprio sistema. (2) Un developer esistente deve testare un nuovo endpoint o una nuova feature. (3) Un developer vuole esplorare il comportamento dell'API con diversi payload senza scrivere codice. (4) Il team di sales vuole fare una demo live dell'API a un prospect. |
| **Preconditions** | La specifica OpenAPI (`api-contract.yaml`) e' aggiornata e servita dall'engine. Il playground e' accessibile senza autenticazione (o con token demo). L'engine di sviluppo/staging e' disponibile come backend. |

**Main Flow:**

1. Alex accede alla sezione "API Playground" della dashboard.
2. Il sistema presenta una lista degli endpoint disponibili, raggruppati per categoria (explain, graph, narrative, what-if, health), basata sulla specifica OpenAPI.
3. Alex seleziona `POST /api/v1/explain`.
4. Il sistema mostra: la descrizione dell'endpoint, il request body schema con campi compilabili (form mode) o editor JSON (raw mode), esempi pre-caricati selezionabili da un dropdown.
5. Alex seleziona l'esempio "Market Regime Score" e il sistema popola il request body con il payload di esempio dalla specifica OpenAPI.
6. Alex modifica alcuni valori (es. cambia `trend_strength.value` da 0.8 a 0.3) e clicca "Send Request".
7. Il sistema invia la richiesta all'engine (staging/development) e mostra la response: status code, headers (incluso X-Request-Id, X-Processing-Time-Ms), response body formattato con syntax highlighting.
8. Alex clicca "Generate Code" e seleziona il linguaggio (Python, Go, JavaScript/TypeScript, cURL, Java). Il sistema genera il codice di integrazione completo per la richiesta appena eseguita.
9. Alex copia il codice generato e lo integra nel proprio progetto.
10. Alex esplora gli altri endpoint: seleziona `GET /api/v1/explain/{id}` e il sistema pre-popola l'`id` con quello restituito dalla richiesta precedente.

**Alternative Flows:**

- **AF-006a: Errore nell'API.** Se l'engine restituisce un errore (400, 422, 500), il sistema mostra il response body di errore con evidenziazione dei campi problematici e suggerimenti per la correzione (basati sul campo `details` dell'ErrorResponse).
- **AF-006b: Chained requests.** Un mode "Workflow" permette di eseguire una sequenza di chiamate (POST explain -> GET explain/{id} -> GET graph -> POST what-if) con auto-popolamento dei parametri tra una chiamata e l'altra.
- **AF-006c: Environment switching.** Alex puo' switchare tra ambienti (localhost, staging, production) da un dropdown. L'URL base cambia automaticamente.
- **AF-006d: Request history.** Le ultime 50 richieste sono salvate in local storage e navigabili per quick replay.

**Postconditions:**

- Alex ha testato l'API interattivamente e ha il codice di integrazione pronto.
- Il tempo di integrazione e' ridotto significativamente.
- Alex non ha avuto bisogno di contattare il team di supporto.

**Business Value:**

- Riduzione del tempo di integrazione del 60% (da 2-4 settimane a 3-7 giorni).
- Riduzione del 70% dei support ticket relativi all'integrazione API.
- Strumento di demo per il team sales (riduzione del ciclo di vendita).
- Miglioramento della developer experience (DX) come vantaggio competitivo.

**Priority:** **Nice-to-have** -- Molto utile ma non bloccante per il lancio. La specifica OpenAPI e' gia' disponibile e tools esterni (Swagger UI, Insomnia) possono coprire parzialmente il bisogno. Diventa should-have quando il numero di clienti in integrazione supera 5.

---

## 3. Personas

### Marco -- Senior Quant

| Campo | Dettaglio |
|-------|-----------|
| **Nome** | Marco Bianchi |
| **Eta'** | 38 |
| **Ruolo** | Senior Quantitative Analyst presso un hedge fund (AUM 2B EUR) |
| **Esperienza** | 12 anni in finanza quantitativa, PhD in matematica applicata |
| **Team** | Quant Research (6 persone), riporta al CIO |

**Goals:**

- Capire immediatamente perche' un segnale di trading ha un valore specifico, senza dover ricostruire il calcolo manualmente.
- Identificare rapidamente quale componente del modello sta "misbehaving" quando i risultati sono inattesi.
- Verificare che le modifiche al modello producano l'effetto atteso prima di metterle in produzione.
- Documentare le investigazioni per il team e per audit futuri.

**Frustrations:**

- "Passo 3 ore a ricostruire un calcolo che il sistema ha fatto in 200ms. E' assurdo."
- "Quando il PM mi chiede 'perche' questo numero?', devo scaricare i dati, aprire un notebook, e ricominciare da zero ogni volta."
- "I log del sistema sono illeggibili. Cerco un ago in un pagliaio di JSON."
- "Non ho un modo veloce per fare 'e se cambiassi questo input, cosa succederebbe?' senza ri-eseguire l'intero pipeline."

**Livello tecnico:** Esperto. Scrive codice in Python e Go quotidianamente. Sa leggere JSON e API docs. Preferisce interfacce pulite e veloci rispetto a quelle "belle". Ama le keyboard shortcuts. Odia le animazioni inutili.

**Come usa ExplainableEngine:** Marco e' il power user primario. Utilizza l'engine sia tramite API (integrato nel pipeline di trading) sia tramite dashboard per debugging interattivo. La dashboard e' il suo "microscopic" per i modelli: naviga il grafo, drilla nei componenti, usa what-if per testare ipotesi. Usa la dashboard 5-10 volte al giorno durante le ore di mercato.

---

### Elena -- Compliance Officer

| Campo | Dettaglio |
|-------|-----------|
| **Nome** | Elena Rossi |
| **Eta'** | 45 |
| **Ruolo** | Head of Compliance presso una banca d'investimento (mid-tier) |
| **Esperienza** | 18 anni in compliance/legal, certificazione CAMS e ACAMS |
| **Team** | Compliance (4 persone), riporta al Chief Risk Officer |

**Goals:**

- Dimostrare all'autorita' di vigilanza che ogni decisione automatizzata e' tracciabile e spiegabile.
- Preparare la documentazione di audit nel minor tempo possibile, senza dipendere dal team IT.
- Monitorare la qualita' delle spiegazioni nel tempo (trend di confidence, missing data).
- Garantire la conformita' al EU AI Act prima dell'entrata in vigore delle sanzioni.

**Frustrations:**

- "Ogni volta che c'e' un'ispezione, devo chiedere all'IT di estrarre i dati. Ci vogliono 3 giorni per avere quello che mi serve."
- "Non capisco i dettagli tecnici dei modelli, ma devo certificare che sono conformi. Ho bisogno di una vista ad alto livello ma verificabile."
- "I report che mi arrivano sono fogli Excel illeggibili. Ho bisogno di qualcosa che un regolatore possa capire."
- "Non ho modo di sapere in anticipo se abbiamo problemi di compliance finche' non arriva l'ispezione."

**Livello tecnico:** Intermedio-basso. Sa usare applicazioni web, Excel avanzato, sistemi GRC. Non scrive codice. Non legge JSON. Ha bisogno di interfacce intuitive con terminologia business, non tecnica. I grafici e le visualizzazioni sono preferiti rispetto alle tabelle di numeri.

**Come usa ExplainableEngine:** Elena usa la dashboard 2-3 volte a settimana per monitoraggio e in modo intensivo (quotidiano) durante i periodi di audit. Si concentra sull'audit trail (UC-001): ricerca, filtro, export PDF. Non interagisce con il grafo tecnico ma legge le narrative e verifica i deterministic hash. Per lei la dashboard e' lo "strumento di certificazione".

---

### Roberto -- Sales Director

| Campo | Dettaglio |
|-------|-----------|
| **Nome** | Roberto Conti |
| **Eta'** | 42 |
| **Ruolo** | Sales Director / Head of Client Relations presso una SGR (AUM 500M EUR) |
| **Esperienza** | 15 anni nel sales istituzionale, CFA charterholder |
| **Team** | Sales & Client Service (8 persone), riporta al CEO |

**Goals:**

- Rispondere alle domande dei clienti sulle performance e sui fattori di rischio senza coinvolgere il team quant.
- Produrre report personalizzati e brandizzati per le client review in autonomia e rapidamente.
- Differenziare il servizio: "Noi vi mostriamo esattamente come calcoliamo i nostri numeri."
- Ridurre il churn rate aumentando la percezione di trasparenza.

**Frustrations:**

- "Il cliente mi chiede 'perche' questo numero?' e io devo dire 'ci torno'. Perdo credibilita' ogni volta."
- "Preparare una client review mi costa 2 giorni di lavoro tra richieste al quant, attesa, formattazione. E' insostenibile."
- "I report tecnici che ricevo dal team non sono presentabili al cliente. Devo riscrivere tutto."
- "Non ho modo di fare una demo convincente del nostro approccio quantitativo durante una riunione commerciale."

**Livello tecnico:** Basso. Sa usare CRM, PowerPoint, Excel. Non sa cosa sia un'API. Non vuole vedere JSON, grafi o codice. Ha bisogno di: narrative in linguaggio naturale, grafici a barre semplici, PDF pronti da inviare. L'interfaccia deve essere "bella" e professionale perche' la condivide con i clienti.

**Come usa ExplainableEngine:** Roberto usa la dashboard 3-5 volte a settimana, quasi esclusivamente per UC-003 (client reporting). Cerca l'explanation, legge la narrative, genera il PDF, lo invia al cliente. Occasionalmente usa la Client View semplificata durante una call con il cliente condividendo lo schermo. Per lui la dashboard e' il "generatore di fiducia".

---

### Sofia -- Risk Manager

| Campo | Dettaglio |
|-------|-----------|
| **Nome** | Sofia Marchetti |
| **Eta'** | 35 |
| **Ruolo** | Senior Risk Manager presso una compagnia assicurativa |
| **Esperienza** | 10 anni nel risk management, FRM certified, background in ingegneria gestionale |
| **Team** | Risk Management (5 persone), riporta al CRO |

**Goals:**

- Quantificare l'impatto di scenari avversi sui risk score prima che si verifichino.
- Produrre stress test documentati e verificabili per il regolatore (IVASS/EIOPA).
- Identificare i fattori di rischio piu' sensibili per concentrare le azioni di mitigazione.
- Confrontare scenari alternativi per raccomandare strategie al CdA.

**Frustrations:**

- "Ogni stress test richiede di ri-eseguire l'intero modello con nuovi parametri. Ci vogliono ore e il team di IT deve essere coinvolto."
- "Non ho un modo rapido per rispondere alla domanda 'e se la volatilita' raddoppiasse?' durante una riunione del comitato rischi."
- "I risultati degli stress test sono numeri su un foglio. Non c'e' tracciabilita' del percorso logico che li ha generati."
- "Confrontare 5 scenari diversi richiede di costruire manualmente una tabella in Excel. Dovrebbe essere automatico."

**Livello tecnico:** Intermedio-alto. Sa usare Python base, R, Excel avanzato. Capisce i concetti statistici e la struttura dei modelli. Non scrive codice di produzione. Apprezza le interfacce interattive (slider, grafici dinamici) ma non ha bisogno di vedere il codice o il JSON raw.

**Come usa ExplainableEngine:** Sofia usa la dashboard quotidianamente per UC-004 (what-if analysis) e settimanalmente per UC-005 (monitoraggio). Gli slider interattivi sono il suo strumento principale. Salva e confronta scenari. Esporta i risultati per i report al CdA e per la documentazione regolamentare. Per lei la dashboard e' il "simulatore di rischio".

---

### Alex -- Backend Developer

| Campo | Dettaglio |
|-------|-----------|
| **Nome** | Alex Ferrara |
| **Eta'** | 29 |
| **Ruolo** | Senior Backend Developer presso una fintech cliente di ExplainableEngine |
| **Esperienza** | 6 anni di sviluppo backend, expertise in Go e distributed systems |
| **Team** | Platform Engineering (12 persone) |

**Goals:**

- Integrare ExplainableEngine nel sistema del proprio datore di lavoro nel minor tempo possibile.
- Capire il comportamento dell'API con diversi payload prima di scrivere codice di produzione.
- Avere code snippets pronti nel proprio linguaggio (Go, Python) per accelerare l'integrazione.
- Troubleshootare rapidamente i problemi di integrazione senza aprire support ticket.

**Frustrations:**

- "La documentazione OpenAPI e' necessaria ma non sufficiente. Ho bisogno di vedere l'API in azione con dati reali."
- "Ogni volta che devo testare un edge case devo scrivere un test ad-hoc. Vorrei un playground per esplorare."
- "Non capisco perche' la mia richiesta restituisce 422. Il messaggio di errore dice 'validation failed' ma non mi dice quale campo e' sbagliato." (Nota: questo e' gia' risolto nel backend con ValidationDetail, ma Alex non lo sa perche' non ha ancora provato.)
- "Devo convincere il mio team lead che l'integrazione e' fattibile. Una demo interattiva vale piu' di 100 slide."

**Livello tecnico:** Esperto. Scrive codice quotidianamente. Legge JSON, YAML, OpenAPI specs. Usa Postman, cURL, IDE con REST client. Preferisce interfacce developer-friendly: syntax highlighting, copy-to-clipboard, code generation. Non ha bisogno di visualizzazioni "belle" -- vuole velocita' e precisione.

**Come usa ExplainableEngine:** Alex usa la dashboard solo nella fase di integrazione (2-4 settimane). Si concentra su UC-006 (API Playground): testa gli endpoint, genera code snippets, esplora gli errori. Dopo l'integrazione, non torna sulla dashboard: interagisce direttamente con l'API dal codice. Per lui la dashboard e' il "ponte verso l'integrazione".

---

## 4. Feature Prioritization Matrix

### Criteri di valutazione

- **Business Value:** impatto su revenue, compliance, retention, efficienza (1-5)
- **Effort:** complessita' di implementazione in sprint (S/M/L/XL)
- **Priority:** MoSCoW (Must / Should / Could / Won't)

| # | Feature | Use Cases | Personas | Business Value | Effort | Priority | Note |
|---|---------|-----------|----------|----------------|--------|----------|------|
| F-001 | **Explanation Explorer (graph + breakdown)** | UC-002, UC-001 | Marco, Elena | 5 | M | **Must-have** | Cuore della dashboard. Visualizza il DAG interattivo e il breakdown numerico. Consuma `GET /explain/{id}` e `GET /explain/{id}/graph`. Senza questa feature la dashboard non ha ragione di esistere. |
| F-002 | **Audit Trail (search, filter, export)** | UC-001 | Elena | 5 | M | **Must-have** | Prerequisito regolamentare. Ricerca per data/target/confidence, paginazione cursor-based, export CSV/JSON. Blocca la vendita a clienti regolamentati se assente. |
| F-003 | **User Authentication** | Tutti | Tutti | 5 | M | **Must-have** | RBAC con ruoli: admin, analyst, compliance, sales, developer. Senza auth, nessun dato sensibile puo' essere esposto nella dashboard. Prerequisito per tutte le altre feature. |
| F-004 | **What-if Simulator with Sliders** | UC-004 | Sofia, Marco | 4 | M | **Should-have** | L'endpoint `POST /explain/{id}/what-if` esiste gia'. Frontend: slider per componente, delta visualization, scenario saving. Alto valore, effort moderato grazie al backend gia' pronto. |
| F-005 | **Narrative Viewer (multi-language)** | UC-003, UC-001 | Roberto, Elena | 4 | S | **Should-have** | Consuma `GET /explain/{id}/narrative`. Frontend semplice: rendering markdown/text, language switcher. Piccolo effort, alto impatto sulla comunicazione client. |
| F-006 | **PDF/Email Report Generation** | UC-003, UC-001 | Roberto, Elena | 4 | L | **Should-have** | Template brandizzato, inclusione selettiva di sezioni, generazione PDF server-side o client-side. Effort elevato per il templating e la qualita' del PDF. Dipende da F-001 e F-005. |
| F-007 | **Live Monitoring Dashboard** | UC-005 | Operations | 3 | L | **Should-have** | Polling periodico, contatori real-time, alert su soglie configurabili. Richiede un flusso continuo di dati. Valore alto ma solo per clienti con trading automatizzato. |
| F-008 | **API Playground / Swagger** | UC-006 | Alex | 3 | M | **Nice-to-have** | Interactive API explorer con code generation. Swagger UI embeddata copre il 70% del bisogno con effort minimo; un playground custom aggiunge il restante 30%. |
| F-009 | **Team/Org Management** | Tutti (admin) | Admin | 2 | M | **Nice-to-have** | Gestione utenti, ruoli, team, permessi. Necessario a lungo termine ma non per l'MVP. Inizialmente, user management puo' essere gestito via config file o admin API. |
| F-010 | **Webhook Notifications** | UC-005, UC-003 | Operations, Roberto | 2 | S | **Nice-to-have** | Notifiche su eventi (alert, report generato) via webhook a sistemi esterni (Slack, email, CRM). Piccolo effort, valore incrementale. Diventa should-have con la crescita della base clienti. |

### Mappa visuale (Value vs Effort)

```
Business Value
     5 |  F-001(M)  F-002(M)  F-003(M)
       |
     4 |  F-005(S)  F-004(M)            F-006(L)
       |
     3 |            F-008(M)            F-007(L)
       |
     2 |  F-010(S)  F-009(M)
       |
     1 |
       +----S---------M---------L---------XL----> Effort
```

### Roadmap consigliata

**Phase 1 -- Foundation (Sprint 1-2):** F-003 (Auth), F-001 (Explorer), F-002 (Audit Trail)
Queste tre feature compongono il prodotto minimo vendibile. Auth e' il prerequisito tecnico. Explorer e' il cuore funzionale. Audit Trail sblocca i clienti regolamentati.

**Phase 2 -- Value Expansion (Sprint 3-4):** F-005 (Narrative), F-004 (What-if), F-006 (PDF Reports)
Queste feature ampliano il pubblico dalla sola audience tecnica (Marco) a quella business (Roberto, Sofia, Elena). La narrative e il PDF sono gli strumenti che Roberto usa quotidianamente. Il what-if simulator e' lo strumento di Sofia.

**Phase 3 -- Scale & Automation (Sprint 5-6):** F-007 (Monitoring), F-008 (API Playground), F-010 (Webhooks)
Queste feature supportano l'operativita' a scala e l'integrazione con ecosistemi esterni. Diventano critiche quando il numero di clienti e il volume di explanation crescono.

**Phase 4 -- Enterprise (Sprint 7+):** F-009 (Team/Org Management)
Gestione multi-tenant, permessi granulari, SSO. Necessario per clienti enterprise con requisiti di governance.

---

## 5. Success Metrics

### KPI primari

| KPI | Definizione | Baseline (pre-dashboard) | Target (6 mesi) | Target (12 mesi) | Come si misura |
|-----|-------------|--------------------------|------------------|-------------------|----------------|
| **DAU** (Daily Active Users) | Utenti unici che accedono alla dashboard in un giorno | 0 | 20 | 50 | Analytics (page views con user ID) |
| **MAU** (Monthly Active Users) | Utenti unici che accedono in un mese | 0 | 50 | 150 | Analytics |
| **DAU/MAU Ratio** | Stickiness: quanto spesso gli utenti tornano | N/A | > 0.3 | > 0.4 | DAU / MAU |
| **Time-to-Insight** | Tempo dal login al primo "aha moment" (prima interazione significativa con un'explanation) | N/A (no dashboard) | < 60 secondi | < 30 secondi | Analytics (tempo tra login e primo click su explanation detail) |
| **Audit Preparation Time** | Tempo per preparare un pacchetto documentale per audit | 3-5 giorni (manuale) | < 4 ore | < 1 ora | Survey + time tracking |
| **Audit Compliance Rate** | Percentuale di audit superati senza rilievi sulla tracciabilita' | Variabile | 95% | 99% | Report di audit |
| **Client Satisfaction (NPS delta)** | Variazione del Net Promoter Score attribuibile alla trasparenza | NPS baseline del cliente | +10 punti | +20 punti | Survey NPS pre/post dashboard |
| **Integration Time** | Tempo dal primo accesso alla documentazione al go-live dell'integrazione API | 2-4 settimane | 1-2 settimane | 3-7 giorni | CRM tracking (data contratto -> data primo API call in produzione) |
| **Support Ticket Volume** | Numero di ticket relativi a "non capisco il risultato" o "serve una spiegazione" | Baseline da misurare | -50% | -80% | Ticketing system |
| **Report Generation** | Numero di report PDF generati e inviati ai clienti | 0 (manuale) | 100/mese | 500/mese | Analytics |

### KPI secondari (health metrics)

| KPI | Definizione | Target | Come si misura |
|-----|-------------|--------|----------------|
| **Page Load Time (P95)** | Tempo di caricamento della pagina al 95esimo percentile | < 2 secondi | RUM (Real User Monitoring) |
| **API Call Success Rate** | Percentuale di chiamate dashboard -> engine che restituiscono 2xx | > 99.5% | Monitoring |
| **Error Rate** | Percentuale di sessioni con almeno un errore visibile all'utente | < 1% | Error tracking (Sentry/equivalent) |
| **Feature Adoption** | Percentuale di utenti che usano ciascuna feature almeno 1x/mese | > 60% per Must-have features | Analytics per feature |
| **Export Completion Rate** | Percentuale di export (PDF/CSV) avviati che completano con successo | > 98% | Analytics |

### Metriche di business (lagging indicators)

| KPI | Definizione | Target (12 mesi) | Come si misura |
|-----|-------------|-------------------|----------------|
| **Client Retention Rate** | Percentuale di clienti che rinnovano, attribuibile a trasparenza | +5% vs baseline | CRM + churn analysis |
| **Upsell Rate** | Percentuale di clienti che passano da API-only a API+Dashboard | 30% | CRM |
| **New Client Acquisition** | Nuovi clienti dove la dashboard e' stata citata come fattore decisionale | 10 clienti/anno | Sales feedback |
| **Regulatory Savings** | Costo evitato per sanzioni e preparazione audit | 500K EUR/anno | Stima basata su tempo risparmiato + sanzioni evitate |

---

## 6. Risks & Mitigations

### R-001: Performance Degradation con Grafi Grandi

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | Il rendering di grafi DAG con 100+ nodi nel browser potrebbe essere lento o inutilizzabile, degradando l'esperienza utente e rendendo inutile la feature core (F-001). |
| **Probabilita'** | Alta -- i modelli quantitativi reali hanno spesso 50-200 componenti con sub-components su 3-5 livelli. |
| **Impatto** | Alto -- se l'Explorer e' lento, i power user (Marco) tornano a cURL e la dashboard perde la sua ragione d'essere. |
| **Mitigazione** | (1) Lazy loading del grafo: renderizzare solo il livello corrente + 1, espandere on-click. (2) Collapsare automaticamente i nodi sotto soglia di contribution_pct (es. < 2%). (3) Usare WebGL-based rendering (es. deck.gl, Sigma.js) invece di SVG per grafi > 50 nodi. (4) Implementare `max_depth` nel parametro della richiesta API per limitare il payload. (5) Testare con dataset realistici (100, 500, 1000 nodi) fin dal primo sprint. |
| **Owner** | Tech Lead Frontend |

### R-002: Disallineamento tra Backend API e Frontend

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | Il backend evolve (nuovi campi, cambiamenti di schema) e il frontend non viene aggiornato, causando errori o dati mancanti nella dashboard. |
| **Probabilita'** | Media -- il backend e' in sviluppo attivo (v0.4.0, 4 sprint completati). |
| **Impatto** | Medio -- errori di rendering o dati incompleti riducono la fiducia nella dashboard. |
| **Mitigazione** | (1) Contract testing: generare i tipi TypeScript direttamente dalla specifica OpenAPI (`api-contract.yaml`) con un code generator (es. openapi-typescript). (2) CI check: se la specifica cambia, il build del frontend fallisce finche' i tipi non vengono rigenerati. (3) Versionare l'API (`/api/v1/`) e supportare backward compatibility. (4) Integration test automatici che verificano la dashboard contro l'engine reale in staging. |
| **Owner** | Tech Lead Backend + Tech Lead Frontend (condiviso) |

### R-003: Adozione Bassa da Parte degli Utenti Non Tecnici

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | Elena (Compliance) e Roberto (Sales) trovano la dashboard troppo tecnica o confusa e continuano a chiedere al team IT di estrarre i dati manualmente. La dashboard esiste ma non viene usata. |
| **Probabilita'** | Media-alta -- la terminologia del dominio (DAG, confidence propagation, deterministic hash) e' intrinsecamente tecnica. |
| **Impatto** | Alto -- se gli utenti non tecnici non adottano, la dashboard copre solo UC-002 (debugging) e non giustifica l'investimento. |
| **Mitigazione** | (1) Design separato per persona: la vista "Compliance" mostra terminologia diversa dalla vista "Analyst". Es. "Affidabilita'" invece di "Confidence Score". (2) User testing con Elena e Roberto (o proxy) prima di ogni release. (3) Onboarding guidato: primo accesso con tour interattivo (5 step max). (4) Vista "Client View" semplificata per UC-003 con solo: valore, top 3 driver, semaforo di confidence, narrative. (5) Investire in narrative quality: se la spiegazione in linguaggio naturale e' chiara, l'utente non ha bisogno di capire il grafo. |
| **Owner** | Product Manager + UX Designer |

### R-004: Sicurezza e Data Leakage

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | Le explanation contengono dati sensibili (modelli proprietari, parametri di trading, dati finanziari dei clienti). Una falla di sicurezza nella dashboard espone questi dati. Un link di condivisione report finisce nelle mani sbagliate. |
| **Probabilita'** | Media -- la dashboard e' un nuovo vettore di attacco che non esisteva con la sola API. |
| **Impatto** | Critico -- data breach in ambito finanziario = sanzioni GDPR + danno reputazionale irreversibile + possibile perdita del vantaggio competitivo (modelli proprietari esposti). |
| **Mitigazione** | (1) Auth robusta fin dal giorno 1 (F-003 e' Must-have Phase 1): OAuth2/OIDC, MFA obbligatorio per ruoli compliance e admin. (2) RBAC granulare: il Sales vede solo le explanation dei suoi clienti, il Compliance vede tutto ma non puo' modificare. (3) Audit log di ogni accesso e ogni export (chi, quando, cosa, da quale IP). (4) Link di condivisione con scadenza, one-time-use, e IP whitelisting opzionale. (5) Penetration test prima del go-live. (6) Encryption at rest e in transit (TLS 1.3, encrypted storage). (7) CORS configurato restrittivamente (middleware gia' presente nel backend). |
| **Owner** | Security Lead + Tech Lead Backend |

### R-005: Scalabilita' del Polling per il Monitoraggio Real-time

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | La dashboard di monitoraggio (UC-005) usa polling HTTP per aggiornare i dati. Con 20+ utenti connessi simultaneamente e polling ogni 10 secondi, il carico sull'engine cresce linearmente (120 req/min per 20 utenti). In caso di incidente, tutti gli utenti si connettono contemporaneamente amplificando il carico. |
| **Probabilita'** | Media -- il monitoraggio e' Phase 3, ma il problema va progettato prima. |
| **Impatto** | Medio -- degradazione dell'engine che impatta tutti gli utenti, non solo la dashboard. |
| **Mitigazione** | (1) Caching layer tra dashboard e engine: la dashboard non interroga l'engine direttamente per il monitoraggio, ma un aggregator service che mantiene lo stato. (2) SSE (Server-Sent Events) come alternativa al polling: una singola connessione per utente, push dal server. (3) Rate limiting per utente sulla dashboard. (4) Circuit breaker: se l'engine non risponde, la dashboard mostra l'ultimo dato noto con un warning "stale data" invece di retry aggressivi. (5) Pre-aggregazione lato backend delle metriche di monitoraggio (contatori, medie) per evitare che la dashboard debba fare N chiamate per calcolare statistiche aggregate. |
| **Owner** | Tech Lead Backend |

### R-006: Scope Creep e Feature Bloat

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | Ogni stakeholder vuole "la sua feature" nella dashboard. Il Sales vuole il CRM integrato, il Quant vuole il notebook integrato, la Compliance vuole il workflow di approvazione. La dashboard diventa un monolite che cerca di fare tutto e non fa niente bene. |
| **Probabilita'** | Alta -- e' il rischio numero uno di qualsiasi prodotto dashboard. |
| **Impatto** | Alto -- ritardo nelle delivery, qualita' mediocre, utenti frustrati da un'interfaccia sovraccarica. |
| **Mitigazione** | (1) MoSCoW rigoroso: questa product strategy definisce cosa e' Must/Should/Could/Won't. Ogni nuova richiesta deve passare per una valutazione formale. (2) Phase gates: non si inizia Phase 2 finche' Phase 1 non e' completata e validata con utenti reali. (3) "One thing well" per release: ogni sprint ha un tema chiaro (es. Sprint 1 = Auth + Explorer, non Auth + Explorer + Reporting + Monitoring). (4) Feedback loop corto: rilasciare feature incrementali e misurare l'adozione prima di aggiungere complessita'. (5) Dire "no" o "not now" esplicitamente e documentare il razionale. |
| **Owner** | Product Manager |

### R-007: Dipendenza dal Backend Go per Feature Frontend

| Campo | Dettaglio |
|-------|-----------|
| **Rischio** | Alcune feature della dashboard richiedono endpoint o funzionalita' che non esistono ancora nel backend. Esempio: ricerca full-text sulle explanation (per audit trail), aggregazione temporale (per monitoraggio), gestione utenti. Il frontend si blocca in attesa del backend. |
| **Probabilita'** | Alta -- il backend e' stato progettato per l'uso API-to-API, non per una dashboard. Mancano endpoint di listing, ricerca, aggregazione. |
| **Impatto** | Medio-alto -- ritardi nella roadmap, compromessi di design (workaround frontend per feature che dovrebbero essere backend). |
| **Mitigazione** | (1) Mappare esplicitamente le dipendenze backend per ogni feature frontend in fase di sprint planning. (2) Prioritizzare lo sviluppo backend delle API mancanti prima del frontend che le consuma. (3) Per il MVP, accettare compromessi ragionevoli: la ricerca puo' essere client-side sulle explanation gia' caricate (non scala, ma funziona per < 1000 explanation). (4) BFF (Backend for Frontend) pattern: un sottile layer intermedio che aggrega chiamate e aggiunge funzionalita' specifiche per la dashboard senza modificare l'engine core. (5) Coordinamento settimanale backend/frontend per allineare le priorita'. |
| **Owner** | Product Manager + Tech Lead Backend |

---

## Appendice A: Mapping Endpoint API -> Feature Dashboard

| Endpoint Backend | Feature Dashboard | Note |
|------------------|-------------------|------|
| `GET /health` | Status indicator nell'header | Indicatore "engine online/offline" |
| `POST /api/v1/explain` | API Playground (F-008) | Usato solo dal playground; il flusso normale e' read-only |
| `GET /api/v1/explain/{id}` | Explorer (F-001), Audit Trail (F-002) | Endpoint piu' usato dalla dashboard |
| `GET /api/v1/explain/{id}/graph` | Explorer - Graph View (F-001) | Rendering DAG interattivo |
| `GET /api/v1/explain/{id}/narrative` | Narrative Viewer (F-005), PDF Report (F-006) | Multi-language support |
| `POST /api/v1/explain/{id}/what-if` | What-if Simulator (F-004) | Slider -> override payload -> delta view |

**Endpoint mancanti nel backend (da sviluppare):**

| Endpoint necessario | Feature Dashboard | Priorita' |
|---------------------|-------------------|-----------|
| `GET /api/v1/explain?target=X&from=DATE&to=DATE&confidence_min=N` | Audit Trail search (F-002) | Must-have Phase 1 |
| `GET /api/v1/explain/stats?from=DATE&to=DATE` | Live Monitoring (F-007) | Should-have Phase 3 |
| `POST /api/v1/auth/*` | Authentication (F-003) | Must-have Phase 1 |
| `GET /api/v1/users/*` | Team Management (F-009) | Nice-to-have Phase 4 |

---

## Appendice B: Glossario

| Termine | Definizione |
|---------|-------------|
| **AIP** | Algorithmic Investment Platform -- piattaforma di trading quantitativo, primo cliente e caso d'uso di ExplainableEngine |
| **AUM** | Assets Under Management -- patrimonio gestito |
| **BFF** | Backend for Frontend -- layer API intermedio specifico per le esigenze della dashboard |
| **DAG** | Directed Acyclic Graph -- struttura dati core dell'engine |
| **DAU/MAU** | Daily/Monthly Active Users |
| **DX** | Developer Experience |
| **EU AI Act** | Regolamento europeo sull'intelligenza artificiale (2024/1689) |
| **GRC** | Governance, Risk, and Compliance |
| **MiFID II** | Markets in Financial Instruments Directive II |
| **MoSCoW** | Must-have, Should-have, Could-have, Won't-have -- metodo di prioritizzazione |
| **NPS** | Net Promoter Score |
| **RBAC** | Role-Based Access Control |
| **RUM** | Real User Monitoring |
| **SGR** | Societa' di Gestione del Risparmio (Italian asset management company) |
| **SSE** | Server-Sent Events |
