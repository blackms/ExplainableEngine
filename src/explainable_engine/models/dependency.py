from __future__ import annotations

from pydantic import BaseModel, Field


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
