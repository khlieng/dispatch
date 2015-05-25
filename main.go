package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/khlieng/name_pending/assets"
	"github.com/khlieng/name_pending/commands"
	"github.com/khlieng/name_pending/storage"
)

func main() {
	initConfig()
	commands.Execute()
}

func initConfig() {
	configPath := path.Join(storage.AppDir, "config.toml")

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
