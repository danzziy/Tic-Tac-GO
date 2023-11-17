package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// TestTicTacGoPublicGame, given a user wants to play a game with a stranger,
// then a room will be created where both players will play tic-tac-toe till a
// winner is chosen.
func TestTicTacGoPublicGame(t *testing.T) {
	// Player1 connects to server via websocket sending Join Room Message
	// Player should then get a waiting for player message
	// Player2 connects to server via websocket sending Join Room Message
	// Player1 and Player2 should now recieve a Start Game Message
	// Make game moves until player1 wins.
	port := 8080
	personA := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))
	personB := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	// Arrange

	// Act
	_ = initializeGameServer()

	// Assert
	personA.GET("/").Expect().Status(http.StatusOK)
	player1 := personA.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
	player1.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Waiting for Player")

	personB.GET("/").Expect().Status(http.StatusOK)
	player2 := personB.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player2.Close()
	player2.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Start Game")

	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Start Game")

}

func initializeGameServer() string {
	return ""
}
