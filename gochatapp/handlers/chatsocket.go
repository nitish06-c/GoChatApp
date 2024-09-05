package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mu sync.Mutex

var clients = make(map[*websocket.Conn]bool)

var broadcast = make(chan Message)

type Message struct {
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// Register the new client
	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	log.Println("New client connected")

	// Read messages from the client
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}
