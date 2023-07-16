package main

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"sync"
)

var webPort = ":80"

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
}

func (application *Config) serve() {
	server := &http.Server{
		Addr:                         fmt.Sprintf("%s", webPort),
		Handler:                      application.routes(),
		DisableGeneralOptionsHandler: false,
	}

	err := server.ListenAndServe()
	if err != nil {
		application.ErrorLog.Println("cannot serve the server :", err)

		log.Fatal("cannot serve the server")
	}

	application.InfoLog.Println("starting the server")
}
