'use client';

import { useState, useMemo } from 'react';
import { ApiExplorer, ENDPOINTS, examplePayloads, resolveApiPath, endpointHasBody } from '@/components/playground/ApiExplorer';
import { CodeSnippet } from '@/components/playground/CodeSnippet';

const BASE_URL = 'https://explainable-engine-516741092583.europe-west1.run.app';

export default function PlaygroundPage() {
  const [endpoint, setEndpoint] = useState('POST /api/v1/explain');
  const [pathParams, setPathParams] = useState<Record<string, string>>({});
  const [queryParams, setQueryParams] = useState<Record<string, string>>({});
  const [body, setBody] = useState(examplePayloads['POST /api/v1/explain'] ?? '');

  const currentEndpoint = ENDPOINTS.find((e) => e.key === endpoint) ?? ENDPOINTS[0];

  const resolvedPath = useMemo(
    () => resolveApiPath(currentEndpoint, pathParams, queryParams),
    [currentEndpoint, pathParams, queryParams]
  );

  const snippetBody = endpointHasBody(endpoint) && body.trim() ? body : undefined;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">API Playground</h1>
        <p className="text-muted-foreground">
          Test API endpoints and generate code snippets
        </p>
      </div>
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <ApiExplorer
          endpoint={endpoint}
          pathParams={pathParams}
          queryParams={queryParams}
          body={body}
          onEndpointChange={setEndpoint}
          onPathParamsChange={setPathParams}
          onQueryParamsChange={setQueryParams}
          onBodyChange={setBody}
        />
        <CodeSnippet
          method={currentEndpoint.method}
          path={resolvedPath}
          body={snippetBody}
          baseUrl={BASE_URL}
        />
      </div>
    </div>
  );
}
