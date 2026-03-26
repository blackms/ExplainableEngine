package api

import (
	"encoding/json"
	"net/http"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// ExplainHandler handles requests for the /api/v1/explain endpoints.
type ExplainHandler struct {
	orchestrator engine.OrchestratorInterface
	store        storage.ExplanationStore
}

// Create handles POST /api/v1/explain.
func (h *ExplainHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.ExplainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Validate required fields.
	if req.Target == "" {
		writeError(w, http.StatusBadRequest, "target is required")
		return
	}
	if len(req.Components) == 0 {
		writeError(w, http.StatusBadRequest, "components must not be empty")
		return
	}

	resp, err := h.orchestrator.Explain(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "explanation failed: "+err.Error())
		return
	}

	if err := h.store.Save(resp); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save explanation: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /api/v1/explain/{id}.
func (h *ExplainHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
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

	writeJSON(w, http.StatusOK, resp)
}
