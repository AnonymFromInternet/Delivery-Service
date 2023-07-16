package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (application *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(application.SessionLoad)

	mux.Get("/", application.handlerGetMainPage)
	mux.Get("/login", application.handlerGetLoginPage)
	mux.Get("/logout", application.handlerGetLogoutPage)
	mux.Get("/register", application.handlerGetRegisterPage)
	mux.Get("/activate-account", application.handlerGetActivateAccount)

	mux.Post("/login", application.handlerPostLoginPage)
	mux.Post("/register", application.handlerPostRegisterPage)

	return mux
}
