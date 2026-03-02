package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server/network"
	st "server/state"
	"time"

	"github.com/gorilla/websocket"
)

func handler(state *st.State, net *network.Networker, conn *network.WSConn, msg network.WsMsg) {
	switch msg.Msg {
	case "new_event": {
		
		event := st.Event{};
		err := json.Unmarshal(msg.Payload, &event);
		if err != nil {
			log := fmt.Sprintf("WS INVALID EVENT JSON from %s", conn.Ip);
			net.Logger.Log(network.INFO_LEVEL, log)

			ok := net.Bad_Behaviour(network.INFRACTION_MALFORMED_DATA, conn.Ip); 
			if !ok { return; }
		}

		state.EventMux.Lock();
		defer state.EventMux.Unlock();

		existing := state.Event_exists(event.Title);
		if existing != nil {
			if existing.Secret != conn.Nonce {
				log := fmt.Sprintf("WS ATTEMPTED ACCESS OF UNOWNED EVENT from %s", conn.Ip);
				net.Logger.Log(network.INFO_LEVEL, log)

				ok := net.Bad_Behaviour(network.INFRACTION_ATTEMPTED_ACCESS, conn.Ip); 
				if !ok { return; }
			}

			event.Secret = conn.Nonce;
			*existing = event;

			err := existing.Sanitize();
			if err != nil {
				log := fmt.Sprintf("WS EVENT from %s FAILED SANITIZATION", conn.Ip);
				net.Logger.Log(network.INFO_LEVEL, log)

				ok := net.Bad_Behaviour(network.INFRACTION_MALFORMED_DATA, conn.Ip); 
				if !ok { return; }
			}
		} else {

			event.Secret = conn.Nonce;
			event.IsOwn = true;

			err := event.Sanitize();
			if err != nil {

				log := fmt.Sprintf("EVENT from %s FAILED SANITIZATION", conn.Ip);
				net.Logger.Log(network.INFO_LEVEL, log)

				net.Bad_Behaviour(network.INFRACTION_MALFORMED_DATA, conn.Ip);
				return;
			}

			state.Events = append(state.Events, event);
		}

	}
	default: {
		_ = net.Bad_Behaviour(network.INFRACTION_MALFORMED_DATA, conn.Ip); 
		return;
	}
	}
}


func Handle_WS(
	net *network.Networker, state *st.State, 
	w http.ResponseWriter, r *http.Request,
) {
	cookie, err := r.Cookie(network.SESSION_COOKIE);
	if err != nil { 
		log := fmt.Sprintf("WS NO COOKIE from %s", network.Get_client_ip(r))
		net.Logger.Log(network.INFO_LEVEL, log);
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return 
	};

	conn, err := net.Upgrader.Upgrade(w, r, http.Header{});
	if err != nil { 
		log := fmt.Sprintf("FAILED WS UPGRADE for %s", network.Get_client_ip(r))
		net.Logger.Log(network.WARNING_LEVEL, log);
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return;
	};

	ctx, cancel := context.WithCancel(r.Context())

	bigConn := &network.WSConn{
		Conn: conn,
		Ctx: ctx,
		Cancel: cancel,
		Out: make(chan network.WsMsg, 5),
		Nonce: network.Get_nonce(cookie),
		Ip: network.Get_client_ip(r),
	};

	net.ConnsMux.Lock();
	net.Conns[bigConn.Nonce] = bigConn;
	net.ConnsMux.Unlock();

	log := fmt.Sprintf("NEW WS CONN from %s", bigConn.Ip);
	net.Logger.Log(network.INFO_LEVEL, log)

	go readLoop(state, net, bigConn)
	writeLoop(net, bigConn)
}

const (
	readWait  = 60 * time.Second  // how long we wait for the next message/pong
	pongWait  = 30 * time.Second
	pingEvery = 20 * time.Second  // must be < pongWait
)

func readLoop(state *st.State, net *network.Networker, conn *network.WSConn) {

	conn.Conn.SetReadLimit(1024 * 1024) // 1MB max message size

	// Reset deadline when we get pong from client
	conn.Conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.Conn.SetPongHandler(func(string) error {
		conn.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-conn.Ctx.Done():
			log := fmt.Sprintf("CLOSED WRITE %s", conn.Ip);
			net.Logger.Log(network.INFO_LEVEL, log)
			return;

		default:
			var msg network.WsMsg
			err := conn.Conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					log := fmt.Sprintf("WS CLOSE ERR for %s :: %s", conn.Ip, err.Error());
					net.Logger.Log(network.INFO_LEVEL, log)
					conn.Cancel();
				} else {
					log := fmt.Sprintf("WS READ_JSON for %s :: %s", conn.Ip, err.Error());
					net.Logger.Log(network.INFO_LEVEL, log)
				}
				return
			}

			handler(state, net, conn, msg)
		}
	}
}

func writeLoop(net *network.Networker, conn *network.WSConn) {
	ticker := time.NewTicker(pingEvery)
	defer func() {
		ticker.Stop();
		conn.Conn.WriteControl(
		  websocket.CloseMessage,
		  websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		  time.Now().Add(time.Second),
		)
		conn.Conn.Close();

		net.ConnsMux.Lock()
		delete(net.Conns, conn.Nonce);
		net.ConnsMux.Unlock()

		close(conn.Out);
	}()

	for {
		select {
		case <- conn.Ctx.Done():
			log := fmt.Sprintf("CLOSED WRITE %s", conn.Ip);
			net.Logger.Log(network.INFO_LEVEL, log)
			return
		case msg, ok := <-conn.Out:
			conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				err := conn.Conn.WriteMessage(websocket.CloseMessage, []byte{}); 
				if err != nil {
					log := fmt.Sprintf("WS WRITE_MSG_ERR for %s :: %s", conn.Ip, err.Error());
					net.Logger.Log(network.INFO_LEVEL, log)
				}
				return;
			}
			err := conn.Conn.WriteJSON(msg); 
			if err != nil {
				log := fmt.Sprintf("WS WRITE_JSON_ERR for %s :: %s", conn.Ip, err.Error());
				net.Logger.Log(network.INFO_LEVEL, log)
				return
			}

		case <-ticker.C:
			// heartbeat ping
			err := conn.Conn.WriteMessage(websocket.PingMessage, nil); 
			if err != nil {
				log := fmt.Sprintf("WS WRITE_PING for %s :: %s", conn.Ip, err.Error());
				net.Logger.Log(network.INFO_LEVEL, log)
				return
			}
		}
	}
}
