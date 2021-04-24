package prometheus

import (
	"context"
	"fmt"
	"github.com/guodongq/jigsaw/pkg/module"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"
)

type Prometheus struct {
	module.DefaultRunProvider
	Config *Config
	srv    *http.Server
}

func New(config *Config) *Prometheus {
	return &Prometheus{Config: config}
}

func (p *Prometheus) Run() error {
	if !p.Config.Enabled {
		logrus.Infof("Prometheus Provider not enabled")
		return nil
	}
	addr := fmt.Sprintf(":%d", p.Config.Port)
	logEntry := logrus.WithFields(logrus.Fields{
		"addr":     addr,
		"endpoint": p.Config.Endpoint,
	})

	mux := http.NewServeMux()
	mux.Handle(p.Config.Endpoint, promhttp.Handler())
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
	if !p.Config.Enabled || p.srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Error while closing Prometheus server")
	}

	return p.DefaultRunProvider.Close()
}
