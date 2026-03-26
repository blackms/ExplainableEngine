'use client';

import { useState, useCallback } from 'react';
import type { ExplainResponse } from '@/lib/api/types';
import type { SavedScenario } from '@/lib/scenarios';
import { WhatIfSimulator } from '@/components/whatif/WhatIfSimulator';
import { ScenarioManager } from '@/components/whatif/ScenarioManager';
import { ScenarioComparison } from '@/components/whatif/ScenarioComparison';

interface WhatIfPageClientProps {
  explanation: ExplainResponse;
}

export function WhatIfPageClient({ explanation }: WhatIfPageClientProps) {
  const [initialModifications, setInitialModifications] = useState<Record<string, number> | undefined>();
  const [comparisonScenarios, setComparisonScenarios] = useState<SavedScenario[]>([]);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleLoad = useCallback((scenario: SavedScenario) => {
    const mods: Record<string, number> = {};
    for (const mod of scenario.modifications) {
      mods[mod.component] = mod.new_value;
    }
    setInitialModifications(mods);
  }, []);

  const handleSelect = useCallback((scenarios: SavedScenario[]) => {
    setComparisonScenarios(scenarios);
  }, []);

  const handleScenarioSaved = useCallback(() => {
    setRefreshKey((k) => k + 1);
  }, []);

  return (
    <div className="space-y-6">
      <WhatIfSimulator
        explanation={explanation}
        initialModifications={initialModifications}
        onScenarioSaved={handleScenarioSaved}
      />
      <ScenarioManager
        explanationId={explanation.id}
        onSelect={handleSelect}
        onLoad={handleLoad}
        refreshKey={refreshKey}
      />
      {comparisonScenarios.length >= 2 && (
        <ScenarioComparison
          scenarios={comparisonScenarios}
          originalValue={explanation.final_value}
        />
      )}
    </div>
  );
}
