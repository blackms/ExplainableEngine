package aip

import (
	"fmt"
	"math"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// TransformInstrumentSentiment converts AIP sentiment data into an ExplainRequest
// that can be processed by the ExplainableEngine orchestrator.
//
// The transformation maps the AIP sentiment components into weighted factors:
//   - sentiment_7d (35%): short-term sentiment, core signal
//   - sentiment_30d (25%): longer-term sentiment, more stable baseline
//   - trend (20%): direction of sentiment change
//   - positive_ratio (10%): ratio of positive articles
//   - news_score (10%): news-specific sentiment score
//
// Confidence is derived from article count and source availability.
func TransformInstrumentSentiment(s *InstrumentSentiment) *models.ExplainRequest {
	// Map sentiment_7d range [-1, 1] to a normalized value [0, 1].
	normalizedScore := (s.Sentiment7D + 1) / 2 // -1->0, 0->0.5, 1->1

	// Confidence based on article count: more articles = higher confidence, capped at 1.0.
	articleConfidence := math.Min(float64(s.ArticleCount7D)/20.0, 1.0)

	components := []models.Component{
		{
			Name:       "sentiment_7d",
			Value:      s.Sentiment7D,
			Weight:     0.35,
			Confidence: articleConfidence,
		},
		{
			Name:       "sentiment_30d",
			Value:      s.Sentiment30D,
			Weight:     0.25,
			Confidence: math.Min(articleConfidence*1.5, 1.0), // 30d window is more stable
		},
		{
			Name:       "trend",
			Value:      s.Trend,
			Weight:     0.20,
			Confidence: articleConfidence,
		},
		{
			Name:       "positive_ratio",
			Value:      s.PositiveRatio,
			Weight:     0.10,
			Confidence: articleConfidence,
		},
		{
			Name:       "news_score",
			Value:      s.NewsSentiment.Score,
			Weight:     0.10,
			Confidence: boolToConfidence(s.NewsSentiment.HasRecentHeadlines),
			Missing:    !s.NewsSentiment.HasRecentHeadlines && s.NewsSentiment.ArticleCount == 0,
		},
	}

	return &models.ExplainRequest{
		Target:     fmt.Sprintf("%s_sentiment", s.Ticker),
		Value:      normalizedScore,
		Components: components,
		Metadata: map[string]string{
			"source":          "aip",
			"ticker":          s.Ticker,
			"sentiment_label": s.SentimentLabel,
			"last_updated":    s.LastUpdated,
		},
	}
}

// TransformMarketMood converts AIP market mood data into an ExplainRequest.
// Each sector becomes a component whose weight is proportional to its article coverage.
func TransformMarketMood(m *MarketMood) *models.ExplainRequest {
	components := make([]models.Component, len(m.Sectors))

	totalArticles := 0
	for _, s := range m.Sectors {
		totalArticles += s.ArticleCount
	}

	// Guard against division by zero when there are no articles.
	if totalArticles == 0 {
		totalArticles = 1
	}

	for i, sector := range m.Sectors {
		weight := float64(sector.ArticleCount) / float64(totalArticles)
		confidence := math.Min(float64(sector.ArticleCount)/100.0, 1.0)

		components[i] = models.Component{
			Name:       sector.Sector,
			Value:      sector.AverageSentiment,
			Weight:     weight,
			Confidence: confidence,
		}
	}

	return &models.ExplainRequest{
		Target:     "market_mood",
		Value:      m.OverallSentiment,
		Components: components,
		Metadata: map[string]string{
			"source":         "aip",
			"total_articles": fmt.Sprintf("%d", m.TotalArticles),
		},
	}
}

// boolToConfidence returns 1.0 for true and 0.0 for false.
func boolToConfidence(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
