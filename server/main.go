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
)

var debug = flag.Bool("debug", false, "enable debug logging")

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

	const LOG_FILE = "/var/ohsal/logs/ohsal.log"
	const KEY_PATH = "/var/ohsal/KEYS.txt"
	net := network.New_Networker(KEY_PATH, LOG_FILE);


	//////////////////////////////////////////////
	// 				PATHS 	  

	dst, ok := os.LookupEnv("dst_dir");
	if !ok { dst = "../ui/"; }

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir(dst)));

	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		handlers.Handle_events(net, state, w, r);
	});
	mux.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
		handlers.Handle_WS(net, state, w, r);
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie := network.Secure_Middleware(net, w, r); cookie == nil { return }
		mux.ServeHTTP(w, r)
	});

	//////////////////////////////////////////////


	server := &http.Server{
		Addr: "0.0.0.0:8080",
		Handler: handler,
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
