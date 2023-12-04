package http_server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"tic-tac-go/pkg/manager"

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
}

func NewHTTPServer(port int, manager manager.Manager) *HTTPServer {
	http.HandleFunc("/", website)
	http.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		regex := regexp.MustCompile("^[012]{9}$")

		var jsonMessage wsMessage
		_ = conn.ReadJSON(&jsonMessage)
		game, _ := manager.StartGame(jsonMessage.Message)

		for _, playerID := range game.PlayerIDs {
			if _, ok := clients[playerID]; !ok {
				clients[playerID] = conn
			}
			clients[playerID].WriteJSON(wsMessage{game.RoomID, playerID, game.Message})
		}

		for {
			_ = conn.ReadJSON(&jsonMessage)

			switch {
			case regex.MatchString(jsonMessage.Message):
				// On successful move send game message to both players using
				// clients[playerID] to do so. On unsuccessful moves, send
				// connection to current player using writeJson
				manager.MakePlayerMove()
			case jsonMessage.Message == "End Game":
				manager.EndGame()
			}
		}
	})

	return &HTTPServer{&http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port)}}
}

func (s *HTTPServer) Start() error {
	s.server.ListenAndServe()
	return nil
}

func (s *HTTPServer) Stop() error {
	s.server.Shutdown(context.Background())
	return nil
}

func website(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type wsMessage struct {
	GameID   string `json:"gameId"`
	PlayerID string `json:"playerId"`
	Message  string `json:"message"`
}
