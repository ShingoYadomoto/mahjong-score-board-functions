package room

import (
	"errors"

	"github.com/ShingoYadomoto/mahjong-score-board/message"
	"github.com/ShingoYadomoto/mahjong-score-board/trace"
)

func init() {
	singletonRoomManager = newRoomManager()
}

const maxRoomNum = 10

type roomID int

type roomManager struct {
	room map[roomID]*room
}

func newRoomManager() *roomManager {
	return &roomManager{room: map[roomID]*room{}}
}

func (rm *roomManager) NewRoom() (*room, error) {
	if len(rm.room) >= maxRoomNum {
		return nil, errors.New("exceed max room num")
	}

	var id roomID
	for {
		if _, exist := rm.room[id]; !exist {
			break
		}
		id++
	}

	r := &room{
		forward: make(chan *message.Message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}

	rm.room[id] = r

	return r, nil
}

// ToDo
func (rm *roomManager) GetRoom() *room {
	for _, r := range rm.room {
		return r
	}

	return nil
}

var singletonRoomManager *roomManager

func GetRoomManager() *roomManager {
	if singletonRoomManager == nil {
		singletonRoomManager = newRoomManager()
	}

	return singletonRoomManager
}
