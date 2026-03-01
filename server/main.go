package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/handlers"
	"server/network"
	"server/state"
	"syscall"

	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"math/big"

	"time"
)

var debug = flag.Bool("debug", false, "enable debug logging")

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

func handleShutdown(onShutdown func()) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigCh
        onShutdown()
        os.Exit(0)
    }()
}


func main() {
	state := state.New_State();
	net := network.New_Networker();


	//////////////////////////////////////////////
	// 				PATHS 	  

	dst, ok := os.LookupEnv("dst_dir");
	if !ok { dst = "../ui/"; }

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir(dst)));
	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		handlers.Handle_events(net, state, w, r);
	})
	mux.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
		handlers.Handle_WS(net, state, w, r);
	})


	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie := network.Secure_Middleware(net, w, r); cookie == nil { return }
		mux.ServeHTTP(w, r)
	});
;

	//////////////////////////////////////////////


	server := &http.Server{
		Addr: "0.0.0.0:80",
		Handler: handler,
		TLSConfig: GenerateTLSConfig(),
	};

	handleShutdown(func() {
		net.Shutdown();
		state.Shutdown();
	})

	net.Logger.Log(network.DEBUG_LEVEL, "STARTED HTTP SERVER.")
	log.Println("Started server!");
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err);
	}
}
