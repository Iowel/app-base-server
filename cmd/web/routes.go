package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/ws", app.WsEndPoint)

	mux.HandleFunc("/ws", app.WsEndPoint)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)

		mux.Get("/books", app.Books)
	})

	mux.Get("/home", app.Home)
	mux.Get("/profile", app.Profile)
	mux.Get("/form", app.Form)
	mux.Get("/users", app.AllUsers)
	mux.Get("/all-users-posts", app.AllUsersposts)
	mux.Get("/edit-profile", app.EditProfile)

	mux.Get("/kino", app.Kino)
	mux.Get("/shop", app.Shop)
	mux.Get("/orders", app.Orders)

	// auth
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Get("/forgot-password", app.ForgotPassword)
	mux.Get("/reset-password", app.ShowResetPassword)

	mux.Get("/verify_email", app.VerifyEmail)

	mux.Get("/all-users", app.AllUsersForAdmin)
	mux.Get("/all-users/{id}", app.OneUser)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
