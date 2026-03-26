from enum import Enum

from pydantic import BaseModel, Field


class NodeType(str, Enum):
    INPUT = "input"
    COMPUTED = "computed"
    OUTPUT = "output"
    MISSING = "missing"


class Node(BaseModel):
    id: str
    label: str
    value: float
    confidence: float = Field(ge=0.0, le=1.0, default=1.0)
    node_type: NodeType = NodeType.INPUT
    metadata: dict | None = None

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "id": "trend_strength",
                    "label": "Trend Strength",
                    "value": 0.8,
                    "confidence": 0.9,
                    "node_type": "input",
                }
            ]
        }
    }
