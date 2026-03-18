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

func Secure_Middleware(
	net *Networker, 
	w http.ResponseWriter, 
	r *http.Request,
) string {

	origin := r.Header.Get("Origin")
	if origin == "http://localhost:5174" || origin == "tauri://localhost" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ip := Get_client_ip(r);
	q := r.URL.Query();
	token := q.Get("token");

	var err error;
	if token == "" {
		token, _, err = generate_signed_cookie_value(net, w);
		if err != nil { 
			log := fmt.Sprintf("GEN SIGNED COOKIE :: %s", err.Error());
			net.Logger.Log(ERROR_LEVEL, log);
			return "";
		}
		
	}

	if err = verify_token(net, &token, w); err != nil { 
		log := fmt.Sprintf("INVALID COOKIE from %s :: %s", ip, err.Error());
		net.Logger.Log(INFO_LEVEL, log);

		return "";
	}

	q.Set("token", token);
	r.URL.RawQuery = q.Encode();

	return token;
}
