package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to the database
	db := initDataBase()
	err := db.Ping()
	if err != nil {
		return
	}

	// create sessions

	// create channels

	// crate wg

	// create and populate app config

	// set up mail

	// listen for web connections
}

func initDataBase() *sql.DB {
	conn := connectToDataBase()
	if conn == nil {
		log.Fatal("cannot connect to database")
	}

	return conn
}

func connectToDataBase() *sql.DB {
	triesAmount := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDataBase(dsn)
		if err != nil {
			log.Println("postgres is not ready yet")
		} else {
			log.Println("connected")
			return connection
		}

		if triesAmount > 10 {
			return nil
		}

		log.Println("waiting for 1 second by the database connection")
		time.Sleep(1 * time.Second)
		triesAmount++
		continue
	}
}

func openDataBase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
