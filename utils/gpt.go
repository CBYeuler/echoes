package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type GPTResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func SendToGPT(userText string) (string, error) {
	// Prepare the request to OpenAI's GPT API
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY is not set")
	}
	// Create the request body
	reqBody := GPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatMessage{
			{Role: "user", Content: userText},
		},
	}
	// Marshal the request body to JSON
	jsonData, err := json.Marshal(reqBody)

	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	// Check if the response status is OK
	var gptResp GPTResponse
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body
	if err := json.NewDecoder(resp.Body).Decode(&gptResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if we received any choices
	if len(gptResp.Choices) == 0 {
		return "", errors.New("no choices received from GPT")
	}

	// Return the content of the first choice's message
	return gptResp.Choices[0].Message.Content, nil
}

// This function can be used in the application to send user text to GPT and receive a reply
// It handles the HTTP request to OpenAI's API, processes the response, and returns the
// GPT reply or an error if something goes wrong
// The function uses the OpenAI API key from the environment variables
// It returns the GPT reply as a string or an error if the request fails
