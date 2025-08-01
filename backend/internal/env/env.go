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
	Api        *Api        `yaml:"api"`
	Data       *Data       `yaml:"data"`
	Upstream   *Upstream   `yaml:"upstream"`
	MapRanking *MapRanking `yaml:"mapRanking"`
	Redis      *Redis      `yaml:"redis"`
	Brawl      *Brawl      `yaml:"brawl"`
}

type Api struct {
	Port   string
	Prefix string
	Cors   bool
	Debug  bool
}

type Data struct {
	RawMapData  string `yaml:"rawMapData"`
	MapRanking  string `yaml:"mapRanking"`
	LastUpdated string `yaml:"lastUpdated"`
	MapNames    string `yaml:"mapNames"`
}

type Upstream struct {
	MatchData *MatchData `yaml:"matchData"`
}

type MatchData struct {
	BasePath            string            `yaml:"basePath"`
	Endpoints           map[string]string `yaml:"endpoints"`
	LastUpdatedEndpoint string            `yaml:"lastUpdatedEndpoint"`
}

type MapRanking struct {
	WinRateWeight float64 `yaml:"winRateWeight"`
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
	QueueLimit         int    `yaml:"queueLimit"`
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
