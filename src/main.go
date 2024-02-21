package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

var clients = make(map[*websocket.Conn]bool)

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// Add the client's WebSocket connection to the clients map
	clients[conn] = true

	for {
		// Read message from the WebSocket connection
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// If there's an error reading the message, remove the client from the map and return
			delete(clients, conn)
			return
		}

		// Send the received message to all connected clients
		for client := range clients {
			if err := client.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	// Define the handler for WebSocket connections
	http.HandleFunc("/ws", handleWebSocketConnection)
	// Handler for the root URL
	http.HandleFunc("/", handleRoot)
	// Print a message indicating the server is running
	fmt.Println("Server is running on http://localhost:8080")
	// Start the HTTP server listening on port 8080
	http.ListenAndServe(":8080", nil)
}
