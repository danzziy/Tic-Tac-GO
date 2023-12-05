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

		_, bytes, _ := conn.ReadMessage()
		game, _ := manager.StartGame(string(bytes))
		for i, playerID := range game.PlayerIDs {
			if _, ok := clients[playerID]; !ok {
				clients[playerID] = conn
			}
			clients[playerID].WriteMessage(websocket.TextMessage, []byte(game.Message[i]))
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
				game, _ := manager.MakePlayerMove(game.RoomID, playerMessage)
				for i, playerID := range game.PlayerIDs {
					clients[playerID].WriteMessage(websocket.TextMessage, []byte(game.Message[i]))
				}
			case playerMessage == "End Game":
				game := manager.EndGame(game.RoomID)
				for i, playerID := range game.PlayerIDs {
					clients[playerID].WriteMessage(websocket.TextMessage, []byte(game.Message[i]))
				}
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
