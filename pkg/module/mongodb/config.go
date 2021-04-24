package mongodb

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/setting"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	URI               string `yaml:"uri"`
	Timeout           int    `yaml:"timeout"`
	MaxPoolSize       int    `yaml:"maxPoolSize"`
	MaxConnIdleTime   int    `yaml:"maxConnIdleTime"`
	HeartBeatInterval int    `yaml:"heartBeatInterval"`
	Database          string `yaml:"database"`
}

const (
	defaultMongoURI          = "mongodb://root:admin@127.0.0.1:27017"
	defaultTimeout           = 20
	defaultMaxPoolSize       = 16
	defaultMaxConnIdleTime   = 30
	defaultHeartBeatInterval = 15
	defaultMongoDatabase     = "test"
)

const (
	mongoURIKey          = "mongodb.uri"
	timeoutKey           = "mongodb.timeout"
	maxPoolSizeKey       = "mongodb.maxPoolSizeKey"
	maxConnIdleTimeKey   = "mongodb.maxConnIdleTime"
	heartBeatIntervalKey = "mongodb.heartBeatInterval"
	mongoDatabaseKey     = "mongodb.database"
)

func NewConfig(s *setting.Setting) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault(mongoURIKey, defaultMongoURI)
	v.SetDefault(timeoutKey, defaultTimeout)
	v.SetDefault(maxPoolSizeKey, defaultMaxPoolSize)
	v.SetDefault(maxConnIdleTimeKey, defaultMaxConnIdleTime)
	v.SetDefault(heartBeatIntervalKey, defaultHeartBeatInterval)
	v.SetDefault(mongoDatabaseKey, defaultMongoDatabase)

	if s.Enable() {
		if err := new(module.Configure).From(v).CfgFile(s.CfgFile).ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		URI:               v.GetString(mongoURIKey),
		Timeout:           v.GetInt(timeoutKey),
		MaxPoolSize:       v.GetInt(maxPoolSizeKey),
		MaxConnIdleTime:   v.GetInt(maxConnIdleTimeKey),
		HeartBeatInterval: v.GetInt(heartBeatIntervalKey),
		Database:          v.GetString(mongoDatabaseKey),
	}, nil
}
