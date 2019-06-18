package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	server "github.com/iot-workshop-2019/iotpanel/cloud"
	"github.com/iot-workshop-2019/iotpanel/dbi"
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
	c.HTML(http.StatusOK, "index", gin.H{"title": "IoT Panel"})
}

// Test ...
func (api *Api) Test(c *gin.Context) {
	c.HTML(http.StatusOK, "test", gin.H{"title": "Terminal"})
}

// Status ...
func (api *Api) Status(c *gin.Context) {
	c.HTML(http.StatusOK, "status", gin.H{})
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

// DevState ...
type DevState struct {
	Node      string `json:"node"`
	Data      string `json:"data"`
	Kind      string `json:"kind"`
	Count     int    `json:"count"`
	Timestamp string `json:"timestamp"`
}

// DevStatus ...
type DevStatus []DevState

// Devicestatus ...
func (api *Api) Devicestatus(c *gin.Context) {
	l := api.Db.StatusNode()
	var v DevStatus
	for _, item := range l {
		v = append(v, DevState{Node: item.Node, Data: item.Data, Kind: item.Kind, Count: item.Count, Timestamp: item.Timestamp.Format("2006-01-02 15:04:05")})
	}
	c.JSON(http.StatusOK, v)
}
