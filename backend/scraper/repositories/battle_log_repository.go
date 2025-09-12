package repositories

import (
	"BrawlPicks/internal/utils"
	"BrawlPicks/scraper/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/apache/arrow/go/v17/arrow/memory"
	"github.com/apache/arrow/go/v17/parquet"
	"github.com/apache/arrow/go/v17/parquet/compress"
	"github.com/apache/arrow/go/v17/parquet/pqarrow"
)

type BattleLogRepositoryInterface interface {
	WriteBattleLog(battle models.Battle) error
	Close() error
}

type BattleLogRepository struct {
	mu          sync.Mutex
	buffer      map[int][]models.Battle
	maxRows     int
	flushTime   time.Duration
	outDir      string
	schema      *arrow.Schema
	flushTicker *time.Ticker
}

func NewBattleLogRepository(maxRows int, flushTime time.Duration, outDir string) *BattleLogRepository {
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "timestamp", Type: arrow.FixedWidthTypes.Timestamp_s},
		{Name: "map_name", Type: arrow.BinaryTypes.String},
		{Name: "rank", Type: arrow.PrimitiveTypes.Int64},
		{Name: "team_W", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
		{Name: "team_L", Type: arrow.ListOf(arrow.PrimitiveTypes.Int64)},
		{Name: "draw_flag", Type: arrow.FixedWidthTypes.Boolean},
	}, nil)
	r := &BattleLogRepository{
		maxRows:     maxRows,
		buffer:      make(map[int][]models.Battle),
		flushTime:   flushTime,
		outDir:      outDir,
		schema:      schema,
		flushTicker: time.NewTicker(flushTime),
	}
	go r.autoFlush()
	return r
}

func (r *BattleLogRepository) WriteBattleLog(battle models.Battle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rank := battle.Rank

	if _, ok := r.buffer[rank]; !ok {
		r.buffer[rank] = make([]models.Battle, 0, r.maxRows)
	}

	r.buffer[rank] = append(r.buffer[rank], battle)
	if len(r.buffer[rank]) >= r.maxRows {
		return r.flushLocked(rank)
	}
	return nil
}

func (r *BattleLogRepository) Close() error {
	r.flushTicker.Stop()
	r.mu.Lock()
	defer r.mu.Unlock()
	for rank, data := range r.buffer {
		if len(data) > 0 {
			if err := r.flushLocked(rank); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *BattleLogRepository) autoFlush() {
	for range r.flushTicker.C {
		r.mu.Lock()
		for rank, data := range r.buffer {
			if len(data) > 0 {
				_ = r.flushLocked(rank)
			}
		}
		r.mu.Unlock()
	}
}

func (r *BattleLogRepository) flushLocked(rank int) (err error) {
	rows := len(r.buffer[rank])
	if rows == 0 {
		return
	}

	pool := memory.DefaultAllocator

	timestampBuilder := array.NewTimestampBuilder(pool, arrow.FixedWidthTypes.Timestamp_s.(*arrow.TimestampType))
	mapNameBuilder := array.NewStringBuilder(pool)
	rankBuilder := array.NewInt64Builder(pool)
	drawFlagBuilder := array.NewBooleanBuilder(pool)

	winTeamBuilder := array.NewListBuilder(pool, arrow.PrimitiveTypes.Int64)
	winTeamMemberBuilder := winTeamBuilder.ValueBuilder().(*array.Int64Builder)

	loseTeamBuilder := array.NewListBuilder(pool, arrow.PrimitiveTypes.Int64)
	loseTeamMemberBuilder := loseTeamBuilder.ValueBuilder().(*array.Int64Builder)

	for _, battle := range r.buffer[rank] {
		timestampBuilder.Append(arrow.Timestamp(battle.Timestamp.Unix()))
		mapNameBuilder.Append(battle.MapName)
		rankBuilder.Append(int64(battle.Rank))
		drawFlagBuilder.Append(battle.Draw)

		winTeamBuilder.Append(true)
		for _, id := range battle.TeamW {
			winTeamMemberBuilder.Append(int64(id))
		}

		loseTeamBuilder.Append(true)
		for _, id := range battle.TeamL {
			loseTeamMemberBuilder.Append(int64(id))
		}
	}

	timestampArray := timestampBuilder.NewArray()
	defer timestampArray.Release()
	mapNameArray := mapNameBuilder.NewArray()
	defer mapNameArray.Release()
	rankArray := rankBuilder.NewArray()
	defer rankArray.Release()
	drawFlagArray := drawFlagBuilder.NewArray()
	defer drawFlagArray.Release()
	winTeamArray := winTeamBuilder.NewArray()
	defer winTeamArray.Release()
	loseTeamArray := loseTeamBuilder.NewArray()
	defer loseTeamArray.Release()

	batch := array.NewRecord(
		r.schema,
		[]arrow.Array{
			timestampArray,
			mapNameArray,
			rankArray,
			winTeamArray,
			loseTeamArray,
			drawFlagArray,
		},
		int64(rows),
	)
	defer batch.Release()

	rankDir := filepath.Join(r.outDir, strconv.Itoa(rank))
	if err := utils.EnsureDir(filepath.Join(rankDir, "placeholder.parquet")); err != nil {
		return err
	}

	t := time.Now().UTC().UnixNano()
	filename := fmt.Sprintf("batch-%d.parquet", t)
	path := filepath.Join(rankDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	pqWriterProps := parquet.NewWriterProperties(
		parquet.WithCompression(compress.Codecs.Snappy),
	)
	parquetWriter, err := pqarrow.NewFileWriter(r.schema, f, pqWriterProps, pqarrow.ArrowWriterProperties{})
	if err != nil {
		return err
	}
	if err := parquetWriter.Write(batch); err != nil {
		parquetWriter.Close()
		return err
	}
	if err := parquetWriter.Close(); err != nil {
		return err
	}

	r.buffer[rank] = r.buffer[rank][:0]
	return nil
}
