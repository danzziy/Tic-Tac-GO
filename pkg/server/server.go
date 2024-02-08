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
				// You will send an termination message after a player has disconnected to both players.
				// Given that a player has disconnected, only one of those messages will go through.
				// The message that fails to send due to a closed socket will result in the closed socket being purged.
				gameInfo, _ := manager.EndGame(game.RoomID)

				for _, player := range gameInfo.Players {
					if _, ok := clients[player.ID]; ok {
						log.Printf("Player Message Opponent Left: %s", "Opponent Left")

						if err := clients[player.ID].WriteMessage(websocket.TextMessage, []byte("Opponent Left")); err != nil {
							delete(clients, player.ID)
							log.Printf("Error Sending Websocket Message: %v", err)
						}
					}
				}
				return
			}
			playerMessage := string(bytes)

			switch {
			case regex.MatchString(playerMessage):
				gameInfo, err := manager.ExecutePlayerMove(game.RoomID, playerMessage)
				if err != nil {
					log.Printf("Failed to Execute Player Move: %v", err)
					return
				}

				sendMessageToClients(gameInfo)
			case playerMessage == "End Game":
				gameInfo, err := manager.EndGame(game.RoomID)
				if err != nil {
					log.Printf("Failed to End Game: %v", err)
					_, _ = manager.EndGame(game.RoomID)
					return
				}
				sendMessageToClients(gameInfo)
			case playerMessage == "Join Room":
				game, err = manager.StartGame(playerMessage)
				if err != nil {
					log.Printf("Error Starting Game: %v", err)
					return
				}

				for _, player := range game.Players {
					if _, ok := clients[player.ID]; !ok {
						clients[player.ID] = conn
					}
					log.Printf("Player Message Join Room: %s", player.Message)

					err = clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message))
					if err != nil {
						log.Printf("Error Sending Websocket Message: %v", err)
						return
					}
				}
			}
		}
	})

	return &HTTPServer{&http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: router}}
}

func (s *HTTPServer) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop() error {
	if err := s.server.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

func sendMessageToClients(game manager.GameRoom) {
	for _, player := range game.Players {
		if _, ok := clients[player.ID]; ok {
			log.Printf("Player Message Send Message: %s", player.Message)

			if err := clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message)); err != nil {
				delete(clients, player.ID)
				log.Printf("Error Sending Websocket Message: %v", err)
			}
		}
	}
}
