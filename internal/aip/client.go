package aip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://aip.aigenconsult.com"
	defaultTimeout = 15 * time.Second
	maxBulkTickers = 50
)

// AIPService defines the interface for interacting with the AIP sentiment API.
// This abstraction allows tests to use mock implementations.
type AIPService interface {
	GetSentiment(ctx context.Context, ticker string) (*InstrumentSentiment, error)
	GetBulkSentiment(ctx context.Context, tickers []string) ([]InstrumentSentiment, error)
	GetMarketMood(ctx context.Context) (*MarketMood, error)
	GetHeadlines(ctx context.Context, ticker string) ([]Headline, error)
	GetHistory(ctx context.Context, ticker string) (*SentimentHistory, error)
}

// Client is the HTTP client for the AIP sentiment API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Verify that Client implements AIPService at compile time.
var _ AIPService = (*Client)(nil)

// NewClient creates a new AIP API client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// GetSentiment fetches sentiment data for a single instrument.
func (c *Client) GetSentiment(ctx context.Context, ticker string) (*InstrumentSentiment, error) {
	url := fmt.Sprintf("%s/api/v1/external/sentiment/instruments/%s", c.baseURL, ticker)

	var envelope struct {
		Data struct {
			Instrument InstrumentSentiment `json:"instrument"`
		} `json:"data"`
	}
	if err := c.doGet(ctx, url, &envelope); err != nil {
		return nil, fmt.Errorf("fetching sentiment for %s: %w", ticker, err)
	}
	result := envelope.Data.Instrument
	return &result, nil
}

// GetBulkSentiment fetches sentiment data for multiple instruments (up to 50).
func (c *Client) GetBulkSentiment(ctx context.Context, tickers []string) ([]InstrumentSentiment, error) {
	if len(tickers) > maxBulkTickers {
		return nil, fmt.Errorf("bulk request exceeds maximum of %d tickers (got %d)", maxBulkTickers, len(tickers))
	}
	if len(tickers) == 0 {
		return nil, fmt.Errorf("tickers list must not be empty")
	}

	url := fmt.Sprintf("%s/api/v1/external/sentiment/instruments/bulk", c.baseURL)
	body := BulkRequest{Tickers: tickers}

	var envelope struct {
		Data struct {
			Instruments []InstrumentSentiment `json:"instruments"`
		} `json:"data"`
	}
	if err := c.doPost(ctx, url, body, &envelope); err != nil {
		return nil, fmt.Errorf("fetching bulk sentiment: %w", err)
	}
	return envelope.Data.Instruments, nil
}

// GetMarketMood fetches the overall market mood and per-sector breakdown.
func (c *Client) GetMarketMood(ctx context.Context) (*MarketMood, error) {
	url := fmt.Sprintf("%s/api/v1/external/sentiment/market-mood", c.baseURL)

	var envelope struct {
		Data MarketMood `json:"data"`
	}
	if err := c.doGet(ctx, url, &envelope); err != nil {
		return nil, fmt.Errorf("fetching market mood: %w", err)
	}
	return &envelope.Data, nil
}

// GetHeadlines fetches news headlines for a single instrument.
func (c *Client) GetHeadlines(ctx context.Context, ticker string) ([]Headline, error) {
	url := fmt.Sprintf("%s/api/v1/external/sentiment/instruments/%s/headlines", c.baseURL, ticker)

	var envelope struct {
		Data struct {
			Headlines []Headline `json:"headlines"`
		} `json:"data"`
	}
	if err := c.doGet(ctx, url, &envelope); err != nil {
		return nil, fmt.Errorf("fetching headlines for %s: %w", ticker, err)
	}
	return envelope.Data.Headlines, nil
}

// GetHistory fetches sentiment history for a single instrument.
func (c *Client) GetHistory(ctx context.Context, ticker string) (*SentimentHistory, error) {
	url := fmt.Sprintf("%s/api/v1/external/sentiment/instruments/%s/history", c.baseURL, ticker)

	var envelope struct {
		Data SentimentHistory `json:"data"`
	}
	if err := c.doGet(ctx, url, &envelope); err != nil {
		return nil, fmt.Errorf("fetching history for %s: %w", ticker, err)
	}
	return &envelope.Data, nil
}

// doGet performs a GET request with the API key header and decodes the response.
func (c *Client) doGet(ctx context.Context, url string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	return c.do(req, dest)
}

// doPost performs a POST request with a JSON body and decodes the response.
func (c *Client) doPost(ctx context.Context, url string, body any, dest any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, dest)
}

// do executes the HTTP request, sets common headers, and decodes the response.
func (c *Client) do(req *http.Request, dest any) error {
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	if err := json.Unmarshal(respBody, dest); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// APIError represents an error response from the AIP API.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("AIP API error (status %d): %s", e.StatusCode, e.Body)
}
