package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir("../")));

	log.Println("Started server!");
	if err := http.ListenAndServe("127.0.0.1:8080", mux); err != nil {
		log.Fatal(err);
	}
}
