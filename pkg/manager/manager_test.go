package manager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartGameWhenNoPublicRoomsAreAvailable(t *testing.T) {
	t.Parallel()
	database := new(mockDatabase)

	// Arrange
	database.On("PublicRoomAvailable").Return(false, nil).Once()
	database.On("CreatePublicRoom", mock.MatchedBy(matchUUID), mock.MatchedBy(matchUUID)).Return(nil)

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	roomID := database.Calls[1].Arguments.String(0)
	playerID := database.Calls[1].Arguments.String(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, GameRoom{roomID, []Player{{playerID, "Waiting for Player"}}}, actualGameRoom)
	database.AssertExpectations(t)
}

func TestStartGameWhenAPublicRoomIsAvailable(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	player2ID := uuid.NewString()

	database := new(mockDatabase)

	expectedGameRoom := GameRoom{roomID, []Player{
		{uuid.NewString(), "Start Game"},
		{player2ID, "Start Game"},
	}}

	// Arrange
	database.On("PublicRoomAvailable").Return(true, nil).Once()
	database.On("JoinPublicRoom", matcher(matchUUID)).Return(roomID, nil).Once()
	database.On("RetrieveGame", matcher(matchUUID)).Return(expectedGameRoom, nil).Once()

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	// Assert
	assert.NoError(t, err)
	database.AssertExpectations(t)
	assert.Equal(t, expectedGameRoom, actualGameRoom)
}

func TestStartGameFailsWhenRetrievingPublicRoomIsAvailable(t *testing.T) {
	t.Parallel()
	database := new(mockDatabase)

	// Arrange
	database.On("PublicRoomAvailable").Return(false, assert.AnError).Once()

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualGameRoom)
	database.AssertExpectations(t)
}

func TestStartGameWhenCreatingPublicRoomFail(t *testing.T) {
	t.Parallel()
	database := new(mockDatabase)

	// Arrange
	database.On("PublicRoomAvailable").Return(false, nil).Once()
	database.On("CreatePublicRoom", mock.MatchedBy(matchUUID), mock.MatchedBy(matchUUID)).Return(assert.AnError)

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualGameRoom)
	database.AssertExpectations(t)
}

func TestStartGameWhenJoiningPublicRoomFail(t *testing.T) {
	t.Parallel()
	database := new(mockDatabase)

	// Arrange
	database.On("PublicRoomAvailable").Return(true, nil).Once()
	database.On("JoinPublicRoom", mock.MatchedBy(matchUUID)).Return("", assert.AnError)

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualGameRoom)
	database.AssertExpectations(t)
}

func TestStartGameWhenRetrievingRoomFails(t *testing.T) {
	t.Parallel()
	database := new(mockDatabase)
	roomID := uuid.NewString()

	// Arrange
	database.On("PublicRoomAvailable").Return(true, nil).Once()
	database.On("JoinPublicRoom", matcher(matchUUID)).Return(roomID, nil).Once()
	database.On("RetrieveGame", matcher(matchUUID)).Return(GameRoom{}, assert.AnError).Once()

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualGameRoom)
	database.AssertExpectations(t)
}

func TestExecutesPlayerMove(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		prevGameState    string
		playerMove       string
		expectedGameRoom GameRoom
	}{
		{
			"000000000", "000010000", GameRoom{uuid.NewString(), []Player{
				{uuid.NewString(), "000010000:Ongoing"},
				{uuid.NewString(), "000010000:Ongoing"},
			}},
		},
		{
			"022110000", "022111000", GameRoom{uuid.NewString(), []Player{
				{uuid.NewString(), "022111000:Win"},
				{uuid.NewString(), "022111000:Lose"},
			}},
		},
		{
			"220110001", "222110001", GameRoom{uuid.NewString(), []Player{
				{uuid.NewString(), "222110001:Lose"},
				{uuid.NewString(), "222110001:Win"},
			}},
		},
	} {
		tc := tc
		t.Run(fmt.Sprintf("with player move %s", tc.playerMove), func(t *testing.T) {
			t.Parallel()
			gameRoom := GameRoom{tc.expectedGameRoom.RoomID, []Player{
				{tc.expectedGameRoom.Players[0].ID, tc.prevGameState},
				{tc.expectedGameRoom.Players[1].ID, tc.prevGameState},
			}}

			database := new(mockDatabase)
			analyzer := new(mockAnalyzer)

			// Arrange
			database.On("RetrieveGame", matcher(matchUUID)).Return(gameRoom, nil).Once()
			analyzer.On("ValidMove", tc.prevGameState, matcher(matchGameBoard)).Return(true).Once()
			database.On("ExecutePlayerMove", matcher(matchUUID), matcher(matchGameBoard)).Return(nil).Once()
			analyzer.On("DetermineWinner", matcher(matchGameBoard), gameRoom.Players).Return(tc.expectedGameRoom.Players).Once()

			// Act
			manager := NewManager(database, analyzer)
			actualGameRoom, err := manager.ExecutePlayerMove(tc.expectedGameRoom.RoomID, tc.playerMove)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedGameRoom, actualGameRoom)
			mock.AssertExpectationsForObjects(t, database, analyzer)
		})
	}
}

func TestEndGame(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	expectedGameRoom := GameRoom{
		roomID, []Player{{player1ID, "Terminate Connection"}, {player2ID, "Terminate Connection"}},
	}

	database := new(mockDatabase)

	// Arrange
	// TODO: You should definitely have a data struct dedicated to retrieving gameRooms.
	// and another for communicating to the clients. Because, it makes no sense for
	// RetrieveGame to return a GameRoom that contains player messages.
	database.On("RetrieveGame", matcher(matchUUID)).Return(GameRoom{roomID, []Player{{player1ID, ""}, {player2ID, ""}}}, nil).Once()
	database.On("DeleteGameRoom", matcher(matchUUID)).Return(nil).Once()

	// Act
	manager := NewManager(database, nil)
	actualGameRoom, err := manager.EndGame(roomID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedGameRoom, actualGameRoom)
	database.AssertExpectations(t)
}

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) PublicRoomAvailable() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *mockDatabase) CreatePublicRoom(roomID string, playerID string) error {
	args := m.Called(roomID, playerID)
	return args.Error(0)
}

func (m *mockDatabase) JoinPublicRoom(playerID string) (string, error) {
	args := m.Called(playerID)
	return args.String(0), args.Error(1)
}

func (m *mockDatabase) RetrieveGame(roomID string) (GameRoom, error) {
	args := m.Called(roomID)
	return args.Get(0).(GameRoom), args.Error(1)
}

func (m *mockDatabase) ExecutePlayerMove(roomID string, playerMove string) error {
	args := m.Called(roomID, playerMove)
	return args.Error(0)
}

func (m *mockDatabase) DeleteGameRoom(roomID string) error {
	args := m.Called(roomID)
	return args.Error(0)
}

type mockAnalyzer struct {
	mock.Mock
}

func (m *mockAnalyzer) ValidMove(prevGameState string, playerMove string) bool {
	args := m.Called(prevGameState, playerMove)
	return args.Bool(0)
}

func (m *mockAnalyzer) DetermineWinner(playerMove string, players []Player) []Player {
	args := m.Called(playerMove, players)
	return args.Get(0).([]Player)
}

func matcher(fn interface{}) interface{} {
	return mock.MatchedBy(fn)
}

func matchUUID(uuid string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(uuid)
}

func matchGameBoard(gamestate string) bool {
	pattern := `^[012]{9}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(gamestate)
}
