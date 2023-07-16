package main

import "net/http"

func (application *Config) SessionLoad(next http.Handler) http.Handler {
	return application.Session.LoadAndSave(next)
}
