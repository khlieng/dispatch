package server

import (
	"crypto/tls"
	"os"

	"github.com/klauspost/cpuid"
)

func getCipherSuites() []uint16 {
	if cpuid.CPU.AesNi() {
		return []uint16{
			tls.TLS_FALLBACK_SCSV,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		}
	}

	return []uint16{
		tls.TLS_FALLBACK_SCSV,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256}
}

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
