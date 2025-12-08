package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/InspectorGadget/realtime-polling-system/db"
	"github.com/InspectorGadget/realtime-polling-system/helpers"
	"github.com/InspectorGadget/realtime-polling-system/models"
	"github.com/InspectorGadget/realtime-polling-system/redis"
	"github.com/gin-gonic/gin"
)

func CreatePoll(c *gin.Context) {
	var req models.CreatePollRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Start a transaction
	tx, _ := db.GetDB().Begin()

	var pollID int
	// 1. Insert Poll
	err := tx.QueryRow("INSERT INTO polls (topic) VALUES ($1) RETURNING id", req.Topic).Scan(&pollID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
		return
	}

	// 2. Insert Options
	for _, optionText := range req.Options {
		// FIX: Changed column 'id' to 'poll_id'
		_, err := tx.Exec(
			"INSERT INTO options (poll_id, text) VALUES ($1, $2)",
			pollID,     // $1
			optionText, // $2
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll options"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll created successfully",
		"poll_id": pollID,
	})
}

func GetPoll(c *gin.Context) {
	pollID := c.Param("id")

	poll, err := helpers.FetchPollFromDB(pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch poll data"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func CastVote(c *gin.Context) {
	var req struct {
		PollID   int `json:"poll_id"`
		OptionID int `json:"option_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := db.GetDB().Exec(
		"UPDATE options SET votes = votes + 1 WHERE id = $1 AND poll_id = $2",
		req.OptionID,
		req.PollID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast vote"})
		return
	}

	// FIX: Use strconv.Itoa or fmt.Sprintf to convert int ID to string for the helper
	pollIDStr := strconv.Itoa(req.PollID)

	// Fetch updated data to broadcast
	poll, err := helpers.FetchPollFromDB(pollIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated poll data"})
		return
	}

	// Marshall json data
	jsonData, _ := json.Marshal(poll)

	// Publish updated poll data to Redis
	redis.GetRedis().Publish(c.Request.Context(), fmt.Sprintf("poll:%d", req.PollID), jsonData)

	c.JSON(http.StatusOK, gin.H{"message": "Vote cast successfully"})
}
