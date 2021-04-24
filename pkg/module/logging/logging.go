package logging

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/sirupsen/logrus"
)

type Logging struct {
	module.DefaultProvider
	Config *Config
}

func New(config *Config) *Logging {
	return &Logging{Config: config}
}

func (p *Logging) Init() error {
	logrus.SetLevel(p.Config.Level)

	if p.Config.Formatter != nil {
		logrus.SetFormatter(p.Config.Formatter)
	}
	if p.Config.Output != nil {
		logrus.SetOutput(p.Config.Output)
	}

	return nil
}
