package probes

import (
	"context"
	"fmt"
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/app"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

type ProbeFunc func() error

type Probes struct {
	module.DefaultRunProvider
	Config          *Config
	livenessProbes  []ProbeFunc
	readinessProbes []ProbeFunc
	appProvider     *app.App
	srv             *http.Server
}

func New(config *Config, appProvider *app.App) *Probes {
	return &Probes{
		Config:      config,
		appProvider: appProvider,
	}
}

func (p *Probes) Run() error {
	if !p.Config.Enabled {
		logrus.Infof("Probes Provider not enabled")
		return nil
	}
	addr := fmt.Sprintf(":%d", p.Config.Port)
	livenessEndpoint := p.appProvider.ParseEndpoint(p.Config.LivenessEndpoint)
	readinessEndpoint := p.appProvider.ParseEndpoint(p.Config.ReadinessEndpoint)
	logEntry := logrus.WithFields(logrus.Fields{
		"addr":              addr,
		"livenessEndpoint":  livenessEndpoint,
		"readinessEndpoint": readinessEndpoint,
	})
	mux := http.NewServeMux()
	mux.HandleFunc(livenessEndpoint, p.livenessHandler)
	mux.HandleFunc(readinessEndpoint, p.readinessHandler)

	p.srv = &http.Server{Addr: addr, Handler: mux}
	p.SetRunning(true)

	logEntry.Info("Probes Provider Launched")
	if err := p.srv.ListenAndServe(); err != http.ErrServerClosed {
		logEntry.WithError(err).Error("Probes Provider launch failed")
	}
	return nil
}

func (p *Probes) Close() error {
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

func (p *Probes) livenessHandler(res http.ResponseWriter, req *http.Request) {
	reqDump, _ := httputil.DumpRequest(req, false)
	logrus.WithField("req", string(reqDump)).Debug("Handling liveness request")
	for _, probe := range p.livenessProbes {
		if err := probe(); err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			if _, err := res.Write([]byte(err.Error())); err != nil {
				logrus.WithError(err).Warnf("Error while writing liveness data")
			}
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (p *Probes) readinessHandler(res http.ResponseWriter, req *http.Request) {
	reqDump, _ := httputil.DumpRequest(req, false)
	logrus.WithField("req", string(reqDump)).Debug("Handling readiness request")
	for _, probe := range p.readinessProbes {
		if err := probe(); err != nil {
			res.WriteHeader(http.StatusServiceUnavailable)
			if _, err := res.Write([]byte(err.Error())); err != nil {
				logrus.WithError(err).Warnf("Error while writing readiness data")
			}
			return
		}
	}
	res.WriteHeader(http.StatusOK)
}

func (p *Probes) AddLivenessProbes(fn ProbeFunc) {
	p.livenessProbes = append(p.livenessProbes, fn)
}

func (p *Probes) AddReadinessProbes(fn ProbeFunc) {
	p.readinessProbes = append(p.readinessProbes, fn)
}
