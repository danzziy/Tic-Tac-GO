package http_server

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(port int) *HTTPServer {
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
