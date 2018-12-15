package commands

import (
	"fmt"
	"runtime"

	"github.com/khlieng/dispatch/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:              "version",
	Short:            "Show version",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("%s\nCommit: %s\nBuild Date: %s\nRuntime: %s\n", version.Tag, version.Commit, version.Date, runtime.Version())
}
