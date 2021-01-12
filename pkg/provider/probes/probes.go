package probes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"go.uber.org/fx"

	"github.com/guodongq/jigsaw/pkg/provider"
	"github.com/guodongq/jigsaw/pkg/provider/app"
	"github.com/guodongq/jigsaw/pkg/provider/settings"
	"github.com/sirupsen/logrus"
)

var Module = func() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, p *Probes) error {
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

type ProbeFunc func() error

type Probes struct {
	provider.DefaultProvider

	Settings        *settings.Settings
	livenessProbes  []ProbeFunc
	readinessProbes []ProbeFunc
	appProvider     *app.App
	srv             *http.Server
}

var probesCfg struct {
	Enabled           bool   `yaml:"enabled"`
	Port              int    `yaml:"port"`
	LivenessEndpoint  string `yaml:"livenessEndpoint"`
	ReadinessEndpoint string `yaml:"readinessEndpoint"`
}

const (
	defaultProbesEnabled           = true
	defaultProbesPort              = 8000
	defaultProbesLivenessEndpoint  = "/healthz"
	defaultProbesReadinessEndpoint = "/ready"
)

func New(settings *settings.Settings, appProvider *app.App) *Probes {
	return &Probes{
		Settings:    settings,
		appProvider: appProvider,
	}
}

func (p *Probes) Init() error {
	probesCfg.Enabled = defaultProbesEnabled
	probesCfg.Port = defaultProbesPort
	probesCfg.LivenessEndpoint = defaultProbesLivenessEndpoint
	probesCfg.ReadinessEndpoint = defaultProbesReadinessEndpoint

	if p.Settings.Enable() {
		if err := p.Settings.Get("probes").Populate(&probesCfg); err != nil {
			return err
		}
	}
	return nil
}

func (p *Probes) Run() error {
	if !probesCfg.Enabled {
		logrus.Infof("Probes Provider not enabled")
		return nil
	}
	addr := fmt.Sprintf(":%d", probesCfg.Port)
	livenessEndpoint := p.appProvider.ParseEndpoint(probesCfg.LivenessEndpoint)
	readinessEndpoint := p.appProvider.ParseEndpoint(probesCfg.ReadinessEndpoint)
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
	if !probesCfg.Enabled || p.srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Error while closing Prometheus server")
	}

	return p.DefaultProvider.Close()
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
