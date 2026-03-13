package auth

import (
	"net/http"
	"strings"
)

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.cfg.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// WebSocket: check token from query param
		if r.Header.Get("Upgrade") == "websocket" {
			token := r.URL.Query().Get("token")
			if token != "" {
				if err := a.ValidateToken(token); err == nil {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// REST: check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if err := a.ValidateToken(token); err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
