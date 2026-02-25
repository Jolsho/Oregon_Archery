package internals

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)


func Handle_events(net *Networker, state *State, cookie *http.Cookie, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query();

	nonce := get_nonce(cookie);

	switch r.Method {
	case "GET": {
		state.EventMux.RLock();
		defer state.EventMux.RUnlock();

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

		state.EventMux.Lock();
		defer state.EventMux.Unlock();

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
		}

		type NewEventMsg struct {
			Msg string `json:"msg"`
			Event Event `json:"event"`
		}
		ws_msg := NewEventMsg{
			Msg: "new_event",
			Event: event,
		};
		for conn_nonce, conn := range net.Conns {
			if conn_nonce != nonce {
				err := conn.WriteJSON(ws_msg)
				if err != nil {
					log.Println("PUT_EVENT::WS::WRITEJSON -> ", err.Error())
				}
			}
		}
	} 

	case "DELETE": {
		title := q.Get("title");

		state.EventMux.Lock();
		defer state.EventMux.Unlock();

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

		type DeleteEventMsg struct {
			Msg string `json:"msg"`
			Title string `json:"title"`
		}
		ws_msg := DeleteEventMsg{
			Msg: "delete_event",
			Title: title,
		};
		for conn_nonce, conn := range net.Conns {
			if conn_nonce != nonce {
				err := conn.WriteJSON(ws_msg)
				if err != nil {
					log.Println("DELETE_EVENT::WS::WRITEJSON -> ", err.Error())
				}
			}
		}
	}

	default:
		http.Error(w, "INVALID METHOD", http.StatusMethodNotAllowed);
		return;
	}
}

