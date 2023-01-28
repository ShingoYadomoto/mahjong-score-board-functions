package server

import (
	"flag"
	"net/http"
	"os"

	"github.com/ShingoYadomoto/mahjong-score-board/player"
	"github.com/ShingoYadomoto/mahjong-score-board/room"
	"github.com/labstack/echo/v4"
)

var commonHeader = map[string]string{
	"Content-Type":                     "application/json",
	"Access-Control-Allow-Origin":      os.Getenv("ALLOW_ORIGIN"),
	"Access-Control-Allow-Headers":     "Content-Type",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Methods":     "GET,POST,OPTIONS",
}

func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			for k, v := range commonHeader {
				res.Header().Set(k, v)
			}

			if c.Request().Method == http.MethodOptions {
				res.WriteHeader(http.StatusOK)
				return nil
			}
			return next(c)
		}
	}
}

func Serve() {
	addr := flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()

	e := echo.New()

	e.Use(CORSMiddleware())

	r := room.NewRoom()
	r.SetTracer(os.Stdout)
	e.GET("/room", r.SyncRoomHandler)
	e.POST("/player", player.CreateHandler)

	go r.Run() // start room

	e.Logger.Fatal(e.Start(*addr))
}
