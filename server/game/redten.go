package game

import (
	"log"
	"net/http"
	"sync"

	ws "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
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
	r := &Redten{lock: sync.Mutex{}}
	r.players = make(map[string]*Player)
	r.rooms = make(map[string]*Room)
	return r
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

//user connect
//when a player connect:
//1.check if this is reconnect
func (s *Redten) NewPlayer(conn *ws.Conn, info *UserInfo) {
	//first check whether user is reconnect
	if p, ok := s.players[info.Id]; ok {
		p.Renew(conn, nil)
		if p.RoomId != "" {
			if r, ok := s.rooms[p.RoomId]; ok {
				msg := r.Join(p)
				_ = p.SendMsg(msg)
			}
		}
		return
	}
	//if user not exit,create it
	p := NewPlayer(conn, info)
	s.lock.Lock()
	s.players[p.Id()] = p
	s.lock.Unlock()

	p.OnClose("redten", s.OnPlayerClose)
	p.OnMsg("redten", s.OnMsg)
	log.Println("online", len(s.players))
	msg := NewMsg("notify").Set("msg", "ok").Set("id", p.Id())
	_ = p.Send(msg.JsonBytes())
}

//
func (s *Redten) OnPlayerClose(id string) {
	log.Println("redten", id, "offline")
}

func (s *Redten) OnMsg(p *Player, m *Msg) {
	log.Println("redten", m.Cmd, string(m.Data))
	switch m.Cmd {
	case "create_room":
		msg := s.CreateRoom(p)
		_ = p.SendMsg(msg)
	case "join_room":
		msg := s.JoinRoom(p, m)
		_ = p.SendMsg(msg)
	case "leave_room":
		log.Println("redten", "leave_room")
		if r, ok := s.rooms[p.RoomId]; ok {
			r.Leave(p)
		}
	case "ping":
		_ = p.Send([]byte("pong"))
	}
}

func (s *Redten) CreateRoom(p *Player) *Msg {
	if p.RoomId != "" {
		return NewMsg("notify").Set("msg", "already in a room").
			Set("room_id", p.RoomId)
	}
	r := NewRoom(p)
	s.rooms[r.ID()] = r
	return NewMsg("notify").Set("msg", "ok").Set("room_id", r.ID())
}

func (s *Redten) JoinRoom(p *Player, m *Msg) *Msg {
	if p.RoomId != "" {
		return NewMsg("notify").Set("msg", "already in a room").
			Set("room_id", p.RoomId)
	}
	id := gjson.Get(string(m.Data), "room_id").Str
	if room, ok := s.rooms[id]; ok {
		log.Println(p.Id(), "joined", id)
		return room.Join(p)
	}
	return NewMsg("notify").Set("msg", "room not found")
}
