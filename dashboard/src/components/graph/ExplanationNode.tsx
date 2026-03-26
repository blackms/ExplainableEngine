'use client';

import { memo } from 'react';
import { Handle, Position, type NodeProps } from '@xyflow/react';

const nodeColors: Record<string, string> = {
  input: '#86efac',    // green-300
  output: '#fda4af',   // rose-300
  computed: '#93c5fd',  // blue-300
  missing: '#d1d5db',   // gray-300
};

interface ExplanationNodeData {
  label: string;
  value: number;
  confidence: number;
  nodeType: string;
  [key: string]: unknown;
}

function ExplanationNodeComponent({ data }: NodeProps) {
  const { label, value, confidence, nodeType } = data as unknown as ExplanationNodeData;
  const bgColor = nodeColors[nodeType] ?? '#e5e7eb';
  const confidencePct = ((confidence ?? 0) * 100).toFixed(0);

  return (
    <div
      className="rounded-lg border border-gray-300 shadow-md px-4 py-3 min-w-[180px] text-xs"
      style={{ backgroundColor: bgColor }}
    >
      <Handle type="target" position={Position.Left} className="!bg-gray-500" />

      <div className="font-semibold text-gray-900 text-sm truncate mb-1">
        {label}
      </div>

      <div className="flex items-center justify-between text-gray-700">
        <span>Value: {Number(value).toFixed(4)}</span>
        <span className="ml-2 text-[10px] uppercase tracking-wide text-gray-500">
          {nodeType}
        </span>
      </div>

      <div className="mt-1.5">
        <div className="flex items-center justify-between text-[10px] text-gray-600 mb-0.5">
          <span>Confidence</span>
          <span>{confidencePct}%</span>
        </div>
        <div className="h-1.5 w-full rounded-full bg-gray-200 overflow-hidden">
          <div
            className="h-full rounded-full transition-all duration-300"
            style={{
              width: `${confidencePct}%`,
              backgroundColor:
                confidence >= 0.8
                  ? '#22c55e'
                  : confidence >= 0.5
                    ? '#eab308'
                    : '#ef4444',
            }}
          />
        </div>
      </div>

      <Handle type="source" position={Position.Right} className="!bg-gray-500" />
    </div>
  );
}

export const ExplanationNode = memo(ExplanationNodeComponent);
