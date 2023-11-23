package http_server

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: If you want to run your tests in parrallel, you have to find available ports and pass that
// 		 into your server.

func TestListensForHTTPConnections(t *testing.T) {
	// Arrange
	server := NewHTTPServer(8080)
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()

	time.Sleep(10 * time.Millisecond)

	// Act
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	// Assert
	assert.NoError(t, err)
	conn.Close()
}

func TestStopsListeningForHTTPConnections(t *testing.T) {
	startChan := make(chan error)
	startErr := assert.AnError

	// Arrange
	server := NewHTTPServer(8080)
	go func() { startChan <- server.Start() }()

	time.Sleep(10 * time.Millisecond)

	// Act
	_ = server.Stop()

	select {
	case startErr = <-startChan:
	case <-time.After(100 * time.Millisecond):
	}

	_, err := net.Dial("tcp", "127.0.0.1:8080")

	// Assert
	assert.Error(t, err)
	assert.NoError(t, startErr)
}
