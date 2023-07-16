package main

import (
	"database/sql"
	"encoding/gob"
	"github.com/AnonymFromInternet/Delivery-Service/data"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db := initDataBase()

	session := initSession()

	infoLog := log.New(os.Stdout, "INFO \t:", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR \t:", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	wg := &sync.WaitGroup{}

	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     wg,
		Models:   data.New(db),
	}

	// set up mail

	app.listenForShutdown()

	app.serve()

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

func initSession() *scs.SessionManager {
	gob.Register(data.User{})

	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
		MaxIdle: 10,
	}
}
