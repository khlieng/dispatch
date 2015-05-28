package commands

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Edit config file",
		Run: func(cmd *cobra.Command, args []string) {
			if editor := findEditor(); editor != "" {
				process := exec.Command(editor, path.Join(appDir, "config.toml"))
				process.Stdin = os.Stdin
				process.Stdout = os.Stdout
				process.Stderr = os.Stderr
				process.Run()
			} else {
				log.Println("Unable to locate editor")
			}
		},
	}

	editors = []string{"nano", "notepad", "vi", "emacs"}
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
