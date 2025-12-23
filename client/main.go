package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Go Chat Client ===")
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		log.Fatal("Username cannot be empty")
	}

	user, err := registerUser(username)
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}

	fmt.Printf("Registered as: %s (ID: %s)\n\n", user.Username, user.ID)

	chatClient := &ChatClient{
		userID:   user.ID,
		username: user.Username,
		done:     make(chan struct{}),
	}

	if err := chatClient.refreshRooms(); err != nil {
		log.Printf("Warning: Failed to get rooms: %v", err)
	}

	if len(chatClient.rooms) == 0 {
		fmt.Println("No rooms available. Use '/create <name>' to create a room.")
	}

	showCommands()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = chatClient.refreshRooms()
			case <-chatClient.done:
				return
			}
		}
	}()

	go chatClient.readInput(reader)

	<-sigChan
	fmt.Println("\nDisconnecting...")
	chatClient.disconnect()
}
