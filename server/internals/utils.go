package internals

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"
)


func Random_8_str() (string, error) {
	rnd_secret := make([]byte, 8);
	n, err := rand.Read(rnd_secret);
	if n < 8 || err != nil { return "", err }
	return base64.RawURLEncoding.EncodeToString(rnd_secret), nil
}

func Rand_8_str_ignored() string {
	sec, _ := Random_8_str();
	return sec;
}

type Task interface {
	get_when() time.Time
}

func Task_Scheduler[T Task](
	tasks *[]T,
	mux *sync.Mutex, 
	default_interval time.Duration, 
	notifier chan struct{},
	done <- chan struct{},
) {
	var timer *time.Timer

	for {
		mux.Lock()
		if len(*tasks) == 0 {
			mux.Unlock()
			select {
			case <-time.After(default_interval):
				continue
			case <-done:
				return
			}
		}

		when := (*tasks)[0].get_when()
		mux.Unlock()

		wait := max(time.Until(when), 0)

		if timer == nil {
			timer = time.NewTimer(wait)
		} else {
			timer.Reset(wait)
		}

		select {
		case <-timer.C:
			select {
			case notifier <- struct{}{}:
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}

func SavePrivateKey(filename string, key *ecdsa.PrivateKey) error {
    // Convert to DER
    der, err := x509.MarshalECPrivateKey(key)
    if err != nil { return err }

    // PEM block
    block := &pem.Block{
        Type:  "EC PRIVATE KEY",
        Bytes: der,
    }

    // Write to file
    f, err := os.Create(filename)
    if err != nil { return err }
    defer f.Close()

    return pem.Encode(f, block)
}

func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
		if os.IsNotExist(err) {
			priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				panic(err)
			}
			return priv, nil;
		}
        return nil, err
    }

    block, _ := pem.Decode(data)
    if block == nil || block.Type != "EC PRIVATE KEY" {
        return nil, fmt.Errorf("invalid PEM block")
    }

    return x509.ParseECPrivateKey(block.Bytes)
}

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
