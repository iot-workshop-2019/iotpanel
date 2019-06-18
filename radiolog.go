package main

import (
	"html/template"
	"os"
	"time"

	handlers "github.com/iot-workshop-2019/iotpanel/api"

	server "github.com/iot-workshop-2019/iotpanel/cloud"
	"github.com/iot-workshop-2019/iotpanel/dbi"

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

	staticFileRoot, ok := os.LookupEnv("STATIC_FILE_ROOT")
	if !ok {
		log.Error("STATIC_FILE_ROOT environment variable required but not set")
		staticFileRoot = "./"
		log.Errorf("Switch on default %v", staticFileRoot)

	}

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

	r.Static("/images", staticFileRoot+"api/views/images")
	r.Static("/static", staticFileRoot+"api/views/static")
	r.Static("/test", staticFileRoot+"api/views/")

	api := &handlers.Api{Db: dbp, Cld: cloudSrv}

	r.GET("/", api.Index)
	r.GET("/pub/:data", api.Publish)
	r.GET("/status", api.Status)
	r.POST("/event", api.Events)
	r.GET("/devst", api.Devicestatus)

	r.Run(":8080")
}
