from __future__ import annotations

from pydantic import BaseModel


class Contribution(BaseModel):
    node_id: str
    label: str
    value: float
    weight: float
    absolute_contribution: float
    percentage: float
    confidence: float
    children: list[Contribution] | None = None
