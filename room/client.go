package room

import (
	"time"

	"github.com/ShingoYadomoto/mahjong-score-board/message"
	"github.com/ShingoYadomoto/mahjong-score-board/player"
	"github.com/gorilla/websocket"
)

const MessageBufferSize = 256

type client struct {
	socket *websocket.Conn
	send   chan *message.Message
	room   *room
	player *player.Player
}

func NewClient(socket *websocket.Conn, room *room, player *player.Player) *client {
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

			c.room.forward <- msg
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
