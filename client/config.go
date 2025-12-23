package main

import (
	"os"
	"strings"
)

var (
	serverURL = getServerURL()
	wsURL     = getWebSocketURL()
)

func getServerURL() string {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		url = "http://localhost:8080"
	}
	return url
}

func getWebSocketURL() string {
	url := os.Getenv("WS_URL")
	if url == "" {
		httpURL := getServerURL()
		httpURL = strings.Replace(httpURL, "http://", "ws://", 1)
		httpURL = strings.Replace(httpURL, "https://", "wss://", 1)
		url = httpURL + "/ws"
	}
	return url
}
