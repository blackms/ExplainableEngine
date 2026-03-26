'use client';

import { useState, useCallback, useMemo } from 'react';

interface BreakdownItem {
  node_id: string;
  label: string;
  value: number;
  weight: number;
  absolute_contribution: number;
  percentage: number;
  confidence: number;
  children?: BreakdownItem[];
}

interface BreakdownChartProps {
  breakdown: BreakdownItem[];
  onDrillDown?: (nodeId: string) => void;
}

interface BreadcrumbEntry {
  label: string;
  items: BreakdownItem[];
}

/**
 * Maps confidence to a blue with varying saturation.
 * High confidence = saturated blue, low = desaturated.
 */
function barColor(confidence: number): string {
  const saturation = Math.round(40 + confidence * 50); // 40-90%
  const lightness = Math.round(55 - confidence * 10);  // 55-45%
  return `hsl(217, ${saturation}%, ${lightness}%)`;
}

function TooltipCard({ item }: { item: BreakdownItem }) {
  return (
    <div className="rounded-lg bg-card px-4 py-3 text-sm ring-1 ring-foreground/10 shadow-lg space-y-1 min-w-[200px]">
      <p className="font-semibold">{item.label}</p>
      <p className="text-muted-foreground">Contribution: {item.absolute_contribution.toFixed(4)}</p>
      <p className="text-muted-foreground">Percentage: {item.percentage.toFixed(1)}%</p>
      <p className="text-muted-foreground">Weight: {item.weight.toFixed(4)}</p>
      <p className="text-muted-foreground">Confidence: {(item.confidence * 100).toFixed(1)}%</p>
      {item.children && item.children.length > 0 && (
        <p className="text-xs text-primary pt-1">Click to expand</p>
      )}
    </div>
  );
}

function BreakdownBar({
  item,
  maxPct,
  index,
  onClick,
}: {
  item: BreakdownItem;
  maxPct: number;
  index: number;
  onClick: (item: BreakdownItem) => void;
}) {
  const [hovered, setHovered] = useState(false);
  const widthPct = maxPct > 0 ? (item.percentage / maxPct) * 100 : 0;
  const labelInside = widthPct > 30;
  const hasChildren = item.children && item.children.length > 0;

  return (
    <div
      className="relative group"
      style={{ animationDelay: `${index * 50}ms` }}
    >
      <div
        className={`flex items-center gap-3 ${hasChildren ? 'cursor-pointer' : ''}`}
        onMouseEnter={() => setHovered(true)}
        onMouseLeave={() => setHovered(false)}
        onClick={() => onClick(item)}
        role={hasChildren ? 'button' : undefined}
        tabIndex={hasChildren ? 0 : undefined}
        onKeyDown={(e) => {
          if (hasChildren && (e.key === 'Enter' || e.key === ' ')) {
            e.preventDefault();
            onClick(item);
          }
        }}
      >
        <span className="text-sm font-medium w-[120px] shrink-0 truncate">
          {item.label}
        </span>
        <div className="flex-1 relative">
          <div
            className="h-8 rounded-full transition-all duration-500 ease-out flex items-center animate-slide-in"
            style={{
              width: `${widthPct}%`,
              backgroundColor: barColor(item.confidence),
              minWidth: '2rem',
              animationDelay: `${index * 50}ms`,
            }}
          >
            {labelInside && (
              <span className="text-xs font-medium text-white pl-3">
                {item.percentage.toFixed(1)}%
              </span>
            )}
          </div>
          {!labelInside && (
            <span
              className="absolute top-1/2 -translate-y-1/2 text-xs font-medium text-muted-foreground ml-2"
              style={{ left: `${widthPct}%` }}
            >
              {item.percentage.toFixed(1)}%
            </span>
          )}
        </div>
      </div>

      {/* Tooltip on hover */}
      {hovered && (
        <div className="absolute left-[120px] bottom-full mb-2 z-50">
          <TooltipCard item={item} />
        </div>
      )}
    </div>
  );
}

export function BreakdownChart({ breakdown, onDrillDown }: BreakdownChartProps) {
  const [drillStack, setDrillStack] = useState<BreadcrumbEntry[]>([]);

  const currentItems = useMemo(() => {
    const items =
      drillStack.length > 0
        ? drillStack[drillStack.length - 1].items
        : breakdown;
    return [...items].sort((a, b) => b.percentage - a.percentage);
  }, [drillStack, breakdown]);

  const maxPct = useMemo(
    () => Math.max(...currentItems.map((i) => i.percentage), 1),
    [currentItems]
  );

  const handleBarClick = useCallback(
    (item: BreakdownItem) => {
      if (item.children && item.children.length > 0) {
        setDrillStack((prev) => [
          ...prev,
          { label: item.label, items: item.children! },
        ]);
        onDrillDown?.(item.node_id);
      }
    },
    [onDrillDown]
  );

  const handleBreadcrumbClick = useCallback((index: number) => {
    if (index < 0) {
      setDrillStack([]);
    } else {
      setDrillStack((prev) => prev.slice(0, index + 1));
    }
  }, []);

  return (
    <div className="space-y-4">
      {/* Breadcrumb navigation */}
      {drillStack.length > 0 && (
        <nav className="flex items-center gap-1 text-sm text-muted-foreground">
          <button
            type="button"
            onClick={() => handleBreadcrumbClick(-1)}
            className="hover:text-foreground transition-colors"
          >
            Root
          </button>
          {drillStack.map((entry, i) => (
            <span key={i} className="flex items-center gap-1">
              <span>/</span>
              <button
                type="button"
                onClick={() => handleBreadcrumbClick(i)}
                className={`hover:text-foreground transition-colors ${
                  i === drillStack.length - 1 ? 'text-foreground font-medium' : ''
                }`}
              >
                {entry.label}
              </button>
            </span>
          ))}
        </nav>
      )}

      {drillStack.length > 0 && (
        <button
          type="button"
          onClick={() => setDrillStack((prev) => prev.slice(0, -1))}
          className="text-sm text-primary hover:underline"
        >
          &larr; Back
        </button>
      )}

      {/* Bars */}
      <div className="space-y-3">
        {currentItems.map((item, index) => (
          <BreakdownBar
            key={item.node_id}
            item={item}
            maxPct={maxPct}
            index={index}
            onClick={handleBarClick}
          />
        ))}
      </div>

      {/* Inline keyframe for staggered animation */}
      <style jsx>{`
        @keyframes slideIn {
          from {
            width: 0%;
            opacity: 0;
          }
          to {
            opacity: 1;
          }
        }
        .animate-slide-in {
          animation: slideIn 0.5s ease-out forwards;
        }
      `}</style>
    </div>
  );
}
