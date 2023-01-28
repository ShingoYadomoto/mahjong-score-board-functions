package server

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ShingoYadomoto/mahjong-score-board/player"
	"github.com/ShingoYadomoto/mahjong-score-board/room"
)

var commonHeader = map[string]string{
	"Content-Type":                     "application/json",
	"Access-Control-Allow-Origin":      os.Getenv("ALLOW_ORIGIN"),
	"Access-Control-Allow-Headers":     "Content-Type",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Methods":     "GET,OPTIONS",
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range commonHeader {
			w.Header().Set(k, v)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func Serve() {
	addr := flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()

	r := room.NewRoom()
	r.SetTracer(os.Stdout)
	http.HandleFunc("/room", CORSMiddleware(r.SyncRoomHandler))
	http.HandleFunc("/player", CORSMiddleware(player.CreateHandler))

	go r.Run() // start room room

	// start web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
