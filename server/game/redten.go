package game

import (
	"net/http"

	ws "github.com/gorilla/websocket"
)

type Redten struct {
}

func NewRedten() *Redten {
	return &Redten{}
}

func (s *Redten) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrade(w, r, nil, 4096, 4096)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//data, _ := json.MarshalIndent(r.Header, "", " ")
	//log.Println("ws header:", string(data))
	c := s.NewPlayer(conn)
	http.SetCookie(w, &http.Cookie{Name: "id", Value: c})
}

func (s *Redten) NewPlayer(conn *ws.Conn) string {
	conn.NextReader()
	return "id=abc123"
}
