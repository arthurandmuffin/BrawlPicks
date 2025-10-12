package repositories

import (
	env "BrawlPicks/scraper/config"
	"BrawlPicks/scraper/services/upstream"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type ProcessedRecordRepositoryInterface interface {
	AddPlayersToQueue(ctx context.Context, tags []string) (err error)
	PopPlayersFromQueue(ctx context.Context, count int) (res []string, err error)
	GetPlayerQueueLength(ctx context.Context) (length int64, err error)
	AddPlayersToBF(ctx context.Context, tags []string) (newTags []string, err error)
	AddGamesToBF(ctx context.Context, games []*upstream.BattleLogBattle) (res []bool, err error)
}

type ProcessedRecordRepository struct {
	client *redis.Client
	e      *env.Env
}

func NewProcessedRecordRepository(client *redis.Client, e *env.Env) *ProcessedRecordRepository {
	return &ProcessedRecordRepository{
		client: client,
		e:      e,
	}
}

func (r *ProcessedRecordRepository) AddPlayersToQueue(ctx context.Context, tags []string) (err error) {
	if len(tags) == 0 {
		return nil
	}

	var (
		queueKey   = r.e.Redis.PlayerQueueName
		queueLimit = r.e.Redis.PlayerQueueLimit
		args       = make([]interface{}, len(tags))
	)

	for i, tag := range tags {
		args[i] = tag
	}

	currentLength, err := r.GetPlayerQueueLength(ctx)
	if err != nil {
		return err
	}

	capacityTrigger := queueLimit * int64(r.e.Redis.CapacityTrigger) / 100
	if currentLength > capacityTrigger {
		return nil
	}

	cmd := r.client.RPush(ctx, queueKey, args...)
	length, err := cmd.Result()
	if err != nil {
		return
	}

	if length > queueLimit {
		logrus.WithFields(logrus.Fields{
			"queue_limit": r.e.Redis.PlayerQueueLimit,
			"trimmed":     length - queueLimit,
		}).Info("queue-truncated-to-limit")
		cmd := r.client.LTrim(ctx, queueKey, 0, queueLimit-1)
		if err = cmd.Err(); err != nil {
			return
		}
	}
	return nil
}

func (r *ProcessedRecordRepository) PopPlayersFromQueue(ctx context.Context, count int) (res []string, err error) {
	cmd := r.client.LPopCount(ctx, r.e.Redis.PlayerQueueName, count)
	res, err = cmd.Result()
	if err == redis.Nil {
		return []string{}, nil
	}
	return res, err
}

func (r *ProcessedRecordRepository) GetPlayerQueueLength(ctx context.Context) (length int64, err error) {
	cmd := r.client.LLen(ctx, r.e.Redis.PlayerQueueName)
	return cmd.Result()
}

func (r *ProcessedRecordRepository) AddPlayersToBF(ctx context.Context, tags []string) (newTags []string, err error) {
	if len(tags) == 0 {
		return nil, nil
	}

	var (
		pipe       = r.client.Pipeline()
		date       = time.Now().UTC().Format("2006-01-02")
		expiryDate = expiryDate(time.Now().UTC(), int(r.e.Redis.PlayerBFTTL))
		key        = r.e.Redis.PlayerBFPrefix + date
		insertCmd  = &redis.BoolSliceCmd{}
	)

	insertCmd = pipe.BFInsert(
		ctx,
		key,
		&redis.BFInsertOptions{
			Capacity: r.e.Redis.PlayerBFCapacity,
			Error:    r.e.Redis.BFErrorRate,
		},
		interfaceSlice(tags)...,
	)
	pipe.ExpireAt(ctx, key, expiryDate)
	if _, err = pipe.Exec(ctx); err != nil {
		return nil, err
	}

	for i, val := range insertCmd.Val() {
		if val {
			newTags = append(newTags, tags[i])
		}
	}
	return
}

// Group all records of same date in same BFInsert?
func (r *ProcessedRecordRepository) AddGamesToBF(ctx context.Context, games []*upstream.BattleLogBattle) (res []bool, err error) {
	if len(games) == 0 {
		return nil, nil
	}

	var (
		pipe       = r.client.Pipeline()
		keys       = make(map[string]struct{})
		daysToLive = r.e.Redis.GameBFTTL
		insertCmds = []*redis.BoolSliceCmd{}
	)

	for _, game := range games {
		if game == nil || !game.IsValid() {
			continue
		}

		date := game.Time.Format("2006-01-02")
		keys[date] = struct{}{}
		insertCmds = append(insertCmds,
			pipe.BFInsert(
				ctx,
				r.e.Redis.GameBFPrefix+date,
				&redis.BFInsertOptions{
					Capacity: r.e.Redis.GameBFCapacity,
					Error:    r.e.Redis.BFErrorRate,
				},
				strconv.FormatInt(game.Time.Unix(), 10)+game.StarPlayerTag,
			),
		)
	}
	if len(insertCmds) == 0 {
		return nil, nil
	}

	for date := range keys {
		datetime, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		expiryDate := expiryDate(datetime, int(daysToLive))
		pipe.ExpireAt(
			ctx,
			r.e.Redis.GameBFPrefix+date,
			expiryDate,
		)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	for _, cmd := range insertCmds {
		res = append(res, cmd.Val()...)
	}
	return
}

func interfaceSlice[T any](in []T) []interface{} {
	out := make([]interface{}, len(in))
	for i, s := range in {
		out[i] = s
	}
	return out
}

func expiryDate(startTime time.Time, daysToLive int) time.Time {
	y, m, d := startTime.Date()
	utcMidnight := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	expiryDate := utcMidnight.AddDate(0, 0, daysToLive)
	return expiryDate
}
