package http_server

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"tic-tac-go/pkg/manager"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: If you want to run your tests in parrallel, you have to find available ports and pass that
// 		 into your server.

func TestListensForHTTPConnections(t *testing.T) {
	// Arrange
	server := NewHTTPServer(8080, nil)
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
	server := NewHTTPServer(8080, nil)
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

// TODO: Actually test that you are providing users with frontend content.
func TestRetrieveFrontendContent(t *testing.T) {
	port := 8080

	// Arrange
	server := NewHTTPServer(8080, nil)
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()

	time.Sleep(10 * time.Millisecond)

	// Act & Assert
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))
	session.GET("/").Expect().Status(http.StatusOK)
}

func TestUpgradesPublicEndpoitToWebsocketConnection(t *testing.T) {
	port := 8080

	// Arrange
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	// Act
	server := NewHTTPServer(8080, nil)
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()

	time.Sleep(10 * time.Millisecond)

	// Assert
	player1 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
}

func TestPublicEndpoitToWebsocketConnection(t *testing.T) {
	port := 8080
	roomID := uuid.NewString()
	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	gameManager := new(mockManager)

	// Arrange
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	gameManager.On("StartGame", "Join Room").Return(
		manager.GameRoom{
			RoomID: roomID, PlayerIDs: []string{player1ID}, Message: "Waiting for Player",
		}, nil,
	).Once()

	gameManager.On("StartGame", "Join Room").Return(
		manager.GameRoom{
			RoomID: roomID, PlayerIDs: []string{player1ID, player2ID}, Message: "Start Game",
		}, nil,
	).Once()

	// Act
	server := NewHTTPServer(8080, gameManager)
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()
	time.Sleep(10 * time.Millisecond)

	// Assert
	player1 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).Websocket()
	defer player1.Close()

	player1.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Waiting for Player")

	player2 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).Websocket()
	defer player2.Close()
	player2.WriteText("Join Room").Expect().TextMessage().Body().IsEqual("Start Game")

	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Start Game")

	gameManager.AssertExpectations(t)
}

type mockManager struct {
	mock.Mock
}

func (m *mockManager) StartGame(message string) (manager.GameRoom, error) {
	args := m.Called(message)

	return args.Get(0).(manager.GameRoom), args.Error(1)
}

func (m *mockManager) MakePlayerMove() {
}

func (m *mockManager) EndGame() {
}
