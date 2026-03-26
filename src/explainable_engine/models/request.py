from __future__ import annotations

import warnings

from pydantic import BaseModel, Field, model_validator


class Component(BaseModel):
    id: str | None = None
    name: str
    value: float
    weight: float = Field(ge=0.0, le=1.0)
    confidence: float = Field(ge=0.0, le=1.0, default=1.0)
    components: list[Component] | None = None

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "name": "trend_strength",
                    "value": 0.8,
                    "weight": 0.4,
                    "confidence": 0.9,
                }
            ]
        }
    }


class ExplainOptions(BaseModel):
    include_graph: bool = True
    include_drivers: bool = True
    max_drivers: int = Field(default=5, ge=1)
    max_depth: int = Field(default=10, ge=1)


class ExplainRequest(BaseModel):
    target: str
    value: float
    components: list[Component]
    options: ExplainOptions = Field(default_factory=ExplainOptions)
    metadata: dict | None = None

    @model_validator(mode="after")
    def warn_if_weights_exceed_one(self) -> ExplainRequest:
        total = sum(c.weight for c in self.components)
        if total > 1.0 + 1e-9:
            warnings.warn(
                f"Component weights sum to {total:.4f}, which exceeds 1.0",
                stacklevel=2,
            )
        return self

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "target": "market_regime_score",
                    "value": 0.72,
                    "components": [
                        {
                            "name": "trend_strength",
                            "value": 0.8,
                            "weight": 0.4,
                            "confidence": 0.9,
                        },
                        {
                            "name": "volatility",
                            "value": 0.5,
                            "weight": 0.3,
                            "confidence": 0.7,
                        },
                        {
                            "name": "momentum",
                            "value": 0.6,
                            "weight": 0.3,
                            "confidence": 0.85,
                        },
                    ],
                }
            ]
        }
    }
