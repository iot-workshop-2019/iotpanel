package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	handlers "github.com/asterix24/radiolog-mqtt/api"
	"github.com/asterix24/radiolog-mqtt/evcal"

	server "github.com/asterix24/radiolog-mqtt/cloud"
	"github.com/asterix24/radiolog-mqtt/dbi"
	gintemplate "github.com/foolin/gin-template"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func schedule(what func(p chan string), delay time.Duration, param chan string) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what(param)
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

	cloudSrv := &server.Server{Db: dbp}
	cloudSrv.Init()

	evcal := &evcal.EvCal{}
	evcal.Init()
	evcal.Events()

	// schedule(cloudSrv.Publish, 20*time.Second)

	r := gin.Default()
	r.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         "api/views",
		Extension:    ".html",
		Master:       "layouts/master",
		Funcs:        template.FuncMap{},
		DisableCache: true,
	})

	api := &handlers.Api{Db: dbp, Cld: cloudSrv, Data: make(chan string)}
	schedule(func(p chan string) {
		fmt.Println("# alive")
		t := time.Now()
		p <- fmt.Sprintf("[%s]: alive", t.String())
	}, 10*time.Second, api.Data)

	r.GET("/pub/:data", api.Publish)
	r.GET("/index", api.Index)
	r.GET("/status", api.Status)
	r.POST("/event", api.Events)
	r.GET("/ws", api.WShandler)

	r.Run(":8080")
}
