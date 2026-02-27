package utils

import (
	"crypto/rand"
	"encoding/base64"
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
