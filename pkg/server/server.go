package http_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"tic-tac-go/pkg/manager"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// TODO: Perhaps consider using go-chi for the server.

var clients = make(map[string]*websocket.Conn)

type HTTPServer struct {
	server *http.Server
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for WebSocket connections
		return true
	},
}

func NewHTTPServer(port int, manager manager.Manager) *HTTPServer {
	router := mux.NewRouter()

	router.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error Upgrading Connection: %v", err)
			return
		}
		defer conn.Close()

		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to Read Websocket Message: %v", err)
			return
		}

		game, err := manager.StartGame(string(bytes))
		if err != nil {
			log.Printf("Error Starting Game: %v", err)
			return
		}
		for _, player := range game.Players {
			if _, ok := clients[player.ID]; !ok {
				clients[player.ID] = conn
			}
			err = clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message))
			if err != nil {
				log.Printf("Error Sending Websocket Message: %v", err)
				return
			}
		}

		regex := regexp.MustCompile("^[012]{9}$")
		for {
			_, bytes, err = conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to Read Websocket Message: %v", err)
				manager.EndGame(game.RoomID)
				return
			}
			playerMessage := string(bytes)
			// On successful move send game message to both players using
			// clients[playerID] to do so. On unsuccessful moves, send
			// connection to current player using writeJson

			switch {
			case regex.MatchString(playerMessage):
				game, err := manager.ExecutePlayerMove(game.RoomID, playerMessage)
				if err != nil {
					log.Printf("Failed to Execute Player Move: %v", err)
					manager.EndGame(game.RoomID)
					sendErrorMessageToClients(game)
					return
				}
				sendMessageToClients(game)
			case playerMessage == "End Game":
				// TODO: Fix logic, after users end the game, they should be able to start a new one.
				game, err := manager.EndGame(game.RoomID)
				if err != nil {
					log.Printf("Failed to End Game: %v", err)
					manager.EndGame(game.RoomID)
					sendErrorMessageToClients(game)
					return
				}
				sendMessageToClients(game)
			case playerMessage == "Join Room":
				game, err := manager.StartGame(playerMessage)
				if err != nil {
					log.Printf("Error Starting Game: %v", err)
					sendErrorMessageToClients(game)
					return
				}
				for _, player := range game.Players {
					if _, ok := clients[player.ID]; !ok {
						clients[player.ID] = conn
					}
					err = clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message))
					if err != nil {
						log.Printf("Error Sending Websocket Message: %v", err)
						return
					}
				}
			}
		}
	})

	return &HTTPServer{&http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: router}}
}

func (s *HTTPServer) Start() error {
	s.server.ListenAndServe()
	return nil
}

func (s *HTTPServer) Stop() error {
	s.server.Shutdown(context.Background())
	return nil
}

func sendMessageToClients(game manager.GameRoom) {
	for _, player := range game.Players {
		if _, ok := clients[player.ID]; ok {
			clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message))
		}
	}
}

func sendErrorMessageToClients(game manager.GameRoom) {
	for _, player := range game.Players {
		if _, ok := clients[player.ID]; ok {
			clients[player.ID].WriteMessage(websocket.TextMessage, []byte("Server Experienced an Error"))
		}
	}
}
