package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func registerUser(username string) (*User, error) {
	url := fmt.Sprintf("%s/api/users/register", serverURL)

	reqBody := map[string]string{"username": username}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf(apiResp.Error)
	}

	userData, _ := json.Marshal(apiResp.Data)
	var user User
	_ = json.Unmarshal(userData, &user)

	return &user, nil
}

func getAllRooms() ([]Room, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/rooms/all", serverURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf(apiResp.Error)
	}

	roomsData, _ := json.Marshal(apiResp.Data)
	var rooms []Room
	_ = json.Unmarshal(roomsData, &rooms)

	return rooms, nil
}

func createRoom(name string) (*Room, error) {
	url := fmt.Sprintf("%s/api/rooms/create", serverURL)

	reqBody := map[string]string{"name": name}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf(apiResp.Error)
	}

	roomData, _ := json.Marshal(apiResp.Data)
	var room Room
	_ = json.Unmarshal(roomData, &room)

	return &room, nil
}

func getMessagesHistory(roomID string, limit, offset int) ([]Message, error) {
	url := fmt.Sprintf("%s/api/messages/history?room_id=%s&limit=%d&offset=%d", serverURL, roomID, limit, offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf(apiResp.Error)
	}

	messagesData, _ := json.Marshal(apiResp.Data)
	var messages []Message
	_ = json.Unmarshal(messagesData, &messages)

	return messages, nil
}
