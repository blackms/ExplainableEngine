package aip

// InstrumentSentiment represents the AIP sentiment response for a single ticker.
type InstrumentSentiment struct {
	Ticker          string        `json:"ticker"`
	Sentiment7D     float64       `json:"sentiment_7d"`
	Sentiment30D    float64       `json:"sentiment_30d"`
	Trend           float64       `json:"trend"`
	ArticleCount7D  int           `json:"article_count_7d"`
	PositiveRatio   float64       `json:"positive_ratio"`
	LastUpdated     string        `json:"last_updated"`
	SentimentLabel  string        `json:"sentiment_label"`
	SentimentScaled int           `json:"sentiment_score_scaled"`
	NewsSentiment   NewsSentiment `json:"news_sentiment"`
	Sources         Sources       `json:"sources"`
}

// NewsSentiment contains the news-specific sentiment data.
type NewsSentiment struct {
	Score              float64 `json:"score"`
	Label              string  `json:"label"`
	ArticleCount       int     `json:"article_count"`
	HasRecentHeadlines bool    `json:"has_recent_headlines"`
}

// Sources describes which data sources contributed to the sentiment.
type Sources struct {
	News    SourceDetail `json:"news"`
	Social  SourceDetail `json:"social"`
	Analyst SourceDetail `json:"analyst"`
}

// AvailableCount returns the number of available sources.
func (s Sources) AvailableCount() int {
	count := 0
	if s.News.Available {
		count++
	}
	if s.Social.Available {
		count++
	}
	if s.Analyst.Available {
		count++
	}
	return count
}

// TotalCount returns the total number of source types.
func (s Sources) TotalCount() int {
	return 3
}

// SourceDetail describes an individual data source.
type SourceDetail struct {
	Available bool   `json:"available"`
	Provider  string `json:"provider,omitempty"`
}

// MarketMood represents the AIP market mood response.
type MarketMood struct {
	OverallSentiment float64           `json:"overall_sentiment"`
	OverallTrend     float64           `json:"overall_trend"`
	TotalArticles    int               `json:"total_articles"`
	Sectors          []SectorSentiment `json:"sectors"`
}

// SectorSentiment represents sentiment data for a single market sector.
type SectorSentiment struct {
	Sector           string  `json:"sector"`
	AverageSentiment float64 `json:"average_sentiment"`
	ArticleCount     int     `json:"article_count"`
	InstrumentCount  int     `json:"instrument_count"`
	Trend            float64 `json:"trend"`
}

// Headline represents a single news headline from AIP.
type Headline struct {
	Title     string  `json:"title"`
	Source    string  `json:"source"`
	URL       string  `json:"url"`
	Sentiment float64 `json:"sentiment"`
	Published string  `json:"published"`
}

// SentimentHistory represents historical sentiment data from AIP.
type SentimentHistory struct {
	Ticker     string           `json:"ticker"`
	DataPoints []HistoryPoint   `json:"data_points"`
}

// HistoryPoint is a single data point in the sentiment history.
type HistoryPoint struct {
	Date      string  `json:"date"`
	Sentiment float64 `json:"sentiment"`
	Articles  int     `json:"articles"`
}

// BulkRequest is the request body for the bulk sentiment endpoint.
type BulkRequest struct {
	Tickers []string `json:"tickers"`
}
