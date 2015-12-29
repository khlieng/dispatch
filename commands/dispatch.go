package commands

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"

	"github.com/khlieng/dispatch/assets"
	"github.com/khlieng/dispatch/server"
	"github.com/khlieng/dispatch/storage"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dispatch",
		Short: "Web-based IRC client in Go.",
		Run: func(cmd *cobra.Command, args []string) {
			storage.Initialize(appDir)
			server.Run(viper.GetInt("port"))
		},

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			appDir = viper.GetString("dir")

			os.Mkdir(appDir, 0777)
			os.Mkdir(path.Join(appDir, "logs"), 0777)

			initConfig()

			viper.SetConfigName("config")
			viper.AddConfigPath(appDir)
			viper.ReadInConfig()
		},
	}

	appDir string
)

func init() {
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.Flags().IntP("port", "p", 1337, "port to listen on")
	rootCmd.PersistentFlags().String("dir", defaultDir(), "directory to store config and data in")

	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))
}

func Execute() {
	rootCmd.Execute()
}

func initConfig() {
	configPath := path.Join(appDir, "config.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config, err := assets.Asset("config.default.toml")
		if err != nil {
			log.Println(err)
			return
		}

		err = ioutil.WriteFile(configPath, config, 0600)
		if err != nil {
			log.Println(err)
		}
	}
}

func defaultDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(dir, ".dispatch")
}
