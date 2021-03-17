package operator

import (
	"fmt"
	"os"

	"github.com/guodongq/jigsaw/internal/operator/serve"

	"github.com/guodongq/jigsaw/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "operator",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	serve.RegisterCommandRecursive(RootCmd)

	RootCmd.AddCommand(version.Version(&version.BuildVersion, &version.BuildGitHash, &version.BuildTime))
}
