package database

import (
	"context"
	"fmt"
	"log"
	"tic-tac-go/pkg/manager"

	"github.com/redis/go-redis/v9"
)

// TODO: You'll want to pass ctx from server layer all the way down here.
var ctx = context.Background()

type database struct {
	redis *redis.Client
}

func NewDatabase(address string, password string) manager.Database {
	return &database{
		redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
		}),
	}
}

func (d *database) PublicRoomAvailable() (bool, error) {
	listLen, _ := d.redis.LLen(ctx, "Public:Rooms:Available").Result()
	if listLen > 0 {
		return true, nil
	}

	return false, nil
}

func (d *database) CreatePublicRoom(roomID string, playerID string) error {
	_ = d.redis.LPush(ctx, "Public:Rooms:Available", roomID)
	_ = d.redis.HSet(ctx, fmt.Sprintf("Room:%s", roomID), "player1ID", playerID)
	return nil
}

func (d *database) JoinPublicRoom(playerID string) (string, string, error) {
	roomID, _ := d.redis.RPop(ctx, "Public:Rooms:Available").Result()
	opponentID, _ := d.redis.HGet(ctx, fmt.Sprintf("Room:%s", roomID), "player1ID").Result()
	d.redis.HSet(ctx, fmt.Sprintf("Room:%s", roomID), "player2ID", playerID, "gameState", "000000000")
	return roomID, opponentID, nil
}

func (d *database) RetrieveGame(roomID string) (manager.GameRoom, error) {
	room, _ := d.redis.HGetAll(ctx, fmt.Sprintf("Room:%s", roomID)).Result()
	log.Printf("Game: %v", room)
	return manager.GameRoom{RoomID: roomID, Players: []manager.Player{
		{ID: room["player1ID"], Message: room["gameState"]}, {ID: room["player2ID"], Message: room["gameState"]},
	}}, nil
}

func (d *database) ExecutePlayerMove(roomID string, playerMove string) error {
	_ = d.redis.HSet(ctx, fmt.Sprintf("Room:%s", roomID), "gameState", playerMove)
	return nil
}

func (d *database) DeleteGameRoom(roomID string) error {
	return nil
}
