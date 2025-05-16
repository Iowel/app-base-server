package auth

import (
	"github.com/Iowel/app-base-server/pkg/helpers"
	"net/http"
)

func Auth(next http.Handler, auth AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := auth.AuthenticateToken(r)
		if err != nil {
			helpers.InvalidCredentials(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
