package main

import (
	"gochatapp/handlers"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handlers.HandleConnections)

	go handlers.HandleMessages()

	log.Println("Server started on :8080") // started server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
