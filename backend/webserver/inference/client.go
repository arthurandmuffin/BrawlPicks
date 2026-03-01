package inference

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	env "BrawlPicks/webserver/config"
)

type Client struct {
	baseURL       string
	recommendPath string
	client        *http.Client
}

type RecommendRequest struct {
	MapName           string `json:"map_name"`
	Mode              string `json:"mode"`
	Rank              int    `json:"rank"`
	AllyBrawlers      []int  `json:"ally_brawlers,omitempty"`
	EnemyBrawlers     []int  `json:"enemy_brawlers,omitempty"`
	CandidateBrawlers []int  `json:"candidate_brawlers,omitempty"`
	BannedBrawlers    []int  `json:"banned_brawlers,omitempty"`
	TopK              *int   `json:"top_k,omitempty"`
}

type RecommendResponse struct {
	ModelID         string           `json:"model_id"`
	Recommendations []Recommendation `json:"recommendations"`
}

type Recommendation struct {
	BrawlerID int     `json:"brawler_id"`
	Score     float64 `json:"score"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("inference returned status %d", e.StatusCode)
}

func IsBadRequest(err error) bool {
	httpErr, ok := err.(*HTTPError)
	return ok && httpErr.StatusCode == http.StatusBadRequest
}

func NewClient(e *env.Env, client *http.Client) *Client {
	return &Client{
		baseURL:       strings.TrimRight(e.Inference.BaseURL, "/"),
		recommendPath: e.Inference.RecommendPath,
		client:        client,
	}
}

func (c *Client) Recommend(ctx context.Context, payload *RecommendRequest) (*RecommendResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.recommendURL(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.decodeHTTPError(resp.StatusCode, raw)
	}

	out := new(RecommendResponse)
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) recommendURL() string {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return c.baseURL + c.recommendPath
	}
	path, err := url.Parse(c.recommendPath)
	if err != nil {
		return c.baseURL + c.recommendPath
	}
	return base.ResolveReference(path).String()
}

func (c *Client) decodeHTTPError(statusCode int, raw []byte) error {
	payload := new(errorResponse)
	if err := json.Unmarshal(raw, payload); err == nil && payload.Error != "" {
		return &HTTPError{
			StatusCode: statusCode,
			Message:    payload.Error,
		}
	}
	return &HTTPError{
		StatusCode: statusCode,
		Message:    strings.TrimSpace(string(raw)),
	}
}
