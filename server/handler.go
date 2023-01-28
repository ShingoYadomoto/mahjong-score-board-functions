package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ShingoYadomoto/mahjong-score-board/data"
	"github.com/ShingoYadomoto/mahjong-score-board/room"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type handler struct{}

func (h *handler) getPlayer(c echo.Context) (*room.Player, error) {
	cookie, err := c.Cookie("playerID")
	if err != nil {
		return nil, err
	}

	idInt, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return nil, err
	}

	return data.GetPlayer(room.PlayerID(idInt))
}

func (h *handler) setCookie(c echo.Context, p *room.Player) {
	c.SetCookie(&http.Cookie{
		Name:     "playerID",
		Value:    fmt.Sprint(p.ID),
		Path:     "/",
		SameSite: http.SameSiteNoneMode, // only dev
		Expires:  time.Now().Add(time.Hour * 72),
		Secure:   true,
		HttpOnly: true,
	})
}

func (h *handler) CreatePlayerHandler(c echo.Context) error {
	var rb struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&rb); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	p, err := data.CreatePlayer(rb.Name)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	h.setCookie(c, p)

	return c.JSON(http.StatusOK, map[string]string{"playerID": fmt.Sprint(p.ID)})
}

func (h *handler) GetPlayerHandler(c echo.Context) error {
	p, err := h.getPlayer(c)
	if err != nil {
		if err == data.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	h.setCookie(c, p)

	return c.JSON(http.StatusOK, map[string]string{"playerID": fmt.Sprint(p.ID)})
}

func (h *handler) CreateRoomHandler(c echo.Context) error {
	p, err := h.getPlayer(c)
	if err != nil {
		if err == data.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	h.setCookie(c, p)

	r, err := data.CreateRoom(p.ID)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{"roomID": fmt.Sprint(r.ID)})
}

func (h *handler) JoinRoomHandler(c echo.Context) error {
	roomIDStr := c.Param("roomID")
	roomIDInt, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	p, err := h.getPlayer(c)
	if err != nil {
		if err == data.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	h.setCookie(c, p)

	err = data.AddPlayerToRoom(room.ID(roomIDInt), p.ID)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) LeaveRoomHandler(c echo.Context) error {
	roomIDStr := c.Param("roomID")
	roomIDInt, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	p, err := h.getPlayer(c)
	if err != nil {
		if err == data.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	h.setCookie(c, p)

	err = data.DeletePlayerFromRoom(room.ID(roomIDInt), p.ID)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) GetRoomHandler(c echo.Context) error {
	p, err := h.getPlayer(c)
	if err != nil {
		if err == data.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	h.setCookie(c, p)

	r, err := data.GetJoinedRoom(p.ID)
	if err != nil {
		if err == data.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{"roomID": fmt.Sprint(r.ID)})
}
