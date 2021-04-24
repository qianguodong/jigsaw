package app

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/setting"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	BasePath string
	Name     string
}

const (
	basePathKey = "app.basePath"
	nameKey     = "app.name"
)

const defaultBasePath = "/"

func NewConfig(s *setting.Setting) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault(basePathKey, defaultBasePath)
	paths := strings.Split(os.Args[0], "/")
	v.SetDefault(nameKey, paths[len(paths)-1])

	if s.Enable() {
		if err := new(module.Configure).From(v).CfgFile(s.CfgFile).ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		BasePath: v.GetString(basePathKey),
		Name:     v.GetString(nameKey),
	}, nil
}
