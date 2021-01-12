package settings

import (
	"github.com/spf13/pflag"
	"go.uber.org/config"
	"go.uber.org/fx"
)

var Module = func(flagSets ...*pflag.FlagSet) fx.Option {
	var cfgFile = ""
	for _, v := range flagSets {
		if flagSet := v.Lookup("config"); flagSet != nil {
			if val := flagSet.Value.String(); len(val) > 0 {
				cfgFile = val
				break
			}
		}
	}
	return fx.Provide(func() (*Settings, error) {
		return New(WithConfigFile(cfgFile))
	})
}

type Settings struct {
	*config.YAML
	configFile string
}

func New(opts ...Option) (*Settings, error) {
	var settings Settings
	for _, opt := range opts {
		opt.Apply(&settings)
	}
	if cfgFile := settings.configFile; len(cfgFile) > 0 {
		yml, err := config.NewYAML(config.File(cfgFile))
		if err != nil {
			return nil, err
		}
		settings.YAML = yml
	}
	return &settings, nil
}

func (s *Settings) Enable() bool {
	return s.YAML != nil
}
