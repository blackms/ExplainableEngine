from enum import Enum

from pydantic import BaseModel, Field


class TransformationType(str, Enum):
    WEIGHTED_SUM = "weighted_sum"
    NORMALIZATION = "normalization"
    THRESHOLD = "threshold"
    CUSTOM = "custom"


class Edge(BaseModel):
    source: str
    target: str
    weight: float = Field(ge=0.0, le=1.0)
    transformation_type: TransformationType = TransformationType.WEIGHTED_SUM
    metadata: dict | None = None

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "source": "trend_strength",
                    "target": "market_regime_score",
                    "weight": 0.4,
                    "transformation_type": "weighted_sum",
                }
            ]
        }
    }
