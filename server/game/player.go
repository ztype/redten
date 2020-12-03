package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type UserInfo struct {
	Id     string
	Name   string
	Avatar string
}

type Player struct {
	info UserInfo
	conn *websocket.Conn
	recv chan Msg
	send chan Msg
}

func NewPlayer(conn *websocket.Conn, info UserInfo) *Player {
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
		timer := time.NewTimer(time.Second)
		select {
		case p.recv <- msg:
			//
		case <-timer.C:
			break
		}
	}
}

func (p *Player) Send(m Msg) error {
	return p.conn.WriteJSON(m)
}
