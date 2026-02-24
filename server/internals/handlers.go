package internals

import (
	"encoding/json"
	"net/http"
	"slices"
)


func Handle_events(state *State, cookie *http.Cookie, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query();

	nonce := get_nonce(cookie);

	switch r.Method {
	case "GET": {
		state.Mux.RLock();
		defer state.Mux.RUnlock();

		event_copy := state.Events;
		for i := range event_copy {
			ev := &event_copy[i];
			ev.IsOwn = ev.Secret == nonce;
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

		event := DefaultEvent();
		event.Kind = "OUTDOOR";
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		state.Mux.Lock();
		defer state.Mux.Unlock();

		existing := state.Event_exists(event.Title);
		if existing != nil {
			if existing.Secret != nonce {
				http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
				return
			}
			event.Secret = nonce;
			*existing = event;

			err := existing.Sanitize();
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest);
				return;
			}

			w.Header().Set("Content-Type", "application/json")
			bytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
				return;
			}
			w.Write(bytes);
			return;

		} else {

			event.Secret = nonce;
			event.IsOwn = true;

			err := event.Sanitize();
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest);
				return;
			}

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
			if existing.Secret != nonce {
				http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
				return;
			}
			state.Events = slices.DeleteFunc(state.Events, func(ev Event) bool {
				return ev.Secret == nonce && ev.Title == title;
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

