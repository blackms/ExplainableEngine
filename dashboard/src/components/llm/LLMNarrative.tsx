'use client';

import { useState, useCallback } from 'react';
import { useLLMNarrative } from '@/lib/api/hooks';
import type { LLMNarrativeRequest } from '@/lib/api/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { SparklesIcon, CopyIcon, RefreshCwIcon } from 'lucide-react';

interface LLMNarrativeProps {
  explanationId: string;
}

const LEVELS = ['basic', 'advanced', 'executive'] as const;
const LANGUAGES = [
  { code: 'en' as const, label: 'EN' },
  { code: 'it' as const, label: 'IT' },
];

export function LLMNarrative({ explanationId }: LLMNarrativeProps) {
  const [level, setLevel] = useState<LLMNarrativeRequest['level']>('basic');
  const [lang, setLang] = useState<LLMNarrativeRequest['lang']>('en');
  const [copied, setCopied] = useState(false);

  const mutation = useLLMNarrative(explanationId);

  const handleGenerate = useCallback(() => {
    mutation.mutate({ level, lang });
  }, [mutation, level, lang]);

  const handleCopy = useCallback(async () => {
    if (!mutation.data?.narrative) return;
    try {
      await navigator.clipboard.writeText(mutation.data.narrative);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API may not be available
    }
  }, [mutation.data?.narrative]);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between gap-2">
          <CardTitle className="flex items-center gap-2">
            <SparklesIcon className="size-4 text-purple-500" />
            AI Narrative
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
        <Tabs
          defaultValue="basic"
          onValueChange={(value) =>
            setLevel(value as LLMNarrativeRequest['level'])
          }
        >
          <TabsList>
            {LEVELS.map((l) => (
              <TabsTrigger key={l} value={l}>
                {l.charAt(0).toUpperCase() + l.slice(1)}
              </TabsTrigger>
            ))}
          </TabsList>

          {LEVELS.map((l) => (
            <TabsContent key={l} value={l}>
              {mutation.isPending ? (
                <div className="space-y-2 pt-2">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <span className="inline-block animate-pulse">
                      Generating
                    </span>
                    <span className="inline-flex gap-0.5">
                      <span className="animate-bounce [animation-delay:0ms]">
                        .
                      </span>
                      <span className="animate-bounce [animation-delay:150ms]">
                        .
                      </span>
                      <span className="animate-bounce [animation-delay:300ms]">
                        .
                      </span>
                    </span>
                  </div>
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-4 w-5/6" />
                </div>
              ) : mutation.isError ? (
                <div className="rounded-md bg-red-500/10 px-3 py-2 text-sm text-red-700 dark:text-red-400">
                  Failed to generate narrative. Please try again.
                </div>
              ) : mutation.data ? (
                <div className="space-y-3 pt-2">
                  <div className="flex items-center gap-2">
                    <span
                      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                        mutation.data.source === 'llm'
                          ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'
                          : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'
                      }`}
                    >
                      {mutation.data.source === 'llm'
                        ? 'AI Generated'
                        : 'Template'}
                    </span>
                    {mutation.data.model && (
                      <span className="text-xs text-muted-foreground">
                        {mutation.data.model}
                      </span>
                    )}
                  </div>
                  <div className="prose prose-sm dark:prose-invert max-w-none">
                    {mutation.data.narrative.split('\n').map((line, i) => (
                      <p key={i} className="mb-1.5 last:mb-0">
                        {line}
                      </p>
                    ))}
                  </div>
                </div>
              ) : (
                <div className="pt-2 text-center">
                  <p className="mb-3 text-sm text-muted-foreground">
                    Generate an AI-powered narrative explanation at the{' '}
                    <span className="font-medium">{l}</span> level.
                  </p>
                </div>
              )}
            </TabsContent>
          ))}
        </Tabs>

        <div className="flex items-center justify-between gap-2">
          <Button
            variant="default"
            size="sm"
            onClick={handleGenerate}
            disabled={mutation.isPending}
          >
            <SparklesIcon className="size-3.5" data-icon="inline-start" />
            {mutation.data ? 'Regenerate' : 'Generate AI Narrative'}
          </Button>

          {mutation.data?.narrative && (
            <div className="flex items-center gap-1.5">
              <Button variant="outline" size="sm" onClick={handleCopy}>
                <CopyIcon className="size-3.5" data-icon="inline-start" />
                {copied ? 'Copied!' : 'Copy'}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleGenerate}
                disabled={mutation.isPending}
              >
                <RefreshCwIcon className="size-3.5" data-icon="inline-start" />
                Regenerate
              </Button>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
