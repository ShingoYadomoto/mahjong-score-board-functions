package server

import (
	"flag"
	"net/http"
	"os"

	"github.com/ShingoYadomoto/mahjong-score-board/deprecated/room"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve() {
	addr := flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()

	var (
		e  = echo.New()
		h  = handler{}
		rm = room.GetRoomManager()
	)

	r, err := rm.NewRoom()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("ALLOW_ORIGIN")},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	r.SetTracer(os.Stdout)

	e.POST("/player", h.CreatePlayerHandler)
	e.POST("/room", h.CreateRoomHandler)
	e.POST("/room/:roomID", h.JoinRoomHandler)
	e.GET("/room/in", h.CheckInRoomHandler)
	e.POST("/room/leave", h.LeaveRoomHandler)
	e.GET("/room", h.RoomSocketHandler)

	e.POST("/room/next", h.NextHandler)

	go r.Run() // start room

	e.Logger.Fatal(e.Start(*addr))
}
