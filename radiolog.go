package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/asterix24/radiolog-mqtt/api"
	"github.com/asterix24/radiolog-mqtt/cloud"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

func initDb() (*sql.DB, error) {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Unable to connect to DB", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error("Unable to talk with DB", err)
		return nil, err
	}
	log.Info("Successfully connected!")
	return db, nil
}

func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	var db_ptr *sql.DB
	var err error

	for {
		defer db.Close()
		db_ptr, err = initDb()
		if err != nil {
			log.Error("DB error", err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	cloud_srv := &server.Server{Db: db_ptr}
	cloud_srv.Init()

	schedule(func() { fmt.Println("# alive") }, 10*time.Second)
	schedule(cloud_srv.Publish, 5*time.Second)

	api := &handlers.Api{Db: db_ptr}
	http.HandleFunc("/index", api.Index)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}
