package commands

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/khlieng/dispatch/storage"
)

var (
	configCmd = &cobra.Command{
		Use:   "config [editor]",
		Short: "Edit config file",
		Run: func(cmd *cobra.Command, args []string) {
			editors = append(args, editors...)

			if editor := findEditor(); editor != "" {
				process := exec.Command(editor, storage.Path.Config())
				process.Stdin = os.Stdin
				process.Stdout = os.Stdout
				process.Stderr = os.Stderr
				process.Run()
			} else {
				log.Println("Unable to locate editor")
			}
		},
	}

	editors = []string{"nano", "code", "vi", "emacs", "notepad"}
)

func findEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		editor, err := exec.LookPath(editor)
		if err == nil {
			return editor
		}
	}

	for _, editor := range editors {
		editor, err := exec.LookPath(editor)
		if err == nil {
			return editor
		}
	}

	return ""
}
