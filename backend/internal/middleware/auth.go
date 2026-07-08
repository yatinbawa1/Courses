package middleware

import (
	"courses/internal/auth"
	"log"
	"net/http"
	"strings"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		} else {
			ctx, err := auth.VerifyToken(r.Context(), authHeader[1])
			log.Printf("Working till here")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized Access"))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
