package internals

import (
	"net/http"
	"strings"
)


const SESSION_COOKIE = "secret";

func Secure_Middleware(net *Networker, w http.ResponseWriter, r *http.Request) *http.Cookie {

	ip := strings.Split(r.RemoteAddr, ":")[0];

	if expires, timedout := net.RateLimiter.is_timed_out(ip); timedout { 
		http.Error(w, "TIMEOUT UNTIL " + expires, http.StatusGatewayTimeout);
		return nil; 
	}

	cookie, err := r.Cookie(SESSION_COOKIE);
	if err != nil {
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusBadRequest)
			return nil;
		}

		timeout, err, status := net.RateLimiter.handle_rates(RATE_SESSION_CREATION, ip);
		if err != nil { 
			net.RateLimiter.handle_timeout(ip, timeout)
			http.Error(w, err.Error(), status);
			return nil;
		};

		value, expire, ok := generate_signed_cookie_value(net, w);
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

	if !verify_cookie(net, cookie, w) { 
		score := bad_behaviour_score(INFRACTION_INVALID_COOKIE);
		timeout, err, status := net.RateLimiter.handle_rates(RATE_BEHAVIOUR, nonce, score);
		if err != nil { 
			net.RateLimiter.handle_timeout(ip, timeout)
			http.Error(w, err.Error(), status);
			return nil; 
		};
		return nil;
	}

	timeout, err, status := net.RateLimiter.handle_rates(RATE_DATA, nonce);
	if err != nil { 
		net.RateLimiter.handle_timeout(ip, timeout)
		http.Error(w, err.Error(), status);
		return nil; 
	};

	return cookie;
}
