package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn        // websocket for client
	send     chan *message          // channel sent messages by others
	room     *room                  // chat room client participate
	userData map[string]interface{} // ユーザーに関する情報を保持
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}

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
