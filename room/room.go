package room

import (
	"io"
	"log"
	"net/http"

	"github.com/ShingoYadomoto/mahjong-score-board/message"
	"github.com/stretchr/objx"

	"github.com/ShingoYadomoto/mahjong-score-board/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	forward chan *message.Message // channel for sending messages to others
	join    chan *client          // channel for client joining room room
	leave   chan *client          // channel for client leaving from room room
	clients map[*client]bool
	tracer  trace.Tracer
}

func NewRoom() *room {
	return &room{
		forward: make(chan *message.Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
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
		case client := <-r.join:
			r.joinClient(client)
		case client := <-r.leave:
			r.leaveClient(client)
		case msg := <-r.forward:
			r.sendToClient(msg)
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var checkOrigin = func(r *http.Request) bool {
	return true
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
	CheckOrigin:     checkOrigin,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました: ", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message.Message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}
