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
	info          *UserInfo
	conn          *websocket.Conn
	recv          chan Msg
	closecallback []func(string)
}

func NewPlayer(conn *websocket.Conn, info *UserInfo) *Player {
	r := make(chan Msg, 1)
	p := &Player{conn: conn, recv: r, info: info}
	conn.SetCloseHandler(p.onClose)
	go p.onMsg()
	return p
}

func (p *Player) Id() string {
	return p.info.Id
}

func (p *Player) OnMsg() <-chan Msg {
	return p.recv
}

func (p *Player) OnClose(f func(string)) {
	p.closecallback = append(p.closecallback, f)
}

func (p *Player) onClose(code int, text string) error {
	close(p.recv)
	log.Println("close:", p.conn.RemoteAddr().String(), "exit")
	for _, callback := range p.closecallback {
		callback(p.Id())
	}
	return nil
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
	log.Println("send:", m)
	return p.conn.WriteJSON(m)
}

// Echo is for testing only
func (p *Player) Echo(c <-chan Msg) {
	for msg := range c {
		_ = p.Send(msg)
	}
}
