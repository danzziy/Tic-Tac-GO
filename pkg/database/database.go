package database

import (
	"context"
	"fmt"
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
	return "", "", nil
}

func (d *database) RetrieveGame(roomID string) (manager.GameRoom, error) {
	return manager.GameRoom{}, nil
}

func (d *database) ExecutePlayerMove(GameRoom string, roomID string) error {
	return nil
}

func (d *database) DeleteGameRoom(roomID string) error {
	return nil
}
