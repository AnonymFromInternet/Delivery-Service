package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Message       string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	// User *data.User
}

func (application *Config) renderTemplate(w http.ResponseWriter, r *http.Request, tN string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	var templates []string
	templates = append(templates, fmt.Sprintf("%s/%s", pathToTemplates, tN))

	for _, partial := range partials {
		templates = append(templates, partial)
	}

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		application.ErrorLog.Println("cannot parse files :", err)

		return
	}

	err = tmpl.Execute(w, application.addDefaultData(td, r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		application.ErrorLog.Println("cannot execute the template :", err)
		return
	}
}

func (application *Config) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Message = application.Session.PopString(r.Context(), "message")
	td.Warning = application.Session.PopString(r.Context(), "warning")
	td.Error = application.Session.PopString(r.Context(), "error")

	if application.isAuthenticated(r) {
		td.Authenticated = true
		// TODO add more info about user
	}

	td.Now = time.Now()

	return td
}

func (application *Config) isAuthenticated(r *http.Request) bool {
	return application.Session.Exists(r.Context(), "userID")
}
