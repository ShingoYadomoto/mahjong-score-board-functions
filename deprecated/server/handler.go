package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ShingoYadomoto/mahjong-score-board/deprecated/message"
	room2 "github.com/ShingoYadomoto/mahjong-score-board/deprecated/room"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize = 1024
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: room2.MessageBufferSize,
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
	p := &room2.Player{
		ID:     playerMap["user_id"].(string),
		Name:   playerMap["name"].(string),
		RoomID: room2.RoomID(playerMap["room_id"].(float64)),
	}

	rm := room2.GetRoomManager()

	r := rm.GetRoom(p.RoomID)
	if r == nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	client := room2.NewClient(socket, r, p)

	r.Join <- client
	defer func() {
		r.Leave <- client
	}()

	go client.Write()
	client.Read()

	return nil
}

func (h *handler) getPlayerMap(req *http.Request) map[string]interface{} {
	authCookie, err := req.Cookie("auth")
	if err != nil {
		return nil
	}

	return objx.MustFromBase64(authCookie.Value)
}

func (h *handler) playerCreated(d map[string]interface{}) bool {
	for _, field := range []string{"user_id", "name"} {
		if _, ok := d[field]; !ok {
			return false
		}
	}
	return true
}

func (h *handler) joinedRoom(d map[string]interface{}) bool {
	for _, field := range []string{"user_id", "name", "room_id"} {
		if _, ok := d[field]; !ok {
			return false
		}
	}
	return true
}

func (h *handler) responseWithCookie(c echo.Context, d map[string]interface{}) error {
	authCookieValue := objx.New(d).MustBase64()

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

func (h *handler) CreatePlayerHandler(c echo.Context) error {
	req := c.Request()

	var p room2.Player
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	m := md5.New()
	if _, err := io.WriteString(m, strings.ToLower(p.Name)); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	p.ID = fmt.Sprintf("%x", m.Sum(nil))

	return h.responseWithCookie(c, map[string]interface{}{
		"user_id": p.ID,
		"name":    p.Name,
	})
}

func (h *handler) CreateRoomHandler(c echo.Context) error {
	var (
		req       = c.Request()
		playerMap = h.getPlayerMap(req)
		rm        = room2.GetRoomManager()
	)

	if !h.playerCreated(playerMap) {
		return c.NoContent(http.StatusBadRequest)
	}

	//r, err := rm.NewRoom()
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, err.Error())
	//}
	r := rm.GetRoom(3)

	return h.responseWithCookie(c, map[string]interface{}{
		"user_id": playerMap["user_id"],
		"name":    playerMap["name"],
		"room_id": r.ID,
	})
}

func (h *handler) JoinRoomHandler(c echo.Context) error {
	var (
		req       = c.Request()
		roomIDStr = c.Param("roomID")
		playerMap = h.getPlayerMap(req)
		rm        = room2.GetRoomManager()
	)

	roomIDInt, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	if !h.playerCreated(playerMap) {
		return c.NoContent(http.StatusBadRequest)
	}

	r := rm.GetRoom(room2.RoomID(roomIDInt))
	if r == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return h.responseWithCookie(c, map[string]interface{}{
		"user_id": playerMap["user_id"],
		"name":    playerMap["name"],
		"room_id": r.ID,
	})
}

func (h *handler) LeaveRoomHandler(c echo.Context) error {
	var (
		req       = c.Request()
		playerMap = h.getPlayerMap(req)
	)

	if !h.playerCreated(playerMap) {
		return c.NoContent(http.StatusBadRequest)
	}

	return h.responseWithCookie(c, map[string]interface{}{
		"user_id": playerMap["user_id"],
		"name":    playerMap["name"],
	})
}

func (h *handler) CheckInRoomHandler(c echo.Context) error {
	var (
		req       = c.Request()
		playerMap = h.getPlayerMap(req)
	)

	if !h.joinedRoom(playerMap) {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) NextHandler(c echo.Context) error {
	var (
		req       = c.Request()
		playerMap = h.getPlayerMap(req)
		rm        = room2.GetRoomManager()
	)

	if !h.joinedRoom(playerMap) {
		return c.NoContent(http.StatusBadRequest)
	}

	r := rm.GetRoom(room2.RoomID(playerMap["room_id"].(float64)))

	r.SendToClient(&message.Message{
		Name:    playerMap["name"].(string),
		Message: "NEXT",
		When:    time.Now(),
	})

	return c.NoContent(http.StatusOK)
}
