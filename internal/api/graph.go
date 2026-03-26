package api

import (
	"net/http"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// GraphHandler handles graph export requests.
type GraphHandler struct {
	store storage.ExplanationStore
}

// Export handles GET /api/v1/explain/{id}/graph.
// It retrieves a stored explanation and serializes its graph in the requested format.
// Supported formats: json (default), dot, mermaid.
func (h *GraphHandler) Export(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	// Validate format before doing any store lookup.
	gf := engine.GraphFormat(format)
	switch gf {
	case engine.FormatJSON, engine.FormatDOT, engine.FormatMermaid:
		// valid
	default:
		writeError(w, http.StatusBadRequest, "unsupported format: "+format)
		return
	}

	resp, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve explanation: "+err.Error())
		return
	}
	if resp == nil {
		writeError(w, http.StatusNotFound, "explanation not found")
		return
	}

	if resp.Graph == nil {
		writeError(w, http.StatusNotFound, "explanation does not contain a graph")
		return
	}

	content, contentType, err := engine.SerializeGraph(resp.Graph, gf)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "serialization failed: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
