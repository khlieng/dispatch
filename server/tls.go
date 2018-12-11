package server

import (
	"os"
)

func (d *Dispatch) certExists() bool {
	cfg := d.Config().HTTPS

	if cfg.Cert == "" || cfg.Key == "" {
		return false
	}

	if _, err := os.Stat(cfg.Cert); err != nil {
		return false
	}
	if _, err := os.Stat(cfg.Key); err != nil {
		return false
	}

	return true
}
