'use client';

import { useEffect, useState } from 'react';
import {
  generateApiKey,
  getApiKeys,
  revokeApiKey,
  type ApiKey,
} from '@/lib/apikeys';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

export function ApiKeyManager() {
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [name, setName] = useState('');
  const [newKey, setNewKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [confirmRevoke, setConfirmRevoke] = useState<string | null>(null);

  useEffect(() => {
    setKeys(getApiKeys());
  }, []);

  const handleGenerate = () => {
    if (!name.trim()) return;
    const key = generateApiKey(name.trim());
    setNewKey(key.key);
    setName('');
    setKeys(getApiKeys());
  };

  const handleCopy = async () => {
    if (!newKey) return;
    await navigator.clipboard.writeText(newKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleRevoke = (id: string) => {
    revokeApiKey(id);
    setKeys(getApiKeys());
    setConfirmRevoke(null);
  };

  return (
    <div className="space-y-6">
      {/* Generate key section */}
      <Card>
        <CardHeader>
          <CardTitle>Generate New Key</CardTitle>
          <CardDescription>
            Create a new API key for programmatic access.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-end gap-3">
            <div className="grid flex-1 gap-2">
              <Label htmlFor="key-name">Key Name</Label>
              <Input
                id="key-name"
                placeholder="e.g. CI Pipeline"
                value={name}
                onChange={(e) => setName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleGenerate()}
              />
            </div>
            <Button onClick={handleGenerate} disabled={!name.trim()}>
              Generate
            </Button>
          </div>

          {newKey && (
            <div className="mt-4 rounded-lg border border-yellow-500/50 bg-yellow-50 p-4 dark:bg-yellow-950/20">
              <p className="mb-2 text-sm font-medium text-yellow-800 dark:text-yellow-200">
                This key won&apos;t be shown again. Copy it now.
              </p>
              <div className="flex items-center gap-2">
                <code className="flex-1 break-all rounded bg-yellow-100 px-2 py-1 text-xs dark:bg-yellow-900/30">
                  {newKey}
                </code>
                <Button variant="outline" size="sm" onClick={handleCopy}>
                  {copied ? 'Copied' : 'Copy'}
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Key list */}
      {keys.length === 0 ? (
        <p className="text-sm text-muted-foreground">No API keys yet.</p>
      ) : (
        <div className="overflow-hidden rounded-lg border">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-muted/50">
                <th className="px-4 py-2 text-left font-medium">Name</th>
                <th className="px-4 py-2 text-left font-medium">Key</th>
                <th className="px-4 py-2 text-left font-medium">Created</th>
                <th className="px-4 py-2 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {keys.map((k) => (
                <tr key={k.id} className="border-b last:border-b-0">
                  <td className="px-4 py-2">{k.name}</td>
                  <td className="px-4 py-2">
                    <code className="text-xs text-muted-foreground">
                      {k.prefix}
                    </code>
                  </td>
                  <td className="px-4 py-2 text-muted-foreground">
                    {new Date(k.createdAt).toLocaleDateString()}
                  </td>
                  <td className="px-4 py-2 text-right">
                    {confirmRevoke === k.id ? (
                      <span className="flex items-center justify-end gap-2">
                        <span className="text-xs text-muted-foreground">
                          Confirm?
                        </span>
                        <Button
                          variant="destructive"
                          size="xs"
                          onClick={() => handleRevoke(k.id)}
                        >
                          Yes
                        </Button>
                        <Button
                          variant="outline"
                          size="xs"
                          onClick={() => setConfirmRevoke(null)}
                        >
                          No
                        </Button>
                      </span>
                    ) : (
                      <Button
                        variant="destructive"
                        size="xs"
                        onClick={() => setConfirmRevoke(k.id)}
                      >
                        Revoke
                      </Button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
