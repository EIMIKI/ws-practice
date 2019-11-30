package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Controller struct {
	Rooms map[string]*Room
}
type Room struct {
	m *melody.Melody
}

func main() {
	r := gin.Default()

	controller := getController()

	r.LoadHTMLGlob("./*.html")

	r.GET("/:room", controller.getRoot)
	r.GET("/:room/ws", controller.getWebSocket)
	r.Run()
}

func getController() *Controller {
	return &Controller{Rooms: map[string]*Room{}}
}

func (controller *Controller) getRoot(c *gin.Context) {
	roomID := c.Param("room")
	room, ok := controller.Rooms[roomID]
	if !ok {
		m := melody.New()
		controller.Rooms[roomID] = &Room{m: m}
		room = controller.Rooms[roomID]
	}

	room.m.HandleMessage(room.handleMessage)
	c.HTML(http.StatusOK, "index.html", nil)
}

func (controller *Controller) getWebSocket(c *gin.Context) {
	roomID := c.Param("room")
	room := controller.Rooms[roomID]
	room.m.HandleRequest(c.Writer, c.Request)
}

func (room *Room) handleMessage(s *melody.Session, msg []byte) {
	room.m.Broadcast(msg)
}
