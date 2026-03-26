'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { useMutation } from '@tanstack/react-query';
import { api } from '@/lib/api/client';
import type { ExplainResponse, Modification, SensitivityResult } from '@/lib/api/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ComponentSlider } from './ComponentSlider';
import { ComparisonView } from './ComparisonView';
import { SensitivityRanking } from './SensitivityRanking';
import { saveScenario } from '@/lib/scenarios';

interface WhatIfSimulatorProps {
  explanation: ExplainResponse;
  initialModifications?: Record<string, number>;
  onScenarioSaved?: () => void;
}

export function WhatIfSimulator({ explanation, initialModifications, onScenarioSaved }: WhatIfSimulatorProps) {
  const [modifications, setModifications] = useState<Record<string, number>>(() => {
    const initial: Record<string, number> = {};
    for (const item of explanation.breakdown) {
      initial[item.label] = initialModifications?.[item.label] ?? item.value;
    }
    return initial;
  });

  const [showSaveDialog, setShowSaveDialog] = useState(false);
  const [scenarioName, setScenarioName] = useState('');

  // Sync when initialModifications changes externally (e.g. loading a scenario)
  useEffect(() => {
    if (initialModifications) {
      setModifications((prev) => {
        const next: Record<string, number> = {};
        for (const item of explanation.breakdown) {
          next[item.label] = initialModifications[item.label] ?? item.value;
        }
        return next;
      });
    }
  }, [initialModifications, explanation.breakdown]);

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const whatIfMutation = useMutation({
    mutationFn: (mods: Modification[]) => api.whatIf(explanation.id, mods),
  });

  const buildModList = useCallback((): Modification[] => {
    return Object.entries(modifications)
      .filter(([name, val]) => {
        const original = explanation.breakdown.find((b) => b.label === name);
        return original && Math.abs(val - original.value) > 0.001;
      })
      .map(([name, val]) => ({ component: name, new_value: val }));
  }, [modifications, explanation.breakdown]);

  useEffect(() => {
    const modList = buildModList();
    if (modList.length === 0) return;

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    debounceRef.current = setTimeout(() => {
      whatIfMutation.mutate(modList);
    }, 300);

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [modifications]);

  const handleSliderChange = useCallback((name: string, value: number) => {
    setModifications((prev) => ({ ...prev, [name]: value }));
  }, []);

  const handleResetAll = useCallback(() => {
    const initial: Record<string, number> = {};
    for (const item of explanation.breakdown) {
      initial[item.label] = item.value;
    }
    setModifications(initial);
  }, [explanation.breakdown]);

  const handleAnalyze = useCallback(() => {
    const modList = buildModList();
    if (modList.length > 0) {
      whatIfMutation.mutate(modList);
    }
  }, [buildModList, whatIfMutation]);

  const hasChanges = buildModList().length > 0;

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>What-if Analysis</CardTitle>
          <p className="text-sm text-muted-foreground">
            Adjust component values to see how they affect the outcome for{' '}
            <span className="font-medium text-foreground">{explanation.target}</span>
          </p>
        </CardHeader>
        <CardContent className="space-y-6">
          {explanation.breakdown.map((item) => {
            const absVal = Math.abs(item.value);
            const minVal = Math.min(0, item.value - absVal * 0.5);
            const maxVal = item.value + absVal * 0.5 || 1;
            const step = Math.max(0.0001, absVal * 0.01);

            return (
              <ComponentSlider
                key={item.node_id}
                name={item.label}
                originalValue={item.value}
                value={modifications[item.label] ?? item.value}
                onChange={(val) => handleSliderChange(item.label, val)}
                min={minVal}
                max={maxVal}
                step={step}
              />
            );
          })}

          <div className="flex items-center gap-3 pt-2">
            <button
              type="button"
              onClick={handleAnalyze}
              disabled={!hasChanges || whatIfMutation.isPending}
              className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {whatIfMutation.isPending ? 'Analyzing...' : 'Analyze'}
            </button>
            <button
              type="button"
              onClick={handleResetAll}
              disabled={!hasChanges}
              className="rounded-md border px-4 py-2 text-sm font-medium hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Reset All
            </button>
            {whatIfMutation.data && (
              <Button
                variant="outline"
                size="default"
                onClick={() => {
                  setScenarioName('');
                  setShowSaveDialog(true);
                }}
              >
                Save Scenario
              </Button>
            )}
            {whatIfMutation.isError && (
              <span className="text-sm text-red-600">
                Error: {whatIfMutation.error?.message ?? 'Analysis failed'}
              </span>
            )}
          </div>

          {showSaveDialog && (
            <div className="flex items-center gap-2 rounded-lg border border-border bg-muted/30 p-3">
              <Input
                placeholder="Scenario name..."
                value={scenarioName}
                onChange={(e) => setScenarioName(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && scenarioName.trim() && whatIfMutation.data) {
                    const modList = buildModList();
                    saveScenario({
                      name: scenarioName.trim(),
                      explanationId: explanation.id,
                      modifications: modList,
                      result: whatIfMutation.data,
                    });
                    setShowSaveDialog(false);
                    setScenarioName('');
                    onScenarioSaved?.();
                  }
                }}
                className="max-w-xs"
                autoFocus
              />
              <Button
                size="sm"
                disabled={!scenarioName.trim()}
                onClick={() => {
                  if (whatIfMutation.data) {
                    const modList = buildModList();
                    saveScenario({
                      name: scenarioName.trim(),
                      explanationId: explanation.id,
                      modifications: modList,
                      result: whatIfMutation.data,
                    });
                    setShowSaveDialog(false);
                    setScenarioName('');
                    onScenarioSaved?.();
                  }
                }}
              >
                Save
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowSaveDialog(false)}
              >
                Cancel
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {whatIfMutation.data && (
        <>
          <ComparisonView original={explanation} result={whatIfMutation.data} />
          <SensitivityRanking ranking={whatIfMutation.data.sensitivity_ranking} />
        </>
      )}
    </div>
  );
}
