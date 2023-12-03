package handler

import (
	"log"
	"user/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *models.Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				return
			}
			c.Conn.WriteJSON(message)
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if messageType == websocket.TextMessage {
			msg := &models.Message{
				Content:  string(data),
				RoomID:   c.RoomID,
				Username: c.Username,
			}
			hub.Broadcast <- msg

		}
	}
}
