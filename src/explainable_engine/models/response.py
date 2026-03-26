from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel, Field


class BreakdownItem(BaseModel):
    node_id: str
    label: str
    value: float
    weight: float
    absolute_contribution: float
    percentage: float
    confidence: float
    children: list[BreakdownItem] | None = None


class DriverItem(BaseModel):
    name: str
    impact: float
    rank: int


class GraphNodeResponse(BaseModel):
    id: str
    label: str
    value: float
    confidence: float
    node_type: str


class GraphEdgeResponse(BaseModel):
    source: str
    target: str
    weight: float
    transformation_type: str


class GraphResponse(BaseModel):
    nodes: list[GraphNodeResponse]
    edges: list[GraphEdgeResponse]


class DependencyNode(BaseModel):
    node_id: str
    label: str
    depth: int
    relation: str | None = None
    children: list[DependencyNode] = Field(default_factory=list)


class DependencyTree(BaseModel):
    root: DependencyNode
    depth: int
    total_nodes: int


class ConfidenceDetail(BaseModel):
    overall: float
    per_node: dict[str, float]


class ExplainMetadata(BaseModel):
    version: str
    created_at: datetime
    deterministic_hash: str
    computation_type: str = "additive"


class ExplainResponse(BaseModel):
    id: str
    target: str
    final_value: float
    confidence: float
    breakdown: list[BreakdownItem]
    top_drivers: list[DriverItem]
    missing_impact: float = 0.0
    graph: GraphResponse | None = None
    dependency_tree: DependencyTree | None = None
    confidence_detail: ConfidenceDetail | None = None
    metadata: ExplainMetadata
