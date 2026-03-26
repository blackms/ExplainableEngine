# Explainable Engine Dashboard — UX Redesign Specification

> Documento di prodotto completo per il redesign dell'Explainable Engine Dashboard.
> Destinatari: designer, frontend developer, product owner.
> Versione: 1.0 — 2026-03-26

---

## 1. PRODUCT IDENTITY

### Nome prodotto

**Explainable Engine**

Il nome resta invariato. E' gia riconosciuto internamente, coerente col backend, e comunica esattamente la value proposition: un motore che rende spiegabili le decisioni. Non servono nomi piu "marketing" — il target e B2B fintech, la credibilita conta piu della creativita.

### Tagline

> **"Ogni numero ha una storia. Noi la rendiamo leggibile."**

Variante tecnica per la documentazione API:
> *"Graph-first explainability for algorithmic decisions."*

### Personalita del prodotto

| Attributo | Descrizione | Cosa significa nel design |
|-----------|-------------|--------------------------|
| **Professionale** | Mai giocoso, mai frivolo — gestiamo decisioni finanziarie | Tipografia pulita (Geist), colori neutri, icone minimali |
| **Avvicinabile** | Non intimidatorio — anche il sales deve capire | Empty state con suggerimenti, tooltip su ogni metrica, linguaggio chiaro |
| **Data-driven** | I numeri sono il cuore, non la decorazione | Numeri grandi e leggibili, chart interattivi, precision nei decimali |
| **Affidabile** | Ogni dato mostrato e verificabile e tracciabile | Badge di confidence ovunque, audit trail accessibile, hash deterministico visibile |

### Sensazione target (Target Feeling)

L'utente che chiude il browser deve pensare:

> **"Capisco questo numero e mi fido."**

Non "che bel design" — ma "ho capito perche il valore e 0.72, so cosa lo guida, e posso spiegarlo a qualcun altro."

### Design Principles

1. **Information density > decorazione**: ogni pixel deve comunicare un dato o un'azione. Zero ornamenti.
2. **Progressive disclosure**: il primo livello mostra il risultato; il secondo il breakdown; il terzo il grafo completo. L'utente sceglie quanto andare in profondita.
3. **Zero dead-end**: ogni schermata suggerisce il prossimo passo. Mai una pagina vuota senza azione.
4. **Speed is a feature**: skeleton loading gia implementato — mantenerlo ovunque. Target: < 200ms per il primo render significativo.

---

## 2. INFORMATION ARCHITECTURE

### Stato attuale (problemi)

La sidebar attuale ha 5 voci piatte: Home, Audit, Monitor, Playground, Settings. I problemi:

- **Home = Create**: la home page e un form di creazione. Non c'e una vera overview.
- **Audit disconnesso**: l'audit trail sembra una feature separata, non il modo naturale di trovare le explanation.
- **Explain detail senza navigazione**: una volta dentro `/explain/[id]`, non c'e modo di navigare tra le sezioni (graph, what-if) senza scroll verticale.
- **Nessuna gerarchia**: tutte le voci hanno lo stesso peso visivo. Un quant e un compliance officer vedono la stessa sidebar.

### Architettura proposta

```
SIDEBAR NAVIGATION (3 livelli visivi)

[Logo] Explainable Engine
─────────────────────────────────
PRINCIPALE
  Dashboard                    /dashboard
  Explanations                 /explanations
    ├─ + Nuova                 /explanations/new
    ├─ Archivio                /explanations       (stessa route, tab "archivio")
    └─ [id] Dettaglio          /explanations/[id]
        ├─ Overview            (tab default)
        ├─ Grafo               (tab)
        ├─ What-If             (tab)
        ├─ AI Insights         (tab)
        └─ Export              (tab)

OPERAZIONI
  Monitoring                   /monitoring
    ├─ Live Feed               (tab default)
    ├─ Alerts                  (tab)
    └─ Statistiche             (tab)

SVILUPPO
  API Playground               /developer/playground
  Documentazione               /developer/docs

─────────────────────────────────
FOOTER SIDEBAR
  Settings                     /settings
  [Avatar] User Menu
```

### Motivazioni delle scelte

1. **"Explanations" come sezione primaria**: e il cuore del prodotto. Raggruppa creazione, ricerca (ex-Audit) e dettaglio sotto un unico concetto. L'utente non pensa "vado in audit" — pensa "cerco una explanation".

2. **Tab nel dettaglio explanation**: invece di una pagina lunghissima con scroll, il dettaglio ha tab orizzontali. Il contenuto attuale di `/explain/[id]` diventa il tab "Overview". Graph e what-if, oggi in sub-route separate, diventano tab nello stesso layout.

3. **Monitoring come sezione separata**: le operazioni in tempo reale hanno una mentalita diversa ("sorveglio") rispetto alle explanation ("analizzo"). Meritano una sezione dedicata.

4. **Developer in fondo**: e una sezione di nicchia. Non deve competere visivamente con le feature principali.

### Routing

| Percorso | Componente | Note |
|----------|-----------|-------|
| `/` | Redirect a `/dashboard` | |
| `/dashboard` | DashboardPage | Overview, metriche, quick actions |
| `/explanations` | ExplanationsPage | Archivio + filtri (ex-Audit) |
| `/explanations/new` | CreateExplanationPage | Form di creazione |
| `/explanations/[id]` | ExplanationDetailPage | Layout con tab |
| `/explanations/[id]?tab=overview` | OverviewTab | Default |
| `/explanations/[id]?tab=graph` | GraphTab | DAG interattivo |
| `/explanations/[id]?tab=whatif` | WhatIfTab | Simulatore |
| `/explanations/[id]?tab=ai` | AIInsightsTab | Narrativa + Q&A + Summary |
| `/explanations/[id]?tab=export` | ExportTab | PDF, share, download |
| `/monitoring` | MonitoringPage | Layout con tab |
| `/monitoring?tab=feed` | LiveFeedTab | Default |
| `/monitoring?tab=alerts` | AlertsTab | |
| `/monitoring?tab=stats` | StatsTab | |
| `/developer/playground` | PlaygroundPage | |
| `/developer/docs` | DocsPage | Swagger/OpenAPI embed |
| `/settings` | SettingsPage | Preferenze + API keys |

### Sidebar: Specifiche di comportamento

- **Desktop** (>= 1024px): sidebar fissa, 256px di larghezza. Collassabile a 64px (solo icone) con toggle button.
- **Tablet** (768-1023px): sidebar collassata di default (solo icone), espandibile al click.
- **Mobile** (< 768px): sidebar nascosta. Hamburger menu in alto a sinistra. Overlay con backdrop scuro.
- **Voce attiva**: background `sidebar-accent`, testo `sidebar-accent-foreground`, bordo sinistro 2px `primary`.
- **Sotto-sezioni**: visibili solo quando la sezione padre e attiva. Animazione: slide-down 150ms ease-out.
- **Badge contatore**: su "Alerts" nel monitoring, mostra il numero di alert non letti. Pallino rosso, testo bianco, min-width 20px.

---

## 3. USER JOURNEYS

### Journey 1: Primo accesso e onboarding

**Persona**: Qualsiasi utente nuovo (quant, analyst, developer)
**Goal**: Capire cosa fa il prodotto e creare la prima explanation
**Entry point**: URL diretto o invito via email

#### Steps

**Step 1 — Welcome Screen**
- L'utente apre `/dashboard` per la prima volta
- Invece del dashboard standard, vede una **welcome card** centrata nel contenuto principale
- Contenuto:
  - Titolo: "Benvenuto in Explainable Engine"
  - Sottotitolo: "Trasforma qualsiasi decisione algoritmica in una spiegazione chiara, verificabile e condivisibile."
  - Tre card orizzontali con icone:
    1. **Crea** (icona: Plus) — "Inserisci i componenti di una decisione"
    2. **Analizza** (icona: BarChart) — "Esplora breakdown, driver e grafo causale"
    3. **Comprendi** (icona: MessageSquare) — "Chiedi all'AI di spiegare in linguaggio naturale"
  - CTA primario: "Crea la tua prima Explanation" (bottone grande, primary color)
  - Link secondario: "Oppure esplora l'API Playground"
- **Condizione di sparizione**: la welcome card scompare dopo che l'utente ha creato almeno 1 explanation. Stato salvato in `localStorage` (key: `ee-onboarding-complete`). Link "Rivedi il tour" nelle Settings per ripristinare.

**Step 2 — Form guidato**
- Click sul CTA porta a `/explanations/new`
- Il form ha lo stesso layout attuale MA con queste aggiunte:
  - **Tooltip su ogni campo**: icona (?) accanto a ogni label. Al hover:
    - Target name: "Il nome dell'output che vuoi spiegare, es. 'portfolio_risk_score'"
    - Value: "Il valore finale calcolato dal tuo sistema"
    - Components: "I fattori che contribuiscono al risultato. Ogni componente ha un nome, valore, peso (0-1) e confidence (0-1)"
  - **Template pre-compilato**: un link "Usa esempio" sopra il form che pre-compila:
    ```
    Target: market_regime_score
    Value: 0.72
    Components:
      - trend_strength: 0.8, weight 0.4, confidence 0.9
      - volatility_index: 0.6, weight 0.35, confidence 0.85
      - momentum_signal: 0.75, weight 0.25, confidence 0.95
    ```
  - **Validazione inline**: bordo rosso + messaggio sotto il campo. Messaggi:
    - Target vuoto: "Inserisci un nome per l'output"
    - Nessun componente: "Aggiungi almeno un componente"
    - Weight fuori range: "Il peso deve essere tra 0 e 1"
    - Confidence fuori range: "La confidence deve essere tra 0 e 1"

**Step 3 — Submit e celebrazione**
- Click su "Submit" → bottone mostra spinner + testo "Analisi in corso..."
- Risposta ricevuta (< 200ms tipicamente) → redirect a `/explanations/[id]`
- **Primo risultato speciale**: se e la prima explanation dell'utente (check `localStorage`):
  - Banner in alto nella pagina dettaglio: sfondo gradient leggero (primary/5%), bordo primary
  - Testo: "La tua prima Explanation! Esplora il breakdown qui sotto, o chiedi all'AI di spiegartela."
  - Bottone: "Prova l'AI" → scroll smooth al tab AI Insights
  - Dismiss: icona X a destra. Una volta chiuso, non riappare.

**Step 4 — Esplorazione guidata**
- La pagina dettaglio mostra il tab Overview di default
- Se e la prima visita, tre **contextual tip** appaiono in sequenza (non tutti insieme):
  1. Sul SummaryCard: "Questo e il riepilogo. Il numero grande e il valore finale, il badge la confidence."
  2. Sul BreakdownChart: "Qui vedi come ogni componente contribuisce al risultato."
  3. Sui tab: "Esplora il Grafo per le dipendenze, What-If per simulazioni, AI per spiegazioni in linguaggio naturale."
- Ogni tip ha: testo, freccia che punta all'elemento, bottone "Capito" per passare al successivo.
- I tip usano un overlay semi-trasparente (backdrop-filter: blur(2px)) che evidenzia solo l'elemento target.

**Success state**: L'utente ha creato un'explanation, ne ha visto il breakdown, e sa dove trovare le feature avanzate.

**Edge cases**:
- API backend non raggiungibile: toast error "Impossibile connettersi al server. Riprova tra qualche secondo." + retry button.
- Utente chiude prima di completare: alla prossima apertura rivede la welcome card.

---

### Journey 2: Quant che fa debug di un segnale

**Persona**: Quant/Trader — utente esperto, ha fretta, conosce i dati
**Goal**: Capire perche un segnale ha un valore inatteso e identificare il componente anomalo
**Entry point**: Link diretto a `/explanations/[id]` (da alert email, Slack, o sistema interno)

#### Steps

**Step 1 — Atterraggio sul dettaglio**
- L'utente clicca un link → `/explanations/abc123`
- **Primo render** (< 100ms): skeleton loading con la struttura della pagina gia visibile
- **Dati caricati** (< 300ms): il tab Overview mostra immediatamente:
  - **Header**: nome target ("portfolio_risk_score") + breadcrumb "Explanations > portfolio_risk_score"
  - **SummaryCard** in alto: valore finale grande (es. "0.72") a sinistra, confidence badge a destra (es. "81% Confidence" con colore semantico: verde >80%, giallo 50-80%, rosso <50%), missing impact sotto ("12% dati mancanti" in arancione se >5%)
  - **Top driver**: sotto il valore, una riga: "Guidato da: trend_strength (impatto 42%)" — click porta al driver nel breakdown

**Step 2 — Analisi del breakdown**
- Sotto il SummaryCard: BreakdownChart (bar chart orizzontale, gia implementato) + DriverRanking (tabella)
- Il quant identifica visivamente il componente anomalo — ad esempio `volatility_index` con un contributo insolitamente alto
- **Interazione**: hover sulla barra → tooltip con dettaglio: "volatility_index: valore 0.92, peso 35%, contributo 0.322 (45% del totale)"
- **Click sulla barra** → highlight del componente + pannello laterale destro (slide-in, 400px) con:
  - Storico del componente (se disponibile)
  - Confidence del singolo nodo
  - Link "Simula variazione" → apre What-If tab con questo componente pre-selezionato

**Step 3 — Esplorazione del grafo**
- L'utente clicca il tab "Grafo"
- **Rendering**: DAG interattivo (ReactFlow, gia implementato)
- Il nodo del componente anomalo e evidenziato automaticamente se l'utente e arrivato dal breakdown (query param `?highlight=volatility_index`)
- **Interazioni**:
  - Zoom: scroll wheel o pinch
  - Pan: drag sullo sfondo
  - Click nodo: seleziona, mostra dettaglio nel pannello laterale
  - Doppio click nodo: espande i figli (se compressi)
- Il quant risale la catena: volatility_index <- raw_vol_data <- market_data_feed
- Identifica che `market_data_feed` ha confidence 0.3 → il problema e nei dati sorgente

**Step 4 — Simulazione What-If**
- Tab "What-If" → WhatIfSimulator (gia implementato)
- Se arrivo dal grafo con un nodo selezionato, il simulatore pre-seleziona quel componente
- L'utente modifica `volatility_index` da 0.92 a 0.60 (valore atteso)
- **Risultato in tempo reale**: il valore finale cambia da 0.72 a 0.61, con diff evidenziata
- DiffTable mostra: ogni componente, valore originale, valore modificato, delta
- SensitivityRanking mostra: quali componenti hanno l'impatto maggiore

**Step 5 — Conclusione**
- Il quant ha identificato il problema (dati di volatilita anomali), verificato l'impatto (0.11 punti), e sa che il segnale corretto sarebbe ~0.61
- Chiude la tab. Tempo totale: < 2 minuti

**Success state**: L'utente ha fatto root-cause analysis completa senza uscire dalla dashboard.

**Edge cases**:
- Explanation non trovata (404): pagina "Explanation non trovata" con campo di ricerca inline + link "Torna all'archivio"
- Grafo troppo grande (>100 nodi): rendering progressivo. Prima i nodi a depth 0-2, poi espansione on-demand. Warning: "Grafo complesso — visualizzazione semplificata. Clicca per espandere."
- What-If timeout: "La simulazione sta impiegando piu del previsto. Il server potrebbe essere sovraccarico." + retry

---

### Journey 3: Sales che prepara un report per il cliente

**Persona**: Analyst/Sales — non tecnico, ha bisogno di spiegazioni semplici e presentabili
**Goal**: Generare una spiegazione comprensibile da un cliente e esportarla come PDF
**Entry point**: Dashboard → cerca l'explanation → AI Insights → Export

#### Steps

**Step 1 — Trovare l'explanation**
- L'utente apre `/explanations` (Archivio)
- Usa la barra di ricerca: digita il nome del target o parte dell'ID
- **Search behavior**: ricerca live con debounce 300ms. I risultati filtrano la tabella in real-time.
- Alternativa: usa i filtri rapidi (chips cliccabili sopra la tabella):
  - "Ultima settimana" | "Ultimo mese" | "Alta confidence" | "Bassa confidence"
- Trova l'explanation → click sulla riga → naviga a `/explanations/[id]`

**Step 2 — Generare Executive Summary**
- Click sul tab "AI Insights"
- Il tab mostra tre sezioni:
  1. **Narrativa AI**: spiegazione generata dall'LLM
  2. **Q&A Chat**: interfaccia conversazionale
  3. **Executive Summary**: generatore di report strutturato
- L'utente clicca su "Genera Executive Summary"
- **Form di configurazione** (inline, non modal):
  - Audience: dropdown con "Board" | "Tecnico" | "Cliente" → seleziona "Cliente"
  - Lingua: toggle "EN" | "IT" → seleziona "IT"
  - Bottone: "Genera Summary"
- **Loading state**: skeleton del summary con shimmer animation. Testo: "L'AI sta analizzando i dati..."
- **Risultato** (3-8 secondi): card strutturata con:
  - **Titolo**: generato dall'AI, es. "Analisi del Risk Score — Portafoglio Alfa"
  - **Summary**: 2-3 paragrafi in linguaggio non tecnico
  - **Key Findings**: lista puntata dei risultati principali
  - **Rischi**: eventuali warning (dati mancanti, bassa confidence)
  - **Raccomandazioni**: suggerimenti actionable
- Ogni sezione ha un'icona di copia a destra (click → "Copiato!" toast)

**Step 3 — Revisione e modifica**
- L'utente legge il summary generato
- Se soddisfatto, procede all'export
- Se vuole modificare: click "Rigenera con istruzioni" → campo di testo: "Rendi il tono piu rassicurante" o "Enfatizza la stabilita del modello"
- Rigenera → nuovo summary con le istruzioni applicate

**Step 4 — Export PDF**
- Tab "Export" (o icona download nel tab AI Insights)
- **Opzioni di export**:
  - Formato: PDF (default) | CSV (solo dati)
  - Contenuto: checkboxes
    - [x] Summary Card
    - [x] Breakdown Chart
    - [x] Driver Ranking
    - [ ] Grafo (opzionale, puo essere pesante)
    - [x] Executive Summary (se generato)
    - [ ] Dati grezzi
  - Branding: "Includi logo aziendale" toggle (logo configurabile nelle Settings)
- Click "Scarica PDF" → generazione client-side → download automatico
- **Toast**: "PDF scaricato: portfolio_risk_score_2026-03-26.pdf"

**Step 5 — Invio al cliente**
- L'utente ha il PDF. Alternativa: click "Copia link" per condividere un link diretto alla explanation (se la condivisione e abilitata nelle Settings).
- Tempo totale: < 3 minuti

**Success state**: PDF professionale pronto per il cliente, con spiegazioni comprensibili e dati verificabili.

**Edge cases**:
- LLM non disponibile: fallback alla narrativa template (gia implementato nel backend). Nota visibile: "Summary generato con template — LLM non disponibile al momento."
- Explanation con dati mancanti: il summary include automaticamente una sezione "Limitazioni" che spiega quali dati mancano e l'impatto sulla reliability.
- PDF generation fallisce: toast error + suggerimento "Prova a deselezionare il Grafo per ridurre le dimensioni del file."

---

### Journey 4: Compliance Officer che conduce un audit

**Persona**: Compliance — meticoloso, ha bisogno di tracciabilita completa, lavora con date range
**Goal**: Revisionare tutte le decisioni di un periodo, identificare quelle a bassa confidence, esportare evidenze
**Entry point**: `/explanations` con filtri

#### Steps

**Step 1 — Impostare i filtri**
- Apre `/explanations` → la pagina mostra l'archivio completo
- **Pannello filtri** (sopra la tabella, collapsabile):
  - **Date range**: due date picker (da / a). Preset rapidi: "Oggi", "Ultima settimana", "Ultimo mese", "Q1 2026", "Q4 2025", "Range personalizzato"
  - **Target**: campo di testo con autocomplete (suggerisce target gia visti)
  - **Confidence**: range slider doppio (min-max), da 0% a 100%
  - **Ricerca libera**: campo full-text che cerca in target, ID, metadata
- L'utente seleziona Q1 2026 (01/01/2026 — 31/03/2026) e imposta confidence 0%-50%
- **Risultati**: la tabella si aggiorna con i filtri applicati. Sopra la tabella: "47 risultati in Q1 2026 con confidence < 50%"

**Step 2 — Ordinare e navigare**
- **Tabella colonne**: ID (troncato), Target, Valore, Confidence, Data creazione, Azioni
- Click sull'header "Confidence" → ordinamento ascendente (freccia su). Primo click: asc, secondo: desc, terzo: reset.
- Le explanation con confidence < 30% hanno un **indicatore visivo**: icona warning arancione + riga con sfondo `destructive/5%`
- **Paginazione**: 20 righe per pagina (configurabile: 20/50/100). Cursor-based (gia implementato). Navigazione prev/next con contatore pagine.

**Step 3 — Review dettagliata**
- Click su una riga → apre `/explanations/[id]` **in una nuova tab** (comportamento: click = stessa tab, Cmd/Ctrl+click = nuova tab, oppure icona "apri in nuova tab" nella colonna Azioni)
- L'utente alterna tra la tabella e il dettaglio
- Nel dettaglio, controlla:
  - Breakdown: i contributi sommano al totale?
  - Confidence: quali nodi hanno confidence bassa?
  - Narrativa: la spiegazione e coerente?
  - Metadata: hash deterministico, timestamp, versione

**Step 4 — Export di massa**
- Torna alla tabella `/explanations`
- Click "Esporta" (bottone in alto a destra, accanto ai filtri)
- **Opzioni**:
  - Formato: CSV (default per audit) | JSON
  - Scope: "Risultati filtrati (47)" | "Selezionati (0)" — se l'utente ha selezionato righe con checkbox
  - Campi: checkboxes per includere/escludere colonne
- Click "Scarica CSV" → download immediato
- **Nome file**: `explainable-engine-audit-Q1-2026-low-confidence.csv`

**Step 5 — Export singole evidenze**
- Per le explanation piu critiche, l'utente genera PDF individuali:
  - Apre il dettaglio → tab Export → seleziona tutti i contenuti → "Scarica PDF"
  - Il PDF include il `deterministic_hash` nel footer come prova di integrita

**Success state**: L'officer ha una lista CSV completa per il report di compliance e PDF individuali come evidenza per i casi critici.

**Edge cases**:
- Troppi risultati (>10.000): la tabella mostra i primi 20, con note "10.247 risultati. Raffina i filtri per una ricerca piu precisa."
- Export CSV di grandi dimensioni: generazione asincrona. Toast: "Preparazione export in corso..." → notifica quando pronto.
- Filtri che non tornano risultati: empty state con illustrazione + "Nessuna explanation trovata per questi filtri. Prova ad ampliare il range di date."

---

### Journey 5: Ops che monitora decisioni in tempo reale

**Persona**: Ops/Risk Manager — sorveglia, reagisce rapidamente, ha bisogno di alert chiari
**Goal**: Monitorare le decisioni in real-time e reagire ad anomalie
**Entry point**: `/monitoring`

#### Steps

**Step 1 — Apertura monitoring**
- Naviga a `/monitoring` → tab "Live Feed" attivo di default
- **Header area**: tre StatsCards in alto:
  1. "Explanation oggi": numero (es. 342), trend (+12% vs ieri), sparkline ultimi 7 giorni
  2. "Confidence media": percentuale (es. 78%), colore semantico, sparkline
  3. "Alert attivi": numero (es. 3), badge rosso se > 0, sparkline

**Step 2 — Live Feed**
- Sotto le stats: feed in real-time
- Ogni entry mostra: timestamp, target name, valore, confidence badge, tipo di computation
- **Auto-refresh**: polling ogni 10 secondi. Indicatore visivo: pallino verde pulsante "Live" in alto a destra del feed
- **Nuove entry**: appaiono in cima con animazione slide-down + flash di sfondo (highlight per 2 secondi poi fade)
- **Pause**: bottone "Pausa" che ferma l'auto-refresh. Utile per leggere senza che il feed scorra.

**Step 3 — Alert notification**
- Un nuovo alert appare: una riga nel feed si evidenzia in rosso
- Contemporaneamente:
  - **Badge sidebar**: il contatore su "Alerts" nella sidebar incrementa
  - **Toast notification**: angolo in basso a destra, "Alert: Low confidence (0.3) su portfolio_risk_score". Rimane 8 secondi. Click → apre il dettaglio.
  - **Browser notification** (se permesso): notifica OS con lo stesso testo

**Step 4 — Investigazione alert**
- Click sull'alert (nel feed o nel toast) → naviga a `/explanations/[id]`
- Il dettaglio si apre con un **banner alert** in cima: sfondo `destructive/10%`, testo "Questa explanation ha generato un alert: Low confidence (0.30)"
- L'utente vede che 2 componenti su 5 hanno `missing: true`
- La confidence e bassa perche i dati mancanti impattano per il 35%

**Step 5 — Tab Alerts**
- Torna a `/monitoring?tab=alerts`
- Lista di tutti gli alert, ordinati per data (piu recenti in cima)
- Ogni alert: timestamp, tipo (low_confidence | missing_data | anomaly), target, valore, status (new | acknowledged | resolved)
- **Azioni per alert**:
  - "Acknowledge" → segna come preso in carico (cambia status, rimuove dal contatore badge)
  - "Apri dettaglio" → naviga all'explanation
  - "Risolvi" → segna come risolto con nota opzionale

**Step 6 — Escalation**
- L'utente copia il link dell'explanation (icona "Copia link" nel dettaglio)
- Lo incolla nel canale Slack del team quant
- (Futura integrazione: bottone "Notifica team" direttamente dalla dashboard)

**Success state**: L'operatore ha visto l'alert in tempo reale, identificato la causa, e notificato il team responsabile.

**Edge cases**:
- Connessione persa: il pallino "Live" diventa arancione con testo "Riconnessione...". Retry automatico ogni 5 secondi. Dopo 30 secondi: "Connessione persa. Verifica la rete."
- Flood di alert (>20 in 1 minuto): raggruppamento automatico. "23 alert simili nell'ultimo minuto" con link per espandere.
- Nessun alert: empty state nel tab Alerts: "Nessun alert attivo. Il sistema funziona regolarmente."

---

### Journey 6: Developer che integra l'API

**Persona**: Developer — vuole codice funzionante, odia la documentazione vaga
**Goal**: Capire l'API, testare una chiamata, ottenere codice pronto all'uso
**Entry point**: `/developer/playground`

#### Steps

**Step 1 — Scegliere l'endpoint**
- Apre `/developer/playground`
- **Layout split-screen**: sinistra = request builder, destra = code snippet (gia implementato)
- **Endpoint selector**: dropdown con tutti gli endpoint disponibili, raggruppati per categoria:
  - Explain: `POST /api/v1/explain`, `GET /api/v1/explain/{id}`, `GET /api/v1/explain`
  - Graph: `GET /api/v1/explain/{id}/graph`
  - Narrative: `GET /api/v1/explain/{id}/narrative`, `POST /api/v1/explain/{id}/narrative-llm`
  - What-If: `POST /api/v1/explain/{id}/whatif`
  - AI: `POST /api/v1/explain/{id}/ask`, `POST /api/v1/explain/{id}/summary`
- L'utente seleziona `POST /api/v1/explain`

**Step 2 — Configurare la request**
- **Body editor**: textarea con syntax highlighting (JSON)
- Pre-compilato con un example payload completo e commentato
- **Path params**: se l'endpoint ha `{id}`, appare un campo per inserirlo. Autocomplete con gli ID delle explanation recenti.
- **Query params**: campi dinamici in base all'endpoint selezionato
- **Headers**: campo per API key (auto-compilato se configurata nelle Settings)

**Step 3 — Eseguire la request**
- Click "Invia" → bottone mostra spinner
- **Response panel** (sotto il body editor):
  - Status code con colore semantico (200 verde, 400 giallo, 500 rosso)
  - Response time in ms
  - Body della risposta con syntax highlighting e folding JSON
  - Headers della risposta (collapsati di default)

**Step 4 — Copiare il codice**
- Panel di destra: code snippet auto-generato
- **Linguaggio selector**: tabs "Python" | "JavaScript" | "Go" | "cURL"
- Il codice include:
  - Import necessari
  - Setup della connessione con API key
  - Chiamata all'endpoint selezionato
  - Parsing della risposta
  - Commenti esplicativi
- Bottone "Copia" → copia negli appunti → toast "Codice copiato!"

**Step 5 — Esplorare la documentazione**
- Link "Documentazione completa" → `/developer/docs`
- Pagina con API reference interattiva (embed Swagger/OpenAPI o documentazione custom):
  - Ogni endpoint con descrizione, parametri, response schema
  - Esempi di request/response
  - Error codes e significato
  - Rate limits e best practices

**Success state**: Il developer ha testato l'API, capito la struttura di request/response, e ha codice pronto da incollare nel suo progetto.

**Edge cases**:
- JSON non valido nel body: bordo rosso sul textarea + messaggio "JSON non valido: Unexpected token at line 5"
- API key mancante: warning "Aggiungi la tua API key nelle Settings per autenticare le richieste"
- Endpoint 500: mostra il body dell'errore + suggerimento "Verifica i parametri della richiesta"

---

## 4. SCREEN-BY-SCREEN SPECS

### Screen 1: Dashboard (Home) — `/dashboard`

**Purpose**: Dare all'utente una vista a colpo d'occhio dell'attivita recente e fornire accesso rapido alle azioni principali.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Breadcrumb: Dashboard]                             │
│                                                      │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────┐ │
│  │ Explain  │  │ Conf.   │  │ Alert   │  │ Ultima │ │
│  │ oggi:342 │  │ media:  │  │ attivi: │  │ expl.: │ │
│  │ +12%     │  │ 78%     │  │ 3       │  │ 2m fa  │ │
│  └─────────┘  └─────────┘  └─────────┘  └────────┘ │
│                                                      │
│  ┌──────────────────────┐  ┌───────────────────────┐ │
│  │  Quick Actions       │  │  Attivita recente     │ │
│  │                      │  │                       │ │
│  │  [+ Nuova Expl.]     │  │  - portfolio_risk...  │ │
│  │  [Apri Archivio]     │  │  - market_regime...   │ │
│  │  [Vai al Monitoring] │  │  - signal_strength... │ │
│  │                      │  │  - credit_score...    │ │
│  └──────────────────────┘  │  - vol_forecast...    │ │
│                            │                       │ │
│                            │  [Vedi tutto ->]      │ │
│                            └───────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Creare una nuova explanation (CTA prominente in Quick Actions).
**Secondary actions**: Navigare alle explanation recenti, controllare gli alert, aprire l'archivio.

**Data shown**:
- 4 StatsCards: explanation create oggi (con trend), confidence media (con trend), alert attivi (con badge), ultima explanation (con timestamp relativo)
- Quick Actions: 3 bottoni per le azioni piu frequenti
- Attivita recente: ultime 10 explanation, con target name, valore, confidence badge, timestamp relativo. Click → dettaglio.

**Empty state** (primo accesso, nessuna explanation):
- Le StatsCards mostrano "0" con testo "Nessun dato ancora"
- Quick Actions resta visibile
- Al posto di "Attivita recente": la Welcome Card di onboarding (Journey 1)

**Loading state**: 4 skeleton card in alto (altezza fissa 80px), skeleton lista a destra (5 righe).

**Error state**: Se le stats API falliscono: le card mostrano "—" con icona warning. Toast: "Impossibile caricare le statistiche. Riprova."

**Responsive**:
- Desktop (>= 1024px): 4 stat card in riga, quick actions + recenti affiancati (2 colonne)
- Tablet (768-1023px): 4 stat card in 2x2 grid, quick actions + recenti impilati
- Mobile (< 768px): stat card in colonna singola (scrollabile orizzontalmente), tutto impilato

---

### Screen 2: Crea Explanation — `/explanations/new`

**Purpose**: Permettere all'utente di inserire i componenti di una decisione e generare un'explanation.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Breadcrumb: Explanations > Nuova]                  │
│                                                      │
│  Crea una nuova Explanation                          │
│  Inserisci i componenti per generare un'analisi      │
│  completa.                    [Usa esempio]          │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  Target name          Value                    │  │
│  │  [________________]   [________]               │  │
│  │                                                │  │
│  │  Componenti                         [+ Aggiungi│  │
│  │  ┌────────────────────────────────────────┐    │  │
│  │  │ Nome          Valore  Peso   Conf.  [x]│    │  │
│  │  │ [__________]  [____]  [____] [____] [x]│    │  │
│  │  │ [__________]  [____]  [____] [____] [x]│    │  │
│  │  │ [__________]  [____]  [____] [____] [x]│    │  │
│  │  └────────────────────────────────────────┘    │  │
│  │                                                │  │
│  │  Opzioni avanzate (collapsato)                 │  │
│  │  ▸ Include grafo, include drivers, max depth   │  │
│  │                                                │  │
│  │  [    Analizza    ]   [Annulla]                │  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Compilare il form e cliccare "Analizza".
**Secondary actions**: Usare il template di esempio, configurare opzioni avanzate, annullare.

**Data shown**: Form fields (target, value, components array), opzioni avanzate (include_graph, include_drivers, max_drivers, max_depth, missing_threshold).

**Opzioni avanzate** (accordion, chiuso di default):
- Include grafo: toggle (default: on)
- Include driver ranking: toggle (default: on)
- Max drivers: number input (default: 5)
- Max depth: number input (default: 10)
- Missing threshold: number input (default: 0.1)

**Empty state**: N/A (la pagina e sempre un form).

**Loading state**: Bottone "Analizza" diventa disabled con spinner + testo "Analisi in corso..."

**Error state**:
- Validazione client-side: bordo rosso sui campi invalidi + messaggio sotto
- Errore server: alert banner sopra il bottone submit con messaggio dell'errore e suggerimento

**Responsive**:
- Desktop: form centrato, max-width 720px
- Tablet: identico, con margini ridotti
- Mobile: campi componente impilati verticalmente (nome occupa tutta la riga, value/weight/confidence in riga sotto)

---

### Screen 3: Explanation Detail — Overview Tab

**Purpose**: Dare una vista sintetica e completa dell'explanation, mostrando valore finale, breakdown, driver e confidence in un'unica schermata.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Breadcrumb: Explanations > portfolio_risk_score]   │
│                                                      │
│  ┌─ Tab Bar ──────────────────────────────────────┐  │
│  │ [Overview]  Grafo  What-If  AI Insights  Export│  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  SUMMARY CARD                                  │  │
│  │  portfolio_risk_score                          │  │
│  │                                                │  │
│  │  0.72           [81% Conf.]    [12% Missing]   │  │
│  │  ───────────────────────────────────────────   │  │
│  │  Top driver: trend_strength (42%)              │  │
│  │  Creato: 26 Mar 2026, 14:32 — v1.0.0          │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌────────────────────┐  ┌───────────────────────┐   │
│  │  BREAKDOWN CHART   │  │  DRIVER RANKING       │   │
│  │  (bar chart oriz.) │  │  1. trend_str.  42%   │   │
│  │                    │  │  2. volatility  35%   │   │
│  │  ████████ 42%      │  │  3. momentum    23%   │   │
│  │  ██████ 35%        │  │                       │   │
│  │  ████ 23%          │  │                       │   │
│  └────────────────────┘  └───────────────────────┘   │
│                                                      │
│  ┌────────────────────┐  ┌───────────────────────┐   │
│  │  CONFIDENCE PANEL  │  │  NARRATIVE VIEWER     │   │
│  │  Overall: 81%      │  │  "The portfolio risk  │   │
│  │  Per-node:         │  │   score of 0.72 is    │   │
│  │  - trend: 90%      │  │   primarily driven    │   │
│  │  - vol: 85%        │  │   by..."              │   │
│  │  - mom: 95%        │  │                       │   │
│  │  Missing: 12%      │  │  [Genera narrativa]   │   │
│  └────────────────────┘  └───────────────────────┘   │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  SENSITIVITY QUICK VIEW (se driver disponibili)│  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Leggere e comprendere l'explanation.
**Secondary actions**: Navigare ai tab specializzati, esportare, copiare link.

**Data shown**: SummaryCard (target, final_value, confidence, missing_impact, top_drivers, metadata), BreakdownChart (breakdown array), DriverRanking (top_drivers), ConfidencePanel (confidence, per_node, missing_impact), NarrativeViewer (narrativa template), SensitivityQuickView (preview della sensitivity).

**Tab Bar behavior**:
- Tab attivo: bordo inferiore 2px `primary`, testo `foreground`, font-weight semibold
- Tab inattivo: testo `muted-foreground`, hover `foreground`
- Transizione: nessuna animazione tra tab (cambio istantaneo — la velocita e piu importante dell'estetica)
- URL: il tab attivo si riflette nel query param `?tab=overview` per permettere deep-linking e condivisione

**Confidence badge — colori semantici**:
- >= 80%: `bg-emerald-500/10 text-emerald-700 border-emerald-200` (verde)
- 50-79%: `bg-amber-500/10 text-amber-700 border-amber-200` (giallo)
- < 50%: `bg-red-500/10 text-red-700 border-red-200` (rosso)

**Empty state**: N/A (se l'explanation esiste, ha sempre dei dati).

**Loading state**: Skeleton che replica esattamente la struttura: un rettangolo per il SummaryCard, due rettangoli affiancati per breakdown+drivers, due per confidence+narrative.

**Error state**: Se il fetch dell'explanation fallisce: pagina centrata con "Explanation non trovata" + campo di ricerca + link "Torna all'archivio".

**Responsive**:
- Desktop: layout come descritto (2 colonne per breakdown/drivers e confidence/narrative)
- Tablet: identico ma con meno padding
- Mobile: tutte le sezioni impilate. Tab bar scrollabile orizzontalmente.

---

### Screen 4: Explanation Detail — Grafo Tab

**Purpose**: Visualizzare le dipendenze causali tra nodi in un DAG interattivo, per capire come il valore finale e stato derivato.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Tab Bar: Overview  [Grafo]  What-If  ...]          │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │                                                │  │
│  │              INTERACTIVE DAG                   │  │
│  │                                                │  │
│  │         [root_node]                            │  │
│  │        /     |      \                          │  │
│  │  [child1] [child2] [child3]                    │  │
│  │     |        |                                 │  │
│  │  [leaf1]  [leaf2]                              │  │
│  │                                                │  │
│  │                                                │  │
│  │  ┌──── Controls ─────┐                         │  │
│  │  │ [Zoom+] [Zoom-]   │                         │  │
│  │  │ [Fit] [Fullscreen] │                         │  │
│  │  │ [Reset]            │                         │  │
│  │  └────────────────────┘                         │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌─ Node Detail Panel (se nodo selezionato) ──────┐  │
│  │  Nome: trend_strength                          │  │
│  │  Valore: 0.80  |  Confidence: 90%             │  │
│  │  Tipo: component  |  Peso: 0.4                │  │
│  │  [Simula variazione]  [Mostra in Breakdown]   │  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Esplorare visivamente le dipendenze tra nodi.
**Secondary actions**: Selezionare un nodo per dettagli, navigare al What-If per simulare, esportare il grafo come immagine.

**Data shown**: Graph nodes (id, label, value, confidence, node_type), Graph edges (source, target, weight, transformation_type).

**Interazioni del grafo**:
- **Pan**: drag sullo sfondo
- **Zoom**: mouse wheel / pinch gesture / bottoni +/-
- **Select node**: click → highlight nodo e edges connessi, dimmerizza il resto (opacity 0.3). Mostra il Node Detail Panel.
- **Double-click node**: espande/collassa i figli
- **Hover node**: tooltip leggero con label e valore
- **Hover edge**: mostra peso e tipo di trasformazione
- **Fit view**: adatta lo zoom per vedere tutti i nodi
- **Fullscreen**: espande il grafo a tutto schermo (overlay)

**Node styling**:
- Forma: rettangolo arrotondato (border-radius 8px)
- Colore bordo: basato sulla confidence (semantica: verde/giallo/rosso)
- Colore sfondo: bianco (light) / card (dark)
- Nodo root: bordo doppio, sfondo primary/5%
- Nodo con `missing: true`: bordo tratteggiato, icona warning

**Edge styling**:
- Spessore: proporzionale al peso (1px-4px)
- Colore: `muted-foreground` (default), `primary` (se connesso al nodo selezionato)
- Frecce: piccole (6px), direzionali

**Empty state**: Se `explanation.graph` e null: "Grafo non disponibile per questa explanation. Ricreala con l'opzione 'Include grafo' attiva."

**Loading state**: Skeleton rettangolo grande con shimmer + testo "Caricamento grafo..."

**Error state**: "Errore nel rendering del grafo. I dati potrebbero essere corrotti." + link "Mostra dati grezzi JSON"

**Responsive**:
- Desktop/Tablet: grafo occupa tutta la larghezza, altezza 60vh. Node Detail Panel sotto.
- Mobile: grafo occupa tutta la larghezza, altezza 50vh. Panel sotto. Messaggio: "Per un'esperienza migliore, ruota il dispositivo."

---

### Screen 5: Explanation Detail — What-If Tab

**Purpose**: Permettere all'utente di simulare variazioni nei componenti e vedere l'impatto sul risultato finale.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Tab Bar: Overview  Grafo  [What-If]  ...]          │
│                                                      │
│  ┌───────────────────────┐  ┌──────────────────────┐ │
│  │  SIMULATORE           │  │  RISULTATO           │ │
│  │                       │  │                      │ │
│  │  trend_strength       │  │  Originale: 0.72     │ │
│  │  [====|======] 0.80   │  │  Simulato:  0.65     │ │
│  │                       │  │  Delta:    -0.07     │ │
│  │  volatility_index     │  │  Variazione: -9.7%   │ │
│  │  [====|======] 0.60   │  │                      │ │
│  │                       │  │  ┌────────────────┐  │ │
│  │  momentum_signal      │  │  │ COMPARISON BAR │  │ │
│  │  [====|======] 0.75   │  │  │ (prima/dopo)   │  │ │
│  │                       │  │  └────────────────┘  │ │
│  │  [Reset]  [Applica]   │  │                      │ │
│  └───────────────────────┘  └──────────────────────┘ │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  SCENARIO MANAGER                              │  │
│  │  [Scenario 1: Base] [Scenario 2: Stress] [+]  │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  DIFF TABLE                                    │  │
│  │  Component    Orig.   Mod.   Delta   %         │  │
│  │  trend_str.   0.80    0.80   0.00    0%        │  │
│  │  volatility   0.60    0.40  -0.20  -33%        │  │
│  │  momentum     0.75    0.75   0.00    0%        │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  SENSITIVITY RANKING                           │  │
│  │  1. volatility_index  — impatto: 0.35          │  │
│  │  2. trend_strength    — impatto: 0.28          │  │
│  │  3. momentum_signal   — impatto: 0.18          │  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Modificare i valori dei componenti tramite slider e vedere come cambia il risultato.
**Secondary actions**: Salvare scenari, confrontare scenari, esportare la simulazione.

**Data shown**: ComponentSlider per ogni componente, risultato simulato (original, modified, delta, delta_percentage), DiffTable (component_diffs), SensitivityRanking.

**Slider behavior**:
- Range: 0 — valore_max * 2 (o 1 se il valore e tra 0 e 1)
- Step: 0.01
- Valore corrente mostrato a destra dello slider
- Colore track: `muted` per la parte inattiva, `primary` per la parte attiva
- Se il valore e stato modificato: label in **bold** + indicatore di delta accanto

**Scenario Manager** (gia implementato):
- Tab orizzontali per ogni scenario salvato
- Bottone "+" per creare nuovo scenario
- Ogni scenario salva la configurazione degli slider
- "Confronta scenari" → mostra ComparisonView con overlay dei risultati

**Calcolo**:
- Click "Applica" → `POST /api/v1/explain/{id}/whatif` con le modifiche
- Risultato aggiorna in tempo reale il pannello destro e la DiffTable
- **Ottimizzazione**: debounce 500ms sullo slider per non fare troppe chiamate API

**Empty state**: N/A (il simulatore mostra sempre i valori originali come punto di partenza).

**Loading state**: Pannello risultato mostra skeleton per valore simulato + spinner discreto.

**Error state**: "Simulazione fallita. Verifica che i valori siano nel range valido." + dettaglio errore in expandable.

**Responsive**:
- Desktop: simulatore e risultato affiancati, DiffTable e Ranking sotto full-width
- Mobile: tutto impilato. Slider occupano tutta la larghezza.

---

### Screen 6: Explanation Detail — AI Insights Tab

**Purpose**: Fornire spiegazioni in linguaggio naturale, permettere domande conversazionali, e generare summary strutturati per diverse audience.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Tab Bar: Overview  Grafo  What-If  [AI Insights]]  │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  NARRATIVA AI                                  │  │
│  │                                                │  │
│  │  Livello: [Basic] [Advanced] [Executive]       │  │
│  │  Lingua:  [EN] [IT]                            │  │
│  │                                                │  │
│  │  "The portfolio risk score of 0.72 indicates   │  │
│  │   a moderate risk level. The primary driver    │  │
│  │   is trend_strength, contributing 42% of..."   │  │
│  │                                                │  │
│  │  Fonte: LLM (GPT-4) | Template               │  │
│  │  [Rigenera]  [Copia]                          │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  Q&A CHAT                                      │  │
│  │                                                │  │
│  │  ┌── Chat History ──────────────────────────┐  │  │
│  │  │  Tu: Perche il valore e cosi alto?       │  │  │
│  │  │  AI: Il valore di 0.72 e influenzato...  │  │  │
│  │  │  Tu: Cosa succederebbe senza volatility? │  │  │
│  │  │  AI: Rimuovendo il contributo di...      │  │  │
│  │  └──────────────────────────────────────────┘  │  │
│  │                                                │  │
│  │  [Digita una domanda...              ] [Invia] │  │
│  │                                                │  │
│  │  Suggerimenti:                                 │  │
│  │  "Perche la confidence e bassa?"               │  │
│  │  "Quali dati mancano?"                         │  │
│  │  "Spiega il driver principale"                 │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │  EXECUTIVE SUMMARY GENERATOR                   │  │
│  │                                                │  │
│  │  Audience: [Board] [Tecnico] [Cliente]         │  │
│  │  Lingua:   [EN] [IT]                           │  │
│  │  [Genera Summary]                              │  │
│  │                                                │  │
│  │  ┌── Risultato (se generato) ──────────────┐  │  │
│  │  │  Titolo: Analisi del Risk Score          │  │  │
│  │  │  Summary: ...                            │  │  │
│  │  │  Key Findings: ...                       │  │  │
│  │  │  Rischi: ...                             │  │  │
│  │  │  Raccomandazioni: ...                    │  │  │
│  │  │  [Copia tutto] [Esporta PDF]            │  │  │
│  │  └──────────────────────────────────────────┘  │  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Ottenere una spiegazione comprensibile della explanation.
**Secondary actions**: Fare domande specifiche via chat, generare summary per audience diverse, copiare/esportare.

**Data shown**: Narrativa (LLM o template), chat history, executive summary (title, summary, key_findings, risks, recommendations).

**Narrativa AI**:
- **Livello**: toggle group (Basic / Advanced / Executive). Cambia la profondita della spiegazione.
  - Basic: 2-3 frasi, linguaggio semplice
  - Advanced: 1-2 paragrafi, include dati tecnici
  - Executive: 1 paragrafo, focus su impatto e azione
- **Lingua**: toggle EN/IT
- **Fonte**: badge che indica se la narrativa e da LLM o da template
- **Auto-load**: la narrativa template si carica automaticamente. La versione LLM richiede click su "Genera con AI".

**Q&A Chat**:
- **Input**: campo di testo con placeholder "Fai una domanda su questa explanation..."
- **Submit**: Enter o click su "Invia"
- **Suggested questions**: 3 chip cliccabili sotto l'input (cambiano in base al contenuto dell'explanation)
- **Chat history**: scrollabile, max 20 messaggi visibili. Messaggi precedenti caricati on-scroll.
- **Risposta AI**: typing indicator (tre puntini animati) durante la generazione. Testo appare tutto insieme quando completo.
- **Errore**: se l'AI non risponde: "Non sono riuscito a generare una risposta. Riprova." in un bubble di errore.

**Executive Summary Generator**:
- **Audience**: radio group (Board, Tecnico, Cliente). Ogni opzione ha una descrizione:
  - Board: "Riepilogo per decisori — focus su rischio e raccomandazioni"
  - Tecnico: "Analisi dettagliata — include metriche e breakdown"
  - Cliente: "Spiegazione accessibile — linguaggio non tecnico"
- **Lingua**: toggle EN/IT
- **Generazione**: 3-8 secondi. Skeleton con shimmer.
- **Risultato**: card strutturata con sezioni separabili. Ogni sezione ha icona "Copia" individuale.
- **Rigenera**: possibilita di aggiungere istruzioni custom prima di rigenerare.

**Empty state**: Narrativa section mostra "Clicca 'Genera' per ottenere una spiegazione in linguaggio naturale." Chat mostra i suggested questions.

**Loading state**: Narrativa: skeleton testo (3 righe). Chat: typing indicator. Summary: skeleton card.

**Error state**: Fallback alla narrativa template con banner "AI non disponibile — narrativa generata con template."

**Responsive**:
- Desktop: tutte le sezioni impilate verticalmente, full-width
- Mobile: identico ma chat input diventa sticky bottom

---

### Screen 7: Archivio Explanations (ex-Audit) — `/explanations`

**Purpose**: Permettere di cercare, filtrare e navigare tutte le explanation generate, con focus su auditabilita e tracciabilita.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Breadcrumb: Explanations]                          │
│                                                      │
│  Archivio Explanations              [+ Nuova] [Exp.] │
│  47 risultati                                        │
│                                                      │
│  ┌── Filtri ──────────────────────────────────────┐  │
│  │  [Ricerca: ___________]                        │  │
│  │  Da: [__/__/____]  A: [__/__/____]            │  │
│  │  Target: [____________]                        │  │
│  │  Confidence: [====|=======|====] 0% - 100%    │  │
│  │                                                │  │
│  │  Quick: [Oggi] [Settimana] [Mese] [Alta conf] │  │
│  │  [Applica filtri]  [Reset]                     │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌── Tabella ─────────────────────────────────────┐  │
│  │  [ ] ID       Target          Value  Conf  Data│  │
│  │  ─────────────────────────────────────────────  │  │
│  │  [ ] abc12..  portfolio_risk  0.72   81%   26/3│  │
│  │  [ ] def34..  market_regime   0.45   62%   25/3│  │
│  │  [ ] ghi56..  signal_str      0.88   93%   25/3│  │
│  │  [ ] jkl78..  credit_score    0.31   28%   24/3│  │
│  │  [ ] mno90..  vol_forecast    0.67   75%   24/3│  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  [< Precedente]  Pagina 1 di 3  [Successivo >]      │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Trovare una explanation specifica e navigare al suo dettaglio.
**Secondary actions**: Filtrare per data/confidence/target, esportare in CSV, selezionare piu righe per azioni batch.

**Data shown**: Lista paginata di explanation con ID (troncato a 8 char), target, final_value, confidence (con badge colore), data creazione.

**Tabella — comportamento colonne**:
- **Checkbox**: seleziona riga per azioni batch (export, delete)
- **ID**: cliccabile, copia ID completo negli appunti. Tooltip con ID intero.
- **Target**: cliccabile, naviga al dettaglio. Testo troncato con ellipsis se troppo lungo.
- **Value**: formattato a 2 decimali
- **Confidence**: badge con colore semantico (verde/giallo/rosso) + percentuale
- **Data**: formato relativo ("2 ore fa") con tooltip che mostra data completa ("26 Mar 2026, 14:32:15 UTC")
- **Ordinamento**: click sull'header per ordinare. Icona freccia indica direzione.

**Filtri — pannello collapsabile**:
- Default: espanso la prima volta, poi ricorda la preferenza dell'utente
- I filtri si applicano automaticamente con debounce 300ms (per la ricerca testuale)
- I filtri quick (chips) aggiornano i campi del form e applicano immediatamente
- I filtri attivi mostrano un contatore: "3 filtri attivi" + bottone "Reset" per pulire tutto
- I filtri si riflettono nell'URL come query params per permettere condivisione e bookmark

**Empty state**: Illustrazione minimale (grafico vuoto) + "Nessuna explanation trovata." + due opzioni: "Rimuovi i filtri" | "Crea la tua prima explanation"

**Loading state**: AuditSkeleton (gia implementato): skeleton per filtri + 5 righe skeleton.

**Error state**: "Errore nel caricamento dell'archivio. Riprova." + retry button.

**Responsive**:
- Desktop: tabella completa con tutte le colonne
- Tablet: nasconde la colonna ID, mostra le altre
- Mobile: tabella diventa lista di card. Ogni card mostra: target (bold), valore + confidence sulla stessa riga, data sotto. Click → dettaglio.

---

### Screen 8: Monitoring — `/monitoring`

**Purpose**: Sorvegliare le decisioni in tempo reale, ricevere alert su anomalie, e avere una visione statistica dell'attivita.

**Layout**: (tre tab: Live Feed, Alerts, Statistiche — struttura gia descritta in Journey 5)

**Tab Live Feed**:
```
┌──────────────────────────────────────────────────────┐
│  [Tab Bar: [Live Feed]  Alerts (3)  Statistiche]     │
│                                                      │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐              │
│  │ Oggi    │  │ Conf.   │  │ Alert   │              │
│  │ 342     │  │ media   │  │ attivi  │              │
│  │ +12%    │  │ 78%     │  │ 3       │              │
│  └─────────┘  └─────────┘  └─────────┘              │
│                                                      │
│  ┌── Live Feed ──── [● Live] [Pausa] ─────────────┐ │
│  │  14:32:15  portfolio_risk   0.72  [81%]         │ │
│  │  14:31:48  market_regime    0.45  [62%]    ⚠    │ │
│  │  14:31:22  signal_strength  0.88  [93%]         │ │
│  │  14:30:55  credit_score     0.31  [28%]    ⚠    │ │
│  │  ...                                            │ │
│  └─────────────────────────────────────────────────┘ │
│                                                      │
│  ┌── Alert Panel ─────────────────────────────────┐  │
│  │  ⚠ Low confidence (0.28) — credit_score        │  │
│  │  ⚠ Missing data (35%) — vol_forecast           │  │
│  │  ⚠ Low confidence (0.31) — market_regime       │  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Tab Alerts**:
- Lista completa di alert con: timestamp, tipo, target, valore, status
- Azioni per riga: Acknowledge, Apri dettaglio, Risolvi
- Filtri: per tipo, per status, per date range

**Tab Statistiche**:
- Grafici temporali: explanation per ora/giorno, confidence media nel tempo, distribuzione confidence
- KPI: totale explanation, media confidence, % con alert, tempo medio di risposta API

**Primary action** (Live Feed): Monitorare il flusso e reagire agli alert.
**Primary action** (Alerts): Gestire e risolvere gli alert.
**Primary action** (Stats): Analizzare trend e pattern.

**Empty state**:
- Live Feed: "In attesa di nuove explanation... Il feed si aggiornera automaticamente."
- Alerts: "Nessun alert attivo. Tutto regolare."
- Stats: "Dati insufficienti per generare statistiche. Le statistiche saranno disponibili dopo le prime 10 explanation."

**Loading state**: StatsCards skeleton + feed skeleton (5 righe).

**Error state**: "Errore nella connessione al feed live." + retry + opzione "Passa a modalita manuale" (refresh esplicito).

**Responsive**:
- Desktop: StatsCards in riga, feed (2/3 larghezza) + alert panel (1/3) affiancati
- Mobile: tutto impilato, alert panel sopra il feed

---

### Screen 9: API Playground — `/developer/playground`

**Purpose**: Permettere agli sviluppatori di testare gli endpoint API in modo interattivo e ottenere codice pronto all'uso.

**Layout**: (gia descritto in Journey 6 — split screen request/code)

**Primary action**: Selezionare un endpoint, configurare la request, inviarla.
**Secondary actions**: Copiare il codice generato, cambiare linguaggio.

**Data shown**: Endpoint selector, request body editor, path/query params, response viewer, code snippet.

**Miglioramenti rispetto all'attuale**:
- **Raggruppamento endpoint**: nel dropdown, endpoint raggruppati per categoria con separatori
- **Response history**: le ultime 5 response restano accessibili (tabs sotto il response viewer)
- **Timing**: tempo di risposta mostrato accanto allo status code
- **JSON folding**: possibilita di collassare/espandere sezioni del JSON nella response
- **Autocomplete ID**: nel path param `{id}`, dropdown con le ultime explanation create

**Empty state**: Response panel: "Seleziona un endpoint e clicca 'Invia' per vedere la risposta."

**Loading state**: Response panel: spinner + "Richiesta in corso..." con timer che conta i secondi.

**Error state**: Response panel mostra status code rosso + body dell'errore formattato.

**Responsive**:
- Desktop: split screen orizzontale (50/50)
- Tablet: split screen con code snippet sotto (stacked)
- Mobile: tutto impilato. Code snippet collapsabile.

---

### Screen 10: Settings — `/settings`

**Purpose**: Configurare le preferenze dell'utente, gestire le API key, e personalizzare l'esperienza.

**Layout**:
```
┌──────────────────────────────────────────────────────┐
│  [Breadcrumb: Settings]                              │
│                                                      │
│  Settings                                            │
│                                                      │
│  ┌── Preferenze ──────────────────────────────────┐  │
│  │  Tema:           [Chiaro] [Scuro] [Sistema]    │  │
│  │  Lingua default: [EN ▾]                        │  │
│  │  Formato numeri: [0.00] [0,00]                 │  │
│  │  Auto-refresh monitoring: [Toggle ON]  10s     │  │
│  │  Notifiche browser: [Toggle OFF]               │  │
│  │  [Salva preferenze]                            │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌── API Keys ────────────────────────────────────┐  │
│  │  Le API key servono per accesso programmatico  │  │
│  │                                                │  │
│  │  Key 1: ee_sk_***************3f2a              │  │
│  │  Creata: 15 Mar 2026  Ultimo uso: Oggi         │  │
│  │  [Copia] [Revoca]                              │  │
│  │                                                │  │
│  │  Key 2: ee_sk_***************8b1c              │  │
│  │  Creata: 01 Feb 2026  Ultimo uso: Mai          │  │
│  │  [Copia] [Revoca]                              │  │
│  │                                                │  │
│  │  [+ Genera nuova API Key]                      │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌── Account ─────────────────────────────────────┐  │
│  │  Email: alessio@example.com                    │  │
│  │  Ruolo: Admin                                  │  │
│  │  Ultimo accesso: 26 Mar 2026, 14:32            │  │
│  │                                                │  │
│  │  [Rivedi tour onboarding]                      │  │
│  │  [Esporta tutti i miei dati]                   │  │
│  │  [Logout]                                      │  │
│  └────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────┘
```

**Primary action**: Modificare le preferenze.
**Secondary actions**: Gestire API key, logout.

**Data shown**: Preferenze utente, lista API key (mascherate), info account.

**Comportamento API Key**:
- **Genera**: click → mostra la key COMPLETA una sola volta in un modal. "Copia e salva in un posto sicuro. Non potrai rivederla."
- **Revoca**: confirmation dialog: "Sei sicuro? Le applicazioni che usano questa key smetteranno di funzionare."
- **Maschera**: mostra solo i primi 5 e ultimi 4 caratteri.

**Salvataggio preferenze**: auto-save con debounce 1s. Toast "Preferenze salvate" quando completato.

**Empty state** (API Keys): "Nessuna API key. Genera la tua prima key per accedere all'API." + bottone.

**Loading state**: Skeleton per le sezioni.

**Error state**: "Impossibile salvare le preferenze. Riprova." accanto al campo che ha fallito.

**Responsive**:
- Desktop: max-width 640px centrato (gia implementato)
- Mobile: full-width con padding standard

---

## 5. MICRO-INTERACTIONS & POLISH

### Transizioni

| Elemento | Animazione | Durata | Easing |
|----------|-----------|--------|--------|
| Cambio pagina (sidebar nav) | Nessuna — cambio istantaneo | 0ms | — |
| Cambio tab (detail page) | Nessuna — cambio istantaneo | 0ms | — |
| Apertura pannello laterale | Slide-in da destra | 200ms | ease-out |
| Chiusura pannello laterale | Slide-out verso destra | 150ms | ease-in |
| Apertura accordion/collapsable | Altezza da 0 a auto | 200ms | ease-out |
| Hover su riga tabella | Background change | 100ms | linear |
| Nuovo elemento nel live feed | Slide-down + highlight flash | 300ms + 2000ms fade | ease-out |
| Modal apertura | Scale 0.95->1 + opacity 0->1 | 200ms | ease-out |
| Modal chiusura | Opacity 1->0 | 150ms | ease-in |

**Principio**: le transizioni servono a orientare, non a intrattenere. Zero animazioni tra pagine. Zero bounce effects.

### Feedback

**Bottoni**:
- Default: sfondo `primary`, testo `primary-foreground`
- Hover: sfondo `primary/90%`
- Active (pressed): sfondo `primary/80%`, scale 0.98
- Disabled: sfondo `muted`, testo `muted-foreground`, cursor not-allowed
- Loading: spinner a sinistra dell'icona, testo "In corso..."

**Form validation**:
- **Inline**: validazione al blur (non al keystroke per non irritare). Bordo rosso + messaggio sotto il campo in `text-destructive`.
- **Submit**: se ci sono errori, scroll al primo campo invalido + focus
- **Success**: campo torna al bordo normale, nessun messaggio di successo per campo (solo toast generale)

**Copy to clipboard**:
- Click icona copia → icona diventa checkmark per 2 secondi → torna a icona copia
- Toast: "Copiato negli appunti"

### Notifiche e toast

**Sistema toast** (angolo in basso a destra, stacked):
- **Success**: bordo sinistro verde, icona checkmark. Durata: 4 secondi.
- **Error**: bordo sinistro rosso, icona X. Durata: 8 secondi (l'utente ha bisogno di piu tempo per leggere errori).
- **Info**: bordo sinistro blu, icona info. Durata: 4 secondi.
- **Warning**: bordo sinistro arancione, icona warning. Durata: 6 secondi.
- Max 3 toast visibili contemporaneamente. I successivi entrano in coda.
- Ogni toast ha una X per dismissione manuale.
- Hover su un toast: pausa il timer di auto-dismiss.

**Messaggi toast standard**:
- "Explanation creata con successo" (success)
- "PDF scaricato: {filename}" (success)
- "Copiato negli appunti" (info)
- "Preferenze salvate" (success)
- "API key revocata" (warning)
- "Errore di connessione. Riprova." (error)
- "Simulazione completata" (success)
- "Summary generato" (success)

### Empty states

Ogni empty state ha tre elementi:
1. **Illustrazione**: icona grande (48px) in `muted-foreground/30%` — non un'illustrazione custom, solo un'icona Lucide ingrandita
2. **Messaggio**: testo descrittivo in `muted-foreground`
3. **Azione**: bottone o link per risolvere lo stato vuoto

Esempi:
- Dashboard senza dati: icona BarChart3 + "Nessuna attivita ancora" + "Crea la tua prima Explanation"
- Archivio filtrato senza risultati: icona Search + "Nessun risultato per questi filtri" + "Rimuovi filtri"
- Chat senza messaggi: icona MessageSquare + "Fai una domanda sull'explanation" + chip suggerimenti

### Loading

**Principio**: skeleton screens ovunque, mai spinner a pagina intera.

- **Pagina**: la struttura (sidebar, header, breadcrumb) e sempre visibile. Solo il contenuto ha skeleton.
- **Card**: rettangolo con border-radius identico alla card reale, sfondo `muted/50%`, shimmer animation.
- **Tabella**: header reale + righe skeleton (altezza identica alle righe reali).
- **Testo**: linee rettangolari con larghezze variabili (70%, 100%, 85%) per simulare testo.
- **Chart**: rettangolo con proporzioni del chart reale.

**Shimmer animation**: gradiente lineare che si muove da sinistra a destra, durata 1.5s, infinite loop.

### Tooltip

- **Trigger**: hover (desktop) / tap (mobile)
- **Delay**: 300ms prima di apparire (evita tooltip accidentali)
- **Posizione**: sopra l'elemento (default), sotto se non c'e spazio
- **Stile**: sfondo `popover`, bordo `border`, ombra sm, border-radius 6px, padding 8px 12px
- **Max width**: 280px
- **Freccia**: triangolo 6px che punta all'elemento trigger

**Dove mettere tooltip**:
- Ogni metrica nella SummaryCard (confidence, missing impact, valore)
- Header delle colonne nella tabella audit
- Icone nella toolbar del grafo
- Parametri nel form di creazione
- Badge di stato (alert types)

### Keyboard shortcuts

| Shortcut | Azione | Contesto |
|----------|--------|----------|
| `Cmd/Ctrl + K` | Apri ricerca globale | Ovunque |
| `Escape` | Chiudi modal/panel/overlay | Quando aperto |
| `Cmd/Ctrl + N` | Nuova explanation | Ovunque |
| `Cmd/Ctrl + E` | Esporta (PDF/CSV) | Detail/Archivio |
| `1-5` | Cambia tab | Detail page |
| `?` | Mostra shortcuts | Ovunque |
| `J / K` | Naviga su/giu nella lista | Archivio/Feed |
| `Enter` | Apri elemento selezionato | Archivio/Feed |

**Ricerca globale** (Cmd+K):
- Modal centrato, input in alto
- Risultati divisi per tipo: Explanation, Pagine, Azioni
- Navigazione con frecce su/giu, Enter per selezionare
- Mostra le ultime 3 ricerche come "Recenti"

### Breadcrumbs

- **Posizione**: sopra il titolo della pagina, sotto l'header
- **Formato**: "Explanations > portfolio_risk_score" (con separator ">")
- **Comportamento**: ogni segmento e cliccabile tranne l'ultimo (pagina corrente)
- **Overflow**: se il breadcrumb e troppo lungo, tronca i segmenti centrali con "..."
- **Stile**: testo `muted-foreground` per i segmenti, `foreground` per l'ultimo

---

## 6. METRICS THAT MATTER

### Dashboard (Home)

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Bounce rate** | % utenti che lasciano la dashboard senza azione | < 20% |
| **Click-through su Quick Actions** | % utenti che usano i quick action buttons | > 60% |
| **Time to first action** | Tempo tra apertura e primo click significativo | < 10s |

### Crea Explanation

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Completion rate** | % di form iniziati vs. form completati | > 80% |
| **Time to submit** | Tempo dall'apertura del form al submit | < 60s |
| **Uso template esempio** | % di utenti che usano "Usa esempio" | Tracciare (alto per nuovi utenti) |
| **Errori di validazione** | Numero medio di errori per submission | < 1 |

### Explanation Detail

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Time to first meaningful insight** | Tempo tra apertura e primo scroll/click significativo | < 15s |
| **Tab navigation rate** | % di utenti che visitano > 1 tab | > 40% |
| **Tab distribution** | Quale tab e piu visitato dopo Overview | Tracciare |
| **Grafo interaction rate** | % di utenti che interagiscono col grafo (zoom, click nodo) | > 50% (tra chi visita il tab) |
| **What-If usage rate** | % di utenti che modificano almeno 1 slider | > 30% (tra chi visita il tab) |
| **AI features adoption** | % di utenti che usano almeno 1 feature AI (narrativa/chat/summary) | > 25% |
| **AI Q&A depth** | Numero medio di domande per sessione | > 2 |
| **Export conversion** | % di visite al dettaglio che risultano in un export (PDF/CSV) | > 10% |

### Archivio (ex-Audit)

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Search-to-click rate** | % di ricerche che portano a un click su una riga | > 70% |
| **Filter usage** | % di sessioni che usano almeno 1 filtro | > 50% |
| **Avg. filtri per sessione** | Numero medio di filtri applicati | 2-3 |
| **Export rate** | % di sessioni archivio che risultano in export CSV | > 15% |
| **Pagine navigate** | Numero medio di pagine nella tabella per sessione | 1-3 |

### Monitoring

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Tempo medio sulla pagina** | Quanto a lungo resta aperto il monitoring | > 5 min |
| **Alert response time** | Tempo tra alert e acknowledge | < 2 min |
| **Alert resolution rate** | % di alert risolti vs. ignorati | > 80% |
| **Feed pause usage** | % di sessioni che usano il bottone Pausa | Tracciare |

### API Playground

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Requests per session** | Numero medio di richieste inviate | > 3 |
| **Code copy rate** | % di sessioni che copiano codice | > 60% |
| **Language distribution** | Quale linguaggio di code snippet e piu copiato | Tracciare |
| **Success rate** | % di richieste che ritornano 2xx | > 90% |

### Cross-platform

| Metrica | Descrizione | Target |
|---------|-------------|--------|
| **Sessioni per utente/settimana** | Frequenza di ritorno | > 3 |
| **Feature discovery** | % di utenti che hanno usato almeno una volta ogni sezione principale | > 50% entro 30 giorni |
| **Onboarding completion** | % di nuovi utenti che completano il primo explanation | > 70% |
| **NPS (Net Promoter Score)** | Survey periodica | > 40 |

### Implementazione tracking

Usare un event bus leggero (custom hook `useTrack`):

```typescript
// Ogni evento ha: name, properties, timestamp (auto)
useTrack('explanation_created', { components_count: 3, used_template: false });
useTrack('tab_switched', { from: 'overview', to: 'graph', explanation_id: 'abc' });
useTrack('export_completed', { format: 'pdf', sections: ['summary', 'breakdown'] });
useTrack('ai_question_asked', { question_length: 45, session_messages: 3 });
useTrack('alert_acknowledged', { alert_type: 'low_confidence', response_time_s: 42 });
```

Backend: gli eventi possono essere inviati a qualsiasi analytics provider (PostHog, Mixpanel, custom). L'interfaccia e agnostica rispetto al provider.

---

## APPENDICE A: Design Tokens

Per mantenere consistenza, questi sono i design token di riferimento (basati sul tema Tailwind/shadcn gia in uso):

| Token | Uso | Valore (light) |
|-------|-----|-----------------|
| `--primary` | CTA, link attivi, bordi focus | Blue 600 |
| `--destructive` | Errori, alert critici | Red 600 |
| `--muted` | Sfondi secondari, disabilitati | Gray 100 |
| `--muted-foreground` | Testo secondario, placeholder | Gray 500 |
| `--accent` | Hover states, selezioni | Gray 100 |
| `--card` | Sfondo card | White |
| `--border` | Bordi, separatori | Gray 200 |
| `--confidence-high` | Badge confidence >= 80% | Emerald 500 |
| `--confidence-medium` | Badge confidence 50-79% | Amber 500 |
| `--confidence-low` | Badge confidence < 50% | Red 500 |

## APPENDICE B: Priorita di implementazione

| Fase | Scope | Effort stimato |
|------|-------|---------------|
| **Fase 1** | Nuova IA (routing, sidebar, breadcrumb) + Dashboard home + Explanation tabs | 2 sprint |
| **Fase 2** | Archivio migliorato (filtri, quick filter, responsive cards) + Onboarding | 1 sprint |
| **Fase 3** | Monitoring tabs + alert management + keyboard shortcuts | 1 sprint |
| **Fase 4** | Ricerca globale (Cmd+K) + micro-interactions polish + analytics tracking | 1 sprint |

---

*Documento redatto il 26 Marzo 2026. Da revisionare con il team design prima dell'implementazione.*
