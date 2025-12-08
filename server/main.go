package main

import (
	"context"

	"github.com/InspectorGadget/realtime-polling-system/controllers"
	"github.com/InspectorGadget/realtime-polling-system/db"
	"github.com/InspectorGadget/realtime-polling-system/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func init() {
	godotenv.Load()

	// Connect to Database
	_, err := db.Connect()
	if err != nil {
		panic(err)
	}

	// Connect to Redis
	if err := redis.Connect(ctx); err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/polls", controllers.CreatePoll)     // Create a new poll
	r.GET("/polls/:id", controllers.GetPoll)     // Get poll details
	r.POST("/vote", controllers.CastVote)        // Cast a vote (triggers update)
	r.GET("/ws/:poll_id", controllers.WsHandler) // Websocket endpoint

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
