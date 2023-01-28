package player

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/objx"
)

type Player struct {
	ID   string
	Name string `json:"name"`
}

func CreateHandler(c echo.Context) error {
	req := c.Request()

	_, err := req.Cookie("auth")
	if err != nil && err != http.ErrNoCookie {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	var p Player
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
