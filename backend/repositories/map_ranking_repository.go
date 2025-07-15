package repositories

import (
	"BrawlPicks/internal/env"
	"BrawlPicks/internal/variables"
	"BrawlPicks/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type MapRankingRepositoryInterface interface {
	NewRawMapData(rank models.Rank, data []*models.MapData) (err error)
	LoadRawMapData(rank models.Rank) (mapData []*models.MapData, err error)
	UpdateRawMapDataSection(rank models.Rank, mapData []*models.MapData) (err error)
	ArchiveRawMapData() (err error)
	NewMapRankings(rank models.Rank, rankings []*models.MapRanking) (err error)
	LoadMapRankings(rank models.Rank) (rankings []*models.MapRanking, err error)
	UpdateMapRankingSection(rank models.Rank, mapRankings []*models.MapRanking) (err error)
	ArchiveMapRankings() (err error)
	LoadLastUpdatedTime() (unixTime int, err error)
	UpdateLastUpdatedTime(unixTime int) (err error)
	AddNewMapName(name string) (err error)
	KnownMapName(mapName string) (known bool, err error)
}

type MapRankingRepository struct {
	e *env.Env
}

func NewMapRankingRepository(e *env.Env) *MapRankingRepository {
	return &MapRankingRepository{
		e: e,
	}
}

func (r *MapRankingRepository) NewRawMapData(rank models.Rank, data []*models.MapData) (err error) {
	var (
		filepath = r.e.Data.RawMapData + string(rank) + variables.JsonExtension
	)

	for _, mapData := range data {
		r.AddNewMapName(mapData.MapName)
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return
	}

	err = ensureDir(filepath)
	if err != nil {
		return
	}

	if err = os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (r *MapRankingRepository) LoadRawMapData(rank models.Rank) (mapData []*models.MapData, err error) {
	data, err := os.ReadFile(r.e.Data.RawMapData + string(rank) + variables.JsonExtension)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}

func (r *MapRankingRepository) UpdateRawMapDataSection(rank models.Rank, newData []*models.MapData) (err error) {
	var (
		filepath = r.e.Data.RawMapData + string(rank) + variables.JsonExtension
		curData  []*models.MapData
	)

	if err := ensureDir(filepath); err != nil {
		return err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &curData); err != nil {
		return err
	}

	for _, newD := range newData {
		known, err := r.KnownMapName(newD.MapName)
		if err != nil {
			return fmt.Errorf("KnownMapName: %s", err)
		} else if !known {
			return fmt.Errorf("map name `%s` not found", newD.MapName)
		}

		for i, curD := range curData {
			if newD.MapName == curD.MapName {
				curData[i] = newD
				break
			}
		}
	}

	jsonData, err := json.MarshalIndent(curData, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (r *MapRankingRepository) ArchiveRawMapData() (err error) {
	// TODO: Needs mail setup or other archive ideas
	return nil
}

func (r *MapRankingRepository) NewMapRankings(rank models.Rank, rankings []*models.MapRanking) (err error) {
	var (
		filepath = r.e.Data.MapRanking + string(rank) + variables.JsonExtension
	)

	jsonData, err := json.MarshalIndent(rankings, "", "\t")
	if err != nil {
		return
	}

	if err := ensureDir(filepath); err != nil {
		return err
	}

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (r *MapRankingRepository) LoadMapRankings(rank models.Rank) (rankings []*models.MapRanking, err error) {
	data, err := os.ReadFile(r.e.Data.MapRanking + string(rank) + variables.JsonExtension)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}

func (r *MapRankingRepository) UpdateMapRankingSection(rank models.Rank, newRankings []*models.MapRanking) (err error) {
	var (
		filepath    = r.e.Data.MapRanking + string(rank) + variables.JsonExtension
		curRankings []*models.MapRanking
	)

	if err := ensureDir(filepath); err != nil {
		return err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &curRankings); err != nil {
		return err
	}

	for _, newRanking := range newRankings {
		known, err := r.KnownMapName(newRanking.Name)
		if err != nil {
			return fmt.Errorf("KnownMapName: %s", err)
		} else if !known {
			return fmt.Errorf("map name `%s` not found", newRanking.Name)
		}

		for i, curRanking := range curRankings {
			if newRanking.Name == curRanking.Name {
				curRankings[i] = newRanking
				break
			}
		}
	}

	jsonData, err := json.MarshalIndent(curRankings, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (r *MapRankingRepository) ArchiveMapRankings() (err error) {
	return nil
}

func (r *MapRankingRepository) LoadLastUpdatedTime() (unixTime int, err error) {
	var (
		filepath    = r.e.Data.LastUpdated + variables.JsonExtension
		lastUpdated = &models.LastUpdated{}
	)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return -1, err
	}

	if err := json.Unmarshal(data, lastUpdated); err != nil {
		return -1, err
	}
	return lastUpdated.Time, nil
}

func (r *MapRankingRepository) UpdateLastUpdatedTime(unixTime int) (err error) {
	var (
		filepath    = r.e.Data.LastUpdated + variables.JsonExtension
		lastUpdated = models.LastUpdated{
			Time: unixTime,
		}
	)

	if err := ensureDir(filepath); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(lastUpdated, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (r *MapRankingRepository) AddNewMapName(name string) (err error) {
	known, err := r.KnownMapName(name)
	if err != nil {
		return err
	} else if known {
		return nil
	}

	var (
		filepath = r.e.Data.MapNames + variables.JsonExtension
		names    []string
	)

	if err := ensureDir(filepath); err != nil {
		return err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &names); err != nil {
		return err
	}

	names = append(names, name)

	jsonData, err := json.MarshalIndent(names, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (r *MapRankingRepository) KnownMapName(mapName string) (known bool, err error) {
	var (
		filepath = r.e.Data.MapNames + variables.JsonExtension
		names    []string
	)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(data, &names); err != nil {
		return false, err
	}

	if slices.Contains(names, mapName) {
		return true, nil
	}
	return false, nil
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}
