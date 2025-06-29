package repositories

import (
	"BrawlPicks/models"
	"encoding/json"
	"net/http"
)

type PLRepository struct{}

func NewPLRepository() *PLRepository {
	return &PLRepository{}
}

func (r *PLRepository) FetchPLResults() (models.PLResults, error) {
	url := "https://storage.googleapis.com/brawlanalyzer-public/pl-results.json.gz"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results models.PLResults
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}
