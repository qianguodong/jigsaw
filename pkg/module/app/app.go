package app

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

type App struct {
	module.DefaultProvider
	Config *Config
}

func New(config *Config) *App {
	return &App{Config: config}
}

func (p *App) Init() error {
	logrus.WithFields(logrus.Fields{
		"name":     p.Name(),
		"basePath": p.ParsePath(),
	}).Info("App Module initialized")
	return nil
}

func (p *App) Name() string {
	return p.Config.Name
}

func (p *App) ParsePath(elem ...string) string {
	res := p.ParseEndpoint(elem...)
	if !strings.HasSuffix(res, "/") {
		res += "/"
	}
	return res
}

func (p *App) ParseEndpoint(elem ...string) string {
	elem = append([]string{p.Config.BasePath}, elem...)
	return path.Join(elem...)
}
