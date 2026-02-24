package internals

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"time"
)

func GenerateTLSConfig() *tls.Config {
    key, _ := rsa.GenerateKey(rand.Reader, 2048)
    tmpl := &x509.Certificate{
        SerialNumber: big.NewInt(1),
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        DNSNames:     []string{"localhost"},
    }
    certDER, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
    cert := tls.Certificate{
        Certificate: [][]byte{certDER},
        PrivateKey:  key,
    }
    return &tls.Config{Certificates: []tls.Certificate{cert}}
}
