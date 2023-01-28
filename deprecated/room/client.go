package room

import (
	"time"

	"github.com/ShingoYadomoto/mahjong-score-board/deprecated/message"
	"github.com/gorilla/websocket"
)

const MessageBufferSize = 256

type Player struct {
	ID     string
	Name   string `json:"name"`
	RoomID RoomID
}

type client struct {
	socket *websocket.Conn
	send   chan *message.Message
	room   *room
	player *Player
}

func NewClient(socket *websocket.Conn, room *room, player *Player) *client {
	return &client{
		socket: socket,
		send:   make(chan *message.Message, MessageBufferSize),
		room:   room,
		player: player,
	}
}

func (c *client) Read() {
	for {
		var msg *message.Message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.player.Name

			c.room.Forward <- msg
		} else {
			break
		}
	}

	c.socket.Close()
}

func (c *client) Write() {
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}

	c.socket.Close()
}
