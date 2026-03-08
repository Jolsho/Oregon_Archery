package network

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func ws_error_handler(
	w http.ResponseWriter, r *http.Request, 
	status int, reason error,
) {
	// TASK_7
	http.Error(w, reason.Error(), status)
};

const READ_BUFFER_SIZE = 1024;
const WRITE_BUFFER_SIZE = 1024;
func New_Upgrader(allowed_origins map[string]struct{}) *websocket.Upgrader {
	check_origin := func (r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Fallback for some older browsers / edge cases
			ref := r.Header.Get("Referer")
			if ref == "" {
				return false
			}

			u, err := url.Parse(ref)
			if err != nil { return false }

			origin = u.Scheme + "://" + u.Host
		}

		_, ok := allowed_origins[origin]
		return ok
	};

	return &websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Duration(4).Seconds()),
		ReadBufferSize: READ_BUFFER_SIZE,
		WriteBufferSize: WRITE_BUFFER_SIZE,
		CheckOrigin: check_origin,
		Error: ws_error_handler,
	};
}
