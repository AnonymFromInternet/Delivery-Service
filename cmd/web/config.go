package main

import (
	"database/sql"
	"fmt"
	"github.com/AnonymFromInternet/Delivery-Service/data"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var webPort = ":80"

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
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

func (application *Config) listenForShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	application.shutdown()
	os.Exit(0)
}

func (application *Config) shutdown() {
	application.InfoLog.Println("would run cleanup tasks")

	application.Wait.Wait()

	application.InfoLog.Println("closing channels and shutting down application...")
}
