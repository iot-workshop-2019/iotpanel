package main

import (
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
		Root:      "api/views",
		Extension: ".html",
		Master:    "layouts/master",
		Funcs: template.FuncMap{
			"seq": func(n int) []int {
				nn := make([]int, n)
				for i := 0; i < n; i++ {
					nn[i] = i
				}
				return nn
			},
		},
		DisableCache: true,
	})

	r.Static("/images", "./api/views/images")
	r.Static("/static", "./api/views/static")
	r.Static("/test", "./api/views/")

	api := &handlers.Api{Db: dbp, Cld: cloudSrv}
	// schedule(func(p chan string) {
	// 	fmt.Println("# alive")
	// 	t := time.Now()
	// 	p <- fmt.Sprintf("[%s]: alive", t.String())
	// }, 10*time.Second, api.Data)

	r.GET("/", api.Index)
	r.GET("/pub/:data", api.Publish)
	r.GET("/status", api.Status)
	r.POST("/event", api.Events)
	r.GET("/devup", api.Devicestatus)

	r.Run(":8080")
}
