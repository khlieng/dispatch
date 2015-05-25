package commands

import (
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/spf13/cobra"

	"github.com/khlieng/name_pending/storage"
)

var (
	editors = []string{"nano", "notepad", "vi", "emacs"}

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Edit config file",
		Run: func(cmd *cobra.Command, _ []string) {
			if editor := findEditor(); editor != "" {
				args := []string{editor, path.Join(storage.AppDir, "config.toml")}
				syscall.Exec(editor, args, os.Environ())
			} else {
				log.Println("Unable to locate editor")
			}
		},
	}
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
