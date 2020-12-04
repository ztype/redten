package game

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	ws "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func NewUid() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}

type Redten struct {
	players map[string]*Player
	rooms   map[string]*Room
	lock    sync.Mutex
}

func NewRedten() *Redten {
	return &Redten{lock: sync.Mutex{}}
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
	//query db and get user info
	return &UserInfo{Id: id}, nil
}

func (s *Redten) NewPlayer(conn *ws.Conn, info *UserInfo) {
	p := NewPlayer(conn, info)
	s.lock.Lock()
	s.players[p.Id()] = p
	s.lock.Unlock()

	p.OnClose(s.OnPlayerClose)
	p.OnMsg(s.OnMsg)
	log.Println("online", len(s.players))

}

func (s *Redten) OnPlayerClose(id string) {
	for _, p := range s.players {
		if p.Id() == id {
			s.lock.Lock()
			delete(s.players, id)
			s.lock.Unlock()
		}
	}
}

func (s *Redten) OnMsg(p *Player, m *Msg) {
	switch m.Cmd {
	case "create_room":
		_ = p.SendMsg(s.CreateRoom(p))
	}
}

func (s *Redten) CreateRoom(p *Player) *Msg {
	r := NewRoom(p)
	data, _ := json.Marshal(map[string]string{"msg": "ok", "room_id": r.ID()})
	return &Msg{Cmd: "notify", Data: string(data)}
}
