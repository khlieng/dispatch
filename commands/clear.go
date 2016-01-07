package commands

import (
	"log"
	"os"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/cobra"

	"github.com/khlieng/dispatch/storage"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear database and message logs",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(storage.Path.Database())
		if err == nil || os.IsNotExist(err) {
			log.Println("Database cleared")
		} else {
			log.Println(err)
		}

		err = os.RemoveAll(storage.Path.Logs())
		if err == nil {
			log.Println("Logs cleared")
		} else {
			log.Println(err)
		}
	},
}
