package database

import (
	"fmt"
	"testing"
	"tic-tac-go/pkg/manager"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestPublicRoomsAvailable(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		numberOfRooms int64
		expctedResult bool
	}{
		{0, false}, // No Rooms Available
		{1, true},  // 1 Room Available
		{2, true},  // 2 Rooms Available
	} {
		tc := tc
		t.Run(fmt.Sprintf("Number of Rooms: %d", tc.numberOfRooms), func(t *testing.T) {
			t.Parallel()
			db, mock := redismock.NewClientMock()
			defer db.Close()

			// Arrange
			mock.ExpectLLen("Public:Rooms:Available").SetVal(tc.numberOfRooms)

			// Act
			database := NewDatabaseTestClient(db)
			roomAvailable, err := database.PublicRoomAvailable()

			// Assert
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expctedResult, roomAvailable)
		})
	}
}

func TestCreatePublicRoom(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHSet(fmt.Sprintf("Room:%s", roomID), "player1ID", playerID).RedisNil()
	mock.ExpectLPush("Public:Rooms:Available", roomID).RedisNil()

	// Act
	database := NewDatabaseTestClient(db)
	err := database.CreatePublicRoom(roomID, playerID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJoinPublicRoom(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	player2ID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectRPop("Public:Rooms:Available").SetVal(roomID)
	mock.ExpectHSet(
		fmt.Sprintf("Room:%s", roomID), "player2ID", player2ID, "gameState", "000000000",
	).RedisNil()

	// Act
	database := NewDatabaseTestClient(db)
	// Player2 are always the ones joining, player1 creates the room.
	actualRoomID, err := database.JoinPublicRoom(player2ID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, roomID, actualRoomID)
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

func TestExecutePlayerMove(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerMove := "000010000"

	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	db := miniredis.RunT(t)
	defer db.Close()

	// Arrange
	db.HSet(fmt.Sprintf("Room:%s", roomID), "player1ID", player1ID, "player2ID", player2ID, "gameState", "000000000")

	// Act
	database := NewDatabase(db.Addr(), "")
	err := database.ExecutePlayerMove(roomID, playerMove)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "player1ID"), player1ID)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "player2ID"), player2ID)
	assert.Equal(t, db.HGet(fmt.Sprintf("Room:%s", roomID), "gameState"), "000010000")
}

func TestDeleteGameRoom(t *testing.T) {
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
	err := database.DeleteGameRoom(roomID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, db.Exists(fmt.Sprintf("Room:%s", roomID)))
}
