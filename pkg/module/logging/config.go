package logging

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/setting"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

const (
	defaultLevel     = "info"
	defaultFormatter = "json"
	defaultOutput    = "stderr"
)

const (
	levelKey     = "logging.level"
	formatterKey = "logging.formatter"
	outputKey    = "logging.output"
)

type Config struct {
	Level     logrus.Level
	Formatter logrus.Formatter
	Output    io.Writer
}

func NewConfig(s *setting.Setting) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault(levelKey, defaultLevel)
	v.SetDefault(formatterKey, defaultFormatter)
	v.SetDefault(outputKey, defaultOutput)

	if s.Enable() {
		if err := new(module.Configure).From(v).CfgFile(s.CfgFile).ReadInConfig(); err != nil {
			return nil, err
		}
	}

	level, err := logrus.ParseLevel(v.GetString(levelKey))
	if err != nil {
		level = logrus.InfoLevel
	}

	var formatter logrus.Formatter
	switch v.GetString(formatterKey) {
	case "text":
		formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	case "json":
		fallthrough
	default:
		formatter = &logrus.JSONFormatter{
			//PrettyPrint: true,
		}
	}

	var output io.Writer
	switch v.GetString(outputKey) {
	case "stdout":
		output = os.Stdout
	case "stderr":
		fallthrough
	default:
		output = os.Stderr
	}
	return &Config{
		Level:     level,
		Formatter: formatter,
		Output:    output,
	}, nil
}
