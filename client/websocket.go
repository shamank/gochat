package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ChatClient struct {
	userID   string
	username string
	roomID   string
	roomName string
	conn     *websocket.Conn
	done     chan struct{}
	rooms    []Room
}

func (c *ChatClient) connect() error {
	u, err := url.Parse(wsURL)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	q := u.Query()
	q.Set("room_id", c.roomID)
	q.Set("user_id", c.userID)
	u.RawQuery = q.Encode()

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	return nil
}

func (c *ChatClient) disconnect() {
	if c.conn != nil {
		_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
		c.conn = nil
	}
}

func (c *ChatClient) readMessages() {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		if msg.UserID != c.userID {
			fmt.Printf("\n[%s]: %s\n", msg.Username, msg.Content)
			if c.roomID == "" {
				fmt.Print("(not in room) > ")
			} else {
				fmt.Printf("[%s] > ", c.roomName)
			}
		}
	}
}

func (c *ChatClient) sendMessage(content string) error {
	url := fmt.Sprintf("%s/api/messages/send?room_id=%s&user_id=%s", serverURL, c.roomID, c.userID)

	reqBody := map[string]string{"content": content}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s", string(body))
	}

	return nil
}
