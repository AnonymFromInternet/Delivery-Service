package main

import (
	"net/http"
)

func (application *Config) handlerGetMainPage(w http.ResponseWriter, r *http.Request) {
	application.renderTemplate(w, r, "home.page.gohtml", nil)
}
