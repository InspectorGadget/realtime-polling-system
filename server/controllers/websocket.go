package controllers

import (
	"log"
	"net/http"

	"github.com/InspectorGadget/realtime-polling-system/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	pollID := c.Param("poll_id")

	// Upgrade HTTP connection to Websocket
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WS Upgrade Error:", err)
		return
	}
	defer ws.Close()

	// Subscribe to Redis channel for the poll
	subscriber := redis.GetRedis().Subscribe(c.Request.Context(), "poll:"+pollID)
	defer subscriber.Close()

	// Go routine to handle incoming messages (if needed)
	ch := subscriber.Channel()
	for msg := range ch {
		err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			log.Println("WS Write Error:", err)
			break
		}
	}
}
