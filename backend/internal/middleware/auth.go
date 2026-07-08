package middleware

import (
	"courses/internal/auth"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		} else {
			ctx, err := auth.VerifyAccessToken(r.Context(), authHeader[1])

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					w.WriteHeader(http.StatusRequestTimeout)
					w.Write([]byte("Auth Token Expired"))
					return
				}

				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized Access"))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
