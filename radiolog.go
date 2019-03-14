package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"./api"
	"./cloud"
	"./evcal"
	"github.com/asterix24/radiolog-mqtt/dbi"
	log "github.com/sirupsen/logrus"
)

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

	dbp := &dbi.DBI{}
	for {
		if err := dbp.Init(); err != nil {
			log.Error("DB error", err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	cloud_srv := &server.Server{Db: dbp}
	cloud_srv.Init()

	evcal := &evcal.EvCal{}
	evcal.Init()
	evcal.Events()

	schedule(func() { fmt.Println("# alive") }, 10*time.Second)
	schedule(cloud_srv.Publish, 5*time.Second)

	api := &handlers.Api{Db: dbp}
	http.HandleFunc("/index", api.Index)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}
