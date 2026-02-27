package network

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"server/utils"
	"strings"
	"time"
)

func generate_cookie_hash(w http.ResponseWriter, nonce, expire string) ([]byte, error) {
	hasher := sha256.New();

	n, err := hasher.Write([]byte(nonce));
	if n != len([]byte(nonce)) || err != nil {
		http.Error(w, "HASHING SECRET", http.StatusInternalServerError);
		if n != len([]byte(nonce)) {
			err =  errors.New("DIDNT WRITE FULL NONCE TO HASH");
		}
		return nil, err;
	}

	n, err = hasher.Write([]byte(expire));
	if n != len([]byte(expire)) || err != nil {
		http.Error(w, "HASHING EXPIRE", http.StatusInternalServerError);
		if n != len([]byte(expire)) {
			err =  errors.New("DIDNT WRITE FULL EXPIRE TO HASH");
		}
		return nil, err;
	}

	return hasher.Sum(nil), nil;
}

func generate_signed_cookie_value(
	net *Networker, w http.ResponseWriter, 
) (string, time.Time, error) {

	expire := time.Now().Add(24 * time.Hour);

	nonce, err := utils.Random_8_str();
	if err != nil {
		http.Error(w, "SECRET GENERATION", http.StatusInternalServerError);
		return "", expire, err;
	};

	expire_bytes, err :=  expire.MarshalBinary();
	if err != nil {
		http.Error(w, "MARSHAL EXPIRE", http.StatusInternalServerError);
		return "", expire, err;
	};
	expire_str := base64.StdEncoding.EncodeToString(expire_bytes);


	hash, err := generate_cookie_hash(w, nonce, expire_str);
	if err != nil { return "", expire, err; }


	sig_bytes, err := ecdsa.SignASN1(rand.Reader, net.PrivKey, hash);
	if err != nil {
		http.Error(w, "SIGNATURE ERROR", http.StatusInternalServerError);
		return "", expire, err;
	}

	cookie_val := base64.StdEncoding.EncodeToString(sig_bytes);
	cookie_val = cookie_val + ":" + nonce + ":" + expire_str;

	return cookie_val, expire, nil
}

func Get_nonce(cookie *http.Cookie) string {
	values := strings.Split(cookie.Value, ":")
	return values[1];
}

func verify_cookie(net *Networker, cookie *http.Cookie, w http.ResponseWriter) error {

	values := strings.Split(cookie.Value, ":")
	if len(values) != 3 {
		http.Error(w, "MALFORMED COOKIE", http.StatusUnauthorized)
		return errors.New("MALFORMED COOKIE");
	}
	sig := values[0];
	nonce := values[1];
	expire_str := values[2];

	expire_bytes, err := base64.StdEncoding.DecodeString(expire_str);
	if err != nil {
		http.Error(w, "DECODE COOKIE EXPIRE", http.StatusInternalServerError);
		return err;
	}
	var expire time.Time
	err = expire.UnmarshalBinary(expire_bytes)
	if err != nil {
		http.Error(w, "UNMARSHAL COOKIE EXPIRE", http.StatusInternalServerError)
		return err;
	}

	hash, err := generate_cookie_hash(w, nonce, expire_str);
	if err != nil { return err; }

	signature, err := base64.StdEncoding.DecodeString(sig);
	if err != nil {
		http.Error(w, "DECODE SIGNATURE", http.StatusUnauthorized)
		return err;
	}
	if !ecdsa.VerifyASN1(net.PubKey, hash, signature) {
		http.Error(w, "INVALID COOKIE", http.StatusUnauthorized);
		return err;
	}


	if expire.Before(time.Now()) {
		value, expire, err := generate_signed_cookie_value(net, w);
		if err != nil { return err; }

		cookie.Value = value;
		cookie.Expires = expire;

		http.SetCookie(w, cookie);
	}

	return nil;
}
