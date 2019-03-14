package dbi

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Load postgress driver
	log "github.com/sirupsen/logrus"
)

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

// DBI  database interface
type DBI struct {
	db *sql.DB
}

// Init init db module interface
func (dbp *DBI) Init() error {
	config := dbConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	var err error
	dbp.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Unable to connect to DB: ", err)
		return err
	}

	defer dbp.db.Close()
	err = dbp.db.Ping()
	if err != nil {
		log.Error("Unable to talk with DB: ", err)
		return err
	}
	log.Info("Successfully connected!")
	return nil
}
