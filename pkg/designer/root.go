package designer

import (
	"fmt"
	"os"

	"github.com/guodongq/jigsaw/pkg/designer/serve"

	"github.com/guodongq/jigsaw/pkg/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "designer",
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
