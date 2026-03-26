package api

import (
	"encoding/json"
	"net/http"

	"github.com/blackms/ExplainableEngine/internal/aip"
	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// AIPHandler handles AIP-powered explanation endpoints.
type AIPHandler struct {
	aipClient    aip.AIPService
	orchestrator engine.OrchestratorInterface
	store        storage.ExplanationStore
}

// ExplainTicker fetches sentiment from AIP for a single ticker, transforms it
// into an ExplainRequest, runs it through the orchestrator, and returns the result.
//
// GET /api/v1/aip/explain/{ticker}
func (h *AIPHandler) ExplainTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	if ticker == "" {
		writeError(w, http.StatusBadRequest, "ticker is required")
		return
	}

	sentiment, err := h.aipClient.GetSentiment(r.Context(), ticker)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to fetch AIP sentiment: "+err.Error())
		return
	}

	req := aip.TransformInstrumentSentiment(sentiment)

	result, err := h.orchestrator.Explain(*req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "explanation failed: "+err.Error())
		return
	}

	if err := h.store.Save(result); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save explanation: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"explanation": result,
		"aip_data":    sentiment,
		"ticker":      ticker,
	})
}

// ExplainMarketMood fetches market mood from AIP, transforms sectors into components,
// runs the explanation, and returns the result.
//
// GET /api/v1/aip/explain/market-mood
func (h *AIPHandler) ExplainMarketMood(w http.ResponseWriter, r *http.Request) {
	mood, err := h.aipClient.GetMarketMood(r.Context())
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to fetch AIP market mood: "+err.Error())
		return
	}

	req := aip.TransformMarketMood(mood)

	result, err := h.orchestrator.Explain(*req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "explanation failed: "+err.Error())
		return
	}

	if err := h.store.Save(result); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save explanation: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"explanation": result,
		"aip_data":    mood,
	})
}

// bulkExplainRequest is the JSON body for the bulk explain endpoint.
type bulkExplainRequest struct {
	Tickers []string `json:"tickers"`
}

// bulkExplainItem holds the result for a single ticker in a bulk request.
type bulkExplainItem struct {
	Ticker      string      `json:"ticker"`
	Explanation any         `json:"explanation,omitempty"`
	AIPData     any         `json:"aip_data,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// ExplainBulk fetches bulk sentiment from AIP and creates an explanation for each ticker.
//
// POST /api/v1/aip/explain/bulk
func (h *AIPHandler) ExplainBulk(w http.ResponseWriter, r *http.Request) {
	var body bulkExplainRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if len(body.Tickers) == 0 {
		writeError(w, http.StatusBadRequest, "tickers list must not be empty")
		return
	}
	if len(body.Tickers) > 50 {
		writeError(w, http.StatusBadRequest, "maximum 50 tickers per request")
		return
	}

	sentiments, err := h.aipClient.GetBulkSentiment(r.Context(), body.Tickers)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to fetch AIP bulk sentiment: "+err.Error())
		return
	}

	results := make([]bulkExplainItem, len(sentiments))
	for i, s := range sentiments {
		sentiment := s // capture loop variable
		req := aip.TransformInstrumentSentiment(&sentiment)

		result, err := h.orchestrator.Explain(*req)
		if err != nil {
			results[i] = bulkExplainItem{
				Ticker: sentiment.Ticker,
				Error:  err.Error(),
			}
			continue
		}

		_ = h.store.Save(result) // best-effort save

		results[i] = bulkExplainItem{
			Ticker:      sentiment.Ticker,
			Explanation: result,
			AIPData:     &sentiment,
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"results": results,
		"count":   len(results),
	})
}
