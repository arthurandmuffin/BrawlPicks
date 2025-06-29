package repositories

import (
	"BrawlPicks/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type BrawlStarsRepository struct {
	apiKey string
}

func NewBrawlStarsRepository(apiKey string) *BrawlStarsRepository {
	return &BrawlStarsRepository{apiKey: apiKey}
}

func (r *BrawlStarsRepository) GetPlayer(tag string) (map[string]interface{}, error) {
	encodedTag := url.QueryEscape("#" + tag) // tag is usually passed without #
	url := fmt.Sprintf("https://api.brawlstars.com/v1/players/%s", encodedTag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *BrawlStarsRepository) GetPlayerData(tag string) (*models.PlayerResponse, error) {
	encodedTag := url.QueryEscape("#" + tag) // tag is usually passed without #
	url := fmt.Sprintf("https://api.brawlstars.com/v1/players/%s", encodedTag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6ImYzOTA3ZjEzLWM1NzAtNDViMS1iYjYwLWRlMzZmNjA4ZWM5MiIsImlhdCI6MTczODIwMjIxNiwic3ViIjoiZGV2ZWxvcGVyLzRhNmMzMjcyLTIzYWMtZDIxYi0zY2NlLTUzYzkxNDNkYjAxNCIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiMTczLjE3Ni4xMzcuMTA1IiwiNTQuODYuNTAuMTM5Il0sInR5cGUiOiJjbGllbnQifV19.ptYJLGyJKL6TMsCK_I-fAT6shsPI-Uf6QUuqXwP5P_dz4pKQOry-khaKxEYUOrKkdhMdjk5lKRee96RtneVtzA")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var player models.PlayerResponse
	if err := json.Unmarshal(bodyBytes, &player); err != nil {
		fmt.Println("Failed to unmarshal:", string(bodyBytes))
		return nil, err
	}

	return &player, nil
}

/*
type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) RepositoryTest() string {
	return "Repository"
}
*/
