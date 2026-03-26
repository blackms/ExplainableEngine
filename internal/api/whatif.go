package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// WhatIfRequest is the request body for what-if analysis.
type WhatIfRequest struct {
	Modifications []engine.Modification `json:"modifications"`
}

// WhatIfHandler handles what-if analysis requests.
type WhatIfHandler struct {
	store        storage.ExplanationStore
	orchestrator engine.OrchestratorInterface
}

// Analyze handles POST /api/v1/explain/{id}/what-if.
func (h *WhatIfHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	// 1. Get original explanation from store.
	original, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve explanation: "+err.Error())
		return
	}
	if original == nil {
		writeError(w, http.StatusNotFound, "explanation not found")
		return
	}

	// 2. Decode request body (modifications).
	var req WhatIfRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// 3. Verify that the original request is available.
	if original.OriginalRequest == nil {
		writeError(w, http.StatusInternalServerError, "original request not available for this explanation")
		return
	}

	// 4. Call AnalyzeSensitivity.
	result, err := engine.AnalyzeSensitivity(
		original.OriginalRequest,
		original,
		req.Modifications,
		h.orchestrator,
	)
	if err != nil {
		var cnfErr *engine.ComponentNotFoundError
		if errors.As(err, &cnfErr) {
			writeError(w, http.StatusUnprocessableEntity, cnfErr.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "what-if analysis failed: "+err.Error())
		return
	}

	// 5. Return SensitivityResult (NOT persisted).
	writeJSON(w, http.StatusOK, result)
}
