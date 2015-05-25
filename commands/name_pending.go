package commands

import (
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/khlieng/name_pending/server"
	"github.com/khlieng/name_pending/storage"
)

var (
	rootCmd = &cobra.Command{
		Use:   "name_pending",
		Short: "Web-based IRC client in Go.",
		Run: func(cmd *cobra.Command, args []string) {
			storage.Initialize()
			server.Run(viper.GetInt("port"))
		},
	}
)

func init() {
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.Flags().IntP("port", "p", 1337, "port to listen on")

	viper.SetConfigName("config")
	viper.AddConfigPath(storage.AppDir)
	viper.ReadInConfig()

	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
}

func Execute() {
	rootCmd.Execute()
}
