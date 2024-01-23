package database

import (
	"fmt"
	"testing"
	"tic-tac-go/pkg/manager"

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
			assert.Equal(t, tc.expctedResult, roomAvailable)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFailedToVerifyPublicRoomAvailability(t *testing.T) {
	t.Parallel()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectLLen("Public:Rooms:Available").SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	roomAvailable, err := database.PublicRoomAvailable()

	// Assert
	assert.Error(t, err)
	assert.False(t, roomAvailable)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePublicRoom(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHSet(fmt.Sprintf("Room:%s", roomID), "player1ID", playerID).SetVal(0)
	mock.ExpectLPush("Public:Rooms:Available", roomID).SetVal(0)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.CreatePublicRoom(roomID, playerID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePublicRoomFailsToCreateRoom(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHSet(fmt.Sprintf("Room:%s", roomID), "player1ID", playerID).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.CreatePublicRoom(roomID, playerID)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePublicRoomFailsToMakeRoomAvailable(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHSet(fmt.Sprintf("Room:%s", roomID), "player1ID", playerID).SetVal(0)
	mock.ExpectLPush("Public:Rooms:Available", roomID).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.CreatePublicRoom(roomID, playerID)

	// Assert
	assert.Error(t, err)
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
	).SetVal(0)

	// Act
	database := NewDatabaseTestClient(db)
	// Player2 are always the ones joining, player1 creates the room.
	actualRoomID, err := database.JoinPublicRoom(player2ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roomID, actualRoomID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJoinPublicRoomWhenRetrievingAvailableRoomsFail(t *testing.T) {
	t.Parallel()
	player2ID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectRPop("Public:Rooms:Available").SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	actualRoomID, err := database.JoinPublicRoom(player2ID)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualRoomID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJoinPublicRoomWhenJoiningRoomFails(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	player2ID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectRPop("Public:Rooms:Available").SetVal(roomID)
	mock.ExpectHSet(
		fmt.Sprintf("Room:%s", roomID), "player2ID", player2ID, "gameState", "000000000",
	).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	actualRoomID, err := database.JoinPublicRoom(player2ID)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualRoomID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrieveGame(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	player1ID := uuid.NewString()
	player2ID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHGetAll(fmt.Sprintf("Room:%s", roomID)).SetVal(map[string]string{
		"player1ID": player1ID, "player2ID": player2ID, "gameState": "000000000",
	})

	// Act
	database := NewDatabaseTestClient(db)
	actualGameRoom, err := database.RetrieveGame(roomID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, manager.GameRoom{RoomID: roomID,
		Players: []manager.Player{{ID: player1ID, Message: "000000000"}, {ID: player2ID, Message: "000000000"}},
	}, actualGameRoom)
}

func TestRetrieveGameFails(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHGetAll(fmt.Sprintf("Room:%s", roomID)).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	actualGameRoom, err := database.RetrieveGame(roomID)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, actualGameRoom)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExecutePlayerMove(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerMove := "000010000"

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHSet(fmt.Sprintf("Room:%s", roomID), "gameState", playerMove).SetVal(0)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.ExecutePlayerMove(roomID, playerMove)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExecutePlayerMoveFails(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()
	playerMove := "000010000"

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectHSet(fmt.Sprintf("Room:%s", roomID), "gameState", playerMove).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.ExecutePlayerMove(roomID, playerMove)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteGameRoom(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectLRem("Public:Rooms:Available", 0, roomID).SetVal(0)
	mock.ExpectDel(fmt.Sprintf("Room:%s", roomID)).SetVal(0)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.DeleteGameRoom(roomID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteGameRoomFailsToRemovePublicRoomFromAvailableList(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectLRem("Public:Rooms:Available", 0, roomID).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.DeleteGameRoom(roomID)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteGameRoomFailsToRemoveRoom(t *testing.T) {
	t.Parallel()
	roomID := uuid.NewString()

	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Arrange
	mock.ExpectLRem("Public:Rooms:Available", 0, roomID).SetVal(0)
	mock.ExpectDel(fmt.Sprintf("Room:%s", roomID)).SetErr(assert.AnError)

	// Act
	database := NewDatabaseTestClient(db)
	err := database.DeleteGameRoom(roomID)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
