package server

import (
	"os"

	"github.com/spf13/viper"
)

func certExists() bool {
	cert := viper.GetString("https.cert")
	key := viper.GetString("https.key")

	if cert == "" || key == "" {
		return false
	}

	if _, err := os.Stat(cert); err != nil {
		return false
	}
	if _, err := os.Stat(key); err != nil {
		return false
	}

	return true
}
