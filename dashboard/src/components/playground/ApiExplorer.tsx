'use client';

import { useState, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

const ENDPOINTS = [
  { key: 'POST /api/v1/explain', method: 'POST', path: '/api/v1/explain' },
  { key: 'GET /api/v1/explain/{id}', method: 'GET', path: '/api/v1/explain/{id}' },
  { key: 'GET /api/v1/explain/{id}/graph', method: 'GET', path: '/api/v1/explain/{id}/graph' },
  { key: 'GET /api/v1/explain/{id}/narrative', method: 'GET', path: '/api/v1/explain/{id}/narrative' },
  { key: 'POST /api/v1/explain/{id}/what-if', method: 'POST', path: '/api/v1/explain/{id}/what-if' },
  { key: 'GET /health', method: 'GET', path: '/health' },
] as const;

const examplePayloads: Record<string, string> = {
  'POST /api/v1/explain': JSON.stringify({
    target: 'market_regime_score',
    value: 0.72,
    components: [
      { name: 'trend_strength', value: 0.8, weight: 0.4, confidence: 0.9 },
      { name: 'volatility', value: 0.5, weight: 0.3, confidence: 0.7 },
      { name: 'momentum', value: 0.6, weight: 0.3, confidence: 0.85 },
    ],
  }, null, 2),
  'POST /api/v1/explain/{id}/what-if': JSON.stringify({
    modifications: [
      { component: 'trend_strength', new_value: 0.95 },
    ],
  }, null, 2),
};

function endpointNeedsId(key: string): boolean {
  return key.includes('{id}');
}

function endpointHasBody(key: string): boolean {
  return key.startsWith('POST');
}

function resolveApiPath(endpoint: typeof ENDPOINTS[number], pathParams: Record<string, string>, queryParams: Record<string, string>): string {
  let resolved: string = endpoint.path;
  if (pathParams.id) {
    resolved = resolved.replace('{id}', pathParams.id);
  }
  const qs = new URLSearchParams(queryParams).toString();
  if (qs) {
    resolved += `?${qs}`;
  }
  return resolved;
}

function resolveBffPath(endpoint: typeof ENDPOINTS[number], pathParams: Record<string, string>, queryParams: Record<string, string>): string {
  let bffPath: string = endpoint.path;

  // Map external API paths to internal BFF routes
  if (bffPath === '/health') {
    bffPath = '/api/health';
  } else {
    bffPath = bffPath.replace('/api/v1/', '/api/');
    // Handle what-if -> whatif mapping
    bffPath = bffPath.replace('/what-if', '/whatif');
  }

  if (pathParams.id) {
    bffPath = bffPath.replace('{id}', pathParams.id);
  }

  const qs = new URLSearchParams(queryParams).toString();
  if (qs) {
    bffPath += `?${qs}`;
  }
  return bffPath;
}

interface ApiExplorerProps {
  endpoint: string;
  pathParams: Record<string, string>;
  queryParams: Record<string, string>;
  body: string;
  onEndpointChange: (endpoint: string) => void;
  onPathParamsChange: (params: Record<string, string>) => void;
  onQueryParamsChange: (params: Record<string, string>) => void;
  onBodyChange: (body: string) => void;
}

export function ApiExplorer({
  endpoint,
  pathParams,
  queryParams,
  body,
  onEndpointChange,
  onPathParamsChange,
  onQueryParamsChange,
  onBodyChange,
}: ApiExplorerProps) {
  const [response, setResponse] = useState<{ status: number; data: unknown; time: number } | null>(null);
  const [loading, setLoading] = useState(false);
  const [bodyError, setBodyError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const currentEndpoint = ENDPOINTS.find((e) => e.key === endpoint) ?? ENDPOINTS[0];

  const handleEndpointChange = useCallback((value: string | null) => {
    if (!value) return;
    onEndpointChange(value);
    onPathParamsChange({});
    onQueryParamsChange({});
    if (examplePayloads[value]) {
      onBodyChange(examplePayloads[value]);
    } else {
      onBodyChange('');
    }
    setResponse(null);
    setBodyError(null);
  }, [onEndpointChange, onPathParamsChange, onQueryParamsChange, onBodyChange]);

  const handleBodyBlur = useCallback(() => {
    if (!endpointHasBody(endpoint) || !body.trim()) {
      setBodyError(null);
      return;
    }
    try {
      JSON.parse(body);
      setBodyError(null);
    } catch {
      setBodyError('Invalid JSON');
    }
  }, [endpoint, body]);

  const handleSend = useCallback(async () => {
    if (endpointNeedsId(endpoint) && !pathParams.id) return;
    if (endpointHasBody(endpoint) && body.trim()) {
      try {
        JSON.parse(body);
      } catch {
        setBodyError('Invalid JSON');
        return;
      }
    }

    setLoading(true);
    setResponse(null);

    const bffPath = resolveBffPath(currentEndpoint, pathParams, queryParams);
    const start = performance.now();

    try {
      const fetchOptions: RequestInit = {
        method: currentEndpoint.method,
        headers: { 'Content-Type': 'application/json' },
      };

      if (endpointHasBody(endpoint) && body.trim()) {
        fetchOptions.body = body;
      }

      const res = await fetch(bffPath, fetchOptions);
      const elapsed = Math.round(performance.now() - start);
      const data = await res.json();

      setResponse({ status: res.status, data, time: elapsed });
    } catch (err) {
      const elapsed = Math.round(performance.now() - start);
      setResponse({
        status: 0,
        data: { error: err instanceof Error ? err.message : 'Network error' },
        time: elapsed,
      });
    } finally {
      setLoading(false);
    }
  }, [endpoint, pathParams, queryParams, body, currentEndpoint]);

  const handleCopyResponse = useCallback(async () => {
    if (!response) return;
    await navigator.clipboard.writeText(JSON.stringify(response.data, null, 2));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [response]);

  const statusColor = response
    ? response.status >= 200 && response.status < 300
      ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
      : response.status >= 400 && response.status < 500
        ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
        : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
    : '';

  const endpointItems = ENDPOINTS.map((e) => ({
    value: e.key,
    label: e.key,
  }));

  return (
    <Card>
      <CardHeader>
        <CardTitle>API Explorer</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Endpoint selector */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Endpoint</label>
          <div className="flex items-center gap-2">
            <span
              className={cn(
                'inline-flex items-center rounded-md px-2 py-0.5 text-xs font-semibold',
                currentEndpoint.method === 'POST'
                  ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
                  : 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
              )}
            >
              {currentEndpoint.method}
            </span>
            <Select
              value={endpoint}
              onValueChange={handleEndpointChange}
              items={endpointItems}
            >
              <SelectTrigger className="flex-1">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {ENDPOINTS.map((e) => (
                  <SelectItem key={e.key} value={e.key}>
                    {e.key}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Path parameters */}
        {endpointNeedsId(endpoint) && (
          <div className="space-y-2">
            <label className="text-sm font-medium">Explanation ID</label>
            <Input
              placeholder="Enter explanation ID"
              value={pathParams.id ?? ''}
              onChange={(e) => onPathParamsChange({ ...pathParams, id: e.target.value })}
            />
          </div>
        )}

        {/* Query parameters for graph endpoint */}
        {endpoint === 'GET /api/v1/explain/{id}/graph' && (
          <div className="space-y-2">
            <label className="text-sm font-medium">Format</label>
            <Select
              value={queryParams.format ?? 'json'}
              onValueChange={(v) => onQueryParamsChange({ ...queryParams, format: v ?? 'json' })}
              items={[
                { value: 'json', label: 'json' },
                { value: 'dot', label: 'dot' },
                { value: 'mermaid', label: 'mermaid' },
              ]}
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="json">json</SelectItem>
                <SelectItem value="dot">dot</SelectItem>
                <SelectItem value="mermaid">mermaid</SelectItem>
              </SelectContent>
            </Select>
          </div>
        )}

        {/* Query parameters for narrative endpoint */}
        {endpoint === 'GET /api/v1/explain/{id}/narrative' && (
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Level</label>
              <Select
                value={queryParams.level ?? 'basic'}
                onValueChange={(v) => onQueryParamsChange({ ...queryParams, level: v ?? 'basic' })}
                items={[
                  { value: 'basic', label: 'basic' },
                  { value: 'advanced', label: 'advanced' },
                ]}
              >
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="basic">basic</SelectItem>
                  <SelectItem value="advanced">advanced</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Language</label>
              <Select
                value={queryParams.lang ?? 'en'}
                onValueChange={(v) => onQueryParamsChange({ ...queryParams, lang: v ?? 'en' })}
                items={[
                  { value: 'en', label: 'en' },
                  { value: 'it', label: 'it' },
                ]}
              >
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="en">en</SelectItem>
                  <SelectItem value="it">it</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        )}

        {/* Request body */}
        {endpointHasBody(endpoint) && (
          <div className="space-y-2">
            <label className="text-sm font-medium">Request Body</label>
            <textarea
              className={cn(
                'h-48 w-full rounded-lg border bg-transparent px-3 py-2 font-mono text-sm outline-none transition-colors',
                'focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50',
                'dark:bg-input/30',
                bodyError
                  ? 'border-destructive'
                  : 'border-input'
              )}
              value={body}
              onChange={(e) => onBodyChange(e.target.value)}
              onBlur={handleBodyBlur}
              spellCheck={false}
            />
            {bodyError && (
              <p className="text-xs text-destructive">{bodyError}</p>
            )}
          </div>
        )}

        {/* Send button */}
        <Button
          onClick={handleSend}
          disabled={loading || (endpointNeedsId(endpoint) && !pathParams.id)}
          className="w-full"
        >
          {loading ? 'Sending...' : 'Send Request'}
        </Button>

        {/* Response panel */}
        {response && (
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className={cn('inline-flex items-center rounded-md px-2 py-0.5 text-xs font-semibold', statusColor)}>
                  {response.status || 'ERR'}
                </span>
                <span className="text-xs text-muted-foreground">
                  {response.time}ms
                </span>
              </div>
              <Button variant="outline" size="xs" onClick={handleCopyResponse}>
                {copied ? 'Copied!' : 'Copy'}
              </Button>
            </div>
            <pre className="max-h-96 overflow-auto rounded-lg bg-muted p-4 font-mono text-xs">
              {JSON.stringify(response.data, null, 2)}
            </pre>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

export { ENDPOINTS, examplePayloads, resolveApiPath, endpointHasBody, endpointNeedsId };
