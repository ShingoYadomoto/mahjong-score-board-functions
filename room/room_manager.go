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

type RoomID int

type roomManager struct {
	room map[RoomID]*room
}

func newRoomManager() *roomManager {
	return &roomManager{room: map[RoomID]*room{}}
}

func (rm *roomManager) NewRoom() (*room, error) {
	if len(rm.room) >= maxRoomNum {
		return nil, errors.New("exceed max room num")
	}

	var id RoomID
	//for {
	//	if _, exist := rm.room[id]; !exist {
	//		break
	//	}
	//	id++
	//}
	id = 3

	r := &room{
		ID:      id,
		forward: make(chan *message.Message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}

	rm.room[id] = r

	return r, nil
}

func (rm *roomManager) GetRoom(id RoomID) *room {
	return rm.room[id]
}

var singletonRoomManager *roomManager

func GetRoomManager() *roomManager {
	if singletonRoomManager == nil {
		singletonRoomManager = newRoomManager()
	}

	return singletonRoomManager
}
