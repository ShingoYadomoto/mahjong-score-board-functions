package server

import (
	"flag"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve() {
	addr := flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()

	var (
		e = echo.New()
		h = handler{}
	)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("ALLOW_ORIGIN")},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	e.POST("/players", h.CreatePlayerHandler)
	e.GET("/player", h.GetPlayerHandler)
	e.POST("/rooms", h.CreateRoomHandler)
	e.POST("/rooms/:roomID/in", h.JoinRoomHandler)
	e.POST("/rooms/:roomID/out", h.LeaveRoomHandler)
	e.GET("/room", h.GetRoomHandler)

	e.Logger.Fatal(e.Start(*addr))
}
