from pydantic import BaseModel

from explainable_engine.models.edge import Edge
from explainable_engine.models.node import Node


class ExplanationGraph(BaseModel):
    nodes: list[Node]
    edges: list[Edge]
    root_node_id: str

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "root_node_id": "market_regime_score",
                    "nodes": [
                        {
                            "id": "market_regime_score",
                            "label": "Market Regime Score",
                            "value": 0.72,
                            "confidence": 0.828,
                            "node_type": "output",
                        },
                        {
                            "id": "trend_strength",
                            "label": "Trend Strength",
                            "value": 0.8,
                            "confidence": 0.9,
                            "node_type": "input",
                        },
                    ],
                    "edges": [
                        {
                            "source": "trend_strength",
                            "target": "market_regime_score",
                            "weight": 0.4,
                            "transformation_type": "weighted_sum",
                        }
                    ],
                }
            ]
        }
    }
