package upstream

import (
	"BrawlPicks/scraper/models"
	"encoding/json"
	"time"
)

const (
	WinResult  = "victory"
	LostResult = "defeat"
	DrawResult = "draw"
)

type BattleLogResponse struct {
	Battles []*BattleLogBattle `json:"items"`
}

type BattleLogBattle struct {
	Time          time.Time
	MapName       string
	Result        string
	StarPlayerTag string
	Teams         []*Team
}

type Team []*Player

type Player struct {
	Tag       string
	BrawlerID int
	Rank      int
}

const BattleTimestampLayout = "20060102T150405.000Z"
const RankedBattleType = "soloRanked"

func (b *BattleLogBattle) UnmarshalJSON(data []byte) error {
	type tmpPlayer struct {
		Tag     string `json:"tag"`
		Details struct {
			Brawler int `json:"id"`
			Rank    int `json:"trophies"`
		} `json:"brawler"`
	}
	var tmpBattle struct {
		BattleTime string `json:"battleTime"`
		Event      struct {
			MapName string `json:"map"`
		} `json:"event"`
		Details struct {
			Type       string `json:"type"`
			Result     string `json:"result"`
			StarPlayer struct {
				Tag string `json:"tag"`
			} `json:"starPlayer"`
			Teams [][]tmpPlayer `json:"teams"`
		} `json:"battle"`
	}

	if err := json.Unmarshal(data, &tmpBattle); err != nil {
		return err
	}

	if tmpBattle.Details.Type != RankedBattleType || len(tmpBattle.Details.Teams) < 2 {
		return nil
	}

	t, err := time.Parse(BattleTimestampLayout, tmpBattle.BattleTime)
	if err != nil {
		return err
	}

	b.Time = t
	b.MapName = tmpBattle.Event.MapName
	b.Result = tmpBattle.Details.Result
	b.StarPlayerTag = tmpBattle.Details.StarPlayer.Tag
	for _, team := range tmpBattle.Details.Teams {
		tmpTeam := Team{}
		for _, player := range team {
			tmpTeam = append(tmpTeam,
				&Player{
					Tag:       player.Tag,
					BrawlerID: player.Details.Brawler,
					Rank:      player.Details.Rank,
				},
			)
		}
		b.Teams = append(b.Teams, &tmpTeam)
	}
	return nil
}

func (b *BattleLogBattle) Transform(tag string) (battle *models.Battle) {
	if !b.IsValid() {
		return nil
	}

	ids := b.getIDTeams()

	battle = &models.Battle{
		Timestamp: b.Time,
		MapName:   b.MapName,
		Rank:      b.getMeanRank(),
		Draw:      b.Result == DrawResult,
	}

	if battle.Draw {
		battle.TeamW, battle.TeamL = ids[0], ids[1]
		return battle
	}

	onTeam0 := false
	for _, p := range *b.Teams[0] {
		if p.Tag == tag {
			onTeam0 = true
			break
		}
	}

	win := WinResult == b.Result
	wIndex := 0
	if onTeam0 != win {
		wIndex = 1
	}
	lIndex := 1 - wIndex

	battle.TeamW = ids[wIndex]
	battle.TeamL = ids[lIndex]
	return
}

func (b *BattleLogBattle) IsValid() bool {
	return b != nil && !b.Time.IsZero() && b.MapName != "" && len(b.Teams) >= 2
}

func (b *BattleLogBattle) getIDTeams() (res [][]int) {
	for _, team := range b.Teams {
		tmp := make([]int, 0)
		for _, player := range *team {
			tmp = append(tmp, player.BrawlerID)
		}
		res = append(res, tmp)
	}
	return res
}

func (b *BattleLogBattle) getMeanRank() int {
	var total, count int
	for _, team := range b.Teams {
		for _, player := range *team {
			count++
			total += player.Rank
		}
	}
	return (total + count/2) / count
}

func (b *BattleLogBattle) GetTags(curTag string) (tags []string) {
	for _, team := range b.Teams {
		for _, player := range *team {
			if player.Tag != curTag {
				tags = append(tags, player.Tag)
			}
		}
	}
	return
}
