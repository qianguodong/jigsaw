package mongodb

import (
	"context"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewConfig),
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, p *MongoDB) error {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return p.Close()
				},
			})
			return p.Init()
		}),
	)
}
