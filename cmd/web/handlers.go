package main

import (
	"fmt"
	"net/http"
)

func (application *Config) handlerGetMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlerGetMainPage()")
	application.renderTemplate(w, r, "home.page.gohtml", nil)
}

func (application *Config) handlerGetLoginPage(w http.ResponseWriter, r *http.Request) {
	application.renderTemplate(w, r, "login.page.gohtml", nil)
}

func (application *Config) handlerPostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = application.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		application.ErrorLog.Println("cannot parse the form on login page :", err)

		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := application.Models.User.GetByEmail(email)
	if err != nil {
		application.Session.Put(r.Context(), "error", "invalid credentials")
		application.ErrorLog.Println("invalid credentials")

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		application.Session.Put(r.Context(), "error", "invalid credentials")
		application.ErrorLog.Println("invalid credentials")

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}
	if !validPassword {
		application.Session.Put(r.Context(), "error", "invalid credentials")
		application.ErrorLog.Println("invalid credentials")

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	application.Session.Put(r.Context(), "userID", user.ID)
	application.Session.Put(r.Context(), "user", user)
	application.Session.Put(r.Context(), "message", "Successful login!")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (application *Config) handlerGetLogoutPage(w http.ResponseWriter, r *http.Request) {
	_ = application.Session.Destroy(r.Context())
	_ = application.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (application *Config) handlerGetRegisterPage(w http.ResponseWriter, r *http.Request) {
	application.renderTemplate(w, r, "register.page.gohtml", nil)
}

func (application *Config) handlerPostRegisterPage(w http.ResponseWriter, r *http.Request) {}

func (application *Config) handlerGetActivateAccount(w http.ResponseWriter, r *http.Request) {}
