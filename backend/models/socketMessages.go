package models

import "time"

type MessageType string

// Define your message types
const (
	MessageTypeChat       = "chat"
	MessageTypeJoin       = "join"
	MessageTypeLeave      = "leave"
	MessageTypeGameUpdate = "game_update"
	MessageTypeError      = "error"
)

type Message struct {
	Type      MessageType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	// PlayerID  string      `json:"playerId,omitempty"` TODO: PlayerIDs
	PlayerName string      `json:"playerName,omitempty"`
	Content    interface{} `json:"content,omitempty"`
}

type PlayerListContent struct {
	PlayerList []string
}
