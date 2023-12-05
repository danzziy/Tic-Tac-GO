package manager

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartGameWhenNoRoomsAvailable(t *testing.T) {
	database := new(mockDatabase)

	// Arrange
	database.On("PublicRoomAvailable").Return(false).Once()
	database.On("CreatePublicRoom", mock.MatchedBy(matchUUID), mock.MatchedBy(matchUUID)).Return(nil)

	// Act
	manager := NewManager(database)
	actualGameRoom, err := manager.StartGame("Join Room")

	roomID := database.Calls[1].Arguments.String(0)
	playerID := database.Calls[1].Arguments.String(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, GameRoom{roomID, []Player{{playerID, "Waiting for Player"}}}, actualGameRoom)
	database.AssertExpectations(t)
}

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) PublicRoomAvailable() bool {
	args := m.Called()

	return args.Bool(0)
}
func (m *mockDatabase) CreatePublicRoom(roomID string, playerID string) error {
	args := m.Called(roomID, playerID)
	return args.Error(0)
}

func matchUUID(uuid string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(uuid)
}
