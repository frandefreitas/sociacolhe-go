package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }

type Hub struct { rooms map[string]map[*websocket.Conn]bool }

func NewHub() *Hub { return &Hub{rooms: map[string]map[*websocket.Conn]bool{}} }

func (h *Hub) Join(c *gin.Context) {
	roomID := c.Param("roomID")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil { return }
	if h.rooms[roomID] == nil { h.rooms[roomID] = map[*websocket.Conn]bool{} }
	h.rooms[roomID][conn] = true
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil { delete(h.rooms[roomID], conn); conn.Close(); break }
		for cli := range h.rooms[roomID] {
			if err := cli.WriteMessage(websocket.TextMessage, msg); err != nil { log.Println("ws write:", err) }
		}
	}
}
