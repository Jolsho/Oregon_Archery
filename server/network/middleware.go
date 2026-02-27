package network

import (
	"fmt"
	"net/http"
	"strings"
)


const SESSION_COOKIE = "secret";

func Get_client_ip(r *http.Request) string {
    // Prefer X-Forwarded-For if present
    clientIP := r.Header.Get("X-Forwarded-For")
    if clientIP != "" {
        // X-Forwarded-For can contain multiple IPs: client, proxy1, proxy2,...
		// first IP is the original client
        clientIP = strings.Split(clientIP, ",")[0] 

        clientIP = strings.TrimSpace(clientIP)
    } else if ip := r.Header.Get("X-Real-IP"); ip != "" {
        clientIP = ip
    } else {
        // fallback to RemoteAddr (will be Nginx proxy IP if you have it)
        clientIP = r.RemoteAddr
    }
	return clientIP;
}

func Secure_Middleware(net *Networker, w http.ResponseWriter, r *http.Request) *http.Cookie {

	ip := Get_client_ip(r);

	if expires, status := net.RateLimiter.Is_timed_out(ip); status != http.StatusOK { 
		http.Error(w, "TIMEOUT UNTIL " + expires, http.StatusTooManyRequests);
		return nil; 
	}

	cookie, err := r.Cookie(SESSION_COOKIE);
	if err != nil {
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusBadRequest)
			return nil;
		}

		value, expire, err := generate_signed_cookie_value(net, w);
		if err != nil { 
			log := fmt.Sprintf("GEN SIGNED COOKIE :: %s", err.Error());
			net.Logger.Log(ERROR_LEVEL, log);
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

	if err = verify_cookie(net, cookie, w); err != nil { 
		log := fmt.Sprintf("INVALID COOKIE from %s :: %s", ip, err.Error());
		net.Logger.Log(INFO_LEVEL, log);

		timeout, err, status := net.RateLimiter.Handle_behaviour(ip, INFRACTION_INVALID_COOKIE);
		if err != nil { 

			log := fmt.Sprintf("TIMEDOUT %s because %s", ip, err.Error());
			net.Logger.Log(INFO_LEVEL, log);

			net.RateLimiter.Handle_timeout(ip, timeout, status)
			http.Error(w, err.Error(), status);
			return nil; 
		};
		return nil;
	}

	timeout, err, status := net.RateLimiter.Handle_rates(ip);
	if err != nil { 
		log := fmt.Sprintf("TIMEDOUT %s because %s", ip, err.Error());
		net.Logger.Log(INFO_LEVEL, log);

		net.RateLimiter.Handle_timeout(ip, timeout, status)
		http.Error(w, err.Error(), status);
		return nil; 
	};

	return cookie;
}
