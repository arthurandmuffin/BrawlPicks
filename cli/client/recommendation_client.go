package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RecommendationClient struct {
	baseURL       string
	recommendPath string
	client        *http.Client
}

type RecommendRequest struct {
	MapName           string `json:"mapName"`
	Mode              string `json:"mode"`
	Rank              int    `json:"rank"`
	AllyBrawlers      []int  `json:"allyBrawlers,omitempty"`
	EnemyBrawlers     []int  `json:"enemyBrawlers,omitempty"`
	CandidateBrawlers []int  `json:"candidateBrawlers,omitempty"`
	BannedBrawlers    []int  `json:"bannedBrawlers,omitempty"`
	TopK              *int   `json:"topK,omitempty"`
}

type RecommendEnvelope struct {
	Code string           `json:"code"`
	Data RecommendPayload `json:"data"`
}

type RecommendPayload struct {
	ModelID         string            `json:"modelId"`
	Recommendations []*Recommendation `json:"recommendations"`
}

type Recommendation struct {
	BrawlerID int     `json:"brawlerId"`
	Score     float64 `json:"score"`
}

type errorEnvelope struct {
	Code string `json:"code"`
	Data any    `json:"data"`
}

func New(baseURL, recommendPath string) *RecommendationClient {
	return &RecommendationClient{
		baseURL:       strings.TrimRight(baseURL, "/"),
		recommendPath: recommendPath,
		client:        &http.Client{},
	}
}

func (c *RecommendationClient) Recommend(payload *RecommendRequest) (*RecommendPayload, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.recommendURL(), bytes.NewReader(body))
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
		return nil, decodeError(raw, resp.StatusCode)
	}

	out := new(RecommendEnvelope)
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *RecommendationClient) recommendURL() string {
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

func decodeError(raw []byte, statusCode int) error {
	payload := new(errorEnvelope)
	if err := json.Unmarshal(raw, payload); err == nil {
		switch value := payload.Data.(type) {
		case string:
			if strings.TrimSpace(value) != "" {
				return fmt.Errorf("http %d: %s", statusCode, value)
			}
		}
	}
	return fmt.Errorf("http %d: %s", statusCode, strings.TrimSpace(string(raw)))
}
