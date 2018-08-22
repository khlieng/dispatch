package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/assets"
	"github.com/khlieng/dispatch/server"
	"github.com/khlieng/dispatch/storage"
	"github.com/khlieng/dispatch/storage/bleve"
	"github.com/khlieng/dispatch/storage/boltdb"
)

const logo = `
    ____   _                     _         _
   |  _ \ (_) ___  _ __    __ _ | |_  ___ | |__
   | | | || |/ __|| '_ \  / _  || __|/ __|| '_ \
   | |_| || |\__ \| |_) || (_| || |_| (__ | | | |
   |____/ |_||___/| .__/  \__,_| \__|\___||_| |_|
                  |_|
                       v0.4
`

var rootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Web-based IRC client in Go.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use == "dispatch" {
			fmt.Println(logo)
		}

		storage.Initialize(viper.GetString("dir"))

		initConfig(storage.Path.Config(), viper.GetBool("reset_config"))

		viper.SetConfigName("config")
		viper.AddConfigPath(storage.Path.Root())
		viper.ReadInConfig()

		viper.WatchConfig()

		prev := time.Now()
		viper.OnConfigChange(func(e fsnotify.Event) {
			now := time.Now()
			// fsnotify sometimes fires twice
			if now.Sub(prev) > time.Second {
				log.Println("New config loaded")
				prev = now
			}
		})
	},

	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("dev") {
			log.Println("Running in development mode, access client at http://localhost:3000")
		}
		log.Println("Storing data at", storage.Path.Root())

		db, err := boltdb.New(storage.Path.Database())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		srv := server.Dispatch{
			Store:        db,
			SessionStore: db,

			GetMessageStore: func(user *storage.User) (storage.MessageStore, error) {
				return boltdb.New(storage.Path.Log(user.Username))
			},
			GetMessageSearchProvider: func(user *storage.User) (storage.MessageSearchProvider, error) {
				return bleve.New(storage.Path.Index(user.Username))
			},
		}

		srv.Run()
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().String("dir", storage.DefaultDirectory(), "directory to store config and data in")
	rootCmd.PersistentFlags().Bool("reset-config", false, "reset to the default configuration, overwriting the current one")
	rootCmd.Flags().StringP("address", "a", "", "interface to which the server will bind")
	rootCmd.Flags().IntP("port", "p", 80, "port to listen on")
	rootCmd.Flags().Bool("dev", false, "development mode")

	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("reset_config", rootCmd.PersistentFlags().Lookup("reset-config"))
	viper.BindPFlag("address", rootCmd.Flags().Lookup("address"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("dev", rootCmd.Flags().Lookup("dev"))

	viper.SetDefault("hexIP", false)
	viper.SetDefault("verify_client_certificates", true)
}

func initConfig(configPath string, overwrite bool) {
	if _, err := os.Stat(configPath); overwrite || os.IsNotExist(err) {
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
