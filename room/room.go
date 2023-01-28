package room

import (
	"io"

	"github.com/ShingoYadomoto/mahjong-score-board/message"
	"github.com/ShingoYadomoto/mahjong-score-board/trace"
)

type room struct {
	ID      RoomID
	forward chan *message.Message // channel for sending messages to others
	Join    chan *client          // channel for client joining room room
	Leave   chan *client          // channel for client leaving from room room
	clients map[*client]bool
	tracer  trace.Tracer
}

func (r *room) SetTracer(w io.Writer) {
	r.tracer = trace.New(w)
}

func (r *room) joinClient(c *client) {
	r.clients[c] = true
	r.tracer.Trace("新たなユーザーが参加しました")
}

func (r *room) leaveClient(c *client) {
	delete(r.clients, c)
	close(c.send)
	r.tracer.Trace("クライアントが退室しました")
}

func (r *room) sendToClient(msg *message.Message) {
	r.tracer.Trace("メッセージを送信しました: ", msg.Message)
	// send messages to all client
	for client := range r.clients {
		select {
		case client.send <- msg:
			// send message
			r.tracer.Trace(" -- クライアントに送信されました")
		default:
			// fail to send message
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace(" -- 送信に失敗しました。 クライアントをクリーンアップします。")
		}
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.joinClient(client)
		case client := <-r.Leave:
			r.leaveClient(client)
		case msg := <-r.forward:
			r.sendToClient(msg)
		}
	}
}
