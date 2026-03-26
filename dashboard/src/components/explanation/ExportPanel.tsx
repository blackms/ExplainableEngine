'use client';
import { Button } from '@/components/ui/button';
import { exportExplanationToPDF } from '@/lib/export/pdf';
import type { ExplainResponse } from '@/lib/api/types';

interface ExportPanelProps {
  explanation: ExplainResponse;
  narrative?: string;
}

export function ExportPanel({ explanation, narrative }: ExportPanelProps) {
  return (
    <div className="flex gap-2">
      <Button variant="outline" size="sm" onClick={() => exportExplanationToPDF(explanation, narrative)}>
        Export PDF
      </Button>
    </div>
  );
}
