package game

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var writeWait = time.Second

type UserInfo struct {
	Id     string
	Uid    string
	Name   string
	Avatar string
}

type Player struct {
	info          *UserInfo
	lock          sync.Mutex
	conn          *websocket.Conn
	closecallback []func(string)
	msgcallback   []func(*Player, *Msg)
	//
	Cards  []*Card
	RoomId string
}

func NewPlayer(conn *websocket.Conn, info *UserInfo) *Player {
	p := &Player{conn: conn, info: info}
	p.Cards = make([]*Card, 0)
	p.lock = sync.Mutex{}
	conn.SetCloseHandler(p.onClose)
	go p.onMsg()
	return p
}

func (p *Player) Renew(conn *websocket.Conn, cards []*Card) {
	p.conn = conn
	p.Cards = cards
}

func (p *Player) Id() string {
	return p.info.Id
}

func (p *Player) OnMsg(f func(p *Player, m *Msg)) {
	p.msgcallback = append(p.msgcallback, f)
}

func (p *Player) OnClose(f func(string)) {
	p.closecallback = append(p.closecallback, f)
}

func (p *Player) onClose(code int, text string) error {
	//log.Println("close:", p.conn.RemoteAddr().String(), "exit")
	for _, callback := range p.closecallback {
		callback(p.Id())
	}
	return nil
}

func (p *Player) onMsg() {
	for {
		_, data, err := p.conn.ReadMessage()
		if err != nil {
			return
		}
		msg := Msg{}
		_ = json.Unmarshal(data, &msg)
		for _, callback := range p.msgcallback {
			callback(p, &msg)
		}
	}
}

func (p *Player) SendMsg(m *Msg) error {
	if err := p.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	//log.Println("send:", m)
	return p.conn.WriteJSON(m)
}
