from pydantic import BaseModel


class PropagationStep(BaseModel):
    node_id: str
    computed_confidence: float
    source_nodes: list[str]
    formula: str


class ConfidenceResult(BaseModel):
    overall_confidence: float
    node_confidences: dict[str, float]
    propagation_path: list[PropagationStep]
