package game

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"sync"
)

const (
	RoomInvalid = iota
	RoomPlaying
	RoomFinished
)

const defaultBase = 10
const defaultRate = 10

type Room struct {
	id      string
	players []*Player
	lock    sync.Mutex
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
	r.players = append(r.players, owner)
	r.lock = sync.Mutex{}
	return r
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) Join(p *Player) error {
	r.lock.Lock()
	r.players = append(r.players, p)
	r.lock.Unlock()
	return nil
}

func (r *Room) Has(id string) bool {
	for _, p := range r.players {
		if p.Id() == id {
			return true
		}
	}
	return false
}

func (r *Room) OnMsg(p *Player, m *Msg) {

}
