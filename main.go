package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

type Controller struct {
	M *melody.Melody
}

func main() {
	r := gin.Default()
	m := melody.New()

	controller := getController(m)

	r.LoadHTMLGlob("./*.html")

	r.GET("/", controller.getRoot)
	r.GET("/ws", controller.getWebSocket)
	m.HandleMessage(controller.HandleMessage)
	r.Run()
}

func getController(m *melody.Melody) *Controller {
	return &Controller{M: m}
}

func (controller *Controller) getRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (controller *Controller) getWebSocket(c *gin.Context) {
	controller.M.HandleRequest(c.Writer, c.Request)
}

func (controller *Controller) HandleMessage(s *melody.Session, msg []byte) {
	controller.M.Broadcast(msg)
}
