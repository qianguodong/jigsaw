package setting

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Module(cmd *cobra.Command) fx.Option {
	return fx.Provide(func() *Setting {
		return New(cmd)
	})
}
