package server

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ShingoYadomoto/mahjong-score-board/room"
)

func Serve() {
	addr := flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()

	r := room.NewRoom()
	r.SetTracer(os.Stdout)
	http.Handle("/room", r)

	go r.Run() // start room room

	// start web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
