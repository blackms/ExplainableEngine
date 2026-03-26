package api

import (
	"net/http"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// NarrativeHandler handles requests for the narrative endpoint.
type NarrativeHandler struct {
	store storage.ExplanationStore
}

// Get handles GET /api/v1/explain/{id}/narrative.
func (h *NarrativeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	level := r.URL.Query().Get("level")
	lang := r.URL.Query().Get("lang")

	// Defaults
	if level == "" {
		level = "basic"
	}
	if lang == "" {
		lang = "en"
	}

	// Validate level
	narrativeLevel := engine.NarrativeLevel(level)
	if narrativeLevel != engine.LevelBasic && narrativeLevel != engine.LevelAdvanced {
		writeError(w, http.StatusBadRequest, "unsupported level: must be 'basic' or 'advanced'")
		return
	}

	// Validate language
	narrativeLang := engine.NarrativeLanguage(lang)
	if narrativeLang != engine.LangEN && narrativeLang != engine.LangIT {
		writeError(w, http.StatusBadRequest, "unsupported lang: must be 'en' or 'it'")
		return
	}

	// Retrieve explanation from store
	resp, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve explanation: "+err.Error())
		return
	}
	if resp == nil {
		writeError(w, http.StatusNotFound, "explanation not found")
		return
	}

	// Generate narrative
	result, err := engine.GenerateNarrative(resp, narrativeLevel, narrativeLang)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate narrative: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
