'use client';

import { useState, useMemo, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';

interface CodeSnippetProps {
  method: string;
  path: string;
  body?: string;
  baseUrl: string;
}

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
    lines.push(`response = requests.post(`);
    lines.push(`    "${url}",`);
    lines.push(`    json=${pythonDict(body)}`);
    lines.push(`)`);
  } else {
    lines.push(`response = requests.get("${url}")`);
  }
  lines.push('print(response.json())');
  return lines.join('\n');
}

function pythonDict(jsonStr: string): string {
  try {
    const obj = JSON.parse(jsonStr);
    return formatPythonValue(obj, 0);
  } catch {
    return jsonStr;
  }
}

function formatPythonValue(value: unknown, indent: number): string {
  const pad = '    '.repeat(indent);
  const innerPad = '    '.repeat(indent + 1);

  if (value === null) return 'None';
  if (typeof value === 'boolean') return value ? 'True' : 'False';
  if (typeof value === 'number') return String(value);
  if (typeof value === 'string') return `"${value}"`;

  if (Array.isArray(value)) {
    if (value.length === 0) return '[]';
    const items = value.map((v) => `${innerPad}${formatPythonValue(v, indent + 1)}`);
    return `[\n${items.join(',\n')}\n${pad}]`;
  }

  if (typeof value === 'object' && value !== null) {
    const entries = Object.entries(value as Record<string, unknown>);
    if (entries.length === 0) return '{}';
    const items = entries.map(
      ([k, v]) => `${innerPad}"${k}": ${formatPythonValue(v, indent + 1)}`
    );
    return `{\n${items.join(',\n')}\n${pad}}`;
  }

  return String(value);
}

function generateGo(method: string, url: string, body?: string): string {
  if (method === 'POST' && body) {
    return `package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

func main() {
    body, _ := json.Marshal(${goValue(body)})
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

function goValue(jsonStr: string): string {
  try {
    const obj = JSON.parse(jsonStr);
    return formatGoValue(obj, 1);
  } catch {
    return `[]byte(\`${jsonStr}\`)`;
  }
}

function formatGoValue(value: unknown, indent: number): string {
  const pad = '\t'.repeat(indent);
  const innerPad = '\t'.repeat(indent + 1);

  if (value === null) return 'nil';
  if (typeof value === 'boolean') return value ? 'true' : 'false';
  if (typeof value === 'number') {
    if (Number.isInteger(value)) return String(value);
    return String(value);
  }
  if (typeof value === 'string') return `"${value}"`;

  if (Array.isArray(value)) {
    if (value.length === 0) return '[]any{}';
    const items = value.map((v) => `${innerPad}${formatGoValue(v, indent + 1)},`);
    return `[]any{\n${items.join('\n')}\n${pad}}`;
  }

  if (typeof value === 'object' && value !== null) {
    const entries = Object.entries(value as Record<string, unknown>);
    if (entries.length === 0) return 'map[string]any{}';
    const items = entries.map(
      ([k, v]) => `${innerPad}"${k}": ${formatGoValue(v, indent + 1)},`
    );
    return `map[string]any{\n${items.join('\n')}\n${pad}}`;
  }

  return String(value);
}

function generateJavaScript(method: string, url: string, body?: string): string {
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

type LanguageKey = typeof LANGUAGES[number]['key'];

export function CodeSnippet({ method, path, body, baseUrl }: CodeSnippetProps) {
  const [copiedLang, setCopiedLang] = useState<string | null>(null);

  const url = `${baseUrl}${path}`;

  const snippets = useMemo<Record<LanguageKey, string>>(() => ({
    curl: generateCurl(method, url, body),
    python: generatePython(method, url, body),
    go: generateGo(method, url, body),
    javascript: generateJavaScript(method, url, body),
  }), [method, url, body]);

  const handleCopy = useCallback(async (lang: string) => {
    await navigator.clipboard.writeText(snippets[lang as LanguageKey]);
    setCopiedLang(lang);
    setTimeout(() => setCopiedLang(null), 2000);
  }, [snippets]);

  return (
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
                  className="absolute right-2 top-2"
                  onClick={() => handleCopy(lang.key)}
                >
                  {copiedLang === lang.key ? 'Copied!' : 'Copy'}
                </Button>
                <pre className="overflow-auto rounded-lg bg-zinc-950 p-4 font-mono text-xs text-zinc-100 dark:bg-zinc-900">
                  {snippets[lang.key]}
                </pre>
              </div>
            </TabsContent>
          ))}
        </Tabs>
      </CardContent>
    </Card>
  );
}
