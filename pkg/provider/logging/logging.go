package logging

import (
	"io"
	"os"

	"github.com/guodongq/jigsaw/pkg/provider"

	"github.com/guodongq/jigsaw/pkg/provider/settings"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = func() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(p *Logging) error {
			return p.Init()
		}),
	)
}

type Logging struct {
	provider.DefaultProvider
	Settings *settings.Settings
}

func New(settings *settings.Settings) *Logging {
	return &Logging{Settings: settings}
}

const (
	defaultLoggingLevel     = "info"
	defaultLoggingFormatter = "json"
	defaultLoggingOutput    = "stderr"
)

func (p *Logging) Init() error {
	var loggingCfg struct {
		Level     string `yaml:"level"`
		Formatter string `yaml:"formatter"`
		Output    string `yaml:"output"`
	}
	loggingCfg.Level = defaultLoggingLevel
	loggingCfg.Formatter = defaultLoggingFormatter
	loggingCfg.Output = defaultLoggingOutput

	if p.Settings.Enable() {
		if err := p.Settings.Get("logging").Populate(&loggingCfg); err != nil {
			return err
		}
	}
	// logging level
	level, err := logrus.ParseLevel(loggingCfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// logging formatter
	var formatter logrus.Formatter
	switch loggingCfg.Formatter {
	case "text":
		formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	case "json":
		fallthrough
	default:
		formatter = &logrus.JSONFormatter{
			PrettyPrint: true,
		}
	}
	logrus.SetFormatter(formatter)

	// logging output
	var output io.Writer
	switch loggingCfg.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		fallthrough
	default:
		output = os.Stderr
	}
	logrus.SetOutput(output)
	return nil
}
