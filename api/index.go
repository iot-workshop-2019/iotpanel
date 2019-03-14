package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asterix24/radiolog-mqtt/dbi"
)

// repository contains the details of a repository
type radiologSummary struct {
	Address   int
	Ntc0      string
	Ntc1      string
	Timestamp time.Time
}

type Data struct {
	Data []radiologSummary
}

type Api struct {
	Db *dbi.DBI
}

// Index ...
func (api *Api) Index(w http.ResponseWriter, r *http.Request) {

	api.Db.Temperature(10)
	//data := Data{}

	/*
		rows, err := api.Db.Query("SELECT address, ntc0, ntc1, timestamp FROM radiologdata")
		if err != nil {
			log.Fatal("Unable to talk with db", err)
			http.Error(w, err.Error(), 500)
			return
		}

		defer rows.Close()
		for rows.Next() {
			summary := radiologSummary{}
			err = rows.Scan(
				&summary.Address,
				&summary.Ntc0,
				&summary.Ntc1,
				&summary.Timestamp)

			if err != nil {
				log.Fatal("Unable to get data", err)
				http.Error(w, err.Error(), 500)
				return
			}
			data.Data = append(data.Data, summary)
		}

		err = rows.Err()
		if err != nil {
			log.Fatal("Unable to get row", err)
			http.Error(w, err.Error(), 500)
			return
		}

		out, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	*/
	fmt.Fprintf(w, string("prova"))
}
