'use client';

import { PreferencesForm } from '@/components/settings/PreferencesForm';
import { ApiKeyManager } from '@/components/settings/ApiKeyManager';
import { Separator } from '@/components/ui/separator';

export default function SettingsPage() {
  return (
    <div className="space-y-8 max-w-2xl">
      <h1 className="text-2xl font-bold">Settings</h1>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold">Preferences</h2>
        <PreferencesForm />
      </section>

      <Separator />

      <section className="space-y-4">
        <h2 className="text-lg font-semibold">API Keys</h2>
        <p className="text-sm text-muted-foreground">
          Manage API keys for programmatic access to the Explainable Engine API.
        </p>
        <ApiKeyManager />
      </section>
    </div>
  );
}
