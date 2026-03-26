'use client';

import { useState, useMemo, useCallback } from 'react';
import { Copy, Check, Send, Terminal } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardAction } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from '@/lib/utils';

const ENDPOINTS = [
  { key: 'POST /api/v1/explain', method: 'POST', path: '/api/v1/explain' },
  { key: 'GET /api/v1/explain/{id}', method: 'GET', path: '/api/v1/explain/{id}' },
  { key: 'GET /api/v1/explain/{id}/graph', method: 'GET', path: '/api/v1/explain/{id}/graph' },
  { key: 'GET /api/v1/explain/{id}/narrative', method: 'GET', path: '/api/v1/explain/{id}/narrative' },
  { key: 'POST /api/v1/explain/{id}/what-if', method: 'POST', path: '/api/v1/explain/{id}/what-if' },
  { key: 'GET /health', method: 'GET', path: '/health' },
] as const;

type Endpoint = (typeof ENDPOINTS)[number];

const examplePayloads: Record<string, string> = {
  'POST /api/v1/explain': JSON.stringify(
    {
      target: 'market_regime_score',
      value: 0.72,
      components: [
        { name: 'trend_strength', value: 0.8, weight: 0.4, confidence: 0.9 },
        { name: 'volatility', value: 0.5, weight: 0.3, confidence: 0.7 },
        { name: 'momentum', value: 0.6, weight: 0.3, confidence: 0.85 },
      ],
    },
    null,
    2,
  ),
  'POST /api/v1/explain/{id}/what-if': JSON.stringify(
    {
      modifications: [{ component: 'trend_strength', new_value: 0.95 }],
    },
    null,
    2,
  ),
};

function endpointNeedsId(key: string): boolean {
  return key.includes('{id}');
}

function endpointHasBody(key: string): boolean {
  return key.startsWith('POST');
}

function resolveApiPath(
  endpoint: Endpoint,
  pathParams: Record<string, string>,
  queryParams: Record<string, string>,
): string {
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

function resolveBffPath(
  endpoint: Endpoint,
  pathParams: Record<string, string>,
  queryParams: Record<string, string>,
): string {
  let bffPath: string = endpoint.path;
  if (bffPath === '/health') {
    bffPath = '/api/health';
  } else {
    bffPath = bffPath.replace('/api/v1/', '/api/');
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

const BASE_URL =
  'https://explainable-engine-516741092583.europe-west1.run.app';

// --- Code snippet generators ---

function generateCurl(method: string, url: string, body?: string): string {
  const lines = [`curl -X ${method} ${url}`];
  if (body) {
    lines[0] += ' \\';
    lines.push('  -H "Content-Type: application/json" \\');
    lines.push(`  -d '${body}'`);
  }
  return lines.join('\n');
}

function generatePython(method: string, url: string, body?: string): string {
  const lines = ['import requests', ''];
  if (method === 'POST' && body) {
    lines.push('response = requests.post(');
    lines.push(`    "${url}",`);
    try {
      const obj = JSON.parse(body);
      lines.push(`    json=${JSON.stringify(obj)}`);
    } catch {
      lines.push(`    json=${body}`);
    }
    lines.push(')');
  } else {
    lines.push(`response = requests.get("${url}")`);
  }
  lines.push('print(response.json())');
  return lines.join('\n');
}

function generateGo(method: string, url: string, body?: string): string {
  if (method === 'POST' && body) {
    return `package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
)

func main() {
    body := []byte(\`${body}\`)
    resp, _ := http.Post("${url}", "application/json", bytes.NewReader(body))
    defer resp.Body.Close()
    data, _ := io.ReadAll(resp.Body)
    fmt.Println(string(data))
}`;
  }
  return `package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, _ := http.Get("${url}")
    defer resp.Body.Close()
    data, _ := io.ReadAll(resp.Body)
    fmt.Println(string(data))
}`;
}

function generateJavaScript(
  method: string,
  url: string,
  body?: string,
): string {
  if (method === 'POST' && body) {
    let formatted: string;
    try {
      formatted = JSON.stringify(JSON.parse(body), null, 2);
    } catch {
      formatted = body;
    }
    return `const response = await fetch("${url}", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(${formatted})
});
const data = await response.json();
console.log(data);`;
  }
  return `const response = await fetch("${url}");
const data = await response.json();
console.log(data);`;
}

const LANGUAGES = [
  { key: 'curl', label: 'curl' },
  { key: 'python', label: 'Python' },
  { key: 'go', label: 'Go' },
  { key: 'javascript', label: 'JavaScript' },
] as const;

type LanguageKey = (typeof LANGUAGES)[number]['key'];

// --- Status badge ---

function StatusBadge({ status }: { status: number }) {
  if (status >= 200 && status < 300) {
    return (
      <Badge className="bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400">
        {status}
      </Badge>
    );
  }
  if (status >= 400 && status < 500) {
    return (
      <Badge className="bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400">
        {status}
      </Badge>
    );
  }
  return (
    <Badge className="bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-400">
      {status || 'ERR'}
    </Badge>
  );
}

export default function PlaygroundPage() {
  const [endpoint, setEndpoint] = useState('POST /api/v1/explain');
  const [pathParams, setPathParams] = useState<Record<string, string>>({});
  const [queryParams, setQueryParams] = useState<Record<string, string>>({});
  const [body, setBody] = useState(examplePayloads['POST /api/v1/explain'] ?? '');
  const [response, setResponse] = useState<{
    status: number;
    data: unknown;
    time: number;
  } | null>(null);
  const [loading, setLoading] = useState(false);
  const [bodyError, setBodyError] = useState<string | null>(null);
  const [copiedResponse, setCopiedResponse] = useState(false);
  const [copiedSnippet, setCopiedSnippet] = useState<string | null>(null);

  const currentEndpoint =
    ENDPOINTS.find((e) => e.key === endpoint) ?? ENDPOINTS[0];

  const resolvedPath = useMemo(
    () => resolveApiPath(currentEndpoint, pathParams, queryParams),
    [currentEndpoint, pathParams, queryParams],
  );

  const snippetBody =
    endpointHasBody(endpoint) && body.trim() ? body : undefined;
  const url = `${BASE_URL}${resolvedPath}`;

  const snippets = useMemo<Record<LanguageKey, string>>(
    () => ({
      curl: generateCurl(currentEndpoint.method, url, snippetBody),
      python: generatePython(currentEndpoint.method, url, snippetBody),
      go: generateGo(currentEndpoint.method, url, snippetBody),
      javascript: generateJavaScript(currentEndpoint.method, url, snippetBody),
    }),
    [currentEndpoint.method, url, snippetBody],
  );

  const handleEndpointChange = useCallback(
    (value: string | null) => {
      if (!value) return;
      setEndpoint(value);
      setPathParams({});
      setQueryParams({});
      if (examplePayloads[value]) {
        setBody(examplePayloads[value]);
      } else {
        setBody('');
      }
      setResponse(null);
      setBodyError(null);
    },
    [],
  );

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
        data: {
          error: err instanceof Error ? err.message : 'Network error',
        },
        time: elapsed,
      });
    } finally {
      setLoading(false);
    }
  }, [endpoint, pathParams, queryParams, body, currentEndpoint]);

  const handleCopyResponse = useCallback(async () => {
    if (!response) return;
    await navigator.clipboard.writeText(
      JSON.stringify(response.data, null, 2),
    );
    setCopiedResponse(true);
    setTimeout(() => setCopiedResponse(false), 2000);
  }, [response]);

  const handleCopySnippet = useCallback(
    async (lang: string) => {
      await navigator.clipboard.writeText(snippets[lang as LanguageKey]);
      setCopiedSnippet(lang);
      setTimeout(() => setCopiedSnippet(null), 2000);
    },
    [snippets],
  );

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">API Playground</h1>
        <p className="text-sm text-muted-foreground">
          Test API endpoints and generate code snippets
        </p>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        {/* Left: Request + Response */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Request</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Endpoint selector */}
              <div className="space-y-2">
                <label className="text-sm font-medium">Endpoint</label>
                <div className="flex items-center gap-2">
                  <span
                    className={cn(
                      'inline-flex items-center rounded-md px-2 py-0.5 text-xs font-semibold shrink-0',
                      currentEndpoint.method === 'POST'
                        ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
                        : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400',
                    )}
                  >
                    {currentEndpoint.method}
                  </span>
                  <Select
                    value={endpoint}
                    onValueChange={handleEndpointChange}
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
                    onChange={(e) =>
                      setPathParams({ ...pathParams, id: e.target.value })
                    }
                  />
                </div>
              )}

              {/* Query parameters for graph endpoint */}
              {endpoint === 'GET /api/v1/explain/{id}/graph' && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">Format</label>
                  <Select
                    value={queryParams.format ?? 'json'}
                    onValueChange={(v) =>
                      setQueryParams({
                        ...queryParams,
                        format: v ?? 'json',
                      })
                    }
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
                      onValueChange={(v) =>
                        setQueryParams({
                          ...queryParams,
                          level: v ?? 'basic',
                        })
                      }
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
                      onValueChange={(v) =>
                        setQueryParams({ ...queryParams, lang: v ?? 'en' })
                      }
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
                      'h-48 w-full rounded-lg border px-3 py-2 font-mono text-sm outline-none transition-colors',
                      'bg-slate-950 text-slate-50',
                      'focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50',
                      bodyError ? 'border-destructive' : 'border-slate-700',
                    )}
                    value={body}
                    onChange={(e) => setBody(e.target.value)}
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
                disabled={
                  loading ||
                  (endpointNeedsId(endpoint) && !pathParams.id)
                }
                className="w-full"
              >
                {loading ? (
                  'Sending...'
                ) : (
                  <>
                    <Send className="size-4" />
                    Send Request
                  </>
                )}
              </Button>
            </CardContent>
          </Card>

          {/* Response */}
          {response && (
            <Card>
              <CardHeader>
                <div className="flex items-center gap-3">
                  <CardTitle>Response</CardTitle>
                  <StatusBadge status={response.status} />
                  <span className="text-xs text-muted-foreground tabular-nums">
                    {response.time}ms
                  </span>
                </div>
                <CardAction>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleCopyResponse}
                  >
                    {copiedResponse ? (
                      <>
                        <Check className="size-3.5" />
                        Copied
                      </>
                    ) : (
                      <>
                        <Copy className="size-3.5" />
                        Copy
                      </>
                    )}
                  </Button>
                </CardAction>
              </CardHeader>
              <CardContent>
                <pre className="max-h-96 overflow-auto rounded-lg bg-slate-950 p-4 font-mono text-xs text-slate-50">
                  {JSON.stringify(response.data, null, 2)}
                </pre>
              </CardContent>
            </Card>
          )}

          {!response && !loading && (
            <div className="flex flex-col items-center justify-center py-16 space-y-4 rounded-xl border border-dashed">
              <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
                <Terminal className="h-6 w-6 text-muted-foreground" />
              </div>
              <div className="text-center space-y-1">
                <h3 className="text-base font-medium">No response yet</h3>
                <p className="text-sm text-muted-foreground max-w-sm">
                  Select an endpoint and click &ldquo;Send Request&rdquo; to see
                  the response.
                </p>
              </div>
            </div>
          )}
        </div>

        {/* Right: Code snippets */}
        <Card>
          <CardHeader>
            <CardTitle>Code Snippets</CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs defaultValue="curl">
              <TabsList>
                {LANGUAGES.map((lang) => (
                  <TabsTrigger key={lang.key} value={lang.key}>
                    {lang.label}
                  </TabsTrigger>
                ))}
              </TabsList>
              {LANGUAGES.map((lang) => (
                <TabsContent key={lang.key} value={lang.key}>
                  <div className="relative mt-2">
                    <Button
                      variant="outline"
                      size="xs"
                      className="absolute right-2 top-2 z-10"
                      onClick={() => handleCopySnippet(lang.key)}
                    >
                      {copiedSnippet === lang.key ? (
                        <>
                          <Check className="size-3" />
                          Copied
                        </>
                      ) : (
                        <>
                          <Copy className="size-3" />
                          Copy
                        </>
                      )}
                    </Button>
                    <pre className="overflow-auto rounded-lg bg-slate-950 p-4 font-mono text-xs text-slate-50">
                      {snippets[lang.key]}
                    </pre>
                  </div>
                </TabsContent>
              ))}
            </Tabs>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
