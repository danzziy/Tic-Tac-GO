package http_server

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"tic-tac-go/pkg/manager"
	"tic-tac-go/pkg/test"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: If you want to run your tests in parrallel, you have to find available ports and pass that
// 		 into your server.

func TestListensForHTTPConnections(t *testing.T) {
	t.Parallel()
	port := test.FindAvailablePort()

	// Arrange
	server := NewHTTPServer(port, nil)
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()

	time.Sleep(10 * time.Millisecond)

	// Act
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	defer func() { _ = conn.Close() }()

	// Assert
	assert.NoError(t, err)
}

func TestStopsListeningForHTTPConnections(t *testing.T) {
	t.Parallel()
	port := test.FindAvailablePort()
	startChan := make(chan error)
	startErr := assert.AnError

	// Arrange
	server := NewHTTPServer(port, nil)
	go func() { startChan <- server.Start() }()

	time.Sleep(10 * time.Millisecond)

	// Act
	_ = server.Stop()

	select {
	case startErr = <-startChan:
	case <-time.After(100 * time.Millisecond):
	}

	_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))

	// Assert
	assert.Error(t, err)
	assert.NoError(t, startErr)
}

func TestUpgradesPublicEndpoitToWebsocketConnection(t *testing.T) {
	t.Parallel()
	port := test.FindAvailablePort()

	// Arrange
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	// Act
	server := NewHTTPServer(port, nil)
	go func() { _ = server.Start() }()
	defer func() { _ = server.Stop() }()

	time.Sleep(10 * time.Millisecond)

	// Assert
	player1 := session.GET("/public").WithWebsocketUpgrade().Expect().
		Status(http.StatusSwitchingProtocols).
		Websocket()
	defer player1.Close()
}

func TestExpectedGameplayForPublicEndpoint(t *testing.T) {
	t.Parallel()
	port := test.FindAvailablePort()
	roomID := uuid.NewString()
	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	gameManager := new(mockManager)

	// Arrange
	session := httpexpect.Default(t, fmt.Sprintf("http://127.0.0.1:%d", port))

	gameManager.On("StartGame", "Join Room").Return(
		manager.GameRoom{
			RoomID: roomID, Players: []manager.Player{
				{ID: player1ID, Message: "Waiting for Player"},
			},
		}, nil,
	).Once()

	gameManager.On("StartGame", "Join Room").Return(
		manager.GameRoom{
			RoomID: roomID, Players: []manager.Player{
				{ID: player1ID, Message: "Start Game"},
				{ID: player2ID, Message: "Start Game"},
			},
		}, nil,
	).Once()

	gameManager.On("ExecutePlayerMove", roomID, "022110000").Return(
		manager.GameRoom{
			RoomID: roomID, Players: []manager.Player{
				{ID: player1ID, Message: "022110000:Ongoing"},
				{ID: player2ID, Message: "022110000:Ongoing"},
			},
		}, nil,
	).Once()

	gameManager.On("ExecutePlayerMove", roomID, "022111000").Return(
		manager.GameRoom{
			RoomID: roomID, Players: []manager.Player{
				{ID: player1ID, Message: "022111000:Win"},
				{ID: player2ID, Message: "022111000:Lose"},
			},
		}, nil,
	).Once()

	gameManager.On("EndGame", roomID).Return(
		manager.GameRoom{
			RoomID: roomID, Players: []manager.Player{
				{ID: player1ID, Message: "Terminate Connection"},
				{ID: player2ID, Message: "Terminate Connection"},
			},
		}, nil,
	).Once()

	// Act
	server := NewHTTPServer(port, gameManager)
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

	player2.WriteText("022110000").Expect().TextMessage().Body().IsEqual("022110000:Ongoing")
	player1.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("022110000:Ongoing")

	player1.WriteText("022111000").Expect().TextMessage().Body().IsEqual("022111000:Win")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("022111000:Lose")

	player1.WriteText("End Game").Expect().TextMessage().Body().IsEqual("Terminate Connection")
	player2.WithoutReadTimeout().Expect().TextMessage().Body().IsEqual("Terminate Connection")

	gameManager.AssertExpectations(t)
}

// TODO: Test that the gameroom does not exist after termination. And a new game room can be started.

type mockManager struct {
	mock.Mock
}

func (m *mockManager) StartGame(message string) (manager.GameRoom, error) {
	args := m.Called(message)
	return args.Get(0).(manager.GameRoom), args.Error(1)
}

func (m *mockManager) ExecutePlayerMove(roomID string, message string) (manager.GameRoom, error) {
	args := m.Called(roomID, message)
	return args.Get(0).(manager.GameRoom), args.Error(1)
}

func (m *mockManager) EndGame(roomID string) (manager.GameRoom, error) {
	args := m.Called(roomID)
	return args.Get(0).(manager.GameRoom), args.Error(1)
}
