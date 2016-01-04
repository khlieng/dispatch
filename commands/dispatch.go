package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/khlieng/dispatch/assets"
	"github.com/khlieng/dispatch/server"
	"github.com/khlieng/dispatch/storage"
)

const logo = `
    ____   _                     _         _
   |  _ \ (_) ___  _ __    __ _ | |_  ___ | |__
   | | | || |/ __|| '_ \  / _  || __|/ __|| '_ \
   | |_| || |\__ \| |_) || (_| || |_| (__ | | | |
   |____/ |_||___/| .__/  \__,_| \__|\___||_| |_|
                  |_|
                       v0.1
`

var rootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Web-based IRC client in Go.",
	Run: func(cmd *cobra.Command, args []string) {
		storage.Initialize()
		server.Run()
	},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use == "dispatch" {
			fmt.Println(logo)
		}

		storage.SetDirectory(viper.GetString("dir"))

		err := os.MkdirAll(storage.Path.Logs(), 0700)
		if err != nil {
			log.Fatal(err)
		}

		initConfig()

		viper.SetConfigName("config")
		viper.AddConfigPath(storage.Path.Root())
		viper.ReadInConfig()
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.Flags().IntP("port", "p", 80, "port to listen on")
	rootCmd.PersistentFlags().String("dir", storage.DefaultDirectory(), "directory to store config and data in")

	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))
}

func initConfig() {
	configPath := storage.Path.Config()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config, err := assets.Asset("config.default.toml")
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Writing default config to", configPath)

		err = ioutil.WriteFile(configPath, config, 0600)
		if err != nil {
			log.Println(err)
		}
	}
}
