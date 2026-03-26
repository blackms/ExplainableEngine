'use client';

import { useState, useCallback } from 'react';
import { useGenerateSummary } from '@/lib/api/hooks';
import type { SummaryRequest, SummaryResponse } from '@/lib/api/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import {
  FileTextIcon,
  BuildingIcon,
  WrenchIcon,
  UsersIcon,
  CheckCircleIcon,
  AlertTriangleIcon,
  ArrowRightIcon,
  CopyIcon,
} from 'lucide-react';

interface SummaryGeneratorProps {
  explanationId: string;
}

const AUDIENCES: {
  value: SummaryRequest['audience'];
  label: string;
  description: string;
  icon: React.ReactNode;
}[] = [
  {
    value: 'board',
    label: 'Board of Directors',
    description: 'High-level business impact',
    icon: <BuildingIcon className="size-5" />,
  },
  {
    value: 'technical',
    label: 'Technical Team',
    description: 'Detailed analysis',
    icon: <WrenchIcon className="size-5" />,
  },
  {
    value: 'client',
    label: 'Client Report',
    description: 'Clear, trustworthy explanation',
    icon: <UsersIcon className="size-5" />,
  },
];

const LANGUAGES = [
  { code: 'en' as const, label: 'EN' },
  { code: 'it' as const, label: 'IT' },
];

function SummarySkeleton() {
  return (
    <div className="space-y-4 pt-2">
      <Skeleton className="h-6 w-2/3" />
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-5/6" />
      <div className="space-y-2">
        <Skeleton className="h-4 w-1/4" />
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-3/4" />
        <Skeleton className="h-3 w-5/6" />
      </div>
      <div className="space-y-2">
        <Skeleton className="h-4 w-1/4" />
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-2/3" />
      </div>
      <div className="space-y-2">
        <Skeleton className="h-4 w-1/3" />
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-4/5" />
      </div>
    </div>
  );
}

function SummaryResult({ data }: { data: SummaryResponse }) {
  return (
    <div className="space-y-4 pt-2">
      <h3 className="text-lg font-semibold">{data.title}</h3>

      <p className="text-sm text-muted-foreground leading-relaxed">
        {data.summary}
      </p>

      {data.key_findings.length > 0 && (
        <div className="space-y-1.5">
          <h4 className="text-sm font-medium">Key Findings</h4>
          <ul className="space-y-1">
            {data.key_findings.map((finding, i) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <CheckCircleIcon className="mt-0.5 size-3.5 shrink-0 text-green-600 dark:text-green-400" />
                <span>{finding}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {data.risks.length > 0 && (
        <div className="space-y-1.5">
          <h4 className="text-sm font-medium">Risks</h4>
          <ul className="space-y-1">
            {data.risks.map((risk, i) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <AlertTriangleIcon className="mt-0.5 size-3.5 shrink-0 text-amber-600 dark:text-amber-400" />
                <span>{risk}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {data.recommendations.length > 0 && (
        <div className="space-y-1.5">
          <h4 className="text-sm font-medium">Recommendations</h4>
          <ul className="space-y-1">
            {data.recommendations.map((rec, i) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <ArrowRightIcon className="mt-0.5 size-3.5 shrink-0 text-blue-600 dark:text-blue-400" />
                <span>{rec}</span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export function SummaryGenerator({ explanationId }: SummaryGeneratorProps) {
  const [audience, setAudience] =
    useState<SummaryRequest['audience']>('board');
  const [lang, setLang] = useState<SummaryRequest['lang']>('en');
  const [copied, setCopied] = useState(false);

  const mutation = useGenerateSummary(explanationId);

  const handleGenerate = useCallback(() => {
    mutation.mutate({ audience, lang });
  }, [mutation, audience, lang]);

  const handleCopy = useCallback(async () => {
    if (!mutation.data) return;
    const d = mutation.data;
    const text = [
      d.title,
      '',
      d.summary,
      '',
      'Key Findings:',
      ...d.key_findings.map((f) => `  - ${f}`),
      '',
      'Risks:',
      ...d.risks.map((r) => `  - ${r}`),
      '',
      'Recommendations:',
      ...d.recommendations.map((r) => `  - ${r}`),
    ].join('\n');
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API may not be available
    }
  }, [mutation.data]);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between gap-2">
          <CardTitle className="flex items-center gap-2">
            <FileTextIcon className="size-4 text-emerald-500" />
            Executive Summary
          </CardTitle>
          <div className="flex items-center gap-1.5">
            {LANGUAGES.map(({ code, label }) => (
              <Button
                key={code}
                variant={lang === code ? 'default' : 'outline'}
                size="xs"
                onClick={() => setLang(code)}
              >
                {label}
              </Button>
            ))}
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Audience selector */}
        <div className="grid grid-cols-3 gap-2">
          {AUDIENCES.map((a) => (
            <button
              key={a.value}
              onClick={() => setAudience(a.value)}
              className={`flex flex-col items-center gap-1.5 rounded-lg border p-3 text-center transition-colors ${
                audience === a.value
                  ? 'border-primary bg-primary/5 text-foreground'
                  : 'border-border bg-background text-muted-foreground hover:border-primary/50 hover:text-foreground'
              }`}
            >
              {a.icon}
              <span className="text-xs font-medium">{a.label}</span>
              <span className="text-[10px] leading-tight opacity-70">
                {a.description}
              </span>
            </button>
          ))}
        </div>

        {/* Result area */}
        {mutation.isPending ? (
          <SummarySkeleton />
        ) : mutation.isError ? (
          <div className="rounded-md bg-red-500/10 px-3 py-2 text-sm text-red-700 dark:text-red-400">
            Failed to generate summary. Please try again.
          </div>
        ) : mutation.data ? (
          <SummaryResult data={mutation.data} />
        ) : null}

        {/* Actions */}
        <div className="flex items-center justify-between gap-2">
          <Button
            variant="default"
            size="sm"
            onClick={handleGenerate}
            disabled={mutation.isPending}
          >
            <FileTextIcon className="size-3.5" data-icon="inline-start" />
            {mutation.data ? 'Regenerate' : 'Generate Summary'}
          </Button>

          {mutation.data && (
            <Button variant="outline" size="sm" onClick={handleCopy}>
              <CopyIcon className="size-3.5" data-icon="inline-start" />
              {copied ? 'Copied!' : 'Copy'}
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
