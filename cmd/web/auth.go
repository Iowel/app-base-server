package main

import (
	"log"
	"net/http"
	"time"
)

func (app *application) AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		expTime := app.Session.GetTime(r.Context(), "tokenExp")
		log.Println("TokenExp from session:", expTime)
		
		if expTime.IsZero() || expTime.Before(time.Now()) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
