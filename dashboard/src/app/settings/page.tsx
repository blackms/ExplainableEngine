'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { Check, Copy, Key, Moon, Monitor, Sun } from 'lucide-react';
import {
  getPreferences,
  savePreferences,
  type UserPreferences,
} from '@/lib/preferences';
import {
  generateApiKey,
  getApiKeys,
  revokeApiKey,
  type ApiKey,
} from '@/lib/apikeys';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from '@/lib/utils';

// --- Theme card component ---

interface ThemeOptionProps {
  value: UserPreferences['theme'];
  label: string;
  icon: React.ReactNode;
  selected: boolean;
  onSelect: () => void;
}

function ThemeOption({ label, icon, selected, onSelect }: ThemeOptionProps) {
  return (
    <button
      type="button"
      onClick={onSelect}
      className={cn(
        'flex flex-col items-center gap-2 rounded-lg border-2 p-4 transition-colors cursor-pointer',
        selected
          ? 'border-primary bg-primary/5'
          : 'border-transparent bg-muted/50 hover:bg-muted',
      )}
    >
      <div
        className={cn(
          'flex h-10 w-10 items-center justify-center rounded-full',
          selected ? 'bg-primary text-primary-foreground' : 'bg-muted',
        )}
      >
        {icon}
      </div>
      <span className="text-sm font-medium">{label}</span>
      {selected && (
        <span className="inline-flex items-center gap-1 text-xs text-primary">
          <Check className="size-3" />
          Active
        </span>
      )}
    </button>
  );
}

export default function SettingsPage() {
  // --- Preferences state ---
  const [prefs, setPrefs] = useState<UserPreferences | null>(null);
  const [saved, setSaved] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // --- API Keys state ---
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [keyName, setKeyName] = useState('');
  const [newKey, setNewKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [confirmRevoke, setConfirmRevoke] = useState<string | null>(null);

  useEffect(() => {
    setPrefs(getPreferences());
    setKeys(getApiKeys());
  }, []);

  const handlePrefChange = useCallback(
    (partial: Partial<UserPreferences>) => {
      const updated = savePreferences(partial);
      setPrefs(updated);
      setSaved(true);
      if (timerRef.current) clearTimeout(timerRef.current);
      timerRef.current = setTimeout(() => setSaved(false), 2000);

      if (partial.theme) {
        const root = document.documentElement;
        if (partial.theme === 'dark') {
          root.classList.add('dark');
        } else if (partial.theme === 'light') {
          root.classList.remove('dark');
        } else {
          const isDark = window.matchMedia(
            '(prefers-color-scheme: dark)',
          ).matches;
          root.classList.toggle('dark', isDark);
        }
      }
    },
    [],
  );

  const handleGenerate = useCallback(() => {
    if (!keyName.trim()) return;
    const key = generateApiKey(keyName.trim());
    setNewKey(key.key);
    setKeyName('');
    setKeys(getApiKeys());
  }, [keyName]);

  const handleCopyKey = useCallback(async () => {
    if (!newKey) return;
    await navigator.clipboard.writeText(newKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [newKey]);

  const handleRevoke = useCallback((id: string) => {
    revokeApiKey(id);
    setKeys(getApiKeys());
    setConfirmRevoke(null);
  }, []);

  if (!prefs) return null;

  return (
    <div className="space-y-8 max-w-2xl mx-auto">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
        <p className="text-sm text-muted-foreground">
          Manage your preferences and API access
        </p>
      </div>

      {/* Preferences */}
      <section className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold">Preferences</h2>
          {saved && (
            <span className="inline-flex items-center gap-1 text-xs text-emerald-600 dark:text-emerald-400 animate-in fade-in">
              <Check className="size-3" />
              Saved
            </span>
          )}
        </div>

        {/* Theme selection - card-based */}
        <Card>
          <CardHeader>
            <CardTitle>Theme</CardTitle>
            <CardDescription>Choose your preferred appearance</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-3">
              <ThemeOption
                value="light"
                label="Light"
                icon={<Sun className="size-5" />}
                selected={prefs.theme === 'light'}
                onSelect={() => handlePrefChange({ theme: 'light' })}
              />
              <ThemeOption
                value="dark"
                label="Dark"
                icon={<Moon className="size-5" />}
                selected={prefs.theme === 'dark'}
                onSelect={() => handlePrefChange({ theme: 'dark' })}
              />
              <ThemeOption
                value="system"
                label="Auto"
                icon={<Monitor className="size-5" />}
                selected={prefs.theme === 'system'}
                onSelect={() => handlePrefChange({ theme: 'system' })}
              />
            </div>
          </CardContent>
        </Card>

        {/* Language */}
        <Card>
          <CardHeader>
            <CardTitle>Language</CardTitle>
            <CardDescription>
              Default language for narratives and exports
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Select
              value={prefs.language}
              onValueChange={(val) =>
                handlePrefChange({
                  language: val as UserPreferences['language'],
                })
              }
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="en">English</SelectItem>
                <SelectItem value="it">Italiano</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        {/* Narrative level */}
        <Card>
          <CardHeader>
            <CardTitle>Narrative Level</CardTitle>
            <CardDescription>
              Default detail level for AI-generated narratives
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Select
              value={prefs.narrativeLevel}
              onValueChange={(val) =>
                handlePrefChange({
                  narrativeLevel: val as UserPreferences['narrativeLevel'],
                })
              }
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="basic">Basic</SelectItem>
                <SelectItem value="advanced">Advanced</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>
      </section>

      <Separator />

      {/* API Keys */}
      <section className="space-y-6">
        <div>
          <h2 className="text-lg font-semibold">API Keys</h2>
          <p className="text-sm text-muted-foreground">
            Manage API keys for programmatic access to the Explainable Engine
            API.
          </p>
        </div>

        {/* Generate new key */}
        <Card>
          <CardHeader>
            <CardTitle>Generate New Key</CardTitle>
            <CardDescription>
              Create a new API key for programmatic access.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-end gap-3">
              <div className="grid flex-1 gap-2">
                <Label htmlFor="key-name">Key Name</Label>
                <Input
                  id="key-name"
                  placeholder="e.g. CI Pipeline"
                  value={keyName}
                  onChange={(e) => setKeyName(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && handleGenerate()}
                />
              </div>
              <Button onClick={handleGenerate} disabled={!keyName.trim()}>
                Generate
              </Button>
            </div>

            {newKey && (
              <div className="rounded-lg border border-amber-500/50 bg-amber-50 p-4 dark:bg-amber-950/20">
                <p className="mb-2 text-sm font-medium text-amber-800 dark:text-amber-200">
                  This key won&apos;t be shown again. Copy it now.
                </p>
                <div className="flex items-center gap-2">
                  <code className="flex-1 break-all rounded bg-amber-100 px-2 py-1 text-xs dark:bg-amber-900/30">
                    {newKey}
                  </code>
                  <Button variant="outline" size="sm" onClick={handleCopyKey}>
                    {copied ? (
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
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Key list */}
        {keys.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 space-y-4 rounded-xl border border-dashed">
            <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
              <Key className="h-6 w-6 text-muted-foreground" />
            </div>
            <div className="text-center space-y-1">
              <h3 className="text-base font-medium">No API keys</h3>
              <p className="text-sm text-muted-foreground max-w-sm">
                Generate your first key to access the API programmatically.
              </p>
            </div>
          </div>
        ) : (
          <div className="overflow-hidden rounded-lg border">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b bg-muted/50">
                  <th className="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wide">
                    Name
                  </th>
                  <th className="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wide">
                    Key
                  </th>
                  <th className="px-4 py-2.5 text-left text-xs font-medium text-muted-foreground uppercase tracking-wide">
                    Created
                  </th>
                  <th className="px-4 py-2.5 text-right text-xs font-medium text-muted-foreground uppercase tracking-wide">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {keys.map((k) => (
                  <tr
                    key={k.id}
                    className="border-b last:border-b-0 transition-colors hover:bg-accent/50"
                  >
                    <td className="px-4 py-2.5 font-medium">{k.name}</td>
                    <td className="px-4 py-2.5">
                      <code className="text-xs text-muted-foreground font-mono">
                        {k.prefix}
                      </code>
                    </td>
                    <td className="px-4 py-2.5 text-muted-foreground">
                      {new Date(k.createdAt).toLocaleDateString()}
                    </td>
                    <td className="px-4 py-2.5 text-right">
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
      </section>
    </div>
  );
}
