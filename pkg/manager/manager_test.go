package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartGameWhenNoRoomsAvailable(t *testing.T) {
	// Arrange

	// Act
	manager := NewManager(nil)
	actualGameRoom, err := manager.StartGame("Join Room")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, GameRoom{}, actualGameRoom)
}
