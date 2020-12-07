package game

import (
	"encoding/json"
	"fmt"
	"log"
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
	closecallback map[string]func(string)
	msgcallback   map[string]func(*Player, *Msg)
	//
	Cards  []*Card
	smsg   chan []byte
	RoomId string
	online bool
}

func NewPlayer(conn *websocket.Conn, info *UserInfo) *Player {
	p := &Player{conn: conn, info: info}
	p.lock = sync.Mutex{}
	p.closecallback = make(map[string]func(string))
	p.msgcallback = make(map[string]func(*Player, *Msg))
	p.smsg = make(chan []byte, 100)
	cs := make([]*Card, 0)
	p.Renew(conn, cs)
	return p
}

//Renew update the player's connection when reconnect
func (p *Player) Renew(conn *websocket.Conn, cards []*Card) {
	p.conn = conn
	if len(cards) > 0 {
		p.Cards = cards
	}
	conn.SetCloseHandler(p.onClose)
	p.online = true
	// if the connection is broken,these 2 goroutine will exit
	go p.onMsg()
	go p.dosend()
}

func (p *Player) Id() string {
	return p.info.Id
}

func (p *Player) IsOnline() bool {
	return p.online
}

func (p *Player) LeaveRoom() {
	p.RoomId = ""
	p.Cards = make([]*Card, 0)

	delete(p.closecallback, "room")
	delete(p.msgcallback, "room")
}

func (p *Player) OnMsg(name string, f func(p *Player, m *Msg)) {
	p.msgcallback[name] = f
}

func (p *Player) OnClose(name string, f func(string)) {
	p.closecallback[name] = f
}

// websocket disconnect callback
func (p *Player) onClose(code int, text string) error {
	//log.Println("close:", p.conn.RemoteAddr().String(), "exit")
	if code == websocket.CloseGoingAway {
		return nil
	}
	p.online = false
	for _, callback := range p.closecallback {
		log.Println("player", "close", code, text)
		callback(p.Id())
	}
	return nil
}

func (p *Player) onMsg() {
	for {
		_, data, err := p.conn.ReadMessage()
		if err != nil {
			//if the connection is broken,the p.onClose() will be
			// called by websocket lib automatically
			return
		}
		msg := NewMsg("")
		_ = json.Unmarshal(data, msg)
		for _, callback := range p.msgcallback {
			callback(p, msg)
		}
	}
}

func (p *Player) SendMsg(m *Msg) error {
	return p.Send(m.JsonBytes())
}

func (p *Player) Send(bs []byte) error {
	timer := time.NewTimer(writeWait)
	select {
	case <-timer.C:
		//p.onClose()
		//if the connection is broken ,the p.onClose
		//will be called by websocket lib automatically
		return fmt.Errorf(p.Id(), "msg send timeout")
	case p.smsg <- bs:
		return nil
	}
}

func (p *Player) dosend() {
	for m := range p.smsg {
		if err := p.conn.WriteMessage(websocket.TextMessage, (m)); err != nil {
			//p.onClose()
			//if the connection is broken,the p.onClose
			// will be called by websocket library automatically
			//exit this loop when error happens,which means
			//the connection is broken
			return
		}
	}
}
