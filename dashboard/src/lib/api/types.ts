// Request types

export interface Component {
  id?: string;
  name: string;
  value: number;
  weight: number;
  confidence: number;
  missing?: boolean;
  components?: Component[];
}

export interface ExplainOptions {
  include_graph: boolean;
  include_drivers: boolean;
  max_drivers: number;
  max_depth: number;
  missing_threshold: number;
}

export interface ExplainRequest {
  target: string;
  value: number;
  components: Component[];
  options?: Partial<ExplainOptions>;
  metadata?: Record<string, string>;
}

// Response types

export interface BreakdownItem {
  node_id: string;
  label: string;
  value: number;
  weight: number;
  absolute_contribution: number;
  percentage: number;
  confidence: number;
  children?: BreakdownItem[];
}

export interface DriverItem {
  name: string;
  impact: number;
  rank: number;
}

export interface GraphNodeResponse {
  id: string;
  label: string;
  value: number;
  confidence: number;
  node_type: string;
}

export interface GraphEdgeResponse {
  source: string;
  target: string;
  weight: number;
  transformation_type: string;
}

export interface GraphResponse {
  nodes: GraphNodeResponse[];
  edges: GraphEdgeResponse[];
}

export interface DependencyNode {
  node_id: string;
  label: string;
  depth: number;
  relation?: string;
  children?: DependencyNode[];
}

export interface DependencyTree {
  root: DependencyNode;
  depth: number;
  total_nodes: number;
}

export interface ConfidenceDetail {
  overall: number;
  per_node: Record<string, number>;
}

export interface ExplainMetadata {
  version: string;
  created_at: string;
  deterministic_hash: string;
  computation_type: string;
}

export interface ExplainResponse {
  id: string;
  target: string;
  final_value: number;
  confidence: number;
  breakdown: BreakdownItem[];
  top_drivers: DriverItem[];
  missing_impact: number;
  graph?: GraphResponse;
  dependency_tree?: DependencyTree;
  confidence_detail?: ConfidenceDetail;
  metadata: ExplainMetadata;
  original_request?: ExplainRequest;
}

// List / Audit

export interface ListOptions {
  cursor?: string;
  limit?: number;
  target?: string;
  min_confidence?: number;
  max_confidence?: number;
  from?: string;
  to?: string;
}

export interface ListResult {
  items: ExplainResponse[];
  next_cursor?: string;
  total: number;
}

// Narrative

export interface NarrativeResult {
  explanation_id: string;
  level: string;
  language: string;
  narrative: string;
  confidence_level: string;
  has_missing_data: boolean;
}

// What-if

export interface Modification {
  component: string;
  new_value: number;
}

export interface ComponentDiff {
  name: string;
  original_value: number;
  modified_value: number;
  delta_value: number;
  delta_percentage: number;
  original_contribution: number;
  modified_contribution: number;
}

export interface SensitivityRanking {
  name: string;
  impact: number;
  rank: number;
}

export interface SensitivityResult {
  original_value: number;
  modified_value: number;
  delta_value: number;
  delta_percentage: number;
  component_diffs: ComponentDiff[];
  sensitivity_ranking: SensitivityRanking[];
}
