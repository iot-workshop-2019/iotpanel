package handlers

import (
	"net/http"

	"github.com/apex/log"
	server "github.com/asterix24/radiolog-mqtt/cloud"
	"github.com/asterix24/radiolog-mqtt/dbi"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Api struct {
	Db  *dbi.DBI
	Cld *server.Server
}

// ButtonJSON stuff ..
type ButtonJSON struct {
	Icons  string `json:"icons"`
	Status string `json:"status"`
}

// Publish ...
func (api *Api) Publish(c *gin.Context) {
	user := c.Params.ByName("name")
	c.JSON(http.StatusOK, gin.H{"user": user, "value": 10})
}

// Index ...
func (api *Api) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{"title": "MQTT Example"})
}

// Status ...
func (api *Api) Status(c *gin.Context) {
	c.HTML(http.StatusOK, "status", gin.H{"url": "ws://" + c.Request.Host + "/devup"})
}

// Events ...
func (api *Api) Events(c *gin.Context) {
	var button ButtonJSON
	c.Bind(&button)
	api.Cld.Publish("show", button.Icons)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

var upgrader = websocket.Upgrader{} // use default options

// Devicestatus ...
func (api *Api) Devicestatus(c *gin.Context) {
	con, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("upgrade:", err)
		return
	}
	defer con.Close()
	for {
		select {
		case d := <-api.Cld.StatusEv:
			err = con.WriteMessage(1, []byte(d.Data))
			if err != nil {
				log.Errorf("write:", err)
				break
			}
		}
		// mt, message, err := con.ReadMessage()
		// if err != nil {
		// 	log.Errorf("read:", err)
		// 	break
		// }
		// log.Infof("recv: %s-%s", mt, message)
		// err = con.WriteMessage(mt, message)
	}
}
