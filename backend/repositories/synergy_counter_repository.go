package repositories

import (
	"BrawlPicks/internal/env"
	"BrawlPicks/internal/utils"
	"BrawlPicks/models"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type SynergyCounterRepositoryInterface interface {
	RecordBattle(battle *models.Battle)
	Close() error
}

type SynergyCounterRepository struct {
	mu            sync.Mutex
	matrices      map[string]map[string]*models.SynergyMatrix
	e             *env.Env
	flushTicker   *time.Ticker
	retentionDays int
	matricesDir   string
}

func NewSynergyCounterRepository(e *env.Env, flushTime time.Duration, retention int, matricesDir string) (*SynergyCounterRepository, error) {
	r := &SynergyCounterRepository{
		matrices:      make(map[string]map[string]*models.SynergyMatrix),
		e:             e,
		flushTicker:   time.NewTicker(flushTime),
		retentionDays: retention,
		matricesDir:   matricesDir,
	}
	if err := utils.EnsureDir(path.Join(r.matricesDir, "placeholder.json")); err != nil {
		return nil, err
	}
	if err := r.loadRecent(); err != nil {
		return nil, err
	}
	go r.autoFlush()
	return r, nil
}

func (r *SynergyCounterRepository) RecordBattle(battle *models.Battle) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dateStr := battle.Timestamp.Format("2006-01-02")
	if _, ok := r.matrices[dateStr]; !ok {
		r.matrices[dateStr] = make(map[string]*models.SynergyMatrix)
	}
	synergyMatrix, ok := r.matrices[dateStr][battle.MapName]
	if !ok || synergyMatrix == nil {
		synergyMatrix = models.NewSynergyMatrix()
		r.matrices[dateStr][battle.MapName] = synergyMatrix
	}

	for _, comb := range battle.SynergyCombinations(true) {
		synergyMatrix.IncrementSynergy(comb, true, battle.Draw)
	}
	for _, comb := range battle.SynergyCombinations(false) {
		synergyMatrix.IncrementSynergy(comb, false, battle.Draw)
	}
	for _, comb := range battle.CounterCombinations() {
		synergyMatrix.IncrementCounter(comb, battle.Draw)
	}
}

func (r *SynergyCounterRepository) loadRecent() (err error) {
	cutoff := r.getCutoffDate()
	files, err := os.ReadDir(r.matricesDir)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() || path.Ext(file.Name()) != ".json" {
			continue
		}
		dateStr := strings.TrimSuffix(file.Name(), ".json")
		dateStr = strings.TrimPrefix(dateStr, "synergy-")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
		if !date.Before(cutoff) {
			b, err := os.ReadFile(path.Join(r.matricesDir, file.Name()))
			if err != nil {
				return err
			}
			var dailyData map[string]*models.SynergyMatrix
			if err = json.Unmarshal(b, &dailyData); err != nil {
				return err
			}
			r.mu.Lock()
			r.matrices[dateStr] = dailyData
			r.mu.Unlock()
		}
	}
	return nil
}

func (r *SynergyCounterRepository) autoFlush() {
	for range r.flushTicker.C {
		r.mu.Lock()
		_ = r.flushLocked()
		r.mu.Unlock()
	}
}

func (r *SynergyCounterRepository) Close() error {
	r.flushTicker.Stop()
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.flushLocked()
}

func (r *SynergyCounterRepository) flushLocked() error {
	cutoff := r.getCutoffDate()
	for dateStr, dailyData := range r.matrices {
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
		if err := r.flushDay(dateStr, dailyData); err != nil {
			return err
		}
		if d.Before(cutoff) {
			delete(r.matrices, dateStr)
		}
	}
	return nil
}

func (r *SynergyCounterRepository) flushDay(dateStr string, dailyData map[string]*models.SynergyMatrix) (err error) {
	data, err := json.MarshalIndent(dailyData, "", " ")
	if err != nil {
		return
	}
	filepath := path.Join(r.matricesDir, fmt.Sprintf("%s.json", dateStr))
	if err := utils.EnsureDir(filepath); err != nil {
		return err
	}
	err = os.WriteFile(filepath, data, 0644)
	return
}

func (r *SynergyCounterRepository) getCutoffDate() time.Time {
	return utils.OffsetDate(time.Now().UTC(), -r.retentionDays)
}
