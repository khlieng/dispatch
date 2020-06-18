package cryptoutil

import (
	"crypto/tls"
	"strings"
)

func DescribeTLS(version, cipherSuite uint16) string {
	cs := tls.CipherSuiteName(cipherSuite)
	cs = cs[:strings.LastIndexByte(cs, '_')]
	cs = cipherSuiteReplacer.Replace(cs)

	return TLSVersionName(version) + " and " + cs
}

func TLSVersionName(version uint16) string {
	return tlsVersionNames[version]
}

var (
	tlsVersionNames = map[uint16]string{
		tls.VersionTLS10: "TLS 1.0",
		tls.VersionTLS11: "TLS 1.1",
		tls.VersionTLS12: "TLS 1.2",
		tls.VersionTLS13: "TLS 1.3",
	}

	cipherSuiteReplacer = strings.NewReplacer(
		"_WITH_", " with ",
		"TLS_", "",
		"_", "-",
	)
)
