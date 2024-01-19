package test

import (
	"log"
	"net"
)

// FindAvailablePort
func FindAvailablePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic("Error finding available port")
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	log.Printf("PORT %d", port)
	return listener.Addr().(*net.TCPAddr).Port
}
