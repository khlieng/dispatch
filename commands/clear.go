package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/khlieng/dispatch/storage"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all user data",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(storage.Path.Database())
		if err == nil || os.IsNotExist(err) {
			log.Println("Database cleared")
		} else {
			log.Println(err)
		}

		err = os.RemoveAll(storage.Path.HMACKey())
		if err == nil || os.IsNotExist(err) {
			log.Println("HMAC key cleared")
		} else {
			log.Println(err)
		}

		err = os.RemoveAll(storage.Path.Users())
		if err == nil {
			log.Println("User data cleared")
		} else {
			log.Println(err)
		}
	},
}
