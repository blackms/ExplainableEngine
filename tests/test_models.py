import warnings

from explainable_engine.models import (
    Component,
    ConfidenceResult,
    Contribution,
    Edge,
    ExplainRequest,
    ExplanationGraph,
    Node,
    NodeType,
    PropagationStep,
    TransformationType,
)


class TestNode:
    def test_create_node(self):
        node = Node(id="x", label="X", value=0.5)
        assert node.id == "x"
        assert node.confidence == 1.0
        assert node.node_type == NodeType.INPUT

    def test_node_roundtrip(self):
        node = Node(id="x", label="X", value=0.5, confidence=0.9, node_type=NodeType.COMPUTED)
        data = node.model_dump_json()
        restored = Node.model_validate_json(data)
        assert restored == node


class TestEdge:
    def test_create_edge(self):
        edge = Edge(source="a", target="b", weight=0.4)
        assert edge.transformation_type == TransformationType.WEIGHTED_SUM

    def test_edge_roundtrip(self):
        edge = Edge(source="a", target="b", weight=0.4)
        restored = Edge.model_validate_json(edge.model_dump_json())
        assert restored == edge


class TestExplanationGraph:
    def test_create_graph(self):
        nodes = [
            Node(id="root", label="Root", value=0.72, node_type=NodeType.OUTPUT),
            Node(id="a", label="A", value=0.8),
        ]
        edges = [Edge(source="a", target="root", weight=0.4)]
        graph = ExplanationGraph(nodes=nodes, edges=edges, root_node_id="root")
        assert len(graph.nodes) == 2
        assert len(graph.edges) == 1

    def test_graph_roundtrip(self):
        nodes = [
            Node(id="root", label="Root", value=0.72, node_type=NodeType.OUTPUT),
            Node(id="a", label="A", value=0.8),
        ]
        edges = [Edge(source="a", target="root", weight=0.4)]
        graph = ExplanationGraph(nodes=nodes, edges=edges, root_node_id="root")
        restored = ExplanationGraph.model_validate_json(graph.model_dump_json())
        assert restored == graph


class TestComponent:
    def test_basic_component(self):
        c = Component(name="trend", value=0.8, weight=0.4)
        assert c.confidence == 1.0
        assert c.components is None

    def test_nested_components(self):
        c = Component(
            name="skills",
            value=82.0,
            weight=0.3,
            components=[
                Component(name="python", value=95.0, weight=0.5),
                Component(name="go", value=70.0, weight=0.5),
            ],
        )
        assert len(c.components) == 2


class TestExplainRequest:
    def test_valid_request(self):
        req = ExplainRequest(
            target="score",
            value=0.72,
            components=[
                Component(name="a", value=0.8, weight=0.4, confidence=0.9),
                Component(name="b", value=0.5, weight=0.3, confidence=0.7),
                Component(name="c", value=0.6, weight=0.3, confidence=0.85),
            ],
        )
        assert req.target == "score"
        assert len(req.components) == 3
        assert req.options.include_graph is True

    def test_weights_warning(self):
        with warnings.catch_warnings(record=True) as w:
            warnings.simplefilter("always")
            ExplainRequest(
                target="score",
                value=0.72,
                components=[
                    Component(name="a", value=0.8, weight=0.6),
                    Component(name="b", value=0.5, weight=0.6),
                ],
            )
            assert len(w) == 1
            assert "exceeds 1.0" in str(w[0].message)

    def test_roundtrip(self):
        req = ExplainRequest(
            target="score",
            value=0.72,
            components=[Component(name="a", value=0.8, weight=0.4)],
        )
        restored = ExplainRequest.model_validate_json(req.model_dump_json())
        assert restored == req


class TestContribution:
    def test_with_children(self):
        c = Contribution(
            node_id="a",
            label="A",
            value=0.8,
            weight=0.4,
            absolute_contribution=0.32,
            percentage=44.4,
            confidence=0.9,
            children=[
                Contribution(
                    node_id="a1",
                    label="A1",
                    value=0.9,
                    weight=0.5,
                    absolute_contribution=0.45,
                    percentage=56.25,
                    confidence=0.95,
                )
            ],
        )
        assert len(c.children) == 1


class TestConfidenceResult:
    def test_create(self):
        result = ConfidenceResult(
            overall_confidence=0.828,
            node_confidences={"a": 0.9, "b": 0.7},
            propagation_path=[
                PropagationStep(
                    node_id="root",
                    computed_confidence=0.828,
                    source_nodes=["a", "b"],
                    formula="(0.4*0.9 + 0.3*0.7) / (0.4+0.3)",
                )
            ],
        )
        assert result.overall_confidence == 0.828
