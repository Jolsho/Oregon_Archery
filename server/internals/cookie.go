package internals

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

func generate_cookie_hash(w http.ResponseWriter, nonce, expire string) ([]byte, bool) {
	hasher := sha256.New();

	n, err := hasher.Write([]byte(nonce));
	if n != len([]byte(nonce)) || err != nil {
		http.Error(w, "HASHING SECRET", http.StatusInternalServerError);
		return nil, false;
	}

	n, err = hasher.Write([]byte(expire));
	if n != len([]byte(expire)) || err != nil {
		http.Error(w, "HASHING EXPIRE", http.StatusInternalServerError);
		return nil, false;
	}

	return hasher.Sum(nil), true;
}

func generate_signed_cookie_value(
	net *Networker, w http.ResponseWriter, 
) (string, time.Time, bool) {

	expire := time.Now().Add(24 * time.Hour);

	nonce, err := Random_8_str();
	if err != nil {
		http.Error(w, "SECRET GENERATION", http.StatusInternalServerError);
		return "", expire, false;
	};

	expire_bytes, err :=  expire.MarshalBinary();
	if err != nil {
		http.Error(w, "MARSHAL EXPIRE", http.StatusInternalServerError);
		return "", expire, false;
	};
	expire_str := base64.StdEncoding.EncodeToString(expire_bytes);


	hash, ok := generate_cookie_hash(w, nonce, expire_str);
	if !ok { return "", expire, false; }


	sig_bytes, err := ecdsa.SignASN1(rand.Reader, net.PrivKey, hash);
	if err != nil {
		http.Error(w, "SIGNATURE ERROR", http.StatusInternalServerError);
		return "", expire, false;
	}

	cookie_val := base64.StdEncoding.EncodeToString(sig_bytes);
	cookie_val = cookie_val + ":" + nonce + ":" + expire_str;

	return cookie_val, expire, true;
}

func get_nonce(cookie *http.Cookie) string {
	values := strings.Split(cookie.Value, ":")
	return values[1];
}

func verify_cookie(net *Networker, cookie *http.Cookie, w http.ResponseWriter) bool {

	values := strings.Split(cookie.Value, ":")
	if len(values) != 3 {
		http.Error(w, "MALFORMED COOKIE", http.StatusUnauthorized)
		return false
	}
	sig := values[0];
	nonce := values[1];
	expire_str := values[2];

	expire_bytes, err := base64.StdEncoding.DecodeString(expire_str);
	if err != nil {
		http.Error(w, "DECODE COOKIE EXPIRE", http.StatusInternalServerError);
		return false;
	}
	var expire time.Time
	err = expire.UnmarshalBinary(expire_bytes)
	if err != nil {
		http.Error(w, "UNMARSHAL COOKIE EXPIRE", http.StatusInternalServerError)
		return false
	}

	hash, ok := generate_cookie_hash(w, nonce, expire_str);
	if !ok { return false; }

	signature, err := base64.StdEncoding.DecodeString(sig);
	if err != nil {
		http.Error(w, "DECODE SIGNATURE", http.StatusUnauthorized)
		return false
	}
	if !ecdsa.VerifyASN1(net.PubKey, hash, signature) {
		http.Error(w, "INVALID COOKIE", http.StatusUnauthorized);
		return false;
	}


	if expire.Before(time.Now()) {
		value, expire, ok := generate_signed_cookie_value(net, w);
		if !ok { return false; }

		cookie.Value = value;
		cookie.Expires = expire;

		http.SetCookie(w, cookie);
	}

	return true;
}
