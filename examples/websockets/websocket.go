package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var counter = 0

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from the client
		_, msg, err := conn.ReadMessage()
		log.Println("stuff: ")

		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// If the message received is "join room", send "waiting for player"
		if string(msg) == "Join Room" {
			log.Printf("counter: %d", counter)
			log.Println("Message: ", string(msg))

			if counter != -4 {
				responseMsg := []byte("Waiting for Player")
				if err := conn.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
					log.Println("Error writing message:", err)
					break
				}
			} else if counter > 0 {
				responseMsg := []byte("Start Game")
				if err := conn.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
					log.Println("Error writing message:", err)
					break
				}
			}
			counter++

		}
	}
}

func main() {
	http.HandleFunc("/", handleWebSocket)

	port := ":8080"
	log.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
