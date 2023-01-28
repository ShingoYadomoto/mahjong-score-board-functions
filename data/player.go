package data

import (
	"github.com/ShingoYadomoto/mahjong-score-board/room"
)

var globalp *room.Player

func GetPlayer(id room.PlayerID) (*room.Player, error) {
	if globalp == nil {
		return nil, ErrNotFound
	}

	return globalp, nil
}

func CreatePlayer(name string) (*room.Player, error) {
	globalp = &room.Player{
		ID:   1,
		Name: name,
	}
	return globalp, nil
}
