package probes

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/setting"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Enabled           bool
	Port              int
	LivenessEndpoint  string
	ReadinessEndpoint string
}

const (
	defaultProbesEnabled           = true
	defaultProbesPort              = 8000
	defaultProbesLivenessEndpoint  = "/healthz"
	defaultProbesReadinessEndpoint = "/ready"
)

const (
	enabledKey           = "probes.enabled"
	portKey              = "probes.port"
	livenessEndpointKey  = "probes.livenessEndpoint"
	readinessEndpointKey = "probes.readinessEndpoint"
)

func NewConfig(s *setting.Setting) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault(enabledKey, defaultProbesEnabled)
	v.SetDefault(portKey, defaultProbesPort)
	v.SetDefault(livenessEndpointKey, defaultProbesLivenessEndpoint)
	v.SetDefault(readinessEndpointKey, defaultProbesReadinessEndpoint)

	if s.Enable() {
		if err := new(module.Configure).From(v).CfgFile(s.CfgFile).ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		Enabled:           v.GetBool(enabledKey),
		Port:              v.GetInt(portKey),
		LivenessEndpoint:  v.GetString(livenessEndpointKey),
		ReadinessEndpoint: v.GetString(readinessEndpointKey),
	}, nil
}
