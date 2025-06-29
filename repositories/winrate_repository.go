package repositories

/*
import (
	"BrawlPicks/models"
	"encoding/json"
	"errors"

	"net/http"
)

type EventResponse struct {
    Active []struct {
        Map struct {
            Name   string `json:"name"`
            ID     int    `json:"id"`
            GameMode struct {
                Name string `json:"name"`
            } `json:"gameMode"`
            Stats []struct {
                BrawlerID int     `json:"brawler"`
                WinRate   float64 `json:"winRate"`
            } `json:"stats"`
        } `json:"map"`
    } `json:"active"`
}

type WinrateRepository interface {
	GetEventWinrates() ([]models.EventWinrate, error)
}

type winrateRepository struct {
	client *http.Client
}

func NewWinrateRepository() WinrateRepository {
	return &winrateRepository{
		client: &http.Client{},
	}
}

func (r *winrateRepository) GetEventWinrates() ([]models.EventWinrate, error) {
	resp, err := r.client.Get("https://api.brawlify.com/v1/events")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Active   []EventData `json:"active"`
		Upcoming []EventData `json:"upcoming"`
	}
= b := range brawlers {
				topBrawlers[i] = models.BrawlerWinrate{
					Name:    b.Name,
					Winrate: b.WinRate,
				}
			}
			result = append(result, models.EventWinrate{
				EventID:     event.Event.ID,
				EventName:   event.Event.Name,
				TopBrawlers: topBrawlers,
			})
		}
	}

	processEvents(raw.Active)
	processEvents(raw.Upcoming)

	if len(result) == 0 {
		return nil, errors.New("no event data found")
	}

	return result, nil
}
*/
