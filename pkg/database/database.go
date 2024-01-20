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

func NewDatabaseTestClient(redis *redis.Client) manager.Database {
	return &database{redis}
}

func (d *database) PublicRoomAvailable() (bool, error) {
	listLen, err := d.redis.LLen(ctx, "Public:Rooms:Available").Result()
	if err != nil {
		return false, err
	}

	if listLen > 0 {
		return true, nil
	}
	return false, nil
}

func (d *database) CreatePublicRoom(roomID string, playerID string) error {
	if err := d.redis.HSet(ctx, fmt.Sprintf("Room:%s", roomID), "player1ID", playerID).Err(); err != nil {
		return err
	}

	if err := d.redis.LPush(ctx, "Public:Rooms:Available", roomID).Err(); err != nil {
		return err
	}
	return nil
}

func (d *database) JoinPublicRoom(playerID string) (string, error) {
	roomID, err := d.redis.RPop(ctx, "Public:Rooms:Available").Result()
	if err != nil {
		return "", err
	}

	if err := d.redis.HSet(
		ctx, fmt.Sprintf("Room:%s", roomID), "player2ID", playerID, "gameState", "000000000",
	).Err(); err != nil {
		return "", err
	}
	return roomID, nil
}

func (d *database) RetrieveGame(roomID string) (manager.GameRoom, error) {
	room, err := d.redis.HGetAll(ctx, fmt.Sprintf("Room:%s", roomID)).Result()
	if err != nil {
		return manager.GameRoom{}, err
	}

	return manager.GameRoom{RoomID: roomID, Players: []manager.Player{
		{ID: room["player1ID"], Message: room["gameState"]}, {ID: room["player2ID"], Message: room["gameState"]},
	}}, nil
}

func (d *database) ExecutePlayerMove(roomID string, playerMove string) error {
	if err := d.redis.HSet(ctx, fmt.Sprintf("Room:%s", roomID), "gameState", playerMove).Err(); err != nil {
		return err
	}
	return nil
}

func (d *database) DeleteGameRoom(roomID string) error {
	if err := d.redis.LRem(ctx, "Public:Rooms:Available", 0, roomID).Err(); err != nil {
		return err
	}
	if err := d.redis.Del(ctx, fmt.Sprintf("Room:%s", roomID)).Err(); err != nil {
		return err
	}
	return nil
}
