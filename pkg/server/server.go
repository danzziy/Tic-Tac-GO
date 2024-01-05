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
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for WebSocket connections
		return true
	},
}

func NewHTTPServer(port int, manager manager.Manager) *HTTPServer {
	http.HandleFunc("/", website)
	http.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()

		_, bytes, _ := conn.ReadMessage()
		game, _ := manager.StartGame(string(bytes))
		for _, player := range game.Players {
			if _, ok := clients[player.ID]; !ok {
				clients[player.ID] = conn
			}
			clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message))
		}

		regex := regexp.MustCompile("^[012]{9}$")
		for {
			_, bytes, _ = conn.ReadMessage()
			playerMessage := string(bytes)
			// On successful move send game message to both players using
			// clients[playerID] to do so. On unsuccessful moves, send
			// connection to current player using writeJson

			switch {
			case regex.MatchString(playerMessage):
				game, _ := manager.ExecutePlayerMove(game.RoomID, playerMessage)
				sendMessageToClients(game)
			case playerMessage == "End Game":
				// TODO: Fix logic after users end the game, they should be able to start a new one.
				game, _ := manager.EndGame(game.RoomID)
				sendMessageToClients(game)
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

func sendMessageToClients(game manager.GameRoom) {
	for _, player := range game.Players {
		clients[player.ID].WriteMessage(websocket.TextMessage, []byte(player.Message))
	}
}
