package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// List handles GET /api/v1/explain — returns a paginated, filtered list of explanations.
func (h *ExplainHandler) List(w http.ResponseWriter, r *http.Request) {
	opts := storage.ListOptions{
		Limit: 20,
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		opts.Limit = n
	}
	if v := r.URL.Query().Get("cursor"); v != "" {
		opts.Cursor = v
	}
	if v := r.URL.Query().Get("target"); v != "" {
		opts.Target = v
	}
	if v := r.URL.Query().Get("min_confidence"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid min_confidence parameter")
			return
		}
		opts.MinConfidence = f
	}
	if v := r.URL.Query().Get("max_confidence"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid max_confidence parameter")
			return
		}
		opts.MaxConfidence = f
	}
	if v := r.URL.Query().Get("from"); v != "" {
		opts.FromTime = v
	}
	if v := r.URL.Query().Get("to"); v != "" {
		opts.ToTime = v
	}

	result, err := h.store.List(opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list explanations: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Stats handles GET /api/v1/stats — returns total explanation count.
func (h *ExplainHandler) Stats(w http.ResponseWriter, r *http.Request) {
	count, err := h.store.Count()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to count explanations: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total_explanations": count,
		"status":             "ok",
	})
}
