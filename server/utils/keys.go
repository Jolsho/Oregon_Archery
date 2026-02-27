package utils
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"encoding/pem"
	"fmt"
	"os"
)


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
