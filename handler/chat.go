package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var rooms = make(map[string][]*websocket.Conn)

func ChatWebSocket(c *gin.Context) {
	orderID := c.Param("orderId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	log.Println("📡 Client connected to order:", orderID)

	rooms[orderID] = append(rooms[orderID], conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Connection closed:", err)
			break
		}

		// Broadcast ke semua koneksi di room (kecuali sender)
		for _, other := range rooms[orderID] {
			if other != conn {
				err := other.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("⚠️ Broadcast failed:", err)
				}
			}
		}
	}
}
