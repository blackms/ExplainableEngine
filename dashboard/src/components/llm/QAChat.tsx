'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import { useAskQuestion } from '@/lib/api/hooks';
import type { ChatMessage } from '@/lib/api/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  MessageCircleIcon,
  SendIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  AlertCircleIcon,
} from 'lucide-react';

interface QAChatProps {
  explanationId: string;
}

const SUGGESTED_QUESTIONS = [
  'Why is this value so high?',
  'Which component has the most impact?',
  'How reliable is this result?',
  'What are the main risk factors?',
];

export function QAChat({ explanationId }: QAChatProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [collapsed, setCollapsed] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  const mutation = useAskQuestion(explanationId);

  const scrollToBottom = useCallback(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, mutation.isPending, scrollToBottom]);

  const handleSend = useCallback(
    (question: string) => {
      if (!question.trim() || mutation.isPending) return;

      const userMessage: ChatMessage = {
        role: 'user',
        content: question.trim(),
      };
      const updatedMessages = [...messages, userMessage];
      setMessages(updatedMessages);
      setInput('');

      mutation.mutate(
        { question: question.trim(), history: messages },
        {
          onSuccess: (data) => {
            setMessages((prev) => [
              ...prev,
              { role: 'assistant', content: data.answer },
            ]);
          },
          onError: () => {
            // Keep user message but show error in the UI
          },
        }
      );
    },
    [messages, mutation]
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        handleSend(input);
      }
    },
    [handleSend, input]
  );

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <MessageCircleIcon className="size-4 text-blue-500" />
            Ask a Question
          </CardTitle>
          <Button
            variant="ghost"
            size="icon-xs"
            onClick={() => setCollapsed(!collapsed)}
          >
            {collapsed ? (
              <ChevronDownIcon className="size-4" />
            ) : (
              <ChevronUpIcon className="size-4" />
            )}
          </Button>
        </div>
      </CardHeader>

      {!collapsed && (
        <CardContent className="space-y-3">
          {/* Chat messages area */}
          <div
            ref={scrollRef}
            className="flex h-64 flex-col gap-2 overflow-y-auto rounded-md border bg-muted/20 p-3"
          >
            {messages.length === 0 && !mutation.isPending ? (
              <div className="flex flex-1 flex-col items-center justify-center gap-3">
                <p className="text-sm text-muted-foreground">
                  Ask anything about this explanation
                </p>
                <div className="flex flex-wrap justify-center gap-1.5">
                  {SUGGESTED_QUESTIONS.map((q) => (
                    <button
                      key={q}
                      onClick={() => handleSend(q)}
                      className="rounded-full border bg-background px-3 py-1 text-xs text-muted-foreground transition-colors hover:border-primary hover:text-foreground"
                    >
                      {q}
                    </button>
                  ))}
                </div>
              </div>
            ) : (
              <>
                {messages.map((msg, i) => (
                  <div
                    key={i}
                    className={`flex ${
                      msg.role === 'user' ? 'justify-end' : 'justify-start'
                    }`}
                  >
                    <div
                      className={`max-w-[80%] rounded-lg px-3 py-2 text-sm ${
                        msg.role === 'user'
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted text-foreground'
                      }`}
                    >
                      {msg.content}
                    </div>
                  </div>
                ))}

                {/* Typing indicator */}
                {mutation.isPending && (
                  <div className="flex justify-start">
                    <div className="rounded-lg bg-muted px-3 py-2 text-sm text-muted-foreground">
                      <span className="inline-flex gap-0.5">
                        <span className="animate-bounce [animation-delay:0ms]">
                          .
                        </span>
                        <span className="animate-bounce [animation-delay:150ms]">
                          .
                        </span>
                        <span className="animate-bounce [animation-delay:300ms]">
                          .
                        </span>
                      </span>
                    </div>
                  </div>
                )}

                {/* Error state */}
                {mutation.isError && (
                  <div className="flex justify-start">
                    <div className="flex items-center gap-2 rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-700 dark:text-red-400">
                      <AlertCircleIcon className="size-3.5 shrink-0" />
                      <span>
                        Failed to get a response. Check that the LLM API key is
                        configured.
                      </span>
                    </div>
                  </div>
                )}
              </>
            )}
          </div>

          {/* Input area */}
          <div className="flex items-center gap-2">
            <Input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask about this explanation..."
              disabled={mutation.isPending}
            />
            <Button
              variant="default"
              size="icon"
              onClick={() => handleSend(input)}
              disabled={!input.trim() || mutation.isPending}
            >
              <SendIcon className="size-4" />
            </Button>
          </div>
        </CardContent>
      )}
    </Card>
  );
}
