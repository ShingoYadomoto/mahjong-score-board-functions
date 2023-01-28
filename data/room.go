package data

import (
	"errors"

	"github.com/ShingoYadomoto/mahjong-score-board/room"
)

var globalr *room.Room

var ErrNotFound = errors.New("data: ErrNotFound")

func GetJoinedRoom(pid room.PlayerID) (*room.Room, error) {
	if globalr.Joined(pid) {
		return globalr, nil
	}
	return nil, ErrNotFound
}

func CreateRoom(pid room.PlayerID) (*room.Room, error) {
	globalr = &room.Room{
		ID:        2,
		PlayerIDs: map[room.PlayerID]struct{}{pid: {}},
	}
	return globalr, nil
}

func AddPlayerToRoom(id room.ID, pid room.PlayerID) error {
	if len(globalr.PlayerIDs) > 4 {
		return errors.New("exceed max member")
	}

	globalr.PlayerIDs[pid] = struct{}{}
	return nil
}

func DeletePlayerFromRoom(id room.ID, pid room.PlayerID) error {
	delete(globalr.PlayerIDs, pid)
	return nil
}
