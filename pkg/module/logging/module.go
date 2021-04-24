package logging

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewConfig),
		fx.Provide(New),
		fx.Invoke(func(p *Logging) error {
			return p.Init()
		}),
	)
}
