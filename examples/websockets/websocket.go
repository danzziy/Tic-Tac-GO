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
	responseMsg := []byte("Start Game")
	if err := conn.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
		log.Println("Error writing message:", err)
	}

	for {
		log.Println("In for loop: ")
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

			responseMsg := []byte("Waiting for Player")
			if err := conn.WriteMessage(websocket.TextMessage, responseMsg); err != nil {
				log.Println("Error writing message:", err)
				break
			}

		}
	}
}

var clients []*websocket.Conn
var num = 0

func publicWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	clients = append(clients, conn)
	num++
	for {
		_, msg, _ := conn.ReadMessage()

		switch string(msg) {
		case "Join Room":
			if num <= 1 {
				conn.WriteMessage(websocket.TextMessage, []byte("Waiting for Player"))
			} else {
				for _, c := range clients {
					c.WriteMessage(websocket.TextMessage, []byte("Start Game"))
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/", handleWebSocket)
	http.HandleFunc("/public", publicWebSocket)

	port := ":8080"
	log.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
