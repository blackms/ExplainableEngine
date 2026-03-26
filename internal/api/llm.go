package api

import (
	"encoding/json"
	"net/http"

	"github.com/blackms/ExplainableEngine/internal/llm"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// LLMHandler handles LLM-powered endpoints for narratives, Q&A, and summaries.
type LLMHandler struct {
	store  storage.ExplanationStore
	llmSvc llm.Service
}

// narrativeRequest is the JSON body for the LLM narrative endpoint.
type narrativeRequest struct {
	Level string `json:"level"` // basic, advanced, executive
	Lang  string `json:"lang"`  // en, it
}

// narrativeResponse is the JSON response for the LLM narrative endpoint.
type narrativeResponse struct {
	Narrative string `json:"narrative"`
	Source    string `json:"source"` // "llm" or "template"
	Model     string `json:"model,omitempty"`
}

// askRequest is the JSON body for the Q&A endpoint.
type askRequest struct {
	Question string        `json:"question"`
	History  []llm.Message `json:"history,omitempty"`
}

// askResponse is the JSON response for the Q&A endpoint.
type askResponse struct {
	Answer string `json:"answer"`
	Model  string `json:"model"`
}

// summaryRequest is the JSON body for the executive summary endpoint.
type summaryRequest struct {
	Audience string `json:"audience"` // board, technical, client
	Lang     string `json:"lang"`     // en, it
}

// GenerateNarrative handles POST /api/v1/explain/{id}/narrative/llm.
func (h *LLMHandler) GenerateNarrative(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var body narrativeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Apply defaults.
	if body.Level == "" {
		body.Level = "basic"
	}
	if body.Lang == "" {
		body.Lang = "en"
	}

	// Validate level.
	switch body.Level {
	case "basic", "advanced", "executive":
		// ok
	default:
		writeError(w, http.StatusBadRequest, "unsupported level: must be 'basic', 'advanced', or 'executive'")
		return
	}

	// Validate language.
	switch body.Lang {
	case "en", "it":
		// ok
	default:
		writeError(w, http.StatusBadRequest, "unsupported lang: must be 'en' or 'it'")
		return
	}

	explanation, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve explanation: "+err.Error())
		return
	}
	if explanation == nil {
		writeError(w, http.StatusNotFound, "explanation not found")
		return
	}

	narrative, err := h.llmSvc.GenerateNarrative(r.Context(), explanation, body.Level, body.Lang)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate narrative: "+err.Error())
		return
	}

	source := "llm"
	model := "claude-sonnet-4-20250514"
	if _, ok := h.llmSvc.(*llm.FallbackService); ok {
		source = "template"
		model = ""
	}

	writeJSON(w, http.StatusOK, narrativeResponse{
		Narrative: narrative,
		Source:    source,
		Model:     model,
	})
}

// AskQuestion handles POST /api/v1/explain/{id}/ask.
func (h *LLMHandler) AskQuestion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var body askRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if body.Question == "" {
		writeError(w, http.StatusBadRequest, "question is required")
		return
	}

	explanation, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve explanation: "+err.Error())
		return
	}
	if explanation == nil {
		writeError(w, http.StatusNotFound, "explanation not found")
		return
	}

	answer, err := h.llmSvc.AnswerQuestion(r.Context(), explanation, body.Question, body.History)
	if err != nil {
		// Return 503 when the LLM service is not available (fallback mode).
		if _, ok := h.llmSvc.(*llm.FallbackService); ok {
			writeError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to answer question: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, askResponse{
		Answer: answer,
		Model:  "claude-sonnet-4-20250514",
	})
}

// GenerateSummary handles POST /api/v1/explain/{id}/summary.
func (h *LLMHandler) GenerateSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var body summaryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Apply defaults.
	if body.Audience == "" {
		body.Audience = "client"
	}
	if body.Lang == "" {
		body.Lang = "en"
	}

	// Validate audience.
	switch body.Audience {
	case "board", "technical", "client":
		// ok
	default:
		writeError(w, http.StatusBadRequest, "unsupported audience: must be 'board', 'technical', or 'client'")
		return
	}

	// Validate language.
	switch body.Lang {
	case "en", "it":
		// ok
	default:
		writeError(w, http.StatusBadRequest, "unsupported lang: must be 'en' or 'it'")
		return
	}

	explanation, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve explanation: "+err.Error())
		return
	}
	if explanation == nil {
		writeError(w, http.StatusNotFound, "explanation not found")
		return
	}

	result, err := h.llmSvc.GenerateSummary(r.Context(), explanation, body.Audience, body.Lang)
	if err != nil {
		// Return 503 when the LLM service is not available (fallback mode).
		if _, ok := h.llmSvc.(*llm.FallbackService); ok {
			writeError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to generate summary: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
