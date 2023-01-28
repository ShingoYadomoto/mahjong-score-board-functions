package data

import (
	"errors"

	"github.com/ShingoYadomoto/mahjong-score-board/room"
)

var globalr *room.Room

func GetRoom(id room.ID) (*room.Room, error) {
	return globalr, nil
}

func CreateRoom(pid room.PlayerID) (*room.Room, error) {
	globalr = &room.Room{
		ID:        2,
		PlayerIDs: map[room.PlayerID]struct{}{pid: {}},
	}
	return globalr, nil
}

func AddPlayerToRoom(pid room.PlayerID) error {
	if len(globalr.PlayerIDs) > 4 {
		return errors.New("exceed max member")
	}

	globalr.PlayerIDs[pid] = struct{}{}
	return nil
}

func DeletePlayerFromRoom(pid room.PlayerID) error {
	delete(globalr.PlayerIDs, pid)
	return nil
}
