'use client';

import { useState } from 'react';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { LLMNarrative } from './LLMNarrative';
import { QAChat } from './QAChat';
import { SummaryGenerator } from './SummaryGenerator';
import {
  SparklesIcon,
  ChevronDownIcon,
  ChevronUpIcon,
} from 'lucide-react';

interface AIAssistantSectionProps {
  explanationId: string;
}

export function AIAssistantSection({
  explanationId,
}: AIAssistantSectionProps) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="flex items-center gap-2 text-base font-semibold">
          <SparklesIcon className="size-4 text-purple-500" />
          AI Assistant
        </h2>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setCollapsed(!collapsed)}
        >
          {collapsed ? (
            <>
              <ChevronDownIcon className="size-4" data-icon="inline-start" />
              Show
            </>
          ) : (
            <>
              <ChevronUpIcon className="size-4" data-icon="inline-start" />
              Hide
            </>
          )}
        </Button>
      </div>

      {!collapsed && (
        <Tabs defaultValue="narrative">
          <TabsList>
            <TabsTrigger value="narrative">AI Narrative</TabsTrigger>
            <TabsTrigger value="chat">Ask a Question</TabsTrigger>
            <TabsTrigger value="summary">Executive Summary</TabsTrigger>
          </TabsList>

          <TabsContent value="narrative">
            <LLMNarrative explanationId={explanationId} />
          </TabsContent>

          <TabsContent value="chat">
            <QAChat explanationId={explanationId} />
          </TabsContent>

          <TabsContent value="summary">
            <SummaryGenerator explanationId={explanationId} />
          </TabsContent>
        </Tabs>
      )}
    </div>
  );
}
