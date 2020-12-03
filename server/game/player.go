package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var writeWait = time.Second
var readWait = time.Second

type UserInfo struct {
	Id     string
	Name   string
	Avatar string
}

type Player struct {
	info *UserInfo
	conn *websocket.Conn
	recv chan Msg
	send chan Msg
}

func NewPlayer(conn *websocket.Conn, info *UserInfo) *Player {
	r := make(chan Msg, 1)
	s := make(chan Msg, 1)
	p := &Player{conn: conn, recv: r, send: s, info: info}
	go p.onMsg()
	return p
}

func (p *Player) OnMsg() <-chan Msg {
	return p.recv
}

func (p *Player) onMsg() {
	for {
		_, data, err := p.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg := Msg{}
		_ = json.Unmarshal(data, &msg)
		timer := time.NewTimer(readWait)
		select {
		case p.recv <- msg:
			//do nothing
			log.Println("recv:", string(data))
		case <-timer.C:
			break
		}
	}
}

func (p *Player) Send(m Msg) error {
	if err := p.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return p.conn.WriteJSON(m)
}

// Echo is for testing only
func (p *Player) Echo(c <-chan Msg) {
	for msg := range c {
		_ = p.Send(msg)
	}
}
