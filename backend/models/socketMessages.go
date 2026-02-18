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
	MessageTypePlayerList = "player_list"
)

type MessageTextContent struct {
	Text string `json:"Text"`
}

type PlayerListMessageContent struct {
	Names []string
}

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

func CreateMessage(Type MessageType, PlayerName string, Content interface{}) Message {
	return Message {
		Type: Type,
		Timestamp: time.Now(),
		PlayerName: PlayerName,
		Content: Content,
	}
}