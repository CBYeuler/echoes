package models

import (
	"time"

	"github.com/CBYeuler/echoes/database"
)

// Message represents a message in the system
// It includes fields for the message ID, username, user text, GPT reply, and timestamps

type Message struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	UserText  string    `json:"user_text"` // User's input text
	GPTReply  string    `json:"gpt_reply"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Message) TableName() string {
	return "messages"
}

// NewMessage creates a new message instance
func NewMessage(username, userText string) *Message {
	now := time.Now()
	return &Message{
		Username:  username,
		UserText:  userText,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateGPTReply updates the GPT reply for the message
func (m *Message) UpdateGPTReply(reply string) {
	m.GPTReply = reply
	m.UpdatedAt = time.Now()
}

// GetFormattedMessage returns a formatted string representation of the message
func (m *Message) GetFormattedMessage() string {
	return m.Username + ": " + m.UserText + " (GPT Reply: " + m.GPTReply + ")"
}

// SaveMessage saves the message to the database
func (m *Message) SaveMessage() error {
	query := `
		INSERT INTO messages (username, user_text, gpt_reply, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := database.DB.Exec(query, m.Username, m.UserText, m.GPTReply, m.CreatedAt, m.UpdatedAt)
	return err
}

// GetMessageByID retrieves a message by its ID
func GetMessageByID(id int) (*Message, error) {
	query := `
		SELECT id, username, user_text, gpt_reply, created_at, updated_at
		FROM messages WHERE id = ?
	`
	row := database.DB.QueryRow(query, id)

	var m Message
	err := row.Scan(&m.ID, &m.Username, &m.UserText, &m.GPTReply, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
