package internals

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const READ_BUFFER_SIZE = 1024;
const WRITE_BUFFER_SIZE = 1024;

func ws_error_handler(
	w http.ResponseWriter, r *http.Request, 
	status int, reason error,
) {
	// TODO -- record errors? to restrict bad behaviour...
	http.Error(w, reason.Error(), status)
};

var allowedOrigins = map[string]struct{}{
	"http://localhost:8080":     {},
}
func check_origin(r *http.Request) bool {
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

    _, ok := allowedOrigins[origin]
    return ok
};

func New_Upgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Duration(4).Seconds()),
		ReadBufferSize: READ_BUFFER_SIZE,
		WriteBufferSize: WRITE_BUFFER_SIZE,
		CheckOrigin: check_origin,
		Error: ws_error_handler,
	};
}

func Handle_WS(state *State, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SESSION_COOKIE);
	if err != nil { 
		http.Error(w, "NO COOKIE", http.StatusUnauthorized);
		return 
	};

	conn, err := state.Upgrader.Upgrade(w, r, http.Header{});
	if err != nil { return };

	state.Conns[get_nonce(cookie)] = conn;

	go func() {
		// TODO -- define WS communication branches
	}();
}

