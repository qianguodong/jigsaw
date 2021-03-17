package serve

import (
	"github.com/guodongq/jigsaw/pkg/provider/app"
	"github.com/guodongq/jigsaw/pkg/provider/logging"
	"github.com/guodongq/jigsaw/pkg/provider/probes"
	"github.com/guodongq/jigsaw/pkg/provider/prometheus"
	"github.com/guodongq/jigsaw/pkg/provider/settings"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the Jigsaw Model Designer server",
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(
			settings.Module(cmd.PersistentFlags()),
			logging.Module(),
			prometheus.Module(),
			app.Module(),
			probes.Module(),
		).Run()
	},
}

func RegisterCommandRecursive(parent *cobra.Command) {
	parent.AddCommand(serveCmd)
}

func init() {
	serveCmd.PersistentFlags().StringP("config", "c", "", "Path to  .yaml config files. Values are loaded in the order provided, meaning that the last config file overwrites values from the previous config file.")
}
