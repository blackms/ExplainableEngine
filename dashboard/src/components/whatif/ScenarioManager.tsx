'use client';

import { useState, useEffect, useCallback } from 'react';
import { getScenarios, deleteScenario, clearScenarios, type SavedScenario } from '@/lib/scenarios';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

interface ScenarioManagerProps {
  explanationId: string;
  onSelect: (scenarios: SavedScenario[]) => void;
  onLoad: (scenario: SavedScenario) => void;
  refreshKey?: number;
}

export function ScenarioManager({ explanationId, onSelect, onLoad, refreshKey }: ScenarioManagerProps) {
  const [scenarios, setScenarios] = useState<SavedScenario[]>([]);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [confirmClear, setConfirmClear] = useState(false);

  const refresh = useCallback(() => {
    setScenarios(getScenarios(explanationId));
  }, [explanationId]);

  useEffect(() => {
    refresh();
  }, [refresh, refreshKey]);

  const handleToggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else if (next.size < 3) {
        next.add(id);
      }
      return next;
    });
  };

  const handleDelete = (scenarioId: string) => {
    deleteScenario(explanationId, scenarioId);
    setSelected((prev) => {
      const next = new Set(prev);
      next.delete(scenarioId);
      return next;
    });
    refresh();
    onSelect([]);
  };

  const handleClearAll = () => {
    if (!confirmClear) {
      setConfirmClear(true);
      return;
    }
    clearScenarios(explanationId);
    setScenarios([]);
    setSelected(new Set());
    setConfirmClear(false);
    onSelect([]);
  };

  const handleCompare = () => {
    const selectedScenarios = scenarios.filter((s) => selected.has(s.id));
    onSelect(selectedScenarios);
  };

  const formatTimestamp = (iso: string) => {
    const date = new Date(iso);
    return date.toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Saved Scenarios</CardTitle>
          <div className="flex items-center gap-2">
            {selected.size >= 2 && (
              <Button size="sm" onClick={handleCompare}>
                Compare Selected ({selected.size})
              </Button>
            )}
            {scenarios.length > 0 && (
              <Button
                variant="destructive"
                size="sm"
                onClick={handleClearAll}
                onBlur={() => setConfirmClear(false)}
              >
                {confirmClear ? 'Confirm Clear All' : 'Clear All'}
              </Button>
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {scenarios.length === 0 ? (
          <p className="text-sm text-muted-foreground py-4 text-center">
            No saved scenarios yet. Run a what-if analysis and save it to start comparing.
          </p>
        ) : (
          <div className="space-y-3">
            {scenarios.map((scenario) => {
              const isSelected = selected.has(scenario.id);
              const deltaSign = scenario.result.delta_value >= 0 ? '+' : '';
              const deltaColor =
                scenario.result.delta_value >= 0 ? 'text-green-600' : 'text-red-600';

              return (
                <div
                  key={scenario.id}
                  className={`flex items-center justify-between rounded-lg border p-3 transition-colors ${
                    isSelected ? 'border-primary bg-primary/5' : 'border-border'
                  }`}
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <input
                      type="checkbox"
                      checked={isSelected}
                      onChange={() => handleToggleSelect(scenario.id)}
                      disabled={!isSelected && selected.size >= 3}
                      className="size-4 shrink-0 rounded border-input accent-primary"
                      aria-label={`Select ${scenario.name} for comparison`}
                    />
                    <div className="min-w-0">
                      <p className="text-sm font-medium truncate">{scenario.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {scenario.modifications.length} modification{scenario.modifications.length !== 1 ? 's' : ''}
                        {' \u00B7 '}
                        {formatTimestamp(scenario.savedAt)}
                        {' \u00B7 '}
                        <span className={deltaColor}>
                          {deltaSign}{scenario.result.delta_percentage.toFixed(1)}%
                        </span>
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    <Button variant="outline" size="sm" onClick={() => onLoad(scenario)}>
                      Load
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => handleDelete(scenario.id)}>
                      Delete
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
