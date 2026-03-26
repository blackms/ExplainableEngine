# Explainable Engine Dashboard -- User Stories

> **Document version**: 1.0
> **Last updated**: 2026-03-26
> **Status**: Draft for Sprint 8-12 planning

---

## Personas

| ID | Persona | Context | Primary need |
|----|---------|---------|--------------|
| P1 | **Quant / Trader** | Uses AIP daily, needs to trust model outputs before acting | Fast, deep decomposition of any score |
| P2 | **Analyst / Sales** | Communicates model results to clients and internal stakeholders | Plain-language explanations, exportable evidence |
| P3 | **Risk / Compliance Officer** | Must prove decisions are auditable and explainable to regulators | Complete audit trail, search, export |
| P4 | **Developer / Integrator** | Builds systems that consume ExplainableEngine API | API exploration, code generation, fast onboarding |
| P5 | **Operations / Desk Lead** | Monitors model health and explanation quality in real time | Live feed, anomaly alerts, summary KPIs |

---

## Use Cases

| ID | Use Case | Primary Persona |
|----|----------|-----------------|
| UC-01 | Understand why a model produced a specific score | P1, P2 |
| UC-02 | Visually trace causal dependencies between components | P1 |
| UC-03 | Assess trustworthiness of an explanation | P1, P3 |
| UC-04 | Communicate an explanation to a non-technical audience | P2 |
| UC-05 | Stress-test a score by modifying inputs | P1 |
| UC-06 | Compare multiple what-if scenarios for a decision | P1, P2 |
| UC-07 | Search and audit historical explanations | P3 |
| UC-08 | Export evidence for regulatory reporting | P3 |
| UC-09 | Monitor explanation quality in real time | P5 |
| UC-10 | Prototype API integrations quickly | P4 |
| UC-11 | Manage access and credentials | P3, P4 |

---

## EPIC A -- Explanation Explorer

**Priority**: Must-have
**Sprint**: 8-9
**Goal**: Deliver the core value proposition -- any user can open an explanation and immediately understand what happened, why, and how trustworthy it is.

---

### US-101: View explanation summary

**Epic**: Explanation Explorer
**Use Case**: UC-01
**Priority**: Must-have

**As a** Quant or Analyst
**I want** to see the target name, final value, overall confidence, and top drivers at a glance when I open an explanation
**So that** I can decide in under 5 seconds whether this score requires further investigation or is safe to act on.

**Business Rationale**:
Every user interaction with the dashboard starts here. If the summary is unclear or slow, users revert to spreadsheets and manual analysis -- eliminating the platform's value proposition. A clear summary directly reduces the average investigation time per score from ~15 minutes (manual) to seconds, which at scale across a trading desk of 20 people translates to measurable productivity gains. It is also the first screen a prospective client sees during a sales demo.

**Acceptance Criteria**:
- [ ] Summary panel displays: target name, final_value (formatted to relevant precision), confidence score, computation_type, and created_at timestamp
- [ ] Top 3 drivers are listed with name, impact score, and rank, pulled from `top_drivers`
- [ ] Confidence is color-coded: green >= 0.8, yellow >= 0.5, red < 0.5
- [ ] Missing impact is displayed with a warning icon when > 0 and color escalation: yellow when missing_impact > 0.10, red when > 0.20
- [ ] Deterministic hash is shown (truncated, expandable) so users can verify reproducibility
- [ ] Page loads and renders summary in under 500ms on a 4G connection
- [ ] All numeric values respect locale formatting (dot vs comma for decimals)

**UI/UX Notes**:
Card-based layout. The summary card sits at the top of the explanation detail page and is always visible without scrolling. Use large typography for final_value and confidence. Top drivers should appear as a compact ranked list with inline spark-bars showing relative impact. The deterministic hash should be in a small monospace font with a copy icon.

**API Dependencies**:
- `GET /api/v1/explain/{id}` -- fields: target, final_value, confidence, top_drivers, missing_impact, metadata.deterministic_hash, metadata.created_at, metadata.computation_type

---

### US-102: Interactive breakdown chart

**Epic**: Explanation Explorer
**Use Case**: UC-01
**Priority**: Must-have

**As a** Quant
**I want** to see a visual chart showing each component's percentage contribution to the final value, and drill into sub-components by clicking
**So that** I can identify exactly which factor is driving a score up or down without manually computing weight * value for every component.

**Business Rationale**:
Component decomposition is the core analytical action. Today, quants export JSON and build ad-hoc spreadsheets to understand contribution breakdowns -- a process that takes 10-20 minutes per explanation and is error-prone. A visual, drill-down chart eliminates this entirely. For compliance, having a standard visual decomposition format also creates a consistent audit artifact across all teams.

**Acceptance Criteria**:
- [ ] Horizontal bar chart displays each component from the `breakdown` array
- [ ] Each bar shows: component name, contribution_pct (as bar width), and absolute contribution value as a label
- [ ] Bars are sorted by contribution_pct descending by default, with option to sort alphabetically
- [ ] Clicking a bar that has `sub_breakdown` data drills down, replacing the chart with the sub-component breakdown and showing a breadcrumb trail for navigation back
- [ ] Hovering over a bar shows a tooltip with: name, value, weight, contribution, contribution_pct, confidence
- [ ] Bar color encodes confidence: same green/yellow/red scale as US-101
- [ ] A "back" button or breadcrumb allows returning to the parent level at any depth
- [ ] Chart is responsive and readable on viewports down to 1024px wide
- [ ] An alternative treemap view is available via a toggle for users who prefer area-based visualization

**UI/UX Notes**:
Default view is horizontal bar chart (familiar to financial users). A toggle in the top-right corner switches to treemap. Breadcrumb trail appears above the chart when drilled into sub-components (e.g., "market_regime_score > trend_strength > ..."). Transitions between levels should animate smoothly (300ms).

**API Dependencies**:
- `GET /api/v1/explain/{id}` -- fields: breakdown (including nested sub_breakdown)

---

### US-103: Dependency graph visualization

**Epic**: Explanation Explorer
**Use Case**: UC-02
**Priority**: Must-have

**As a** Quant
**I want** to see an interactive directed acyclic graph (DAG) showing how components feed into the target, with nodes colored by type and edges weighted visually
**So that** I can trace the full causal chain of a score and quickly identify which upstream inputs are responsible for an unexpected result.

**Business Rationale**:
When a score behaves unexpectedly, quants need to trace causality. Without a graph, this means reading raw JSON and mentally reconstructing the tree -- a process that took 30-60 minutes per incident in pre-platform workflows. During production incidents (e.g., a model outputs an anomalous trading signal), every minute of delay translates to potential financial exposure. The graph visualization reduces root-cause identification to under 2 minutes.

**Acceptance Criteria**:
- [ ] Graph is rendered as an interactive DAG using the data from the graph endpoint
- [ ] Node colors distinguish type: target nodes (blue), component nodes (green), sub_component nodes (gray)
- [ ] Edge thickness is proportional to edge weight
- [ ] Edge labels show the weight value (e.g., "w=0.4")
- [ ] Clicking a node opens a detail panel showing: id, label, type, value, and (if component) its contribution_pct and confidence
- [ ] Zoom (scroll wheel or pinch), pan (drag), and fit-to-screen (button) are supported
- [ ] User can toggle between three layout algorithms: hierarchical (top-down), force-directed, and left-to-right
- [ ] Graph supports export to PNG and SVG via a download button
- [ ] For graphs with > 50 nodes, automatic clustering groups sub_components under their parent, expandable on click
- [ ] Graph format selector allows switching between JSON, DOT, and Mermaid raw output views (for developers)

**UI/UX Notes**:
The graph occupies the full width of the content area and has a minimum height of 400px. A floating toolbar in the top-left provides: zoom in, zoom out, fit-to-screen, layout toggle, export. Node hover shows a subtle glow effect. The detail panel slides in from the right when a node is clicked, similar to a property inspector.

**API Dependencies**:
- `GET /api/v1/explain/{id}` -- fields: graph (nodes, edges)
- `GET /api/v1/explain/{id}/graph?format=json` -- full graph data
- `GET /api/v1/explain/{id}/graph?format=dot` -- for export
- `GET /api/v1/explain/{id}/graph?format=mermaid` -- for export

---

### US-104: Confidence indicator panel

**Epic**: Explanation Explorer
**Use Case**: UC-03
**Priority**: Must-have

**As a** Quant or Compliance Officer
**I want** to see an overall confidence gauge and a per-component confidence breakdown, with clear warnings when missing data degrades reliability
**So that** I can immediately assess whether an explanation is trustworthy enough to act on or whether I need to flag it for manual review.

**Business Rationale**:
"Can I trust this number?" is the single most frequent question asked by every persona. In regulated financial services, acting on a low-confidence automated decision without documentation is a compliance violation. This panel provides the trust signal that determines whether a user acts on the score, escalates it, or rejects it. Without it, users either over-trust (risk) or under-trust (missed opportunity) every output. Regulators (MiFID II, AIFMD) increasingly require documented confidence assessments for algorithmic decisions.

**Acceptance Criteria**:
- [ ] A radial gauge or semi-circular meter displays overall confidence (0.0 to 1.0) with color zones: green (>= 0.8), yellow (>= 0.5), red (< 0.5)
- [ ] Below the gauge, a ranked list shows per-component confidence: component name, confidence value, and a small inline bar
- [ ] Components with confidence < 0.5 are highlighted in red with a warning icon
- [ ] A "Missing Data" section is visible when missing_impact > 0, showing: missing_impact value as percentage, severity label (Minor < 10%, Moderate 10-20%, Severe > 20%), and a brief explanation of what missing data means
- [ ] When missing_impact > 0.20, a prominent red banner appears at the top of the entire explanation page: "High missing data impact -- this explanation may not be reliable"
- [ ] Clicking on any component in the confidence list navigates to that component in the breakdown chart (US-102)
- [ ] Confidence values show two decimal places

**UI/UX Notes**:
Panel sits to the right of the breakdown chart (desktop) or below it (mobile). The gauge should be visually prominent -- it is the trust signal. Use a semi-circular gauge with a needle animation on load. The missing data banner is a full-width alert that cannot be dismissed (but can be collapsed).

**API Dependencies**:
- `GET /api/v1/explain/{id}` -- fields: confidence, breakdown[].confidence, missing_impact

---

### US-105: Narrative explanation viewer

**Epic**: Explanation Explorer
**Use Case**: UC-04
**Priority**: Must-have

**As an** Analyst or Sales person
**I want** to read a plain-language narrative of the explanation, toggle between basic and advanced detail levels, and switch between English and Italian
**So that** I can communicate model results to clients and stakeholders without translating numbers into words myself, saving 15-30 minutes per client communication.

**Business Rationale**:
The sales team sends explanation summaries to clients daily. Before this feature, an analyst manually writes a 2-3 paragraph summary for each explanation -- a process that takes 15-30 minutes and introduces human error and inconsistency. Automated narratives ensure consistency, reduce preparation time, and allow Sales to service more clients per day. Italian language support is required because the primary client base operates in Italian financial markets. The copy-to-clipboard feature supports rapid insertion into emails and reports.

**Acceptance Criteria**:
- [ ] Narrative text is displayed in a readable panel with proper typography (16px+ body, generous line-height)
- [ ] A toggle switch allows selecting "Basic" (non-technical summary) or "Advanced" (includes technical details, specific values, confidence analysis)
- [ ] A language selector provides "EN" and "IT" options
- [ ] Changing level or language fetches the new narrative without a full page reload (smooth transition)
- [ ] A "Copy to clipboard" button copies the narrative as plain text
- [ ] A "Copy as Markdown" button copies the narrative with Markdown formatting preserved
- [ ] Loading state shows a skeleton text placeholder while the narrative endpoint responds
- [ ] If the narrative endpoint returns an error, a fallback message is shown: "Narrative unavailable. Please view the data breakdown above."
- [ ] Narrative panel is printable (renders cleanly in print CSS)

**UI/UX Notes**:
Full-width panel below the graph. The toggle (Basic/Advanced) and language selector (EN/IT) appear as pill-shaped buttons in the panel header. Copy buttons are in the top-right corner of the panel. Text should feel like reading a report paragraph, not a data table.

**API Dependencies**:
- `GET /api/v1/explain/{id}/narrative?level=basic&lang=en`
- `GET /api/v1/explain/{id}/narrative?level=advanced&lang=en`
- `GET /api/v1/explain/{id}/narrative?level=basic&lang=it`
- `GET /api/v1/explain/{id}/narrative?level=advanced&lang=it`

---

## EPIC B -- What-if Simulator

**Priority**: Must-have
**Sprint**: 9-10
**Goal**: Enable users to modify component values and immediately see how the target score would change, supporting stress testing and decision preparation.

---

### US-201: Modify input values with sliders

**Epic**: What-if Simulator
**Use Case**: UC-05
**Priority**: Must-have

**As a** Quant or Portfolio Manager
**I want** to adjust any component's value using a slider and immediately see the predicted impact on the target score
**So that** I can stress-test a model output against hypothetical market conditions (e.g., "what if volatility doubles?") without rebuilding the full model, enabling faster risk assessment during volatile markets.

**Business Rationale**:
Stress testing is a daily activity on trading desks and a quarterly regulatory requirement. Currently, running a what-if scenario requires modifying input data, re-running the model, and comparing outputs -- a cycle that takes 30+ minutes per scenario. Sliders reduce this to real-time interaction (< 2 seconds per scenario). During market volatility events, the ability to rapidly assess "what if X changes by Y%" is the difference between a proactive hedge and a reactive loss.

**Acceptance Criteria**:
- [ ] Each component in the breakdown is listed with a horizontal slider
- [ ] Slider range is current value +/- 50% by default, with the option to type a custom value outside this range
- [ ] As the user drags a slider, a real-time delta indicator shows: new value, absolute change, and percentage change from original
- [ ] Sliders for weight and confidence are also available (collapsible "advanced" section)
- [ ] A "Reset All" button restores all sliders to original values
- [ ] A "Submit What-If" button sends the modified values to the what-if endpoint and displays the new result
- [ ] Modified components are visually highlighted (blue border or background) so the user can see at a glance what changed
- [ ] Debouncing ensures the delta preview updates at most every 200ms during drag (to avoid flicker)
- [ ] At least 10 components can be displayed without layout issues

**UI/UX Notes**:
Sliders appear in a left panel alongside the breakdown chart. Each slider row shows: component name, current value (muted), slider, new value (bold if changed). The delta indicator appears inline next to the new value as a colored pill: green for positive delta, red for negative. The "Submit What-If" button is sticky at the bottom of the panel.

**API Dependencies**:
- `POST /api/v1/explain/{id}/what-if` -- request body contains modified component values
- `GET /api/v1/explain/{id}` -- to load original values for slider initialization

---

### US-202: Side-by-side comparison

**Epic**: What-if Simulator
**Use Case**: UC-05
**Priority**: Must-have

**As a** Quant
**I want** to see the original explanation and the what-if result side by side, with deltas highlighted
**So that** I can make an informed decision by visually comparing the before and after states, rather than trying to remember numbers from one view while looking at another.

**Business Rationale**:
Decision-making requires comparison. Showing the original and modified states side by side eliminates the cognitive overhead of switching between views and reduces errors in manual comparison. This directly supports the decision workflow: adjust, compare, decide. Portfolio managers use this in meetings to justify allocation changes -- "here is the current state, here is what happens if we adjust."

**Acceptance Criteria**:
- [ ] After a what-if submission, the screen splits into two columns: "Original" (left) and "What-if" (right)
- [ ] Both columns show: final_value, confidence, and breakdown chart
- [ ] A "Delta" column or overlay highlights the differences: green for improvement (value moves in favorable direction), red for degradation
- [ ] Delta values show both absolute and percentage change
- [ ] The breakdown chart bars are overlaid or mirrored to visually show which components shifted
- [ ] The overall confidence change is displayed prominently at the top of the comparison
- [ ] A "Back to Editor" button returns to the slider view (US-201) with the current modifications preserved
- [ ] Comparison view is printable for meeting handouts

**UI/UX Notes**:
Two-column layout on desktop (minimum 1280px viewport). On narrower screens, stack vertically with original on top and what-if below. Delta badges appear next to each value. Consider a butterfly chart (back-to-back bars) for the breakdown comparison to make differences immediately scannable.

**API Dependencies**:
- `POST /api/v1/explain/{id}/what-if` -- response contains full explanation structure for the modified scenario
- `GET /api/v1/explain/{id}` -- original explanation for left column

---

### US-203: Sensitivity ranking

**Epic**: What-if Simulator
**Use Case**: UC-05
**Priority**: Must-have

**As a** Quant
**I want** to see a ranked list of which input components have the most influence on the output, quantified by sensitivity
**So that** I can focus my attention on the variables that matter most, rather than wasting time investigating components that have negligible effect.

**Business Rationale**:
Attention is the scarcest resource on a trading desk. With explanations that may contain 10-50 components, knowing where to focus is worth significant time savings. Sensitivity analysis also feeds into risk management: the components with highest sensitivity are the components to hedge or monitor. This ranking turns the what-if simulator from a toy into a professional risk tool.

**Acceptance Criteria**:
- [ ] A ranked list shows all components ordered by sensitivity (highest impact first)
- [ ] Each row displays: rank number, component name, sensitivity score, and a horizontal bar proportional to sensitivity
- [ ] Sensitivity is computed by running what-if with +/-10% perturbation on each component individually and measuring the output delta
- [ ] A horizontal bar chart visualizes the sensitivity distribution across all components
- [ ] Clicking a component in the ranking auto-scrolls to that component's slider in the what-if editor (US-201)
- [ ] The ranking updates after each what-if submission to reflect the new baseline
- [ ] A "Top N" filter allows viewing only the top 3, 5, or 10 most sensitive components

**UI/UX Notes**:
This appears as a collapsible section below the what-if sliders or as a tab alongside the slider panel. The bar chart uses a single color (blue) with opacity gradient to distinguish magnitude. The top 3 components have a subtle highlight or badge ("High Sensitivity").

**API Dependencies**:
- `POST /api/v1/explain/{id}/what-if` -- multiple calls with individual component perturbations to compute sensitivity (may be batched client-side or require a dedicated sensitivity endpoint in a future API version)

---

### US-204: Save and compare scenarios

**Epic**: What-if Simulator
**Use Case**: UC-06
**Priority**: Should-have

**As a** Portfolio Manager
**I want** to save a what-if scenario with a custom name and compare up to 3 saved scenarios side by side
**So that** I can prepare for board meetings and investment committee reviews by presenting "here are the 3 scenarios we evaluated and why we chose option B."

**Business Rationale**:
Investment decisions are rarely made in real time. The typical workflow is: explore multiple scenarios over several hours, save the interesting ones, then present the top 3 in a meeting. Without save/compare, users screenshot results and paste into PowerPoint -- an error-prone process that takes 30+ minutes per meeting. Saved scenarios also create a decision audit trail: "we considered these alternatives before choosing this course of action," which is valuable for compliance and post-mortem analysis.

**Acceptance Criteria**:
- [ ] After a what-if result is displayed, a "Save Scenario" button allows the user to enter a custom name (max 100 characters) and save
- [ ] Saved scenarios are stored in browser localStorage (MVP) with the explanation ID, scenario name, modified component values, and result summary
- [ ] A "Saved Scenarios" panel lists all scenarios for the current explanation, showing: name, timestamp, final_value delta, confidence delta
- [ ] User can select 2 or 3 scenarios and click "Compare" to open a multi-column comparison view
- [ ] Comparison view shows: scenario name as column header, final_value, confidence, and per-component breakdown in each column
- [ ] A "Delete" button on each saved scenario allows removal
- [ ] Maximum of 10 saved scenarios per explanation (localStorage constraint)
- [ ] Export comparison as CSV (one row per component, one column per scenario)

**UI/UX Notes**:
The "Saved Scenarios" panel appears as a sidebar or a drawer accessible via a button. The comparison view extends the two-column layout from US-202 to three columns. If only 2 scenarios are selected, the third column shows the original. Color-coded deltas in each column relative to the original.

**API Dependencies**:
- `POST /api/v1/explain/{id}/what-if` -- to regenerate full results when reopening a saved scenario
- No backend storage required for MVP (localStorage)

---

## EPIC C -- Audit Trail

**Priority**: Should-have
**Sprint**: 10-11
**Goal**: Provide compliance officers and analysts with a complete, searchable, exportable log of all explanations produced by the system.

---

### US-301: List all explanations

**Epic**: Audit Trail
**Use Case**: UC-07
**Priority**: Should-have

**As a** Compliance Officer
**I want** to see a paginated table of all explanations the system has produced, sorted by date
**So that** I can verify that every automated decision has a corresponding explanation on record, satisfying regulatory audit requirements.

**Business Rationale**:
Financial regulations (MiFID II Article 17, AIFMD) require firms using algorithmic decision-making to maintain logs of decisions and their rationale. This table is the primary interface for demonstrating compliance during audits. Without it, the compliance team would need to query the API directly or request database exports from engineering -- a dependency that adds days to audit preparation and creates operational risk.

**Acceptance Criteria**:
- [ ] A table displays explanations with columns: target, final_value, confidence, computation_type, created_at, deterministic_hash (truncated)
- [ ] Table supports server-side pagination using cursor-based pagination (20 items per page default, configurable: 20, 50, 100)
- [ ] Each column header is clickable for sorting (ascending/descending toggle)
- [ ] Default sort is created_at descending (most recent first)
- [ ] Confidence column cells are color-coded (green/yellow/red per the same scale)
- [ ] Rows with missing_impact > 0.10 display a warning icon
- [ ] Clicking a row navigates to the full explanation detail view (US-303)
- [ ] The total count of explanations is displayed at the top of the table (approximate is acceptable)
- [ ] Table renders within 1 second for the first page

**UI/UX Notes**:
Standard data table with sticky header. Alternating row colors for readability. Pagination controls at the bottom with page size selector. The table occupies the full content width. A subtle row hover highlights the clickable row.

**API Dependencies**:
- `GET /api/v1/explain` (list endpoint -- requires backend implementation of a list/search endpoint with cursor pagination; currently not in the API contract)
- This story requires a **new backend endpoint**: `GET /api/v1/explain?cursor=X&limit=N&sort=field:asc|desc`

---

### US-302: Search and filter explanations

**Epic**: Audit Trail
**Use Case**: UC-07
**Priority**: Should-have

**As a** Compliance Officer
**I want** to filter the explanations list by date range, target name, confidence range, and missing data presence
**So that** when an auditor asks "show me all low-confidence decisions in Q1," I can produce the answer in under 30 seconds instead of requesting a custom database query.

**Business Rationale**:
Audit requests are unpredictable in scope. "Show me all decisions with confidence below 0.6 in the past 90 days" is a real question regulators ask. The ability to answer this instantly -- rather than filing a Jira ticket for engineering -- reduces audit response time from days to seconds. This is also the mechanism by which internal risk teams proactively monitor decision quality: daily review of low-confidence explanations.

**Acceptance Criteria**:
- [ ] A filter bar above the table provides: date range picker (from/to), target name text search (substring match, case-insensitive), confidence range slider (min/max, 0.0-1.0), missing data toggle (Any / With Missing Data / Without Missing Data)
- [ ] Filters are applied server-side (not client-side filtering) for performance on large datasets
- [ ] Active filters are shown as removable chips/pills above the table
- [ ] A "Clear All Filters" button resets to default view
- [ ] URL query parameters reflect the current filter state (enabling bookmarkable/shareable filtered views)
- [ ] Filter changes update the table without full page reload
- [ ] Full-text search on target names returns results within 500ms
- [ ] Empty state shows a clear "No explanations match your filters" message with a suggestion to broaden criteria

**UI/UX Notes**:
Horizontal filter bar with collapsible advanced options. The date range picker uses a calendar widget. The confidence slider is a dual-thumb range slider. Filter chips appear between the filter bar and the table. Consider a "Save Filter" feature in a future iteration for compliance officers who run the same queries repeatedly.

**API Dependencies**:
- `GET /api/v1/explain?cursor=X&limit=N&target=X&confidence_min=X&confidence_max=X&from=X&to=X&has_missing_data=bool` (requires backend implementation of query parameters on the list endpoint)

---

### US-303: Explanation detail view

**Epic**: Audit Trail
**Use Case**: UC-07
**Priority**: Should-have

**As a** Compliance Officer
**I want** to click on any explanation in the audit list and see the full detail view with all panels (summary, breakdown, graph, narrative)
**So that** I can drill down from a high-level audit list into complete evidence for any individual decision without switching tools.

**Business Rationale**:
Audit workflows require both breadth (scanning the list) and depth (examining individual decisions). This story connects the audit list (US-301) to the full explanation view (Epic A), creating a complete audit workflow within a single tool. Without this, compliance officers would need to copy explanation IDs and navigate manually -- a workflow that breaks concentration and increases the chance of reviewing the wrong record.

**Acceptance Criteria**:
- [ ] Clicking a row in the explanation list (US-301) navigates to a detail page at `/explanations/{id}`
- [ ] The detail page includes all panels from Epic A: summary (US-101), breakdown chart (US-102), dependency graph (US-103), confidence panel (US-104), narrative viewer (US-105)
- [ ] A "Back to List" button returns to the audit list with filters and scroll position preserved
- [ ] The browser URL updates to include the explanation ID (deep-linkable)
- [ ] A breadcrumb shows: Explanations > {target name} for navigation context
- [ ] Page loads all panels within 2 seconds (parallel API calls)

**UI/UX Notes**:
Same layout as the standalone explanation view from Epic A, but with the addition of a breadcrumb and "Back to List" navigation. All panels load in parallel with independent loading states (skeleton screens). The page should feel like opening a detailed record from a master list.

**API Dependencies**:
- `GET /api/v1/explain/{id}` -- full explanation data
- `GET /api/v1/explain/{id}/narrative?level=basic&lang=en` -- default narrative
- `GET /api/v1/explain/{id}/graph?format=json` -- graph data

---

### US-304: Export to PDF and CSV

**Epic**: Audit Trail
**Use Case**: UC-08
**Priority**: Should-have

**As a** Compliance Officer
**I want** to export a single explanation as a PDF (including graph, breakdown, and narrative) and export a filtered list as CSV
**So that** I can attach evidence to regulatory filings, share with external auditors who do not have dashboard access, and archive records in the firm's document management system.

**Business Rationale**:
Regulators accept PDF and CSV as evidence formats. External auditors do not have (and should not have) access to internal tools. Compliance teams routinely attach explanation evidence to regulatory filings (e.g., RTS 25 reporting, transaction reporting). The alternative is manual screenshot-and-paste into Word -- a process that takes 20+ minutes per explanation and produces inconsistent, unprofessional output. PDF export also supports the firm's legal obligation to maintain records for 5+ years in a portable format.

**Acceptance Criteria**:
- [ ] On the explanation detail page (US-303), an "Export PDF" button generates a PDF containing: explanation summary, breakdown table, dependency graph (as rendered image), confidence panel data, narrative text (both basic and advanced if available), metadata (ID, hash, timestamp, computation_type)
- [ ] PDF includes a header with: firm logo placeholder, "Explainable Engine Report", date generated, explanation ID
- [ ] PDF is generated client-side (no server dependency) using the data already loaded
- [ ] On the explanation list page (US-301), an "Export CSV" button exports the currently filtered results
- [ ] CSV includes columns: id, target, final_value, confidence, missing_impact, computation_type, created_at, deterministic_hash
- [ ] CSV export handles up to 10,000 rows (paginating through the API if necessary)
- [ ] Export buttons show a loading indicator during generation
- [ ] File names follow the pattern: `explanation_{id}_{date}.pdf` and `explanations_export_{date}.csv`

**UI/UX Notes**:
Export buttons are in the top-right action bar of both the detail and list views. PDF generation should use a print-optimized layout (no interactive elements, static graph rendering). Consider a brief "Generating..." modal for CSV exports that fetch multiple pages.

**API Dependencies**:
- `GET /api/v1/explain/{id}` -- all data for PDF
- `GET /api/v1/explain/{id}/graph?format=json` -- for rendering graph in PDF
- `GET /api/v1/explain/{id}/narrative?level=basic&lang=en` and `?level=advanced&lang=en` -- for PDF narrative sections
- `GET /api/v1/explain?cursor=X&limit=100` -- paginated fetching for CSV export

---

## EPIC D -- Live Monitoring

**Priority**: Should-have
**Sprint**: 11
**Goal**: Give operations teams real-time visibility into the stream of explanations being produced, with alerts for anomalous results.

---

### US-401: Recent explanations feed

**Epic**: Live Monitoring
**Use Case**: UC-09
**Priority**: Should-have

**As a** Desk Lead
**I want** to see a live-updating feed of the most recent explanations as they are created
**So that** I can monitor model activity in real time and notice unusual patterns (e.g., sudden spike in low-confidence explanations) before they become incidents.

**Business Rationale**:
Trading desks operate in real time. Models may produce hundreds of explanations per hour during market hours. Without a live feed, the desk lead learns about problems only when a trader reports a bad signal -- at which point the damage is done. A live feed enables proactive monitoring: "I see 5 consecutive low-confidence explanations for volatility -- something may be wrong with the data feed." Early detection of model degradation prevents costly incorrect trading signals.

**Acceptance Criteria**:
- [ ] A live feed page displays the 50 most recent explanations in reverse chronological order
- [ ] The feed auto-updates via polling every 10 seconds (with a visible "Last updated: X seconds ago" indicator)
- [ ] Each feed item shows: target name, final_value, confidence (color-coded), missing_impact warning icon (if > 0.10), created_at (relative time, e.g., "2m ago")
- [ ] New items appear at the top with a brief highlight animation (fade-in)
- [ ] When new items arrive while the user is scrolled down, a "N new explanations" banner appears at the top (clicking scrolls to top)
- [ ] A "Pause Feed" toggle stops polling (useful during investigation)
- [ ] Clicking a feed item navigates to the explanation detail view (US-303)
- [ ] A connection status indicator shows whether polling is active, paused, or failing

**UI/UX Notes**:
Card list layout, similar to a news feed or log viewer. Each card is compact (single line or two lines). The page should feel like a terminal or monitoring dashboard -- information-dense, minimal decoration. Consider a dark theme option for this page (common in trading desk tools).

**API Dependencies**:
- `GET /api/v1/explain?limit=50&sort=created_at:desc` (list endpoint with sort)
- Polling interval: 10 seconds. Future enhancement: WebSocket/SSE for push-based updates.

---

### US-402: Anomaly alerts

**Epic**: Live Monitoring
**Use Case**: UC-09
**Priority**: Should-have

**As a** Desk Lead
**I want** to receive visual alerts when an explanation has confidence below 0.5 or missing data impact above 20%
**So that** I can immediately investigate potentially unreliable model outputs before they influence trading decisions, reducing the firm's exposure to model risk.

**Business Rationale**:
A model producing a low-confidence signal that goes undetected can result in a bad trade. In a worst case, this means direct financial loss. Alerting on confidence and missing data thresholds is the minimum viable monitoring required by internal model risk policies. Most firms' model risk frameworks (SR 11-7) require documented monitoring of model output quality. This feature satisfies that requirement and provides the evidence trail.

**Acceptance Criteria**:
- [ ] Anomaly conditions are defined as: confidence < 0.5, missing_impact > 0.20 (thresholds configurable in Settings)
- [ ] When a feed item (US-401) matches an anomaly condition, it is highlighted with a red border and a warning badge
- [ ] A notification badge appears on the monitoring page tab/nav item showing the count of unacknowledged anomalies
- [ ] An optional browser notification (with user permission) fires for each anomaly, showing: target name and the anomaly reason
- [ ] An optional audio alert (subtle chime) plays for anomalies (toggle in Settings)
- [ ] An "Anomalies Only" filter on the feed page shows only flagged explanations
- [ ] User can "Acknowledge" an anomaly (marks it as reviewed, removes it from the badge count)
- [ ] Anomaly acknowledgment is persisted in localStorage (MVP) with timestamp and user action

**UI/UX Notes**:
Anomaly items in the feed have a red-tinted background and a pulsing dot icon. The notification badge on the nav item uses the standard red circle with count convention. Browser notifications use the Notifications API with a fallback message if permissions are denied. The "Anomalies Only" filter is a toggle pill above the feed.

**API Dependencies**:
- Same as US-401 (list endpoint with polling)
- Anomaly detection logic is client-side based on confidence and missing_impact thresholds from the explanation data

---

### US-403: Summary statistics

**Epic**: Live Monitoring
**Use Case**: UC-09
**Priority**: Should-have

**As a** Desk Lead or Operations Manager
**I want** to see summary KPI cards (total explanations today, average confidence, average response time) and a simple time-series chart of explanation volume
**So that** I can track operational health at a glance and detect trends (e.g., declining confidence over the day) that indicate systemic issues.

**Business Rationale**:
Individual explanations are tactical; aggregate statistics are strategic. A declining average confidence trend over the day may indicate data feed degradation that no single explanation would reveal. Volume spikes may indicate unexpected model behavior. These KPIs are standard operational metrics that the desk lead reports to management weekly. Without automated computation, this requires manual SQL queries against the database -- a task that takes 30 minutes and is performed by an engineer, not the person who needs the information.

**Acceptance Criteria**:
- [ ] Three KPI cards are displayed at the top of the monitoring page: "Explanations Today" (count), "Avg Confidence" (0.00 format, color-coded), "Anomalies Today" (count, red if > 0)
- [ ] A time-series line chart shows explanation volume per hour for the current day (24-hour view)
- [ ] A second overlay line on the chart shows average confidence per hour
- [ ] KPI cards update on the same polling interval as the feed (10 seconds)
- [ ] Hovering over a chart data point shows: hour, explanation count, avg confidence for that hour
- [ ] Chart Y-axes: left for count, right for confidence (0-1 scale)
- [ ] A date selector allows viewing statistics for previous days (not just today)
- [ ] Empty state for days with no data shows "No explanations recorded for this date"

**UI/UX Notes**:
KPI cards appear as a horizontal row above the feed (US-401). The time-series chart sits between the KPI cards and the feed. Cards use large typography for the number, small label below. The chart should be minimal and clean -- no gridlines, light axis labels. Use the same color coding for confidence on the chart line (segment color changes as confidence crosses thresholds).

**API Dependencies**:
- `GET /api/v1/explain?from=X&to=X&limit=100` -- paginated fetching for aggregation
- Note: A dedicated statistics/aggregation endpoint would significantly improve performance. Consider `GET /api/v1/stats?date=YYYY-MM-DD` as a future backend enhancement.

---

## EPIC E -- API Playground

**Priority**: Nice-to-have
**Sprint**: 12
**Goal**: Accelerate developer onboarding and integration by providing an interactive API exploration tool within the dashboard.

---

### US-501: Interactive API explorer

**Epic**: API Playground
**Use Case**: UC-10
**Priority**: Nice-to-have

**As a** Developer integrating ExplainableEngine into a downstream system
**I want** to build API requests using a form, send them, and see the response with syntax highlighting
**So that** I can understand the API behavior through hands-on experimentation without setting up curl, Postman, or writing throwaway code, reducing my integration onboarding time from days to hours.

**Business Rationale**:
Every new API consumer goes through an exploration phase: reading docs, trying requests, understanding error cases. This typically takes 1-2 days with traditional documentation. An interactive playground in the dashboard reduces this to 1-2 hours by providing immediate feedback. Faster integration means faster time-to-revenue for new clients that consume the API programmatically. It also reduces support burden on the engineering team (fewer "how do I use the API?" questions).

**Acceptance Criteria**:
- [ ] A form-based interface allows building requests for each endpoint: POST /explain, GET /explain/{id}, GET /graph, GET /narrative, POST /what-if
- [ ] Endpoint selection via dropdown populates the form with the correct parameters, headers, and body template
- [ ] The request body editor supports JSON with syntax highlighting and basic validation (brackets matching, required fields)
- [ ] A "Send" button executes the request and displays: HTTP status code, response headers, response body (JSON with syntax highlighting and collapsible sections)
- [ ] Request and response are displayed side by side (or top/bottom on narrow screens)
- [ ] Error responses are displayed with the error code and message highlighted
- [ ] Request history (last 20 requests) is stored in localStorage and accessible via a dropdown
- [ ] A "Load Example" button populates the form with a working example for the selected endpoint
- [ ] The base URL is configurable (to support localhost, staging, production environments)

**UI/UX Notes**:
Split-pane layout: left for request building, right for response viewing. The form should feel like Postman or Swagger UI but simpler. JSON editor should use a monospace font with line numbers. Response body should be collapsible at each object level. Use color-coded status badges (2xx green, 4xx yellow, 5xx red).

**API Dependencies**:
- All existing API endpoints (the playground calls them directly)
- CORS must be configured to allow requests from the dashboard origin

---

### US-502: Code snippet generator

**Epic**: API Playground
**Use Case**: UC-10
**Priority**: Nice-to-have

**As a** Developer
**I want** to generate ready-to-use code snippets (curl, Python, Go, JavaScript) for the current API request
**So that** I can copy working code directly into my integration codebase instead of translating API documentation into code manually, eliminating transcription errors and saving 30+ minutes per endpoint.

**Business Rationale**:
The last mile of API integration is writing the actual HTTP client code. Developers universally prefer copy-paste working code over writing from documentation. Code snippets in 4 languages cover the primary integration scenarios: curl for testing, Python for data science consumers, Go for backend services, JavaScript for frontend/Node consumers. Reducing integration friction directly accelerates client onboarding and reduces time-to-revenue.

**Acceptance Criteria**:
- [ ] After building a request in the API explorer (US-501), a "Code" tab shows generated code for the current request
- [ ] Language selector provides: curl, Python (requests library), Go (net/http), JavaScript (fetch API)
- [ ] Generated code includes: correct HTTP method, URL with path parameters filled, headers (Content-Type, X-Request-Id), request body (formatted), basic error handling
- [ ] Code is syntax-highlighted in the appropriate language
- [ ] A "Copy to Clipboard" button copies the code as plain text
- [ ] Generated code is functional (can be run as-is after setting the base URL)
- [ ] Python snippet includes a comment showing how to install the requests library
- [ ] Go snippet includes proper import statements
- [ ] Code updates automatically when the request form changes

**UI/UX Notes**:
The "Code" tab appears alongside the "Response" tab in the right pane of the API explorer. Language selector is a set of pill buttons (curl | Python | Go | JS). Copy button is in the top-right corner of the code block. Syntax highlighting should use a developer-friendly dark theme (e.g., Monokai or similar).

**API Dependencies**:
- None (code generation is purely client-side based on the current request configuration)

---

## EPIC F -- Settings and Authentication

**Priority**: Should-have
**Sprint**: 10
**Goal**: Secure the dashboard with authentication, provide API key management for programmatic access, and allow user personalization.

---

### US-601: User login

**Epic**: Settings and Authentication
**Use Case**: UC-11
**Priority**: Should-have

**As a** Compliance Officer
**I want** to log in with my Google account or email/password
**So that** access to explanation data is authenticated and attributed, creating an audit trail of who viewed which explanations and when, as required by the firm's information security policy.

**Business Rationale**:
Uncontrolled access to model explanation data is a compliance and security risk. Explanations may reveal proprietary trading strategies. An audit trail of access is required by SOC 2 controls, GDPR (data access logging), and internal information barrier policies. Authentication is also the prerequisite for role-based access control in future iterations (e.g., traders see their own models, compliance sees everything).

**Acceptance Criteria**:
- [ ] Login page offers two options: "Sign in with Google" (OAuth 2.0) and email/password form
- [ ] Google OAuth follows the standard redirect flow and creates a session
- [ ] Email/password login validates credentials and creates a session
- [ ] Session is stored as an HttpOnly secure cookie with a 24-hour expiry
- [ ] Unauthenticated access to any dashboard page redirects to the login page
- [ ] A "Sign Out" button in the user menu clears the session and redirects to login
- [ ] Failed login attempts show a clear error message (no information leakage about which field is wrong)
- [ ] After 5 failed login attempts from the same IP, a 60-second lockout is applied
- [ ] The user's display name and avatar (from Google, or initials for email users) appear in the top-right navigation

**UI/UX Notes**:
Clean, centered login page with the ExplainableEngine logo. Google button follows Google's branding guidelines. Email/password form is minimal: two fields + submit button. The page should inspire confidence and professionalism. No "Sign Up" flow in MVP -- users are provisioned by an admin.

**API Dependencies**:
- Requires backend authentication service (new): Google OAuth callback, session management, user store
- This is a **new backend requirement** not currently in the API contract

---

### US-602: API key management

**Epic**: Settings and Authentication
**Use Case**: UC-11
**Priority**: Should-have

**As a** Developer
**I want** to generate and revoke API keys from the dashboard
**So that** I can set up secure machine-to-machine integration without sharing personal credentials, and revoke access immediately if a key is compromised.

**Business Rationale**:
Programmatic API consumers (backend services, data pipelines) need API keys rather than user sessions. Key management in the dashboard eliminates the manual process of requesting keys from engineering (which today involves Slack messages and manual database insertions). Revocation capability is a security requirement: if a key is leaked, it must be disableable in seconds, not hours. API keys also enable usage tracking per integration, supporting future billing and capacity planning.

**Acceptance Criteria**:
- [ ] A "Settings > API Keys" page lists all API keys for the current user: key name, created date, last used date, status (active/revoked)
- [ ] A "Generate New Key" button opens a form: key name (required), expiry (optional: 30d, 90d, 1y, never)
- [ ] After generation, the full API key is displayed once with a copy button and a warning: "This key will not be shown again"
- [ ] The key list shows only the last 4 characters of each key (masked)
- [ ] A "Revoke" button on each key immediately disables it with a confirmation dialog
- [ ] Revoked keys remain visible in the list (grayed out) with revocation timestamp
- [ ] Maximum of 5 active keys per user
- [ ] Generated keys follow the format `ee_live_{32 random alphanumeric characters}`

**UI/UX Notes**:
Standard settings page layout. The key list is a simple table. The generation flow is a modal dialog. The "shown once" warning should be prominent (yellow background, warning icon). Consider a "Created successfully" success state that stays visible until the user dismisses it, ensuring they have time to copy the key.

**API Dependencies**:
- Requires backend API key management endpoints (new): POST /api/v1/keys, GET /api/v1/keys, DELETE /api/v1/keys/{id}
- This is a **new backend requirement** not currently in the API contract

---

### US-603: User preferences

**Epic**: Settings and Authentication
**Use Case**: UC-11
**Priority**: Nice-to-have

**As a** regular dashboard user
**I want** to set my default language (EN/IT), default narrative level (basic/advanced), and visual theme (light/dark)
**So that** the dashboard opens with my preferred configuration every time, eliminating the repetitive task of toggling settings on each visit.

**Business Rationale**:
Personalization drives adoption. Users who must reconfigure the tool on every visit experience friction that pushes them toward alternative workflows (spreadsheets, manual processes). Italian-speaking users who must switch to IT every time are especially affected, as it is their primary working language. Dark mode is not cosmetic -- trading desk environments often use low-light setups where dark themes reduce eye strain during 8+ hour sessions. Higher adoption means higher platform ROI.

**Acceptance Criteria**:
- [ ] A "Settings > Preferences" page allows configuring: default language (EN or IT), default narrative level (Basic or Advanced), visual theme (Light, Dark, or System)
- [ ] Changes are saved immediately (auto-save, no submit button required)
- [ ] Preferences are stored server-side (tied to user account) so they persist across devices
- [ ] When the user opens any explanation, the narrative defaults to the preferred language and level (overridable per-view)
- [ ] Theme changes apply immediately without page reload
- [ ] Dark theme applies to all dashboard pages consistently (no unstyled components)
- [ ] A "Reset to Defaults" button restores: EN, Basic, Light
- [ ] Preferences are loaded at login and applied before first render (no flash of wrong theme)

**UI/UX Notes**:
Simple form with radio buttons or dropdown selectors. Auto-save with a brief "Saved" confirmation toast. Theme preview: a small card showing how the selected theme looks. Preferences page is accessible from the user menu dropdown in the top-right navigation.

**API Dependencies**:
- Requires backend user preferences storage (new): GET /api/v1/users/me/preferences, PUT /api/v1/users/me/preferences
- Fallback to localStorage if backend is not available

---

## Summary

| Epic | Stories | Priority | Sprint | New Backend Required |
|------|---------|----------|--------|---------------------|
| A -- Explanation Explorer | US-101 to US-105 | Must-have | 8-9 | No |
| B -- What-if Simulator | US-201 to US-204 | Must-have (204 Should-have) | 9-10 | No |
| C -- Audit Trail | US-301 to US-304 | Should-have | 10-11 | Yes (list/search endpoint) |
| D -- Live Monitoring | US-401 to US-403 | Should-have | 11 | Yes (list endpoint, ideally stats endpoint) |
| E -- API Playground | US-501, US-502 | Nice-to-have | 12 | No |
| F -- Settings and Auth | US-601 to US-603 | Should-have (603 Nice-to-have) | 10 | Yes (auth, keys, preferences) |

### New Backend Endpoints Required

These user stories surface the need for backend endpoints not currently in the API contract:

1. **`GET /api/v1/explain`** (list/search) -- Required by US-301, US-302, US-401, US-403. Must support: cursor pagination, sort, date range filter, target search, confidence range filter, missing data filter.
2. **`GET /api/v1/stats`** (aggregation) -- Strongly recommended for US-403. Without it, the dashboard must paginate through all explanations client-side to compute daily statistics.
3. **`POST /api/v1/auth/login`**, **`POST /api/v1/auth/google`**, **`POST /api/v1/auth/logout`** -- Required by US-601.
4. **`POST /api/v1/keys`**, **`GET /api/v1/keys`**, **`DELETE /api/v1/keys/{id}`** -- Required by US-602.
5. **`GET /api/v1/users/me/preferences`**, **`PUT /api/v1/users/me/preferences`** -- Required by US-603.

### Story Dependency Map

```
US-101 (summary)  ─┐
US-102 (breakdown) ─┤
US-103 (graph)     ─┼─> US-303 (detail view) ─> US-304 (export)
US-104 (confidence)─┤
US-105 (narrative) ─┘

US-201 (sliders) ──> US-202 (comparison) ──> US-204 (save scenarios)
US-201 (sliders) ──> US-203 (sensitivity)

US-301 (list) ──> US-302 (search) ──> US-303 (detail)
US-301 (list) ──> US-304 (CSV export)

US-301 (list) ──> US-401 (feed) ──> US-402 (alerts)
US-401 (feed) ──> US-403 (stats)

US-501 (explorer) ──> US-502 (snippets)

US-601 (login) ──> US-602 (API keys)
US-601 (login) ──> US-603 (preferences)
```
