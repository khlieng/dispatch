package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/khlieng/dispatch/assets"
	"github.com/khlieng/dispatch/config"
	"github.com/khlieng/dispatch/server"
	"github.com/khlieng/dispatch/storage"
	"github.com/khlieng/dispatch/storage/bleve"
	"github.com/khlieng/dispatch/storage/boltdb"
	"github.com/khlieng/dispatch/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const logo = `
    ____   _                     _         _
   |  _ \ (_) ___  _ __    __ _ | |_  ___ | |__
   | | | || |/ __|| '_ \  / _  || __|/ __|| '_ \
   | |_| || |\__ \| |_) || (_| || |_| (__ | | | |
   |____/ |_||___/| .__/  \__,_| \__|\___||_| |_|
                  |_|

   %s
   Commit: %s
   Build Date: %s
   Runtime: %s

`

var rootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Web-based IRC client in Go.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("version") {
			printVersion()
			os.Exit(0)
		}

		if cmd == cmd.Root() {
			fmt.Printf(logo, version.Tag, version.Commit, version.Date, runtime.Version())
		}

		storage.Initialize(viper.GetString("dir"), viper.GetString("data"), viper.GetString("conf"))

		err := initConfig(storage.Path.Config(), viper.GetBool("reset-config"))
		if err != nil {
			log.Fatal(err)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("dev") {
			log.Println("Running in development mode, access client at http://localhost:3000")
		}
		log.Println("Storing data at", storage.Path.DataRoot())

		db, err := boltdb.New(storage.Path.Database())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		storage.GetMessageStore = func(user *storage.User) (storage.MessageStore, error) {
			return boltdb.New(storage.Path.Log(user.Username))
		}

		storage.GetMessageSearchProvider = func(user *storage.User) (storage.MessageSearchProvider, error) {
			return bleve.New(storage.Path.Index(user.Username))
		}

		cfg, cfgUpdated := config.LoadConfig()
		dispatch := server.New(cfg)

		go func() {
			for {
				dispatch.SetConfig(<-cfgUpdated)
				log.Println("New config loaded")
			}
		}()

		dispatch.Store = db
		dispatch.SessionStore = db

		dispatch.Run()
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().String("data", storage.DefaultDirectory(), "directory to store data in")
	rootCmd.PersistentFlags().String("conf", storage.DefaultDirectory(), "directory to store configuration in")
	rootCmd.PersistentFlags().String("dir", storage.DefaultDirectory(), "directory to store config and data in")
	rootCmd.PersistentFlags().Bool("reset-config", false, "reset to the default configuration, overwriting the current one")
	rootCmd.Flags().StringP("address", "a", "", "interface to which the server will bind")
	rootCmd.Flags().IntP("port", "p", 80, "port to listen on")
	rootCmd.Flags().Bool("dev", false, "development mode")
	rootCmd.Flags().BoolP("version", "v", false, "show version")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())

	viper.SetDefault("verify_certificates", true)
	viper.SetDefault("https.enabled", true)
	viper.SetDefault("https.port", 443)
	viper.SetDefault("auth.anonymous", true)
	viper.SetDefault("auth.login", true)
	viper.SetDefault("auth.registration", true)
	viper.SetDefault("dcc.enabled", true)
	viper.SetDefault("dcc.autoget.delete", true)
}

func initConfig(configPath string, overwrite bool) error {
	if _, err := os.Stat(configPath); overwrite || os.IsNotExist(err) {
		config, err := assets.Asset("config.default.toml")
		if err != nil {
			return err
		}

		log.Println("Writing default config to", configPath)

		return ioutil.WriteFile(configPath, config, 0600)
	}
	return nil
}
