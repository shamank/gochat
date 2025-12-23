package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func showCommands() {
	fmt.Println("Available commands:")
	fmt.Println("  /rooms              - Show all rooms")
	fmt.Println("  /create <name>      - Create a new room")
	fmt.Println("  /join <number>      - Join a room by number")
	fmt.Println("  /leave              - Leave current room")
	fmt.Println("  /history [limit]    - Show message history (default: 10)")
	fmt.Println("  /exit               - Quit application")
	fmt.Println()
}

func (c *ChatClient) readInput(reader *bufio.Reader) {
	for {
		if c.roomID == "" {
			fmt.Print("(not in room) > ")
		} else {
			fmt.Printf("[%s] > ", c.roomName)
		}

		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Error reading input: %v", err)
			continue
		}

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		if strings.HasPrefix(text, "/") {
			if err := c.handleCommand(text, reader); err != nil {
				log.Printf("Command error: %v", err)
			}
			continue
		}

		if text == "exit" || text == "/exit" {
			os.Exit(0)
		}

		if c.roomID == "" {
			fmt.Println("You are not in any room. Use '/join <number>' to join a room first.")
			continue
		}

		if err := c.sendMessage(text); err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			if c.conn == nil {
				fmt.Printf("[%s]: %s\n", c.username, text)
			}
		}
	}
}
