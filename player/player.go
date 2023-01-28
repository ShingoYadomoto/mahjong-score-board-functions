package player

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stretchr/objx"
)

type Player struct {
	ID   string
	Name string `json:"name"`
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p Player
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := md5.New()
	if _, err := io.WriteString(m, strings.ToLower(p.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = fmt.Sprintf("%x", m.Sum(nil))

	authCookieValue := objx.New(map[string]interface{}{
		"userid": p.ID,
		"name":   p.Name,
	}).MustBase64()

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}
