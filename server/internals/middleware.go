package internals

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func is_timed_out(w http.ResponseWriter, state *State, ip string) (is_timedout bool) {
	state.TimeoutMux.RLock();
	if timeout, exists := state.TimeOut[ip]; exists {
		if timeout.Before(time.Now()) {
			delete(state.TimeOut, ip);
		} else {
			error_str := fmt.Sprintf("TIMEOUT::%s", 
				timeout.UTC().Format(time.RFC1123),
			);
			http.Error(w, error_str, http.StatusGatewayTimeout);
			return true;
		}
	}
	state.TimeoutMux.Unlock();

	return false;
}

const SESSION_COOKIE = "secret";

func Secure_Middleware(state *State, w http.ResponseWriter, r *http.Request) *http.Cookie {

	ip := strings.Split(r.RemoteAddr, ":")[0];
	if is_timed_out(w, state, ip) { 
		return nil; 
	}

	cookie, err := r.Cookie(SESSION_COOKIE);
	if err != nil {
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusBadRequest)
			return nil;
		}

		timeout, err, status := state.handle_rates(RATE_SESSION_CREATION, ip);
		if err != nil { 
			state.TimeOut[ip] = time.Now().Add(timeout);
			http.Error(w, err.Error(), status);
			return nil;
		};

		value, expire, ok := generate_signed_cookie_value(state, w);
		if !ok { 
			return nil;
		}

		cookie = &http.Cookie{
			Name:     SESSION_COOKIE,
			Value:    value,
			Path:     "/",                 // cookie is valid for all paths
			HttpOnly: true,                // inaccessible to JS (prevents XSS)
			Secure:   false,               // true in HTTPS
			Expires:  expire, 			   // expiration
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie);
	}

	nonce := get_nonce(cookie);

	if !verify_cookie(state, cookie, w) { 
		score := bad_behaviour_score(INFRACTION_INVALID_COOKIE);
		timeout, err, status := state.handle_rates(RATE_BEHAVIOUR, nonce, score);
		if err != nil { 
			state.TimeOut[ip] = time.Now().Add(timeout);
			http.Error(w, err.Error(), status);
			return nil; 
		};
		return nil;
	}

	timeout, err, status := state.handle_rates(RATE_DATA, nonce);
	if err != nil { 
		state.TimeOut[ip] = time.Now().Add(timeout);
		http.Error(w, err.Error(), status);
		return nil; 
	};

	return cookie;
}
