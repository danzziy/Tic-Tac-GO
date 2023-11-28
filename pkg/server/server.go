package http_server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: Perhaps consider using go-chi for the server.

type HTTPServer struct {
	server *http.Server
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHTTPServer(port int) *HTTPServer {
	http.HandleFunc("/", website)
	http.HandleFunc("/public", publicWebSocket)

	return &HTTPServer{
		&http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port)},
	}
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

func publicWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
}
