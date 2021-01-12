package app

import (
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"go.uber.org/fx"

	"github.com/guodongq/jigsaw/pkg/provider"
	"github.com/guodongq/jigsaw/pkg/provider/settings"
)

var Module = func() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(p *App) error {
			return p.Init()
		}),
	)
}

type App struct {
	provider.DefaultProvider
	Settings *settings.Settings
	name     string
	basePath string
}

func New(settings *settings.Settings) *App {
	return &App{Settings: settings}
}

const defaultBasePath = "/"

var appCfg struct {
	BasePath string `yaml:"basePath"`
	Name     string `yaml:"name"`
}

func (p *App) Init() error {
	appCfg.BasePath = defaultBasePath
	paths := strings.Split(os.Args[0], "/")
	appCfg.Name = paths[len(paths)-1]

	if p.Settings.Enable() {
		if err := p.Settings.Get("app").Populate(&appCfg); err != nil {
			return err
		}
	}
	logrus.WithFields(logrus.Fields{
		"name":     p.Name(),
		"basePath": p.ParsePath(),
	}).Info("App Provider initialized")
	return nil
}

func (p *App) Name() string {
	return appCfg.Name
}

func (p *App) ParsePath(elem ...string) string {
	res := p.ParseEndpoint(elem...)
	if !strings.HasSuffix(res, "/") {
		res += "/"
	}
	return res
}

func (p *App) ParseEndpoint(elem ...string) string {
	elem = append([]string{appCfg.BasePath}, elem...)
	return path.Join(elem...)
}
