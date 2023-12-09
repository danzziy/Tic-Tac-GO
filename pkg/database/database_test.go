package database

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestPublicRoomsAvailable(t *testing.T) {
	t.Parallel()

	db := miniredis.RunT(t)

	// Arrange
	db.Lpush("Public:Rooms:Available", uuid.NewString())
	defer db.Close()

	// Act
	database := NewDatabase(db.Addr(), "")
	roomAvailable, err := database.PublicRoomAvailable()

	// Assert
	assert.NoError(t, err)
	assert.True(t, roomAvailable)
}

func TestPublicRoomsAreNotAvailable(t *testing.T) {
	t.Parallel()

	db := miniredis.RunT(t)
	defer db.Close()

	// Act
	database := NewDatabase(db.Addr(), "")
	roomAvailable, err := database.PublicRoomAvailable()

	// Assert
	assert.NoError(t, err)
	assert.False(t, roomAvailable)
}

func TestCreatePublicRoom(t *testing.T) {
	t.Parallel()

	roomID := uuid.NewString()
	playerID := uuid.NewString()
	db := miniredis.RunT(t)
	defer db.Close()

	// Act
	database := NewDatabase(db.Addr(), "")
	err := database.CreatePublicRoom(roomID, playerID)

	// Assert
	assert.NoError(t, err)
	db.CheckList(t, "Public:Rooms:Available", roomID)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "player1ID"), playerID)
}
