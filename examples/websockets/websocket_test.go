package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// TestWebsockets serves as an example of how to test websockets.
func TestWebsockets(t *testing.T) {
	// Act
	port := 8080
	personA := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	// Assert
	player1 := personA.GET("/").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
	player1.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Waiting for Player")

}
