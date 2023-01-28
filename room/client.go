package room

import (
	"time"

	"github.com/ShingoYadomoto/mahjong-score-board/message"
	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message.Message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message.Message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()

			c.room.forward <- msg
		} else {
			break
		}
	}

	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}

	c.socket.Close()
}
