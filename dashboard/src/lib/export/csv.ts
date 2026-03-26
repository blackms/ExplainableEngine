import type { ExplainResponse } from '@/lib/api/types';

export function exportListToCSV(items: ExplainResponse[], filename: string = 'explanations.csv') {
  const headers = ['ID', 'Target', 'Value', 'Confidence', 'Missing Impact', 'Top Driver', 'Created At'];
  const rows = items.map(item => [
    item.id,
    item.target,
    item.final_value.toFixed(4),
    item.confidence.toFixed(4),
    item.missing_impact.toFixed(4),
    item.top_drivers[0]?.name || 'N/A',
    item.metadata.created_at,
  ]);

  const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
  downloadFile(csv, filename, 'text/csv');
}

function downloadFile(content: string, filename: string, mimeType: string) {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
