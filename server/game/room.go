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
const playersNum = 2

const defaultBase = 10
const defaultRate = 10

type Room struct {
	id      string
	players []*Player
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
	r.players = make([]*Player, 0)
	r.owner = owner
	r.Join(owner)
	return r
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) Join(p *Player) *Msg {
	//1.check whether room is full already
	if len(r.players) == playersNum {
		return NewMsg("notify").Set("msg", "room is full")
	}
	//2.join room
	r.lock.Lock()
	p.RoomId = r.ID()
	r.players = append(r.players, p)
	r.lock.Unlock()
	p.OnMsg("room", r.OnMsg)
	//3.start the game,if player is sufficient
	if len(r.players) == playersNum {
		r.Broadcast(nil, NewMsg("game_start"))
		time.AfterFunc(time.Second, r.DealCards)
	}
	//4.when a new player joins room,notify other players
	idx := len(r.players) - 1
	r.Broadcast(p, NewMsg("join_room").
		Set("id", p.Id()).Set("seat_index", idx))
	return NewMsg("notify").Set("msg", "ok").
		Set("room_id", r.ID()).
		Set("seat_index", idx)
}

func (r *Room) Leave(p *Player) {
	r.lock.Lock()
	p.LeaveRoom()
	r.players = []*Player{}
	r.lock.Unlock()
	r.Broadcast(p, NewMsg("leave_room").Set("id", p.Id()))
}

//
func (r *Room) Dismiss() {
	for _, p := range r.players {
		p.LeaveRoom()
	}
	r.players = []*Player{}
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
	if len(offline) > 0 {
		m.Set("offlines", offline)
	}
	for _, pp := range r.players {
		if p != nil && pp.Id() == p.Id() {
			continue
		}
		_ = pp.SendMsg(m)
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
