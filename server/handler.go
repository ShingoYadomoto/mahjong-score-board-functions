package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ShingoYadomoto/mahjong-score-board/player"
	"github.com/ShingoYadomoto/mahjong-score-board/room"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize = 1024
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: room.MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type handler struct{}

func (h *handler) RoomSocketHandler(c echo.Context) error {
	req := c.Request()

	socket, err := upgrader.Upgrade(c.Response(), req, nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	playerMap := objx.MustFromBase64(authCookie.Value)
	p := &player.Player{
		ID:   playerMap["userid"].(string),
		Name: playerMap["name"].(string),
	}

	rm := room.GetRoomManager()

	r := rm.GetRoom()

	client := room.NewClient(socket, r, p)

	r.Join <- client
	defer func() {
		r.Leave <- client
	}()

	go client.Write()
	client.Read()

	return nil
}

func (h *handler) CreatePlayerHandler(c echo.Context) error {
	req := c.Request()

	_, err := req.Cookie("auth")
	if err != nil && err != http.ErrNoCookie {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	var p player.Player
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	m := md5.New()
	if _, err := io.WriteString(m, strings.ToLower(p.Name)); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	p.ID = fmt.Sprintf("%x", m.Sum(nil))

	authCookieValue := objx.New(map[string]interface{}{
		"userid": p.ID,
		"name":   p.Name,
	}).MustBase64()

	c.SetCookie(&http.Cookie{
		Name:     "auth",
		Value:    authCookieValue,
		Path:     "/",
		SameSite: http.SameSiteNoneMode, // only dev
		Expires:  time.Now().Add(time.Hour * 72),
		Secure:   true,
		HttpOnly: true,
	})

	return c.NoContent(http.StatusOK)
}
