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
	p := &room.Player{
		ID:     playerMap["user_id"].(string),
		Name:   playerMap["name"].(string),
		RoomID: room.RoomID(playerMap["room_id"].(float64)),
	}

	rm := room.GetRoomManager()

	r := rm.GetRoom(p.RoomID)
	if r == nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

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

	var p room.Player
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	m := md5.New()
	if _, err := io.WriteString(m, strings.ToLower(p.Name)); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	p.ID = fmt.Sprintf("%x", m.Sum(nil))

	authCookieValue := objx.New(map[string]interface{}{
		"user_id": p.ID,
		"name":    p.Name,
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

func (h *handler) CreateRoomHandler(c echo.Context) error {
	req := c.Request()

	authCookie, err := req.Cookie("auth")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	playerMap := objx.MustFromBase64(authCookie.Value)
	for _, field := range []string{"user_id", "name"} {
		if _, ok := playerMap[field]; !ok {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	rm := room.GetRoomManager()

	//r, err := rm.NewRoom()
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, err.Error())
	//}
	r := rm.GetRoom(3)

	authCookieValue := objx.New(map[string]interface{}{
		"user_id": playerMap["user_id"],
		"name":    playerMap["name"],
		"room_id": r.ID,
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

func (h *handler) JoinRoomHandler(c echo.Context) error {
	req := c.Request()

	roomIDStr := c.Param("roomID")
	roomIDInt, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	playerMap := objx.MustFromBase64(authCookie.Value)

	rm := room.GetRoomManager()

	r := rm.GetRoom(room.RoomID(roomIDInt))
	if r == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	playerMap["room_id"] = r.ID
	authCookieValue := objx.New(playerMap).MustBase64()

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

func (h *handler) LeaveRoomHandler(c echo.Context) error {
	req := c.Request()

	authCookie, err := req.Cookie("auth")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	playerMap := objx.MustFromBase64(authCookie.Value)

	delete(playerMap, "room_id")

	authCookieValue := objx.New(playerMap).MustBase64()

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

func (h *handler) CheckInRoomHandler(c echo.Context) error {
	req := c.Request()

	authCookie, err := req.Cookie("auth")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	playerMap := objx.MustFromBase64(authCookie.Value)
	for _, field := range []string{"user_id", "name", "room_id"} {
		if _, ok := playerMap[field]; !ok {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	return c.NoContent(http.StatusOK)
}
