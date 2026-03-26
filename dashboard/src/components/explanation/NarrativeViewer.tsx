'use client';

import { useState, useCallback } from 'react';
import { useNarrative } from '@/lib/api/hooks';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';

interface NarrativeViewerProps {
  explanationId: string;
}

const LEVELS = ['basic', 'advanced'] as const;
const LANGUAGES = [
  { code: 'en', label: 'EN' },
  { code: 'it', label: 'IT' },
] as const;

export function NarrativeViewer({ explanationId }: NarrativeViewerProps) {
  const [level, setLevel] = useState<string>('basic');
  const [lang, setLang] = useState<string>('en');
  const [copied, setCopied] = useState(false);

  const { data, isLoading, isError } = useNarrative(explanationId, level, lang);

  const handleCopy = useCallback(async () => {
    if (!data?.narrative) return;
    try {
      await navigator.clipboard.writeText(data.narrative);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API may not be available
    }
  }, [data?.narrative]);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between gap-2">
          <CardTitle>Narrative Explanation</CardTitle>
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
        <Tabs defaultValue="basic" onValueChange={(value) => setLevel(value as string)}>
          <TabsList>
            {LEVELS.map((l) => (
              <TabsTrigger key={l} value={l}>
                {l.charAt(0).toUpperCase() + l.slice(1)}
              </TabsTrigger>
            ))}
          </TabsList>

          {LEVELS.map((l) => (
            <TabsContent key={l} value={l}>
              {isLoading ? (
                <div className="space-y-2 pt-2">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-4 w-5/6" />
                </div>
              ) : isError ? (
                <div className="rounded-md bg-red-500/10 px-3 py-2 text-sm text-red-700 dark:text-red-400">
                  Failed to load narrative. Please try again.
                </div>
              ) : (
                <div className="prose prose-sm dark:prose-invert max-w-none pt-2">
                  {l === 'advanced' ? (
                    data?.narrative.split('\n').map((line, i) => (
                      <p key={i} className="mb-1.5 last:mb-0">
                        {line}
                      </p>
                    ))
                  ) : (
                    <p>{data?.narrative}</p>
                  )}
                </div>
              )}
            </TabsContent>
          ))}
        </Tabs>

        {data?.narrative && (
          <div className="flex justify-end">
            <Button variant="outline" size="sm" onClick={handleCopy}>
              {copied ? 'Copied!' : 'Copy text'}
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
