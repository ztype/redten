package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/ztype/redten/server/game"
)

var fAddr = flag.String("addr", ":1010", "listen addr")
var fDir = flag.String("dir", "./static/html", "html file dir")

func home(dir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "" {
			if c, err := r.Cookie("id"); err != nil {
				uid := game.NewUid()
				http.SetCookie(w, &http.Cookie{Name: "id", Value: uid})
				log.Println("new id:", uid)
			} else {
				log.Println("get id:", c.Value)
			}
			http.ServeFile(w, r, filepath.Join(dir, "./index.html"))
			return
		}
		http.Error(w, "404 not found", http.StatusNotFound)
	}
}

func main() {
	app := game.NewRedten()

	r := mux.NewRouter()

	r.HandleFunc("/", home(*fDir))
	r.HandleFunc("/ws", app.ServeWs)
	log.Println("service start at", *fAddr)
	err := http.ListenAndServe(*fAddr, r)
	log.Println(err)
}
