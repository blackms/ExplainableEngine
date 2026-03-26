'use client';

import { useMemo, useCallback, useState } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  type Node,
  type Edge,
  type NodeMouseHandler,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import dagre from 'dagre';
import type { GraphResponse } from '@/lib/api/types';
import { ExplanationNode } from './ExplanationNode';

const nodeTypes = { explanation: ExplanationNode };

const NODE_WIDTH = 200;
const NODE_HEIGHT = 80;

function buildFlowGraph(graph: GraphResponse): { nodes: Node[]; edges: Edge[] } {
  const g = new dagre.graphlib.Graph();
  g.setDefaultEdgeLabel(() => ({}));
  g.setGraph({ rankdir: 'LR', nodesep: 80, ranksep: 120 });

  graph.nodes.forEach((n) => {
    g.setNode(n.id, { width: NODE_WIDTH, height: NODE_HEIGHT });
  });

  graph.edges.forEach((e) => {
    g.setEdge(e.source, e.target);
  });

  dagre.layout(g);

  const nodes: Node[] = graph.nodes.map((n) => {
    const pos = g.node(n.id);
    return {
      id: n.id,
      position: { x: pos.x - NODE_WIDTH / 2, y: pos.y - NODE_HEIGHT / 2 },
      data: {
        label: n.label,
        value: n.value,
        confidence: n.confidence,
        nodeType: n.node_type,
      },
      type: 'explanation',
    };
  });

  const edges: Edge[] = graph.edges.map((e, i) => ({
    id: `e-${i}`,
    source: e.source,
    target: e.target,
    label: `w=${e.weight.toFixed(2)}`,
    animated: true,
    style: { strokeWidth: Math.max(1, e.weight * 4) },
  }));

  return { nodes, edges };
}

interface NodeDetailPopover {
  x: number;
  y: number;
  label: string;
  value: number;
  confidence: number;
  nodeType: string;
}

interface ExplanationGraphProps {
  graph: GraphResponse;
  className?: string;
}

export function ExplanationGraph({ graph, className }: ExplanationGraphProps) {
  const { nodes, edges } = useMemo(() => buildFlowGraph(graph), [graph]);
  const [popover, setPopover] = useState<NodeDetailPopover | null>(null);

  const onNodeClick: NodeMouseHandler = useCallback((event, node) => {
    const data = node.data as Record<string, unknown>;
    const rect = (event.target as HTMLElement).getBoundingClientRect();
    setPopover({
      x: rect.right + 8,
      y: rect.top,
      label: data.label as string,
      value: data.value as number,
      confidence: data.confidence as number,
      nodeType: data.nodeType as string,
    });
  }, []);

  const onPaneClick = useCallback(() => {
    setPopover(null);
  }, []);

  return (
    <div className={className ?? 'h-[500px] w-full'}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        onNodeClick={onNodeClick}
        onPaneClick={onPaneClick}
        fitView
        fitViewOptions={{ padding: 0.2 }}
        minZoom={0.2}
        maxZoom={2}
      >
        <Background gap={16} size={1} />
        <Controls position="bottom-left" />
        <MiniMap
          position="bottom-right"
          nodeStrokeWidth={3}
          zoomable
          pannable
        />
      </ReactFlow>

      {popover && (
        <div
          className="fixed z-50 rounded-lg border bg-card p-3 shadow-lg text-sm ring-1 ring-foreground/10"
          style={{ left: popover.x, top: popover.y }}
        >
          <div className="font-semibold mb-1">{popover.label}</div>
          <div className="text-muted-foreground space-y-0.5">
            <p>Value: {popover.value.toFixed(4)}</p>
            <p>Confidence: {(popover.confidence * 100).toFixed(1)}%</p>
            <p>Type: {popover.nodeType}</p>
          </div>
        </div>
      )}
    </div>
  );
}
