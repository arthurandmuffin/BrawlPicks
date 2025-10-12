package env

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Env struct {
	Redis   *Redis   `yaml:"redis"`
	Brawl   *Brawl   `yaml:"brawl"`
	Scraper *Scraper `yaml:"scraper"`
	Monitor *Monitor `yaml:"monitor"`
	Storage *Storage `yaml:"storage"`
}

type Redis struct {
	Credentials      *RedisCredentials `yaml:"credentials"`
	PlayerQueueName  string            `yaml:"queueName"`
	PlayerQueueLimit int64             `yaml:"queueLimit"`
	PlayerBFPrefix   string            `yaml:"playerBFPrefix"`
	PlayerBFCapacity int64             `yaml:"playerBFCapacity"`
	PlayerBFTTL      int64             `yaml:"playerBFTTL"`
	GameBFPrefix     string            `yaml:"gameBFPrefix"`
	GameBFCapacity   int64             `yaml:"gameBFCapacity"`
	GameBFTTL        int64             `yaml:"gameBFTTL"`
	BFErrorRate      float64           `yaml:"bfErrorRate"`
}

type RedisCredentials struct {
	Address    string `yaml:"address"`
	Password   string `yaml:"password"`
	MasterName string `yaml:"masterName"`
}

type Brawl struct {
	TopPlayersEndpoint string `yaml:"topPlayersEndpoint"`
	BattleLogEndpoint  string `yaml:"battleLogEndpoint"`
	Key                string `yaml:"key"`
}

type Scraper struct {
	Seeding   *ScraperSeeding   `yaml:"seeding"`
	RateLimit *ScraperRateLimit `yaml:"rateLimit"`
	Workers   *ScraperWorkers   `yaml:"workers"`
	Queue     *ScraperQueue     `yaml:"queue"`
}

type ScraperSeeding struct {
	Threshold       int64 `yaml:"threshold"`
	CooldownSeconds int   `yaml:"cooldownSeconds"`
}

type ScraperRateLimit struct {
	QPS   int `yaml:"qps"`
	Burst int `yaml:"burst"`
}

type ScraperWorkers struct {
	IO  int `yaml:"io"`
	CPU int `yaml:"cpu"`
}

type ScraperQueue struct {
	Batch           int `yaml:"batch"`
	Low             int `yaml:"low"`
	High            int `yaml:"high"`
	ChannelSize     int `yaml:"channelSize"`
	CapacityTrigger int `yaml:"capacityTrigger"`
}

type Monitor struct {
	HeartbeatInterval int `yaml:"heartbeatInterval"`
}

type Storage struct {
	BattleLog *BattleLogStorage `yaml:"battleLog"`
	Synergy   *SynergyStorage   `yaml:"synergy"`
}

type BattleLogStorage struct {
	Dir          string `yaml:"dir"`
	MaxRows      int    `yaml:"maxRows"`
	FlushSeconds int    `yaml:"flushSeconds"`
}

type SynergyStorage struct {
	Dir           string `yaml:"dir"`
	RetentionDays int    `yaml:"retentionDays"`
	FlushSeconds  int    `yaml:"flushSeconds"`
}

func Get(path string) (env *Env, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	env = &Env{}
	if err = yaml.Unmarshal(data, env); err != nil {
		return
	}
	return
}

func (e *Env) Info() {
	strs := []string{}

	v := reflect.ValueOf(*e)
	for i := 0; i < v.NumField(); i++ {
		strs = append(strs, fmt.Sprintf("%+v", v.Field(i).Interface()))
	}

	logrus.Info(strings.Join(strs, ", "))
}
