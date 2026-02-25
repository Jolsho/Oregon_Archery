package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internals"
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

// TODO --> NEED TO IMPLEMENT LOGGING...
// like for bad actors... 
// make sure we can find their IPs later on.
// just output to a log file...

func main() {
	state := internals.New_State();
	net := internals.New_Networker();


	//////////////////////////////////////////////
	// 				PATHS 	  

	dst, ok := os.LookupEnv("dst_dir");
	if !ok { dst = "../ui/"; }

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir(dst)));

	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		if cookie := internals.Secure_Middleware(net, w, r); cookie != nil {
			internals.Handle_events(net, state, cookie, w, r);
		}
	})

	mux.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
		if cookie := internals.Secure_Middleware(net, w, r); cookie != nil {
			internals.Handle_WS(net, state, cookie, w, r);
		}
	})

	//////////////////////////////////////////////


	server := &http.Server{
		Addr: "0.0.0.0:8080",
		Handler: mux,
		TLSConfig: internals.GenerateTLSConfig(),
	};

	handleShutdown(func() {
		net.Shutdown();
	})

	log.Println("Started server!");
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err);
	}
}
