package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"server/data"
	"slices"
	"sync"
	"time"
)

func handle_events(state *data.State, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query();
	cookie, err := r.Cookie("secret");
	if err != nil {
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		secret, err := data.Random_secret();
		if err != nil {
			http.Error(w, "RANDOM_GENERATION_ERROR", http.StatusInternalServerError);
			return;
		}

		cookie = &http.Cookie{
			Name:     "secret",
			Value:    secret,
			Path:     "/",                 // cookie is valid for all paths
			HttpOnly: true,                // inaccessible to JS (prevents XSS)
			Secure:   false,               // true in HTTPS
			Expires:  time.Now().Add(24 * time.Hour), // expiration
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie);
	}

	switch r.Method {
	case "GET": {
		state.Mux.RLock();
		defer state.Mux.RUnlock();

		event_copy := state.Events;
		for i := range event_copy {
			ev := &event_copy[i];
			ev.IsOwn = ev.Secret == cookie.Value;
		}

		bytes, err := json.Marshal(event_copy)
		if err != nil {
			http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
			return;
		}
		w.Write(bytes);
		return;
	}
	case "POST", "PUT": {

		var event data.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		state.Mux.Lock();
		defer state.Mux.Unlock();

		existing := state.Event_exists(event.Title);
		if existing != nil {
			if existing.Secret != cookie.Value {
				http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
				return
			}
			event.Secret = cookie.Value;
			*existing = event;

			w.Header().Set("Content-Type", "application/json")
			bytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
				return;
			}
			w.Write(bytes);
			return;

		} else {

			event.Divisions = data.DIVISIONS;
			event.Secret = cookie.Value;
			event.IsOwn = true;

			w.Header().Set("Content-Type", "application/json")
			bytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
				return;
			}
			w.Write(bytes);

			state.Events = append(state.Events, event);
			return;
		}
	} 
	case "DELETE": {
		title := q.Get("title");

		state.Mux.Lock();
		defer state.Mux.Unlock();

		existing := state.Event_exists(title);
		if existing != nil {
			if existing.Secret != cookie.Value {
				http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
				return;
			}
			state.Events = slices.DeleteFunc(state.Events, func(ev data.Event) bool {
				return ev.Secret == cookie.Value && ev.Title == title;
			}); 

		} else {
			http.Error(w, "EVENT NOT EXISTS", http.StatusBadRequest);
			return;
		}
	}
	default:
		http.Error(w, "INVALID METHOD", http.StatusMethodNotAllowed);
		return;
	}
}

func generateTLSConfig() *tls.Config {
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


func main() {
	state := data.State{
		Events: data.Events,
		Mux: &sync.RWMutex{},
	};

	dst, ok := os.LookupEnv("dst_dir");
	if !ok {
		dst = "../";
	}

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir(dst)));

	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		handle_events(&state, w, r);
	})

	server := &http.Server{
		Addr: "0.0.0.0:8080",
		Handler: mux,
		TLSConfig: generateTLSConfig(),
	};

	log.Println("Started server!");
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err);
	}
}
