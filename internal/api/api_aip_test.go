package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/aip"
	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// mockAIPService implements aip.AIPService for testing.
type mockAIPService struct {
	sentiments  map[string]*aip.InstrumentSentiment
	marketMood  *aip.MarketMood
	errSentiment error
	errMood      error
}

func (m *mockAIPService) GetSentiment(_ context.Context, ticker string) (*aip.InstrumentSentiment, error) {
	if m.errSentiment != nil {
		return nil, m.errSentiment
	}
	s, ok := m.sentiments[ticker]
	if !ok {
		return nil, fmt.Errorf("ticker %s not found", ticker)
	}
	return s, nil
}

func (m *mockAIPService) GetBulkSentiment(_ context.Context, tickers []string) ([]aip.InstrumentSentiment, error) {
	if m.errSentiment != nil {
		return nil, m.errSentiment
	}
	var results []aip.InstrumentSentiment
	for _, t := range tickers {
		s, ok := m.sentiments[t]
		if !ok {
			continue
		}
		results = append(results, *s)
	}
	return results, nil
}

func (m *mockAIPService) GetMarketMood(_ context.Context) (*aip.MarketMood, error) {
	if m.errMood != nil {
		return nil, m.errMood
	}
	return m.marketMood, nil
}

func (m *mockAIPService) GetHeadlines(_ context.Context, _ string) ([]aip.Headline, error) {
	return nil, nil
}

func (m *mockAIPService) GetHistory(_ context.Context, _ string) (*aip.SentimentHistory, error) {
	return nil, nil
}

func newAIPTestRouter(svc aip.AIPService) http.Handler {
	store := storage.NewInMemoryStore()
	orch := &mockOrchestrator{}
	return api.NewRouter(store, orch, api.WithAIPClient(svc))
}

func sampleSentiment(ticker string) *aip.InstrumentSentiment {
	return &aip.InstrumentSentiment{
		Ticker:         ticker,
		Sentiment7D:    0.158,
		Sentiment30D:   0.270,
		Trend:          -0.112,
		ArticleCount7D: 6,
		PositiveRatio:  0.33,
		LastUpdated:    "2026-03-26T12:00:00Z",
		SentimentLabel: "Neutral",
		NewsSentiment: aip.NewsSentiment{
			Score:              0.158,
			Label:              "Neutral",
			ArticleCount:       6,
			HasRecentHeadlines: true,
		},
		Sources: aip.Sources{
			News:    aip.SourceDetail{Available: true},
			Social:  aip.SourceDetail{Available: false},
			Analyst: aip.SourceDetail{Available: false},
		},
	}
}

func TestExplainTicker_Success(t *testing.T) {
	svc := &mockAIPService{
		sentiments: map[string]*aip.InstrumentSentiment{
			"AAPL": sampleSentiment("AAPL"),
		},
	}
	router := newAIPTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/aip/explain/AAPL", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["ticker"] != "AAPL" {
		t.Errorf("ticker = %v, want AAPL", resp["ticker"])
	}
	if resp["explanation"] == nil {
		t.Error("response should include explanation")
	}
	if resp["aip_data"] == nil {
		t.Error("response should include aip_data")
	}
}

func TestExplainTicker_AIPError(t *testing.T) {
	svc := &mockAIPService{
		sentiments:   map[string]*aip.InstrumentSentiment{},
		errSentiment: fmt.Errorf("connection refused"),
	}
	router := newAIPTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/aip/explain/AAPL", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
}

func TestExplainMarketMood_Success(t *testing.T) {
	svc := &mockAIPService{
		sentiments: map[string]*aip.InstrumentSentiment{},
		marketMood: &aip.MarketMood{
			OverallSentiment: 0.107,
			OverallTrend:     0.043,
			TotalArticles:    65043,
			Sectors: []aip.SectorSentiment{
				{Sector: "Energy", AverageSentiment: 0.203, ArticleCount: 274, InstrumentCount: 22, Trend: 0.125},
				{Sector: "Technology", AverageSentiment: 0.146, ArticleCount: 726, InstrumentCount: 30, Trend: 0.050},
			},
		},
	}
	router := newAIPTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/aip/explain/market-mood", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["explanation"] == nil {
		t.Error("response should include explanation")
	}
	if resp["aip_data"] == nil {
		t.Error("response should include aip_data")
	}
}

func TestExplainMarketMood_AIPError(t *testing.T) {
	svc := &mockAIPService{
		sentiments: map[string]*aip.InstrumentSentiment{},
		errMood:    fmt.Errorf("timeout"),
	}
	router := newAIPTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/aip/explain/market-mood", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
}

func TestExplainBulk_Success(t *testing.T) {
	svc := &mockAIPService{
		sentiments: map[string]*aip.InstrumentSentiment{
			"AAPL": sampleSentiment("AAPL"),
			"MSFT": sampleSentiment("MSFT"),
		},
	}
	router := newAIPTestRouter(svc)

	body, _ := json.Marshal(map[string]any{"tickers": []string{"AAPL", "MSFT"}})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/aip/explain/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	count, ok := resp["count"].(float64)
	if !ok || count != 2 {
		t.Errorf("count = %v, want 2", resp["count"])
	}
}

func TestExplainBulk_EmptyTickers(t *testing.T) {
	svc := &mockAIPService{sentiments: map[string]*aip.InstrumentSentiment{}}
	router := newAIPTestRouter(svc)

	body, _ := json.Marshal(map[string]any{"tickers": []string{}})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/aip/explain/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestExplainBulk_TooManyTickers(t *testing.T) {
	svc := &mockAIPService{sentiments: map[string]*aip.InstrumentSentiment{}}
	router := newAIPTestRouter(svc)

	tickers := make([]string, 51)
	for i := range tickers {
		tickers[i] = fmt.Sprintf("T%d", i)
	}

	body, _ := json.Marshal(map[string]any{"tickers": tickers})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/aip/explain/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAIPRoutesNotRegistered_WithoutClient(t *testing.T) {
	// When no AIP client is provided, AIP routes should return 404.
	store := storage.NewInMemoryStore()
	orch := &mockOrchestrator{}
	router := api.NewRouter(store, orch) // no WithAIPClient

	req := httptest.NewRequest(http.MethodGet, "/api/v1/aip/explain/AAPL", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d (routes should not be registered without AIP client)", rec.Code, http.StatusNotFound)
	}
}
