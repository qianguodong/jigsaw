package prometheus

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/setting"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Port     int
	Enabled  bool
	Endpoint string
}

const (
	portKey     = "prometheus.port"
	enabledKey  = "prometheus.enabled"
	endpointKey = "prometheus.endpoint"
)

const (
	defaultPort     = 9090
	defaultEndpoint = "/metrics"
	defaultEnable   = true
)

func NewConfig(s *setting.Setting) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault(portKey, defaultPort)
	v.SetDefault(enabledKey, defaultEnable)
	v.SetDefault(endpointKey, defaultEndpoint)

	if s.Enable() {
		if err := new(module.Configure).From(v).CfgFile(s.CfgFile).ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		Port:     v.GetInt(portKey),
		Enabled:  v.GetBool(enabledKey),
		Endpoint: v.GetString(endpointKey),
	}, nil
}
