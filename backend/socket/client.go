package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adettinger/go-quizgame/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	ID       uuid.UUID
	Conn     *websocket.Conn
	Manager  *Manager
	Send     chan models.Message
	UserData UserData
}

type UserData struct {
	PlayerId uuid.UUID
	Name     string
}

func (client Client) Logf(message string, args ...interface{}) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %v", client.ID, message))
	for _, arg := range args {
		switch v := arg.(type) {
		case error:
			sb.WriteString(": " + v.Error())
		case fmt.Stringer:
			sb.WriteString(v.String())
		default:
			sb.WriteString(fmt.Sprintf("%v", v))
		}
	}
	log.Print(sb.String())
}

func (client *Client) ErrorAndKill(errorMsg string) {
	errorToSend := models.Message{
		Type:       models.MessageTypeError,
		Timestamp:  time.Now(),
		PlayerName: "System",
		Content:    models.MessageTextContent{Text: errorMsg},
	}
	jsonMessage, err := json.Marshal(errorToSend)
	if err != nil {
		client.Logf("Error marshaling message: ", err)
		close(client.Send)
		client.Conn.Close()
		return
	}
	if err := client.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		client.Logf("Error writing message", err)
		close(client.Send)
		client.Conn.Close()
		return
	}
	close(client.Send)
	client.Conn.Close()
}
