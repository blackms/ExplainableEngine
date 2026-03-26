'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Plus, Trash2, Send, Clock } from 'lucide-react';
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
import { useCreateExplanation } from '@/lib/api/hooks';
import type { Component as ComponentType } from '@/lib/api/types';

interface RecentExplanation {
  id: string;
  target: string;
  value: number;
  confidence: number;
  created_at: string;
}

const STORAGE_KEY = 'explainable-engine-recent';

function loadRecent(): RecentExplanation[] {
  if (typeof window === 'undefined') return [];
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

function saveRecent(items: RecentExplanation[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items.slice(0, 20)));
}

function emptyComponent(): ComponentType {
  return { name: '', value: 0, weight: 1, confidence: 1 };
}

export default function HomePage() {
  const router = useRouter();
  const createMutation = useCreateExplanation();

  const [target, setTarget] = useState('');
  const [value, setValue] = useState<number>(0);
  const [components, setComponents] = useState<ComponentType[]>([
    emptyComponent(),
  ]);
  const [recent, setRecent] = useState<RecentExplanation[]>([]);
  const [formOpen, setFormOpen] = useState(false);

  useEffect(() => {
    setRecent(loadRecent());
  }, []);

  const updateComponent = useCallback(
    (index: number, field: keyof ComponentType, val: string | number) => {
      setComponents((prev) => {
        const next = [...prev];
        next[index] = { ...next[index], [field]: val };
        return next;
      });
    },
    []
  );

  const addComponent = useCallback(() => {
    setComponents((prev) => [...prev, emptyComponent()]);
  }, []);

  const removeComponent = useCallback((index: number) => {
    setComponents((prev) => prev.filter((_, i) => i !== index));
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validComponents = components.filter((c) => c.name.trim() !== '');
    if (!target.trim() || validComponents.length === 0) return;

    try {
      const result = await createMutation.mutateAsync({
        target: target.trim(),
        value,
        components: validComponents,
        options: {
          include_graph: true,
          include_drivers: true,
        },
      });

      const entry: RecentExplanation = {
        id: result.id,
        target: result.target,
        value: result.final_value,
        confidence: result.confidence,
        created_at: result.metadata.created_at,
      };

      const updated = [entry, ...recent.filter((r) => r.id !== result.id)];
      setRecent(updated);
      saveRecent(updated);

      router.push(`/explain/${result.id}`);
    } catch {
      // Error is available via createMutation.error
    }
  };

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      {/* New Explanation */}
      <Card>
        <CardHeader>
          <CardTitle>New Explanation</CardTitle>
          <CardDescription>
            Submit components to generate an explainability breakdown
          </CardDescription>
        </CardHeader>
        <CardContent>
          {!formOpen ? (
            <Button onClick={() => setFormOpen(true)} size="lg">
              <Plus className="size-4" data-icon="inline-start" />
              Create Explanation
            </Button>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="target">Target name</Label>
                  <Input
                    id="target"
                    placeholder="e.g. Final Score"
                    value={target}
                    onChange={(e) => setTarget(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="value">Value</Label>
                  <Input
                    id="value"
                    type="number"
                    step="any"
                    value={value}
                    onChange={(e) => setValue(parseFloat(e.target.value) || 0)}
                    required
                  />
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label>Components</Label>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={addComponent}
                  >
                    <Plus className="size-3.5" data-icon="inline-start" />
                    Add
                  </Button>
                </div>

                <div className="space-y-2">
                  {components.map((comp, i) => (
                    <div
                      key={i}
                      className="grid grid-cols-[1fr_80px_80px_80px_auto] items-end gap-2"
                    >
                      <div className="space-y-1">
                        {i === 0 && (
                          <Label className="text-xs text-muted-foreground">
                            Name
                          </Label>
                        )}
                        <Input
                          placeholder="Component name"
                          value={comp.name}
                          onChange={(e) =>
                            updateComponent(i, 'name', e.target.value)
                          }
                        />
                      </div>
                      <div className="space-y-1">
                        {i === 0 && (
                          <Label className="text-xs text-muted-foreground">
                            Value
                          </Label>
                        )}
                        <Input
                          type="number"
                          step="any"
                          value={comp.value}
                          onChange={(e) =>
                            updateComponent(
                              i,
                              'value',
                              parseFloat(e.target.value) || 0
                            )
                          }
                        />
                      </div>
                      <div className="space-y-1">
                        {i === 0 && (
                          <Label className="text-xs text-muted-foreground">
                            Weight
                          </Label>
                        )}
                        <Input
                          type="number"
                          step="any"
                          min="0"
                          max="1"
                          value={comp.weight}
                          onChange={(e) =>
                            updateComponent(
                              i,
                              'weight',
                              parseFloat(e.target.value) || 0
                            )
                          }
                        />
                      </div>
                      <div className="space-y-1">
                        {i === 0 && (
                          <Label className="text-xs text-muted-foreground">
                            Conf.
                          </Label>
                        )}
                        <Input
                          type="number"
                          step="any"
                          min="0"
                          max="1"
                          value={comp.confidence}
                          onChange={(e) =>
                            updateComponent(
                              i,
                              'confidence',
                              parseFloat(e.target.value) || 0
                            )
                          }
                        />
                      </div>
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => removeComponent(i)}
                        disabled={components.length === 1}
                        className={i === 0 ? 'mt-5' : ''}
                      >
                        <Trash2 className="size-3.5 text-muted-foreground" />
                      </Button>
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex items-center gap-2 pt-2">
                <Button
                  type="submit"
                  disabled={createMutation.isPending}
                  size="lg"
                >
                  <Send className="size-4" data-icon="inline-start" />
                  {createMutation.isPending ? 'Submitting...' : 'Submit'}
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="lg"
                  onClick={() => setFormOpen(false)}
                >
                  Cancel
                </Button>
              </div>

              {createMutation.isError && (
                <p className="text-sm text-destructive">
                  {createMutation.error.message}
                </p>
              )}
            </form>
          )}
        </CardContent>
      </Card>

      {/* Recent Explanations */}
      {recent.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Recent Explanations</CardTitle>
            <CardDescription>
              Previously generated explanations
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ul className="divide-y divide-border">
              {recent.map((item) => (
                <li key={item.id}>
                  <button
                    type="button"
                    onClick={() => router.push(`/explain/${item.id}`)}
                    className="flex w-full items-center justify-between gap-4 py-3 text-left transition-colors hover:bg-muted/50 rounded-lg px-2 -mx-2"
                  >
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium">
                        {item.target}
                      </p>
                      <p className="flex items-center gap-1.5 text-xs text-muted-foreground">
                        <Clock className="size-3" />
                        {new Date(item.created_at).toLocaleString()}
                      </p>
                    </div>
                    <div className="flex shrink-0 items-center gap-4 text-sm">
                      <span className="text-muted-foreground">
                        Value:{' '}
                        <span className="font-medium text-foreground">
                          {item.value.toFixed(2)}
                        </span>
                      </span>
                      <span className="text-muted-foreground">
                        Conf:{' '}
                        <span className="font-medium text-foreground">
                          {(item.confidence * 100).toFixed(0)}%
                        </span>
                      </span>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
