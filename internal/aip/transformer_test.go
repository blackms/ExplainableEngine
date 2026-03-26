package aip

import (
	"math"
	"testing"
)

func TestTransformInstrumentSentiment(t *testing.T) {
	s := &InstrumentSentiment{
		Ticker:         "NVDA",
		Sentiment7D:    0.158,
		Sentiment30D:   0.270,
		Trend:          -0.112,
		ArticleCount7D: 6,
		PositiveRatio:  0.33,
		LastUpdated:    "2026-03-26T12:00:00Z",
		SentimentLabel: "Neutral",
		NewsSentiment: NewsSentiment{
			Score:              0.158,
			Label:              "Neutral",
			ArticleCount:       6,
			HasRecentHeadlines: true,
		},
		Sources: Sources{
			News:    SourceDetail{Available: true},
			Social:  SourceDetail{Available: false},
			Analyst: SourceDetail{Available: false},
		},
	}

	req := TransformInstrumentSentiment(s)

	// Target should include the ticker.
	if req.Target != "NVDA_sentiment" {
		t.Errorf("target = %q, want %q", req.Target, "NVDA_sentiment")
	}

	// Value should be normalized: (0.158 + 1) / 2 = 0.579.
	expectedValue := (0.158 + 1) / 2
	if math.Abs(req.Value-expectedValue) > 1e-9 {
		t.Errorf("value = %f, want %f", req.Value, expectedValue)
	}

	// Should have 5 components.
	if len(req.Components) != 5 {
		t.Fatalf("components count = %d, want 5", len(req.Components))
	}

	// Verify component names and weights.
	expectedNames := []string{"sentiment_7d", "sentiment_30d", "trend", "positive_ratio", "news_score"}
	expectedWeights := []float64{0.35, 0.25, 0.20, 0.10, 0.10}
	for i, c := range req.Components {
		if c.Name != expectedNames[i] {
			t.Errorf("component[%d].Name = %q, want %q", i, c.Name, expectedNames[i])
		}
		if math.Abs(c.Weight-expectedWeights[i]) > 1e-9 {
			t.Errorf("component[%d].Weight = %f, want %f", i, c.Weight, expectedWeights[i])
		}
	}

	// Article confidence: 6 / 20 = 0.3.
	expectedConf := 6.0 / 20.0
	if math.Abs(req.Components[0].Confidence-expectedConf) > 1e-9 {
		t.Errorf("sentiment_7d confidence = %f, want %f", req.Components[0].Confidence, expectedConf)
	}

	// 30d confidence: min(0.3 * 1.5, 1.0) = 0.45.
	expectedConf30 := math.Min(expectedConf*1.5, 1.0)
	if math.Abs(req.Components[1].Confidence-expectedConf30) > 1e-9 {
		t.Errorf("sentiment_30d confidence = %f, want %f", req.Components[1].Confidence, expectedConf30)
	}

	// news_score should not be missing since HasRecentHeadlines is true.
	if req.Components[4].Missing {
		t.Error("news_score should not be marked missing when HasRecentHeadlines is true")
	}

	// Metadata should contain source info.
	if req.Metadata["source"] != "aip" {
		t.Errorf("metadata source = %q, want %q", req.Metadata["source"], "aip")
	}
	if req.Metadata["ticker"] != "NVDA" {
		t.Errorf("metadata ticker = %q, want %q", req.Metadata["ticker"], "NVDA")
	}
}

func TestTransformInstrumentSentiment_ZeroArticles(t *testing.T) {
	s := &InstrumentSentiment{
		Ticker:         "XYZ",
		Sentiment7D:    0.0,
		ArticleCount7D: 0,
		NewsSentiment: NewsSentiment{
			Score:              0,
			ArticleCount:       0,
			HasRecentHeadlines: false,
		},
	}

	req := TransformInstrumentSentiment(s)

	// With 0 articles, confidence should be 0.
	if req.Components[0].Confidence != 0.0 {
		t.Errorf("confidence with 0 articles = %f, want 0.0", req.Components[0].Confidence)
	}

	// news_score should be missing.
	if !req.Components[4].Missing {
		t.Error("news_score should be missing when no headlines and 0 articles")
	}
}

func TestTransformInstrumentSentiment_HighArticleCount(t *testing.T) {
	s := &InstrumentSentiment{
		Ticker:         "AAPL",
		Sentiment7D:    0.5,
		ArticleCount7D: 100,
		NewsSentiment: NewsSentiment{
			Score:              0.5,
			ArticleCount:       100,
			HasRecentHeadlines: true,
		},
	}

	req := TransformInstrumentSentiment(s)

	// Confidence should be capped at 1.0 for high article counts.
	if req.Components[0].Confidence != 1.0 {
		t.Errorf("confidence with 100 articles = %f, want 1.0", req.Components[0].Confidence)
	}
}

func TestTransformMarketMood(t *testing.T) {
	m := &MarketMood{
		OverallSentiment: 0.107,
		OverallTrend:     0.043,
		TotalArticles:    65043,
		Sectors: []SectorSentiment{
			{Sector: "Energy", AverageSentiment: 0.203, ArticleCount: 274, InstrumentCount: 22, Trend: 0.125},
			{Sector: "Technology", AverageSentiment: 0.146, ArticleCount: 726, InstrumentCount: 30, Trend: 0.050},
		},
	}

	req := TransformMarketMood(m)

	if req.Target != "market_mood" {
		t.Errorf("target = %q, want %q", req.Target, "market_mood")
	}

	if math.Abs(req.Value-0.107) > 1e-9 {
		t.Errorf("value = %f, want %f", req.Value, 0.107)
	}

	if len(req.Components) != 2 {
		t.Fatalf("components count = %d, want 2", len(req.Components))
	}

	// Weights should be proportional to article counts.
	totalArticles := 274.0 + 726.0
	expectedEnergyWeight := 274.0 / totalArticles
	if math.Abs(req.Components[0].Weight-expectedEnergyWeight) > 1e-9 {
		t.Errorf("Energy weight = %f, want %f", req.Components[0].Weight, expectedEnergyWeight)
	}

	expectedTechWeight := 726.0 / totalArticles
	if math.Abs(req.Components[1].Weight-expectedTechWeight) > 1e-9 {
		t.Errorf("Technology weight = %f, want %f", req.Components[1].Weight, expectedTechWeight)
	}

	// Sector names.
	if req.Components[0].Name != "Energy" {
		t.Errorf("component[0].Name = %q, want %q", req.Components[0].Name, "Energy")
	}
	if req.Components[1].Name != "Technology" {
		t.Errorf("component[1].Name = %q, want %q", req.Components[1].Name, "Technology")
	}

	// Confidence: 274/100 = 2.74, capped at 1.0.
	if req.Components[0].Confidence != 1.0 {
		t.Errorf("Energy confidence = %f, want 1.0", req.Components[0].Confidence)
	}
}

func TestTransformMarketMood_EmptySectors(t *testing.T) {
	m := &MarketMood{
		OverallSentiment: 0.0,
		Sectors:          []SectorSentiment{},
	}

	req := TransformMarketMood(m)

	if len(req.Components) != 0 {
		t.Errorf("components count = %d, want 0", len(req.Components))
	}
}

func TestTransformMarketMood_ZeroArticles(t *testing.T) {
	m := &MarketMood{
		OverallSentiment: 0.0,
		Sectors: []SectorSentiment{
			{Sector: "Energy", AverageSentiment: 0.1, ArticleCount: 0},
		},
	}

	req := TransformMarketMood(m)

	// With 0 total articles, weight should be 0/1 = 0 (guarded by totalArticles=1).
	if req.Components[0].Weight != 0.0 {
		t.Errorf("weight with 0 articles = %f, want 0.0", req.Components[0].Weight)
	}
	if req.Components[0].Confidence != 0.0 {
		t.Errorf("confidence with 0 articles = %f, want 0.0", req.Components[0].Confidence)
	}
}

func TestBoolToConfidence(t *testing.T) {
	if boolToConfidence(true) != 1.0 {
		t.Error("boolToConfidence(true) should be 1.0")
	}
	if boolToConfidence(false) != 0.0 {
		t.Error("boolToConfidence(false) should be 0.0")
	}
}
