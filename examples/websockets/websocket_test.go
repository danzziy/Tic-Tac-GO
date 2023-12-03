package main_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// TestWebsockets serves as an example of how to test websockets.
func TestWebsockets(t *testing.T) {
	// Act
	port := 8080
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	// Assert
	player1 := session.GET("/").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
	// player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Start Game")

	player1.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Waiting for Player")

	player2 := session.GET("/").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player2.Close()

	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Start Game")
	log.Printf("P1: %v", player1)
	log.Printf("P2: %v", player2)

}

func TestPublicEndpoitToWebsocketConnection(t *testing.T) {
	port := 8080

	// Act & Assert
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))
	player1 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
	player1.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Waiting for Player")

	session.GET("/").Expect().Status(http.StatusOK)
	player2 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player2.Close()
	player2.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Start Game")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Start Game")

	player1.WriteText("000010000").Expect().TextMessage().Body().IsEqual("000010000")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("000010000")

	player2.WriteText("200010000").Expect().TextMessage().Body().IsEqual("200010000")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("200010000")

	player1.WriteText("200110000").Expect().TextMessage().Body().IsEqual("200110000")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("200110000")

	player2.WriteText("220110000").Expect().TextMessage().Body().IsEqual("220110000")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("220110000")

	player1.WriteText("220111000").Expect().TextMessage().Body().IsEqual("220111000")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("220111000")
}
