package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/guodongq/jigsaw/pkg/provider"
	"github.com/guodongq/jigsaw/pkg/provider/settings"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = func() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, p *Prometheus) error {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go p.Run()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return p.Close()
				},
			})
			return p.Init()
		}),
	)
}

type Prometheus struct {
	provider.DefaultProvider
	Settings *settings.Settings
	srv      *http.Server
}

func New(settings *settings.Settings) *Prometheus {
	return &Prometheus{Settings: settings}
}

const (
	defaultPort     = 9090
	defaultEndpoint = "/metrics"
	defaultEnable   = true
)

var prometheusCfg struct {
	Port     int    `yaml:"port"`
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

func (p *Prometheus) Init() error {
	prometheusCfg.Port = defaultPort
	prometheusCfg.Endpoint = defaultEndpoint
	prometheusCfg.Enabled = defaultEnable

	if p.Settings.Enable() {
		if err := p.Settings.Get("prometheus").Populate(&prometheusCfg); err != nil {
			return err
		}
	}
	return nil
}

func (p *Prometheus) Run() error {
	if !prometheusCfg.Enabled {
		logrus.Infof("Prometheus Provider not enabled")
		return nil
	}
	addr := fmt.Sprintf(":%d", prometheusCfg.Port)
	logEntry := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": prometheusCfg.Endpoint,
	})

	mux := http.NewServeMux()
	mux.Handle(prometheusCfg.Endpoint, promhttp.Handler())
	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Info("Prometheus Provider launched")
	if err := p.srv.ListenAndServe(); err != http.ErrServerClosed {
		logEntry.WithError(err).Error("Prometheus Provider launch failed")
		return err
	}

	return nil
}

func (p *Prometheus) Close() error {
	if !prometheusCfg.Enabled || p.srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Error while closing Prometheus server")
	}

	return p.DefaultProvider.Close()
}
