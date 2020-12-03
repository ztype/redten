package game

import (
	"net/http"

	ws "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func NewUid() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}

type Redten struct {
	players []*Player
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
	info, err := s.getUserInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.NewPlayer(conn, info)
}

func (s *Redten) getUserInfo(r *http.Request) (*UserInfo, error) {
	ck, err := r.Cookie("id")
	if err == nil {
		return s.getUinfo(ck.Value)
	}
	// create a new userinfo
	u := &UserInfo{
		Id:   NewUid(),
		Name: "jack",
	}
	return u, nil
}

func (s *Redten) getUinfo(id string) (*UserInfo, error) {
	return nil, nil
}

func (s *Redten) NewPlayer(conn *ws.Conn, info *UserInfo) {
	p := NewPlayer(conn, info)
	s.players = append(s.players, p)
	//testing
	go p.Echo(p.OnMsg())
}
