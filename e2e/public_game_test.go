package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"tic-tac-go/pkg/analyzer"
	"tic-tac-go/pkg/config"
	"tic-tac-go/pkg/database"
	"tic-tac-go/pkg/manager"
	game "tic-tac-go/pkg/server"
	"tic-tac-go/pkg/test"

	"github.com/alicebob/miniredis/v2"
	"github.com/gavv/httpexpect/v2"
)

// TODO: Check the database has been updated.

// TestTicTacGoPublicGame, given a user wants to play a game with a stranger,
// then a room will be created where both players will play tic-tac-toe till a
// winner is chosen.
func TestTicTacGoPublicGame(t *testing.T) {
	t.Parallel()
	port := test.FindAvailablePort()
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))
	db := miniredis.RunT(t)
	defer db.Close()

	// Arrange

	// Act
	server := initializeGameServer(port, db.Addr())
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()

	// Assert
	player1 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
	player1.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Waiting for Player")

	player2 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player2.Close()

	player2.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Start Game")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Start Game")

	player1.WriteText("000010000").Expect().TextMessage().Body().IsEqual("000010000:Ongoing")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("000010000:Ongoing")

	player2.WriteText("200010000").Expect().TextMessage().Body().IsEqual("200010000:Ongoing")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("200010000:Ongoing")

	player1.WriteText("200110000").Expect().TextMessage().Body().IsEqual("200110000:Ongoing")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("200110000:Ongoing")

	player2.WriteText("220110000").Expect().TextMessage().Body().IsEqual("220110000:Ongoing")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("220110000:Ongoing")

	player1.WriteText("220111000").Expect().TextMessage().Body().IsEqual("220111000:Win")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("220111000:Lose")

	player1.WriteText("End Game")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Terminate Connection")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Terminate Connection")
}

func initializeGameServer(listeningPort int, dbAddr string) *game.HTTPServer {
	env := config.NewConfig(
		[]string{fmt.Sprintf("LISTENING_PORT=%d", listeningPort), fmt.Sprintf("DATABASE_HOST=%s", dbAddr), "DATABASE_PASSWORD=something"},
	)
	port, _ := env.ListeningPort()
	databaseHost, _ := env.DatabaseHost()
	databasePassword, _ := env.DatabasePassword()
	return game.NewHTTPServer(port, manager.NewManager(
		database.NewDatabase(databaseHost, databasePassword), analyzer.NewAnalyzer()),
	)
}
