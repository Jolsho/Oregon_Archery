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


func main() {
	state := internals.New_State();

	dst, ok := os.LookupEnv("dst_dir");
	if !ok {
		dst = "../ui/";
	}

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir(dst)));

	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		if cookie := internals.Secure_Middleware(state, w, r); cookie != nil {
			internals.Handle_events(state, cookie, w, r);
		}
	})

	mux.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
		if cookie := internals.Secure_Middleware(state, w, r); cookie != nil {
			internals.Handle_WS(state, w, r);
		}
	})

	server := &http.Server{
		Addr: "0.0.0.0:8080",
		Handler: mux,
		TLSConfig: internals.GenerateTLSConfig(),
	};

	handleShutdown(func() {
		state.Persist();
	})

	log.Println("Started server!");
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err);
	}
}
