package commands

import (
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/spf13/cobra"

	"github.com/khlieng/name_pending/server"
	"github.com/khlieng/name_pending/storage"
)

var (
	rootCmd = &cobra.Command{
		Use:  "name_pending",
		Long: "Web-based IRC client in Go.",
		Run: func(cmd *cobra.Command, args []string) {
			storage.Initialize(development)
			server.Run(port, development)
		},
	}

	development bool
	port        int
)

func init() {
	rootCmd.Flags().BoolVarP(&development, "dev", "d", false, "development mode")
	rootCmd.Flags().IntVarP(&port, "port", "p", 1337, "port to listen on")
}

func Execute() {
	rootCmd.Execute()
}
