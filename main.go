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

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", controller.getRoot)
	r.GET("/:room", controller.getRoom)
	r.GET("/:room/ws", controller.getWebSocket)
	r.POST("/create", controller.postCreate)
	r.Run()
}

func getController() *Controller {
	return &Controller{Rooms: map[string]*Room{}}
}

func (controller Controller) getRoot(c *gin.Context) {
	rooms := []string{}
	for key, _ := range controller.Rooms {
		rooms = append(rooms, key)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"rooms": rooms})
}

func (controller *Controller) getRoom(c *gin.Context) {
	roomID := c.Param("room")
	if _, ok := controller.Rooms[roomID]; !ok {
		c.HTML(http.StatusNotFound, "404.html", nil)
	} else {
		c.HTML(http.StatusOK, "room.html", nil)
	}
}

func (controller *Controller) getWebSocket(c *gin.Context) {
	roomID := c.Param("room")
	room := controller.Rooms[roomID]
	room.m.HandleRequest(c.Writer, c.Request)
}

func (controller Controller) postCreate(c *gin.Context) {
	roomID := c.PostForm("id")
	room, ok := controller.Rooms[roomID]
	if !ok {
		m := melody.New()
		controller.Rooms[roomID] = &Room{m: m}
		room = controller.Rooms[roomID]
	}
	room.m.HandleMessage(room.handleMessage)
	c.HTML(http.StatusCreated, "404.html", nil)
}

func (room *Room) handleMessage(s *melody.Session, msg []byte) {
	room.m.Broadcast(msg)
}
