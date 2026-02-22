package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/data"
)

func handle_events(state *data.State, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query();
	title := q.Get("title");

	if title != "" {
		switch r.Method {
		case "POST", "PUT": {

			var event data.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}


			state.Mux.Lock();
			defer state.Mux.Unlock();
			if e, exists := state.Events[title]; exists {
				if e.Secret != event.Secret {
					http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
					return
				}
			}

			if event.Title != title {
				delete(state.Events, title)
			}
			state.Events[event.Title] = event;

		} 
		case "DELETE": {
			secret := q.Get("secret");

			state.Mux.Lock();
			defer state.Mux.Unlock();
			if event, ok := state.Events[title]; ok {
				if event.Secret != secret {
					http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
					return;
				}
				delete(state.Events, title)
			}
		}
		default:
			http.Error(w, "INVALID METHOD", http.StatusMethodNotAllowed);
			return;
		}
	} else {
		state.Mux.RLock();
		defer state.Mux.RUnlock();

		bytes, err := json.Marshal(state.Events)
		if err != nil {
			http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
			return;
		}
		w.Write(bytes);
	}
}


func main() {
	state := data.State{
		Events: map[string]data.Event{},
	};

	mux := http.NewServeMux();
	mux.Handle("/", http.FileServer(http.Dir("../")));

	mux.HandleFunc("/events", func (w http.ResponseWriter, r *http.Request) {
		handle_events(&state, w, r);
	})

	mux.HandleFunc("/divisions", func (w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(data.DIVISIONS)
		if err != nil {
			http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
			return;
		}
		w.Write(bytes);
	});


	log.Println("Started server!");
	if err := http.ListenAndServe("127.0.0.1:8080", mux); err != nil {
		log.Fatal(err);
	}
}
