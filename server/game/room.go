package game

import (
	"encoding/hex"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	RoomInvalid = iota
	RoomPlaying
	RoomFinished
)
const playersNum = 4

const defaultBase = 10
const defaultRate = 10

type Room struct {
	id      string
	players map[string]*Player
	lock    sync.Mutex
	owner   *Player
	//
	base int
	rate int
}

//4 digit hex string
func NewRoomId() string {
	//rand.Seed(time.Now().UnixNano())
	id := ""
	for i := 0; i < 4; i++ {
		n := rand.Intn(16)
		id += hex.EncodeToString([]byte{byte(n)})[1:]
	}
	return strings.ToUpper(id)
}

func NewRoom(owner *Player) *Room {
	r := &Room{}
	r.id = NewRoomId()
	r.base = defaultBase
	r.rate = defaultRate
	r.lock = sync.Mutex{}
	r.players = make(map[string]*Player)
	r.owner = owner
	r.Join(owner)
	return r
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) Join(p *Player) *Msg {
	if len(r.players) == 4 {
		return NewMsg("notify").Set("msg", "room is full")
	}
	r.lock.Lock()
	p.RoomId = r.ID()
	r.players[p.Id()] = p
	r.lock.Unlock()
	p.OnMsg("room", r.OnMsg)
	//if player is sufficient,start the game
	if len(r.players) == playersNum {
		time.AfterFunc(time.Second, r.DealCards)
	}
	return NewMsg("notify").Set("msg", "ok").Set("room_id", r.ID())
}

func (r *Room) Leave(p *Player) {
	r.lock.Lock()
	p.LeaveRoom()
	delete(r.players, p.Id())
	r.lock.Unlock()
}

//
func (r *Room) Dismiss() {
	for id, p := range r.players {
		p.LeaveRoom()
		delete(r.players, id)
	}
}

func (r *Room) OnMsg(p *Player, m *Msg) {
	log.Println("room", r.ID(), m.Cmd, string(m.Data))
	switch m.Cmd {
	case "notify":
		r.Broadcast(p, m)
	case "play_card":
	case "choose_show":
	case "choose_to_play":
	}
}

func (r *Room) OnPlayerClose(id string) {

}

func (r *Room) Broadcast(p *Player, m *Msg) {
	var offline []string
	for _, pp := range r.players {
		if !pp.IsOnline() {
			offline = append(offline, pp.Id())
		}
	}
	for _, pp := range r.players {
		if p != nil && pp.Id() == p.Id() {
			continue
		}
		_ = pp.SendMsg(m.Set("offlines", offline))
	}
}

// === game area ===

//deal card to each player
func (r *Room) DealCards() {
	cards := NewCardSet()
	ShuffleCards(cards)
	cardsets := CutCards(cards, len(r.players))
	i := 0
	for _, p := range r.players {
		set := (cardsets[i])
		SortCards(set)
		msg := NewMsg("deal_cards").Set("cards", set)
		_ = p.SendMsg(msg)
		i++
	}
}
