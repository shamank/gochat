package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func (c *ChatClient) handleCommand(cmd string, reader *bufio.Reader) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "/rooms":
		c.showRooms()
		return nil

	case "/create":
		if len(parts) < 2 {
			fmt.Println("Usage: /create <room_name>")
			return nil
		}
		roomName := strings.Join(parts[1:], " ")
		return c.createNewRoom(roomName)

	case "/join":
		if len(parts) < 2 {
			fmt.Println("Usage: /join <room_number>")
			c.showRooms()
			return nil
		}
		var roomIndex int
		if _, err := fmt.Sscanf(parts[1], "%d", &roomIndex); err != nil {
			return fmt.Errorf("invalid room number")
		}
		return c.joinRoom(roomIndex)

	case "/leave":
		return c.leaveRoom()

	case "/history":
		limit := 10
		if len(parts) >= 2 {
			if parsedLimit, err := fmt.Sscanf(parts[1], "%d", &limit); err != nil || parsedLimit == 0 {
				limit = 10
			}
		}
		return c.showHistory(limit)

	case "/exit":
		os.Exit(0)
		return nil

	case "/help":
		showCommands()
		return nil

	default:
		fmt.Printf("Unknown command: %s\n", parts[0])
		fmt.Println("Type '/help' to see available commands")
		return nil
	}
}

func (c *ChatClient) showRooms() {
	if err := c.refreshRooms(); err != nil {
		fmt.Printf("Failed to refresh rooms: %v\n", err)
		return
	}

	if len(c.rooms) == 0 {
		fmt.Println("No rooms available. Use '/create <name>' to create one.")
		return
	}

	fmt.Println("\nAvailable rooms:")
	for i, room := range c.rooms {
		current := ""
		if room.ID == c.roomID {
			current = " (current)"
		}
		fmt.Printf("  %d. %s%s\n", i+1, room.Name, current)
	}
	fmt.Println()
}

func (c *ChatClient) createNewRoom(roomName string) error {
	if roomName == "" {
		return fmt.Errorf("room name cannot be empty")
	}

	newRoom, err := createRoom(roomName)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	if err := c.refreshRooms(); err != nil {
		log.Printf("Failed to refresh rooms: %v", err)
	}

	fmt.Printf("Room '%s' created successfully! Use '/join %d' to join it.\n", newRoom.Name, len(c.rooms))
	return nil
}

func (c *ChatClient) joinRoom(roomIndex int) error {
	if err := c.refreshRooms(); err != nil {
		return fmt.Errorf("failed to refresh rooms: %w", err)
	}

	if roomIndex < 1 || roomIndex > len(c.rooms) {
		return fmt.Errorf("invalid room number. Use '/rooms' to see available rooms")
	}

	newRoom := &c.rooms[roomIndex-1]
	if newRoom.ID == c.roomID {
		fmt.Println("You are already in this room.")
		return nil
	}

	if c.conn != nil {
		c.disconnect()
	}

	c.roomID = newRoom.ID
	c.roomName = newRoom.Name

	if err := c.connect(); err != nil {
		fmt.Printf("Warning: Failed to connect to WebSocket: %v\n", err)
		fmt.Println("You can still send messages, but won't receive real-time updates.")
	} else {
		go c.readMessages()
	}

	fmt.Printf("Joined room: %s\n", newRoom.Name)

	messages, err := getMessagesHistory(newRoom.ID, 10, 0)
	if err == nil && len(messages) > 0 {
		fmt.Println("\n--- Recent Messages ---")
		for _, msg := range messages {
			fmt.Printf("[%s]: %s\n", msg.Username, msg.Content)
		}
		fmt.Println("--- End History ---")
	}

	return nil
}

func (c *ChatClient) leaveRoom() error {
	if c.roomID == "" {
		fmt.Println("You are not in any room.")
		return nil
	}

	fmt.Printf("Leaving room: %s\n", c.roomName)

	if c.conn != nil {
		c.disconnect()
	}

	c.roomID = ""
	c.roomName = ""
	fmt.Println("Left the room. Use '/join <number>' to join another room.")
	return nil
}

func (c *ChatClient) showHistory(limit int) error {
	if c.roomID == "" {
		fmt.Println("You are not in any room. Use '/join <number>' to join a room first.")
		return nil
	}

	messages, err := getMessagesHistory(c.roomID, limit, 0)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(messages) == 0 {
		fmt.Println("No messages in this room.")
		return nil
	}

	fmt.Printf("\n--- Message History (last %d) ---\n", len(messages))
	for _, msg := range messages {
		fmt.Printf("[%s]: %s\n", msg.Username, msg.Content)
	}
	fmt.Println("--- End History ---")
	return nil
}

func (c *ChatClient) refreshRooms() error {
	rooms, err := getAllRooms()
	if err != nil {
		return err
	}
	c.rooms = rooms
	return nil
}
