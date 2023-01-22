package chat

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ShingoYadomoto/mahjong-score-board/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

var avatar Avatar = TryAvatars{
	UseFileServerAvatar,
	UseAuthAvatar,
	UseGravatarAvatar,
}

func Serve() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	// Gominiauthのセットアップ
	authCallBackEndpoint := "http://localhost:8080/auth/callback/"
	gomniauth.SetSecurityKey(secretKey)
	gomniauth.WithProviders(
		//facebook.New(clientID, , authCallBackEndpoint+"facebook"),
		github.New(githubClientID, githubSecret, authCallBackEndpoint+"github"),
		google.New(googleClientID, googleSecret, authCallBackEndpoint+"google"),
	)

	//r := newRoom(UseAuthAvatar)
	//r := newRoom(UseGravatarAvatar)
	//r := newRoom(UseFileServerAvatar)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/auth/", LoginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	go r.run() // start chat room

	// start web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		path := filepath.Join("templates", t.filename)
		t.templ = template.Must(template.ParseFiles(path))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}
