# Explainable Engine Dashboard -- Design System

> Version 1.0 | Last updated: 2026-03-26
>
> Stack: Next.js + shadcn/ui + Tailwind CSS v4 + Geist Sans/Mono
>
> This document is the single source of truth for all visual decisions.
> A developer should be able to implement any screen from this spec alone.

---

## 1. Design Principles

1. **Clarity over decoration.** Every element earns its place. If removing something changes nothing, remove it.
2. **Data density done right.** Show a lot of information without feeling overwhelming. Use whitespace and hierarchy to create breathing room within dense layouts.
3. **Progressive disclosure.** Summary first, detail on demand. Top-level views show aggregates; clicks and hovers reveal depth.
4. **Trust through design.** Confident typography, precise numbers, clear color-coded indicators. Users make decisions from this data -- it must feel authoritative.
5. **Consistent rhythm.** Predictable spacing, aligned grids, uniform patterns. Once a user learns one screen, every other screen should feel familiar.

---

## 2. Color System

All colors are defined as CSS custom properties in `globals.css` and consumed via Tailwind utilities. The system uses the existing shadcn variable architecture (`--primary`, `--background`, etc.) but remaps values to the palette below.

### 2.1 Core Palette

| Token               | Tailwind Class     | Hex       | Usage                                    |
|----------------------|--------------------|-----------|------------------------------------------|
| `--foreground`       | `text-foreground`  | `#0f172a` | Headers, primary text, emphasis (slate-900) |
| `--muted-foreground` | `text-muted-foreground` | `#64748b` | Secondary text, labels, captions (slate-500) |
| `--background`       | `bg-background`    | `#ffffff` | Page background                          |
| `--card`             | `bg-card`          | `#f8fafc` | Cards, panels, surface areas (slate-50)  |
| `--border`           | `border-border`    | `#e2e8f0` | Subtle dividers, card borders (slate-200)|
| `--primary`          | `bg-primary`       | `#2563eb` | Links, primary buttons, selected states (blue-600) |
| `--primary-foreground` | `text-primary-foreground` | `#ffffff` | Text on primary backgrounds        |
| `--accent`           | `bg-accent`        | `#eff6ff` | Hover backgrounds, subtle highlights (blue-50) |
| `--accent-foreground`| `text-accent-foreground` | `#1e40af` | Text on accent backgrounds (blue-800) |
| `--destructive`      | `bg-destructive`   | `#ef4444` | Delete actions, error states (red-500)   |
| `--ring`             | `ring-ring`        | `#2563eb` | Focus rings (matches primary)            |

### 2.2 Confidence Colors (Semantic)

These are NOT mapped to shadcn variables. Use direct Tailwind classes.

| Level                | Condition   | Color          | Hex       | Tailwind          |
|----------------------|-------------|----------------|-----------|-------------------|
| High confidence      | >= 0.8      | emerald-500    | `#10b981` | `text-emerald-500`, `bg-emerald-500` |
| High (muted bg)      | >= 0.8      | emerald-50     | `#ecfdf5` | `bg-emerald-50`   |
| Moderate confidence  | >= 0.5      | amber-500      | `#f59e0b` | `text-amber-500`, `bg-amber-500` |
| Moderate (muted bg)  | >= 0.5      | amber-50       | `#fffbeb` | `bg-amber-50`     |
| Low confidence       | < 0.5       | rose-500       | `#f43f5e` | `text-rose-500`, `bg-rose-500` |
| Low (muted bg)       | < 0.5       | rose-50        | `#fff1f2` | `bg-rose-50`      |

Helper function pattern (use across all components):

```tsx
function confidenceColor(value: number) {
  if (value >= 0.8) return { text: "text-emerald-500", bg: "bg-emerald-500", bgMuted: "bg-emerald-50", label: "High confidence" };
  if (value >= 0.5) return { text: "text-amber-500", bg: "bg-amber-500", bgMuted: "bg-amber-50", label: "Moderate confidence" };
  return { text: "text-rose-500", bg: "bg-rose-500", bgMuted: "bg-rose-50", label: "Low -- review recommended" };
}
```

### 2.3 Chart Palette (Sequential)

For bar charts, area charts, and multi-series visualizations. Ordered from most prominent to least.

| Index | Color     | Hex       | Tailwind    |
|-------|-----------|-----------|-------------|
| 1     | blue-500  | `#3b82f6` | `fill-blue-500` |
| 2     | blue-400  | `#60a5fa` | `fill-blue-400` |
| 3     | blue-300  | `#93c5fd` | `fill-blue-300` |
| 4     | sky-400   | `#38bdf8` | `fill-sky-400`  |
| 5     | sky-300   | `#7dd3fc` | `fill-sky-300`  |
| 6     | cyan-300  | `#67e8f9` | `fill-cyan-300` |

Map to CSS variables:

```css
:root {
  --chart-1: #3b82f6;
  --chart-2: #60a5fa;
  --chart-3: #93c5fd;
  --chart-4: #38bdf8;
  --chart-5: #7dd3fc;
  --chart-6: #67e8f9;
}
```

### 2.4 Dark Mode

| Token               | Light         | Dark            |
|----------------------|---------------|-----------------|
| `--background`       | `#ffffff`     | `#020617` (slate-950) |
| `--card`             | `#f8fafc`     | `#0f172a` (slate-900) |
| `--border`           | `#e2e8f0`     | `#1e293b` (slate-800) |
| `--foreground`       | `#0f172a`     | `#f1f5f9` (slate-100) |
| `--muted-foreground` | `#64748b`     | `#94a3b8` (slate-400) |
| `--primary`          | `#2563eb`     | `#3b82f6` (blue-500, slightly lighter) |
| `--accent`           | `#eff6ff`     | `#1e293b` (slate-800) |
| `--ring`             | `#2563eb`     | `#3b82f6` |

Chart colors in dark mode: use the same hues but increase lightness by one stop (e.g., blue-500 becomes blue-400). Confidence colors remain unchanged -- they are already designed for both backgrounds.

---

## 3. Typography

Font stack: **Geist Sans** (already bundled with Next.js via `next/font`). Monospace: **Geist Mono**.

### 3.1 Type Scale

| Token    | Tailwind Classes                              | Size   | Weight     | Usage                                |
|----------|-----------------------------------------------|--------|------------|--------------------------------------|
| Display  | `text-3xl font-bold tracking-tight`           | 30px   | 700        | Page titles ("Dashboard", "Audit Log") |
| Heading  | `text-xl font-semibold`                       | 20px   | 600        | Section headers ("Recent Explanations") |
| Subhead  | `text-base font-medium`                       | 16px   | 500        | Card titles, dialog headings         |
| Body     | `text-sm font-normal`                         | 14px   | 400        | Default paragraph text, descriptions |
| Caption  | `text-xs font-normal text-muted-foreground`   | 12px   | 400        | Timestamps, metadata, secondary info |
| Metric   | `text-4xl font-bold tabular-nums`             | 36px   | 700        | KPI numbers, stat card values        |
| Metric-sm| `text-2xl font-bold tabular-nums`             | 24px   | 700        | Inline metrics, secondary numbers    |
| Code     | `text-sm font-mono`                           | 14px   | 400        | Hashes, IDs, API paths, JSON         |

### 3.2 Rules

- All numeric data uses `tabular-nums` (via `font-variant-numeric: tabular-nums`) to prevent layout shifts when values change.
- Line height: use Tailwind defaults (`leading-normal` for body, `leading-tight` for headings).
- No text larger than `text-4xl` anywhere in the application.
- Truncate long text with `truncate` (single line) or `line-clamp-2` (two lines). Never let text overflow a container.

---

## 4. Spacing and Layout

### 4.1 Spacing Scale

| Context           | Tailwind   | Value  |
|-------------------|------------|--------|
| Page padding      | `p-6`     | 24px   |
| Card padding      | `p-5`     | 20px   |
| Section gap       | `space-y-6` | 24px |
| Card gap (grid)   | `gap-4`   | 16px   |
| Inner card gap    | `space-y-3` | 12px |
| Tight inner gap   | `space-y-2` | 8px  |
| Element gap       | `gap-2`   | 8px    |
| Icon-to-text gap  | `gap-2`   | 8px    |

### 4.2 Page Structure

```
max-w-7xl mx-auto p-6
```

Maximum content width: **1280px**, centered horizontally. Padding: 24px on all sides. On screens wider than 1440px, content remains centered with equal margins.

### 4.3 Grid System

| Layout         | Tailwind Classes                       | Usage                      |
|----------------|----------------------------------------|----------------------------|
| Full width     | `col-span-full`                        | Hero sections, tables      |
| Two column     | `grid grid-cols-2 gap-4`              | Side-by-side panels        |
| Three column   | `grid grid-cols-3 gap-4`              | Stat card rows             |
| Four column    | `grid grid-cols-4 gap-4`              | Stat cards on desktop      |
| Dashboard      | `grid grid-cols-12 gap-4`             | Flexible mixed layouts     |

### 4.4 Border Radius

Use the shadcn radius scale consistently:

| Element        | Tailwind     | Computed  |
|----------------|-------------|-----------|
| Cards          | `rounded-lg` | 10px     |
| Buttons        | `rounded-md` | 8px      |
| Badges         | `rounded-full`| pill     |
| Inputs         | `rounded-md` | 8px      |
| Progress bars  | `rounded-full`| pill     |
| Avatars        | `rounded-full`| circle   |

---

## 5. Component Specifications

### 5.1 Navigation Sidebar

**Dimensions:**
- Width: `w-60` (240px) when expanded
- Width: `w-16` (64px) when collapsed (icon-only, desktop)
- Height: full viewport (`h-screen`, `sticky top-0`)

**Visual:**
- Background: `bg-card` (slate-50 light, slate-900 dark)
- Right border: `border-r border-border`
- No shadow

**Logo area (top):**
- Padding: `px-4 py-5`
- Icon: 24x24px product icon
- Text: "Explainable Engine" in `text-sm font-semibold`
- Collapsed: icon only, centered

**Navigation items:**
- Padding: `py-2 px-3`
- Border radius: `rounded-md`
- Icon: 18x18px (Lucide icons), `text-muted-foreground`
- Label: `text-sm font-normal text-muted-foreground`
- Gap between icon and label: `gap-3`
- Vertical gap between items: `space-y-1`

**States:**
- Default: transparent background
- Hover: `bg-accent` (blue-50 light, slate-800 dark)
- Active/selected: `bg-primary text-primary-foreground font-medium`
- Active icon: `text-primary-foreground`

**Section groups:**
- Divider: `<Separator />` with `my-3`
- Section label (optional): `text-xs font-medium text-muted-foreground uppercase tracking-wider px-3 mb-1`

**Bottom section:**
- Separator above
- User row: avatar (32px, `rounded-full`) + name (`text-sm`) + email (`text-xs text-muted-foreground`)
- Settings gear icon button aligned right

**Mobile behavior (< 1024px):**
- Sidebar hidden by default
- Hamburger button in top-left of page header
- Opens as slide-over from left with `fixed inset-y-0 left-0 z-50`
- Backdrop: `bg-black/50` with click-to-close
- Transition: `translate-x` 250ms ease-out

---

### 5.2 Dashboard Home (Page: `/`)

**Page header:**
```
<div class="flex items-center justify-between">
  <div>
    <h1 class="text-3xl font-bold tracking-tight">Welcome back, {name}</h1>
    <p class="text-sm text-muted-foreground mt-1">Here's what's happening with your explanations.</p>
  </div>
  <div class="flex items-center gap-2">
    <Button variant="outline" size="sm">
      <Search class="h-4 w-4 mr-2" /> Search  <kbd>Cmd+K</kbd>
    </Button>
    <Button size="sm">
      <Plus class="h-4 w-4 mr-2" /> New Explanation
    </Button>
  </div>
</div>
```

**Stat cards row:**
- Layout: `grid grid-cols-4 gap-4` (desktop), `grid-cols-2` (tablet), `grid-cols-1` (mobile)
- Four cards: Total Explanations, Avg Confidence, Active Alerts, Last Updated

Each stat card (see 5.3 for full spec):
```
┌─────────────────────────┐
│ Total Explanations       │  <- caption (text-xs text-muted-foreground)
│ 142                      │  <- metric (text-4xl font-bold tabular-nums)
│ +12 from last week  ↑    │  <- trend (text-xs, green if positive)
└─────────────────────────┘
```

**Recent Explanations section:**
- Header: `flex justify-between items-center`
  - Left: "Recent Explanations" in heading style
  - Right: "View all" link (`text-sm text-primary hover:underline`)
- List: card containing rows, each row is a link

Each row:
```
<div class="flex items-center justify-between py-3 px-4 hover:bg-accent/50 rounded-md cursor-pointer transition-colors">
  <div class="flex items-center gap-3">
    <span class="text-sm font-medium font-mono">market_regime</span>
    <span class="text-sm tabular-nums text-muted-foreground">0.72</span>
  </div>
  <div class="flex items-center gap-4">
    <ConfidenceDot value={0.85} />           <!-- colored dot + percentage -->
    <span class="text-xs text-muted-foreground w-16 text-right">2m ago</span>
    <ChevronRight class="h-4 w-4 text-muted-foreground" />
  </div>
</div>
```
- Rows separated by `divide-y divide-border`

**Quick Actions section:**
- Layout: `grid grid-cols-3 gap-4` (desktop), `grid-cols-1` (mobile)
- Each action: card with icon (24px), title (`text-sm font-medium`), description (`text-xs text-muted-foreground`)
- Hover: `ring-1 ring-primary/20 shadow-sm` transition
- Actions: "New Explanation", "Search History", "Monitor Dashboard"

---

### 5.3 Stat Cards

**Container:**
- Background: `bg-card`
- Border: `border border-border`
- Radius: `rounded-lg`
- Padding: `p-5`
- Min height: none (content-driven, typically ~96px)
- Hover: `hover:shadow-sm transition-shadow duration-150`

**Content layout (vertical stack, `space-y-2`):**
1. **Label**: `text-xs font-medium text-muted-foreground uppercase tracking-wide`
2. **Value**: `text-4xl font-bold tabular-nums text-foreground`
3. **Trend** (optional): `text-xs` with colored indicator
   - Positive: `text-emerald-500` with `TrendingUp` icon (12px)
   - Negative: `text-rose-500` with `TrendingDown` icon (12px)
   - Neutral: `text-muted-foreground`
   - Format: "+12 from last week"

**Confidence stat card variant:**
- Value: percentage number
- Below value: confidence dot + label ("High confidence")
- Mini progress bar (h-1.5 rounded-full) in the confidence color

**Responsive:**
- Desktop (>1024px): `grid-cols-4`
- Tablet (640-1024px): `grid-cols-2`
- Mobile (<640px): `grid-cols-1`

---

### 5.4 Explanation Detail -- Overview Tab (Page: `/explain/[id]`)

**Page header:**
```
<div class="flex items-center justify-between">
  <div class="flex items-center gap-3">
    <Button variant="ghost" size="icon" asChild>
      <Link href="/"><ArrowLeft class="h-4 w-4" /></Link>
    </Button>
    <div>
      <h1 class="text-xl font-semibold font-mono">market_regime_score</h1>
      <p class="text-xs text-muted-foreground">Created 2 hours ago by API</p>
    </div>
  </div>
  <ExportDropdown />  <!-- shadcn DropdownMenu: "Export JSON", "Export CSV", "Copy Link" -->
</div>
```

**Top metrics row:**
- Layout: `grid grid-cols-2 gap-4`

Card 1 -- Final Value:
```
┌────────────────────────────────┐
│ Final Value                     │  <- text-xs text-muted-foreground uppercase
│ 0.72                            │  <- text-4xl font-bold tabular-nums
│ Computed from 3 components      │  <- text-xs text-muted-foreground
└────────────────────────────────┘
```

Card 2 -- Confidence:
```
┌────────────────────────────────┐
│ Confidence                      │  <- text-xs text-muted-foreground uppercase
│ 82.5%                           │  <- text-4xl font-bold tabular-nums text-emerald-500
│ ████████████████░░░░  High      │  <- h-2 rounded-full progress bar
└────────────────────────────────┘
```

**Component Breakdown section:**
- Section header: "Component Breakdown" in heading style
- Full-width card containing horizontal bar chart (see 5.5)

**Bottom row:**
- Layout: `grid grid-cols-2 gap-4`

Card 1 -- Top Drivers:
```
┌────────────────────────────────┐
│ Top Drivers                     │
│                                 │
│ 1  trend_strength   1.00       │  <- rank number (text-xs text-muted-foreground),
│ 2  momentum         0.53       │     name (text-sm font-mono), value (text-sm tabular-nums)
│ 3  volatility       0.36       │
└────────────────────────────────┘
```
- Each row: `flex items-center gap-3 py-2`
- Rank: `w-6 text-center text-xs text-muted-foreground`
- Name: `flex-1 text-sm font-mono`
- Value: `text-sm tabular-nums text-right`
- Separator between rows: `divide-y divide-border`

Card 2 -- Data Quality:
```
┌────────────────────────────────┐
│ Data Quality                    │
│                                 │
│ ✓ No missing data              │  <- CheckCircle icon (emerald-500) + text-sm
│ ✓ All inputs validated         │
│ ✓ 3/3 components computed      │
└────────────────────────────────┘
```
- Each row: icon (16px) + label (`text-sm`)
- Success: `CheckCircle` in emerald-500
- Warning: `AlertTriangle` in amber-500
- Error: `XCircle` in rose-500

**Tab bar** (below header, above content):
- Uses shadcn `<Tabs>` component
- Tabs: Overview | Graph | What-If | Narrative
- Active tab: underline style (`border-b-2 border-primary text-primary`)
- Inactive: `text-muted-foreground hover:text-foreground`
- Full-width bottom border: `border-b border-border`

---

### 5.5 Breakdown Chart

**Container:** full-width within its parent card, padding `p-0` (card padding handles outer spacing).

**Bars:**
- Layout: vertical stack, `space-y-3`
- Each bar row: label + bar + percentage

```
<div class="space-y-3">
  {components.map((c, i) => (
    <div key={c.name} class="space-y-1">
      <div class="flex items-center justify-between">
        <span class="text-sm font-mono">{c.name}</span>
        <span class="text-sm tabular-nums font-medium">{c.percentage}%</span>
      </div>
      <div class="h-8 w-full bg-muted rounded-full overflow-hidden">
        <div
          class="h-full bg-primary rounded-full transition-all duration-300 ease-out"
          style={{ width: `${c.percentage}%`, transitionDelay: `${i * 50}ms` }}
        />
      </div>
    </div>
  ))}
</div>
```

**Bar visual spec:**
- Track: `h-8 bg-muted rounded-full`
- Fill: `h-8 bg-primary rounded-full`
- Color intensity: full `bg-primary` for the top component, `bg-primary/80` for second, `bg-primary/60` for third, etc. (opacity decreases by 20% per rank, minimum 40%)
- Corner radius: `rounded-full` on both track and fill (pill shape)

**Hover behavior:**
- Bar fill: lighten slightly (`bg-primary/90`)
- Show tooltip card positioned above the bar:
  - Content: component name, raw value, percentage, confidence
  - Style: `bg-popover text-popover-foreground shadow-md rounded-lg p-3 text-sm`

**Click/expand behavior:**
- Clicking a bar row expands it to show child components (if any)
- Children indented with `pl-6`
- Child bars use next chart color in sequence
- Expand/collapse: `ChevronDown` icon rotates, content slides down 200ms ease-out

**Animation on first render:**
- Bar widths animate from 0% to final value
- Duration: 300ms ease-out
- Stagger: 50ms delay per bar (first bar starts immediately, second at 50ms, etc.)

---

### 5.6 Confidence Gauge

**DO NOT** use a circular/radial gauge. Use a horizontal progress bar for cleanliness.

**Layout (vertical stack, `space-y-1`):**

```
<div class="space-y-1">
  <div class="flex items-baseline gap-2">
    <span class="text-4xl font-bold tabular-nums {confidenceColor(value).text}">
      {(value * 100).toFixed(1)}%
    </span>
    <span class="text-sm text-muted-foreground">{confidenceColor(value).label}</span>
  </div>
  <div class="h-3 w-full bg-muted rounded-full overflow-hidden">
    <div
      class="h-full rounded-full transition-all duration-500 ease-out {confidenceColor(value).bg}"
      style={{ width: `${value * 100}%` }}
    />
  </div>
</div>
```

**Spec:**
- Progress track: `h-3 bg-muted rounded-full`
- Progress fill: `h-3 rounded-full` in the confidence color
- Number: `text-4xl font-bold tabular-nums` in the confidence color
- Label: `text-sm text-muted-foreground`
  - >= 0.8: "High confidence"
  - >= 0.5: "Moderate confidence"
  - < 0.5: "Low -- review recommended"

**Compact variant** (for use in tables and list rows):
- Dot only: `<span class="inline-block h-2 w-2 rounded-full {confidenceColor(value).bg}" />`
- Dot + number: dot followed by `text-xs tabular-nums` percentage

---

### 5.7 Graph View (Page: `/explain/[id]/graph`)

Uses `@xyflow/react` (already installed).

**Canvas:**
- Background: white with subtle dot grid
  - Dot color: `slate-200` (light) / `slate-700` (dark)
  - Dot spacing: 20px
  - Dot size: 1px
- Full height: `calc(100vh - 64px)` (minus header)
- Full width: 100% of content area (sidebar excluded)

**Nodes (custom node component):**
```
┌──────────────────────────┐
│ trend_strength           │  <- text-sm font-mono font-semibold, p-3
│──────────────────────────│
│ Value: 1.00    85% ●     │  <- text-xs tabular-nums, confidence dot
│ ██████████████████░░     │  <- h-1.5 confidence bar
└──────────────────────────┘
```

- Size: min-width 200px, auto-height
- Background: `bg-card`
- Border: `border-2` colored by node type:
  - Source/input: `border-emerald-300`
  - Intermediate: `border-blue-300`
  - Output/root: `border-violet-300`
  - Missing/error: `border-rose-300`
- Radius: `rounded-lg`
- Shadow: `shadow-sm`
- Selected state: `ring-2 ring-primary shadow-md`

**Edges:**
- Style: bezier curves (default xyflow behavior)
- Stroke color: `slate-300` (light) / `slate-600` (dark)
- Stroke width: `1 + (weight * 3)` px (min 1px, max 4px based on edge weight)
- Animated: false (static lines for clarity)
- Label on hover: weight value in a small tooltip (`bg-popover text-xs rounded px-2 py-1`)
- Arrow: small arrowhead at target end

**Controls panel (bottom-left):**
- Position: `absolute bottom-4 left-4`
- Background: `bg-card/90 backdrop-blur-sm`
- Border: `border border-border rounded-lg shadow-sm`
- Buttons: zoom in, zoom out, fit view, reset
- Each button: `p-2` icon-only, `hover:bg-accent`
- Separator between zoom and layout buttons

**Minimap (bottom-right):**
- Position: `absolute bottom-4 right-4`
- Size: 160x120px
- Background: `bg-card/80 backdrop-blur-sm`
- Border: `border border-border rounded-lg`
- Viewport indicator: `border-2 border-primary/50`

---

### 5.8 What-If Simulator (Page: `/explain/[id]/whatif`)

**Layout:**
- Two columns: `grid grid-cols-12 gap-6`
- Left panel (inputs): `col-span-5`
- Right panel (results): `col-span-7`
- Mobile: single column, inputs above results

**Left panel -- Input Sliders:**

Card title: "Adjust Inputs" in subhead style.

Each slider group:
```
<div class="space-y-2 py-3 border-b border-border last:border-0">
  <div class="flex items-center justify-between">
    <label class="text-sm font-medium font-mono">{inputName}</label>
    <Button variant="ghost" size="icon" class="h-6 w-6">
      <RotateCcw class="h-3 w-3" />  <!-- reset to original -->
    </Button>
  </div>
  <Slider value={current} min={min} max={max} step={step} />
  <div class="flex items-center justify-between">
    <span class="text-sm tabular-nums">{currentValue}</span>
    <DeltaDisplay original={original} current={current} />
  </div>
</div>
```

**Delta display:**
- Positive change: `text-emerald-500 text-xs font-medium` showing "+0.15 (+18.7%)"
- Negative change: `text-rose-500 text-xs font-medium` showing "-0.10 (-12.5%)"
- No change: `text-muted-foreground text-xs` showing "No change"

**Slider styling:**
- Track: `h-2 bg-muted rounded-full`
- Fill: `bg-primary`
- Thumb: `h-4 w-4 rounded-full bg-primary border-2 border-background shadow-sm`
- Focus: `ring-2 ring-ring ring-offset-2`

**Right panel -- Results:**

Card title: "Simulation Results" in subhead style.

Before/After comparison:
```
┌──────────────────────────────────────────┐
│ Predicted Output                          │
│                                          │
│ 0.72  →  0.68                            │  <- text-2xl font-bold tabular-nums
│ Original   Simulated                      │  <- text-xs text-muted-foreground
│                                          │
│ Delta: -0.04 (-5.6%)                     │  <- rose-500 text-sm font-medium
└──────────────────────────────────────────┘
```

Arrow between values: `ArrowRight` icon (16px, `text-muted-foreground`)

Below: diff table showing each component's before/after:
```
Component        Original    New        Delta
─────────────────────────────────────────────
trend_strength   0.49        0.44       -0.05
momentum         0.28        0.28        0.00
volatility       0.23        0.26       +0.03
```
- Table: shadcn `<Table>` component
- Delta column: colored (emerald/rose/muted-foreground)
- Changed rows: subtle `bg-amber-50/50` highlight

**"Run Simulation" button:**
- Position: bottom of left panel, sticky
- Style: `Button` primary, full width
- Loading state: spinner icon + "Simulating..."

---

### 5.9 AI Chat Panel

**Container:**
- Position: fixed right side panel
- Width: `w-96` (384px)
- Height: full viewport
- Background: `bg-background`
- Border: `border-l border-border`
- Shadow: `shadow-xl` (prominent, floating feel)
- Z-index: `z-40`

**Open/close:**
- Toggle button: floating bottom-right of page, `bg-primary text-primary-foreground rounded-full p-3 shadow-lg`
- Icon: `MessageSquare` (20px)
- Panel slides in from right: `translate-x-full` to `translate-x-0`, 250ms ease-out

**Header:**
```
<div class="flex items-center justify-between px-4 py-3 border-b border-border">
  <div class="flex items-center gap-2">
    <Sparkles class="h-4 w-4 text-primary" />
    <span class="text-sm font-semibold">AI Assistant</span>
  </div>
  <Button variant="ghost" size="icon" class="h-8 w-8">
    <X class="h-4 w-4" />
  </Button>
</div>
```

**Message area:**
- Scrollable: `flex-1 overflow-y-auto p-4 space-y-4`

User message:
```
<div class="flex justify-end">
  <div class="max-w-[85%] bg-primary text-primary-foreground rounded-2xl rounded-br-sm px-4 py-2.5">
    <p class="text-sm">{message}</p>
    <span class="text-xs text-primary-foreground/60 mt-1 block text-right">10:32 AM</span>
  </div>
</div>
```

AI message:
```
<div class="flex justify-start">
  <div class="max-w-[85%] bg-card border border-border rounded-2xl rounded-bl-sm px-4 py-2.5">
    <p class="text-sm text-foreground">{message}</p>
    <span class="text-xs text-muted-foreground mt-1 block">10:32 AM</span>
  </div>
</div>
```

**Suggested questions** (shown when chat is empty or after AI response):
- Layout: flex wrap, `gap-2`, above input area
- Each chip: `bg-accent text-accent-foreground text-xs px-3 py-1.5 rounded-full border border-border cursor-pointer hover:bg-accent/80 transition-colors`
- Examples: "Why is confidence low?", "Explain the top driver", "What if momentum doubles?"

**Input area:**
```
<div class="border-t border-border p-3">
  <div class="flex items-end gap-2">
    <Textarea
      class="flex-1 min-h-[40px] max-h-[120px] resize-none text-sm"
      placeholder="Ask about this explanation..."
    />
    <Button size="icon" class="h-9 w-9 shrink-0">
      <Send class="h-4 w-4" />
    </Button>
  </div>
</div>
```

---

### 5.10 Browse/Audit Table (Page: `/audit`)

**Page header:**
- Title: "Audit Log" in display style
- Subtitle: "Browse and search all explanations" in `text-sm text-muted-foreground`

**Filters row:**
- Inline above table (NOT in a separate panel)
- Layout: `flex items-center gap-3 flex-wrap`
- Components:
  - Search input: `<Input>` with search icon, `w-64`
  - Confidence filter: `<Select>` with options "All", "High", "Moderate", "Low"
  - Date range: `<Select>` with "Last 24h", "Last 7d", "Last 30d", "All time"
  - Clear filters: `<Button variant="ghost" size="sm">` shown only when filters are active
- Margin below: `mb-4`

**Table:**

```
<Table>
  <TableHeader>
    <TableRow class="hover:bg-transparent">
      <TableHead class="text-xs font-medium text-muted-foreground uppercase tracking-wide">
        Name
      </TableHead>
      ...
    </TableRow>
  </TableHeader>
</Table>
```

Columns:
| Column     | Width   | Alignment | Content                                    |
|------------|---------|-----------|---------------------------------------------|
| Name       | flex-1  | left      | `font-mono text-sm font-medium`            |
| Value      | w-24    | right     | `tabular-nums text-sm`                     |
| Confidence | w-28    | left      | colored dot (h-2 w-2) + `text-sm tabular-nums` percentage |
| Components | w-20    | right     | count (`text-sm text-muted-foreground`)    |
| Created    | w-32    | right     | relative time (`text-sm text-muted-foreground`) |
| Action     | w-10    | center    | `ChevronRight` icon (16px)                 |

**Row behavior:**
- Entire row is clickable: `cursor-pointer`
- Hover: `hover:bg-accent/50 transition-colors`
- Click navigates to `/explain/[id]`
- Focus: `focus-visible:bg-accent/50 focus-visible:outline-none`

**Timestamps:**
- Display: relative ("2h ago", "3d ago")
- Tooltip (shadcn `<Tooltip>`): absolute datetime ("March 26, 2026 at 14:32:05 UTC")

**Empty state:** See section 5.12.

**Pagination:**
```
<div class="flex items-center justify-between pt-4">
  <span class="text-sm text-muted-foreground">Showing 1-20 of 142 explanations</span>
  <div class="flex items-center gap-2">
    <Button variant="outline" size="sm" disabled={page === 1}>Previous</Button>
    <Button variant="outline" size="sm" disabled={page === lastPage}>Next</Button>
  </div>
</div>
```

---

### 5.11 Toast Notifications

Uses shadcn `<Sonner>` (already in the project or add via `shadcn add sonner`).

**Position:** bottom-right (`bottom-right` in Sonner config).

**Visual:**
- Background: `bg-card`
- Border: `border border-border`
- Radius: `rounded-lg`
- Shadow: `shadow-lg`
- Padding: `p-4`
- Max width: 360px
- Left accent border (4px) by type:
  - Success: `border-l-4 border-l-emerald-500`
  - Error: `border-l-4 border-l-rose-500`
  - Info: `border-l-4 border-l-blue-500`
  - Warning: `border-l-4 border-l-amber-500`

**Content:**
- Title: `text-sm font-medium`
- Description: `text-sm text-muted-foreground`
- Action button (optional): `text-sm font-medium text-primary hover:underline`

**Behavior:**
- Auto-dismiss: 4000ms
- Progress bar at bottom: thin (h-0.5), shrinks from full width to 0 over the dismiss duration
- Hover pauses auto-dismiss
- Close button: `X` icon, top-right, `text-muted-foreground`
- Stack: max 3 visible, newer on top

---

### 5.12 Empty States

**Layout:** centered within the content area.

```
<div class="flex flex-col items-center justify-center py-16 space-y-4">
  <div class="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
    <FileQuestion class="h-6 w-6 text-muted-foreground" />
  </div>
  <div class="text-center space-y-1">
    <h3 class="text-base font-medium">No explanations yet</h3>
    <p class="text-sm text-muted-foreground max-w-sm">
      Create your first explanation to see it here.
    </p>
  </div>
  <Button size="sm">
    <Plus class="h-4 w-4 mr-2" /> Create Explanation
  </Button>
</div>
```

**Contextual variants:**

| Page            | Icon           | Title                      | Description                                    |
|-----------------|----------------|----------------------------|------------------------------------------------|
| Dashboard       | `BarChart3`    | "No explanations yet"      | "Create your first explanation to get started." |
| Audit           | `FileSearch`   | "No results found"         | "Try adjusting your filters or search terms."  |
| Graph           | `GitBranch`    | "Graph not available"      | "This explanation has no graph data."          |
| What-If         | `SlidersH`     | "No simulation data"       | "Run a simulation to see results."             |
| Monitor (empty) | `Activity`     | "No alerts"                | "Everything looks good. No active alerts."     |
| Search (no results) | `Search`   | "No matches"               | "We couldn't find anything matching your search." |

---

## 6. Animation and Transitions

All animations use CSS transitions or Tailwind's `transition-*` utilities. No external animation libraries required.

| Element              | Property          | Duration | Easing      | Tailwind / CSS                              |
|----------------------|-------------------|----------|-------------|----------------------------------------------|
| Page fade-in         | opacity           | 200ms    | ease-out    | `animate-in fade-in duration-200`            |
| Tab content switch   | opacity + translateY | 150ms | ease-out    | `animate-in fade-in slide-in-from-bottom-1 duration-150` |
| Sidebar slide (mobile)| translateX       | 250ms    | ease-out    | `transition-transform duration-250 ease-out` |
| Chat panel slide     | translateX        | 250ms    | ease-out    | `transition-transform duration-250 ease-out` |
| Card hover shadow    | box-shadow        | 150ms    | ease        | `transition-shadow duration-150`             |
| Button press         | transform         | 100ms    | ease        | `active:scale-[0.98] transition-transform duration-100` |
| Breakdown bars       | width             | 300ms    | ease-out    | Inline style: `transition: width 300ms ease-out` + staggered `transitionDelay` |
| Confidence bar       | width             | 500ms    | ease-out    | Inline style: `transition: width 500ms ease-out` |
| Number count-up      | --                | 400ms    | ease-out    | Use `framer-motion`'s `useSpring` or a simple `requestAnimationFrame` counter |
| Expand/collapse      | height + opacity  | 200ms    | ease-out    | CSS `grid-template-rows: 0fr/1fr` transition |
| Skeleton loading     | opacity pulse     | 2000ms   | ease-in-out | shadcn `<Skeleton>` (built-in pulse)         |
| Toast enter          | translateY + opacity | 200ms | ease-out    | Sonner default animation                     |
| Tooltip              | opacity + scale   | 100ms    | ease-out    | shadcn `<Tooltip>` default                   |

**Rules:**
- Never animate layout properties (`width`, `height` of containers). Use `transform` and `opacity` wherever possible.
- Exception: chart bars and progress bars may animate `width` since they are visual indicators, not layout containers.
- Respect `prefers-reduced-motion`: wrap all non-essential animations in `@media (prefers-reduced-motion: no-preference)`.
- No animation should exceed 500ms. The UI must never feel sluggish.

---

## 7. Responsive Breakpoints

Using Tailwind's default breakpoint system:

| Breakpoint | Width    | Prefix | Layout behavior                                |
|------------|----------|--------|------------------------------------------------|
| Mobile     | < 640px  | (none) | Single column. Sidebar collapsed (overlay). Stacked cards. Full-width tables (horizontal scroll). |
| Tablet     | >= 640px | `sm:`  | Two-column grids. Sidebar still collapsed. Stat cards in 2-column grid. |
| Desktop    | >= 1024px| `lg:`  | Full layout with sidebar visible. 3-4 column grids. Side panels. |
| Wide       | >= 1440px| `xl:`  | `max-w-7xl` container centered. No layout change, just more breathing room. |

**Specific responsive rules:**

| Component        | Mobile              | Tablet              | Desktop              |
|------------------|---------------------|----------------------|----------------------|
| Sidebar          | Hidden, overlay     | Hidden, overlay      | Visible, w-60        |
| Stat cards       | 1 col               | 2 cols               | 4 cols               |
| Dashboard grid   | Stack               | 2 cols               | Full 12-col grid     |
| What-If layout   | Stack (inputs above) | Stack               | 2 cols (5/7 split)   |
| Chat panel       | Full screen overlay  | w-80 side panel      | w-96 side panel      |
| Table             | Horizontal scroll   | Horizontal scroll    | Full width            |
| Graph             | Full width, no minimap | Full, minimap     | Full, minimap + controls |

---

## 8. Dark Mode

Dark mode is toggled via a class on `<html>` (`class="dark"`). The existing `@custom-variant dark (&:is(.dark *))` in `globals.css` handles scoping.

### 8.1 CSS Variable Overrides

All dark mode values are defined in the `.dark` selector in `globals.css`. See section 2.4 for the full mapping.

### 8.2 Component Adjustments

| Component       | Light                        | Dark                              |
|-----------------|------------------------------|-----------------------------------|
| Cards           | `bg-card` (slate-50)         | `bg-card` (slate-900)            |
| Borders         | `border-border` (slate-200)  | `border-border` (slate-800)      |
| Chart bars      | `bg-primary` (blue-600)      | `bg-primary` (blue-500)          |
| Graph dot grid  | slate-200 dots               | slate-700 dots                   |
| Graph nodes     | white bg, colored border     | slate-900 bg, colored border     |
| Confidence dots | Same colors (designed for both) | Same colors                    |
| Shadows         | `shadow-sm`                  | `shadow-sm shadow-black/20` (slightly stronger) |
| Code blocks     | `bg-slate-100`               | `bg-slate-800`                   |

### 8.3 Toggle Control

- Location: bottom of sidebar (or in settings)
- Component: shadcn `<Switch>` or a three-state toggle (light / system / dark)
- Icon: `Sun` / `Moon` / `Monitor`
- Persist preference in `localStorage` under key `theme`
- Default: `system` (follows OS preference via `prefers-color-scheme`)

---

## 9. Implementation Priority

### Phase 1 -- Foundation (Week 1)

| # | Task                              | Files affected                          |
|---|-----------------------------------|-----------------------------------------|
| 1 | Update CSS variables in globals.css | `src/app/globals.css`                 |
| 2 | Configure typography (Geist Sans) | `src/app/layout.tsx`                   |
| 3 | Redesign navigation sidebar       | `src/components/app-sidebar.tsx` (or new) |
| 4 | Dashboard home page               | `src/app/page.tsx`                     |
| 5 | Stat card component               | `src/components/dashboard/StatCard.tsx` |
| 6 | Explanation detail overview tab   | `src/app/explain/[id]/page.tsx`        |
| 7 | Breakdown chart redesign          | `src/components/explanation/BreakdownChart.tsx` |
| 8 | Confidence gauge redesign         | `src/components/explanation/ConfidenceGauge.tsx` (new or refactor) |

### Phase 2 -- Feature Pages (Week 2)

| # | Task                              | Files affected                          |
|---|-----------------------------------|-----------------------------------------|
| 9 | Graph view polish                 | `src/app/explain/[id]/graph/page.tsx`  |
| 10| What-If simulator redesign        | `src/app/explain/[id]/whatif/`         |
| 11| AI Chat panel                     | `src/components/chat/ChatPanel.tsx` (new) |
| 12| Audit/Browse table                | `src/app/audit/page.tsx`               |
| 13| Monitor page                      | `src/app/monitor/page.tsx`             |
| 14| Empty states (all pages)          | Per-page components                    |
| 15| Toast notification system         | `src/components/ui/sonner.tsx` + layout |

### Phase 3 -- Polish (Week 3)

| # | Task                              | Files affected                          |
|---|-----------------------------------|-----------------------------------------|
| 16| Entrance animations               | Components from Phase 1-2              |
| 17| Number count-up animations        | StatCard, metrics                      |
| 18| Dark mode full audit              | `globals.css` + all components         |
| 19| Mobile responsive pass            | All page and component files           |
| 20| Keyboard shortcuts (Cmd+K, etc.) | `src/components/CommandPalette.tsx` (new) |
| 21| Loading states (skeletons)        | All `loading.tsx` files                |
| 22| Accessibility audit (a11y)        | All interactive components             |

---

## 10. Accessibility Requirements

These are non-negotiable and apply to every component:

1. **Color contrast:** All text meets WCAG 2.1 AA (4.5:1 for body text, 3:1 for large text). The confidence colors on white/dark backgrounds have been verified.
2. **Focus indicators:** Every interactive element has a visible focus ring (`ring-2 ring-ring ring-offset-2`). Never remove outlines.
3. **Keyboard navigation:** All functionality reachable via keyboard. Tab order follows visual order. Escape closes modals/panels.
4. **Screen readers:** Use semantic HTML (`<nav>`, `<main>`, `<table>`, `<h1>`-`<h6>`). Add `aria-label` to icon-only buttons. Use `role="status"` for live-updating metrics.
5. **Reduced motion:** Wrap animations in `@media (prefers-reduced-motion: no-preference)`. Provide instant transitions as fallback.
6. **Touch targets:** Minimum 44x44px for all interactive elements on mobile.

---

## 11. File Naming Conventions

```
src/
  components/
    ui/              <- shadcn primitives (do not modify unless necessary)
    dashboard/       <- dashboard-specific components
      StatCard.tsx
      RecentList.tsx
      QuickActions.tsx
    explanation/     <- explanation detail components
      BreakdownChart.tsx
      ConfidenceGauge.tsx
      DriverRanking.tsx
      DataQuality.tsx
    chat/            <- AI chat components
      ChatPanel.tsx
      ChatMessage.tsx
      SuggestedQuestions.tsx
    graph/           <- graph view components
      CustomNode.tsx
      GraphControls.tsx
    monitor/         <- monitoring components
      AlertPanel.tsx
      LiveFeed.tsx
      StatsCards.tsx
    layout/          <- layout components
      AppSidebar.tsx
      PageHeader.tsx
      EmptyState.tsx
      CommandPalette.tsx
```

All component files: PascalCase. All utility/hook files: camelCase. One component per file.
