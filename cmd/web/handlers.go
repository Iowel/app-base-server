package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Iowel/app-base-server/pkg/encryption"
	"github.com/Iowel/app-base-server/pkg/urlsigner"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "login", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	//save userID exp to session
	id, err := app.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userID", id)

	//save userRole exp to session
	user, err := app.GetOneUser(id)
	if err != nil {
		app.errorLog.Println("ошибка получения пользователя:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userRole", user.Role)

	// save token exp to session
	app.Session.Put(r.Context(), "tokenExp", time.Now().Add(100*time.Hour))
	// app.Session.Put(r.Context(), "tokenExp", time.Now().Add(10*time.Second))

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) Books(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "books", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "register", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "forgot-password", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ShowResetPassword(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")

	// validate url
	theURL := r.RequestURI
	testURL := fmt.Sprintf("%s%s", app.config.frontend, theURL)

	signer := urlsigner.Signer{
		Secret: []byte(app.config.secretkey),
	}

	valid := signer.VerifyToken(testURL)

	if !valid {
		app.errorLog.Println("Invalid URL - detected tampering")
		return
	}

	// check expired link
	expired := signer.Expired(testURL, 60)
	if expired {
		app.errorLog.Println("Link expired")
		return
	}

	// crypt mail
	encryptor := encryption.Encryption{
		Key: []byte(app.config.secretkey),
	}

	encryptedEmail, err := encryptor.Encrypt(email)
	if err != nil {
		app.errorLog.Println("Encryption failed")
		return
	}

	data := make(map[string]interface{})
	data["email"] = encryptedEmail

	if err := app.renderTemplate(w, r, "reset-password", &templateData{Data: data}); err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "profile", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Form(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "form", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) AllUsersForAdmin(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "all-users", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "users", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) AllUsersposts(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "all-users-posts", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "one-user", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) OneProfile(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "one-profile", &templateData{}); err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) EditProfile(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "editprofile", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Kino(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "kino", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "verifyemail", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Shop(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "shop", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Orders(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "orders", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}
