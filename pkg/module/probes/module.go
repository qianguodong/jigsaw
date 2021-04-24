package probes

import (
	"context"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewConfig),
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, p *Probes) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go p.Run()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return p.Close()
				},
			})
		}),
	)
}
