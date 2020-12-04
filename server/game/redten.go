package game

import (
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
	players []*Player
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
	s.players = append(s.players, p)
	log.Println("online", len(s.players))
	p.OnClose(s.OnPlayerClose)
	s.lock.Unlock()
	//testing
	go p.Echo(p.OnMsg())
}

func (s *Redten) OnPlayerClose(id string) {
	for i, p := range s.players {
		if p.Id() == id {
			s.lock.Lock()
			s.players = append(s.players[:i], s.players[i+1:]...)
			s.lock.Unlock()
		}
	}
}
