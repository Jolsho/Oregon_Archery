package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/handlers"
	"server/network"
	"server/state"
	"syscall"
)

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

	DST := "/app/ui"; 
	LOG_FILE := "/var/ohsal/logs/ohsal.log";
	KEY_PATH := "/var/ohsal/KEYS.txt";

	LISTEN_IP := "0.0.0.0";
	LISTEN_PORT := 80;

	allowedOrigins := map[string]struct{}{
		"https://testohsal.com":     	{},
		"https://www.testohsal.com":    {},
	}

	test := flag.Bool("test", false, "run in test mode")
	flag.Parse()

	if *test {
		LOG_FILE = "../mnt/ohsal/logs/ohsal.log";
		KEY_PATH = "../mnt/ohsal/KEYS.txt";

		LISTEN_IP = "127.0.0.1";
		LISTEN_PORT = 8080;
		DST = "../ui/dist";

		allowedOrigins["http://localhost:8080"] = struct{}{};
		allowedOrigins["tauri://localhost"] = struct{}{};
		allowedOrigins["http://localhost:5174"] = struct{}{};
	}

	state := state.New_State();
	net := network.New_Networker(KEY_PATH, LOG_FILE, allowedOrigins);

	//////////////////////////////////////////////
	// 				PATHS 	  

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir(DST)));

	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		handlers.Handle_events(net, state, w, r);
	});
	mux.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
		handlers.Handle_WS(net, state, w, r);
	});
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie := network.Secure_Middleware(net, w, r); cookie == "" {
			return
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return ;
		}

		mux.ServeHTTP(w, r)
	})

	//////////////////////////////////////////////


	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", LISTEN_IP, LISTEN_PORT),
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
