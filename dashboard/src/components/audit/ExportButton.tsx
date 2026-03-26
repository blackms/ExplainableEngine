'use client';
import { Button } from '@/components/ui/button';
import { exportListToCSV } from '@/lib/export/csv';
import type { ExplainResponse } from '@/lib/api/types';

interface ExportButtonProps {
  items: ExplainResponse[];
}

export function ExportButton({ items }: ExportButtonProps) {
  return (
    <Button variant="outline" size="sm" onClick={() => exportListToCSV(items)}>
      Export CSV
    </Button>
  );
}
