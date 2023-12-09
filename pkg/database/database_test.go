package database

import (
	"fmt"
	"testing"
	"tic-tac-go/pkg/manager"

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

func TestJoinPublicRoom(t *testing.T) {
	t.Parallel()

	roomID := uuid.NewString()
	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	db := miniredis.RunT(t)
	defer db.Close()

	// Arrange
	db.Lpush("Public:Rooms:Available", roomID)
	// Player2 are always the ones joining, player1 creates the room.
	db.HSet(fmt.Sprintf("Room:%s", roomID), "player1ID", player1ID)

	// Act
	database := NewDatabase(db.Addr(), "")
	// Player2 are always the ones joining, player1 creates the room.
	actualRoomID, actualPlayer1ID, err := database.JoinPublicRoom(player2ID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, db.Exists("Public:Rooms:Available"))
	assert.Equal(t, roomID, actualRoomID)
	assert.Equal(t, player1ID, actualPlayer1ID)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "player1ID"), player1ID)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "player2ID"), player2ID)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "gameState"), "000000000")
}

func TestRetrieveGame(t *testing.T) {
	t.Parallel()

	roomID := uuid.NewString()
	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	db := miniredis.RunT(t)
	defer db.Close()

	// Arrange
	db.HSet(fmt.Sprintf("Room:%s", roomID), "player1ID", player1ID, "player2ID", player2ID, "gameState", "000000000")

	// Act
	database := NewDatabase(db.Addr(), "")
	actualGameRoom, err := database.RetrieveGame(roomID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, manager.GameRoom{RoomID: roomID,
		Players: []manager.Player{{ID: player1ID, Message: "000000000"}, {ID: player2ID, Message: "000000000"}},
	}, actualGameRoom)
}
