package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/network"
	state_p "server/state"
	"slices"
)

type EventRes struct {
	Events []state_p.Event `json:"events"`
	Token string `json:"token"`
}

func Handle_events(net *network.Networker, state *state_p.State, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query();

	token := q.Get("token");
	if token == "" {
		http.Error(w, "NO TOKEN", http.StatusBadRequest);
		return;
	}

	nonce := network.Get_nonce(token);
	ip := network.Get_client_ip(r);

	switch r.Method {
	case "GET": {
		state.EventMux.RLock();
		defer state.EventMux.RUnlock();

		res := EventRes{
			Events: state.Events,
			Token: token,
		};
		res.Events = state.Events;
		for i := range res.Events {
			ev := &res.Events[i];
			ev.IsOwn = ev.Secret == nonce;
		}

		bytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
			return;
		}

		w.Write(bytes);
		return;
	}

	case "POST", "PUT": {

		event := state_p.DefaultEvent();
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

				log := fmt.Sprintf("ATTEMPTED ACCESS OF UNOWNED EVENT from %s", ip);
				net.Logger.Log(network.INFO_LEVEL, log)


				http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
				return
			}
			event.Secret = nonce;
			*existing = event;

			err := existing.Sanitize();
			if err != nil {

				log := fmt.Sprintf("EVENT from %s FAILED SANITIZATION", ip);
				net.Logger.Log(network.INFO_LEVEL, log)


				http.Error(w, err.Error(), http.StatusBadRequest);
				return;
			}

			w.Header().Set("Content-Type", "application/json")
			bytes, err := json.Marshal(event)
			if err != nil {
				net.Logger.Log(network.WARNING_LEVEL, "EVENT JSON MARSHAL FAILED...")

				http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
				return;
			}
			w.Write(bytes);

		} else {
			if state.TooManyEvents(ip) {
				log := fmt.Sprintf("CREATED TOO MANY EVENTS %s", ip);
				net.Logger.Log(network.INFO_LEVEL, log)


				http.Error(w, "TOO MANY EVENTS", http.StatusUnauthorized);
				return
			}

			event.Secret = nonce;
			event.IsOwn = true;

			err := event.Sanitize();
			if err != nil {

				log := fmt.Sprintf("EVENT from %s FAILED SANITIZATION", ip);
				net.Logger.Log(network.INFO_LEVEL, log)


				http.Error(w, err.Error(), http.StatusBadRequest);
				return;
			}

			w.Header().Set("Content-Type", "application/json")
			bytes, err := json.Marshal(event)
			if err != nil {
				net.Logger.Log(network.WARNING_LEVEL, "EVENT JSON MARSHAL FAILED...")
				http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
				return;
			}
			w.Write(bytes);

			state.Events = append(state.Events, event);
		}

		newEvent := struct {
			Event state_p.Event `json:"event"`
		}{Event: event};
		payload, err := json.Marshal(&newEvent);
		if err != nil {
			log := fmt.Sprintf("WS MARSHAL NEW EVENT:: %s", err.Error());
			net.Logger.Log(network.WARNING_LEVEL, log);
			http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
			return;
		}

		ws_msg := network.WsMsg{
			Msg: "new_event",
			Payload: payload,
		};
		for _, conn := range net.Conns {
			if conn.Nonce != nonce {
				conn.Out <- ws_msg;
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

				log := fmt.Sprintf("ATTEMPTED ACCESS OF UNOWNED EVENT from %s", ip);
				net.Logger.Log(network.INFO_LEVEL, log)

				http.Error(w, "INVALID SECRET", http.StatusUnauthorized);
				return;
			}

			state.EventPerUser[ip]--;

			state.Events = slices.DeleteFunc(state.Events, func(ev state_p.Event) bool {
				return ev.Secret == nonce && ev.Title == title;
			}); 

		} else {
			http.Error(w, "EVENT NOT EXISTS", http.StatusBadRequest);
			return;
		}

		payload_bytes, err := json.Marshal(
			struct {
				Title string `json:"title"`
			}{
				Title: title,
			},
		);
		if err != nil {
			log := fmt.Sprintf("WS MARSHAL DELETE EVENT :: %s", err.Error());
			net.Logger.Log(network.WARNING_LEVEL, log);
			http.Error(w, "EVENT MARSHAL FAILED", http.StatusInternalServerError);
			return;
		}

		ws_msg := network.WsMsg{
			Msg: "delete_event",
			Payload: payload_bytes,
		};
		for conn_nonce, conn := range net.Conns {
			if conn_nonce != nonce {
				conn.Out <- ws_msg;
			}
		}
	}

	default:
		log := fmt.Sprintf("INVALID METHOD from %s", ip);
		net.Logger.Log(network.INFO_LEVEL, log)

		http.Error(w, "INVALID METHOD", http.StatusMethodNotAllowed);
		return;
	}
}

