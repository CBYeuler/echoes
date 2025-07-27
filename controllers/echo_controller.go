package controllers

import (
	"log"
	"net/http"

	//"github.com/CBYeuler/echoes/database"
	"github.com/CBYeuler/echoes/models"
	"github.com/gin-gonic/gin"
)

// HandleEcho processes the user's input text and returns a response from GPT
// It expects a JSON body with user_text field
// This function is called when the user sends a message to the Echoes API
// It retrieves the username from the context, creates a new message instance,

func HandleEcho(c *gin.Context) {
	var input struct {
		UserText string `json:"user_text"`
	}
	// HandleEcho processes the user's input text and returns a response from GPT
	// It expects a JSON body with user_text field

	// Bind the JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Failed to bind JSON input:", err) // Log binding errors
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// Retrieve the username from the context
	// This assumes the AuthMiddleware has set the username in the context
	username, exists := c.Get("username")
	if !exists {
		log.Println("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create a new message instance
	// This message will be saved to the database
	msg := models.NewMessage(username.(string), input.UserText)

	reply, err := utils.sendToGPT(input.UserText)
	if err := msg.SaveMessage(); err != nil {
		log.Println("Failed to save message:", err) // Log database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}
	msg.UpdateGPTReply(reply)
	if err := msg.SaveMessage(); err != nil {
		log.Println("Failed to update message:", err) // Log database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
		// Simulate a GPT reply (in a real application, this would call the OpenAI API)
		//gptReply := "Echo: " + input.UserText // This is a placeholder for the
	}
	c.JSON(http.StatusOK, gin.H{
		"username":          username,
		"user_text":         input.UserText,
		"gpt_reply":         reply,
		"formatted_message": msg.GetFormattedMessage(),
	})
}
