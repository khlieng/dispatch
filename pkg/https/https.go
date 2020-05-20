package https

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/khlieng/dispatch/pkg/netutil"
	"github.com/klauspost/cpuid"
)

type Config struct {
	Addr      string
	PortHTTP  string
	PortHTTPS string
	HTTPOnly  bool

	StoragePath string
	Domain      string
	Email       string

	Cert string
	Key  string
}

func Serve(handler http.Handler, cfg Config) error {
	errCh := make(chan error, 1)

	httpSrv := &http.Server{
		Addr: net.JoinHostPort(cfg.Addr, cfg.PortHTTP),
	}

	if !cfg.HTTPOnly {
		httpSrv.ReadTimeout = 5 * time.Second
		httpSrv.WriteTimeout = 5 * time.Second

		httpsSrv := &http.Server{
			Addr:        net.JoinHostPort(cfg.Addr, cfg.PortHTTPS),
			ReadTimeout: 5 * time.Second,
			IdleTimeout: 120 * time.Second,
			Handler:     handler,
		}

		redirect := HTTPSRedirect(cfg.PortHTTPS, handler)

		if cfg.Cert != "" || cfg.Key != "" {
			httpSrv.Handler = redirect
			httpsSrv.TLSConfig = TLSConfig(nil)

			go func() {
				errCh <- httpSrv.ListenAndServe()
			}()

			go func() {
				errCh <- httpsSrv.ListenAndServeTLS(cfg.Cert, cfg.Key)
			}()
		} else {
			if cfg.StoragePath != "" {
				certmagic.Default.Storage = &certmagic.FileStorage{
					Path: cfg.StoragePath,
				}
			}

			certmagic.Default.MustStaple = true

			magic := certmagic.NewDefault()

			acme := certmagic.NewACMEManager(magic, certmagic.ACMEManager{
				Agreed: true,
				Email:  cfg.Email,
			})

			magic.Issuer = acme

			domains := []string{cfg.Domain}
			if cfg.Domain == "" {
				domains = []string{}
				magic.OnDemand = maxObtain(3)
			}

			err := magic.ManageSync(domains)
			if err != nil {
				return err
			}

			httpSrv.Handler = acme.HTTPChallengeHandler(redirect)
			httpsSrv.TLSConfig = TLSConfig(magic.TLSConfig())

			go func() {
				errCh <- httpSrv.ListenAndServe()
			}()

			go func() {
				errCh <- httpsSrv.ListenAndServeTLS("", "")
			}()
		}
	} else {
		httpSrv.ReadTimeout = 5 * time.Second
		httpSrv.IdleTimeout = 120 * time.Second
		httpSrv.Handler = handler

		return httpSrv.ListenAndServe()
	}

	return <-errCh
}

func HTTPSRedirect(portHTTPS string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}

		if fallback != nil && netutil.IsPrivate(host) {
			fallback.ServeHTTP(w, r)
			return
		}

		u := url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(host, portHTTPS),
			Path:   r.RequestURI,
		}

		w.Header().Set("Connection", "close")
		w.Header().Set("Location", u.String())
		w.WriteHeader(http.StatusMovedPermanently)
	}
}

func TLSConfig(tlsConfig *tls.Config) *tls.Config {
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.CipherSuites = defaultCipherSuites()
	tlsConfig.CurvePreferences = []tls.CurveID{
		tls.X25519,
		tls.CurveP256,
	}
	tlsConfig.PreferServerCipherSuites = true

	return tlsConfig
}

func defaultCipherSuites() []uint16 {
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

func maxObtain(limit int) *certmagic.OnDemandConfig {
	requested := []string{}

	return &certmagic.OnDemandConfig{
		DecisionFunc: func(name string) error {
			for _, n := range requested {
				if name == n {
					return nil
				}
			}

			if len(requested) == limit {
				return fmt.Errorf("OnDemand cert limit reached")
			}

			requested = append(requested, name)
			return nil
		},
	}
}
