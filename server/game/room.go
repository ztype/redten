package game

import "sync"

const (
	RoomInvalid = iota
	RoomPlaying
	RoomFinished
)

const defaultBase = 10
const defaultRate = 10

type Room struct {
	players []*Player
	lock    sync.Mutex
	//
	base int
	rate int
}

func NewRoom(owner *Player) *Room {
	r := &Room{}
	r.base = defaultBase
	r.rate = defaultRate
	r.players = append(r.players, owner)
	r.lock = sync.Mutex{}
	return r
}

func (r *Room) Join(p *Player) *Room {
	r.lock.Lock()
	r.players = append(r.players, p)
	r.lock.Unlock()
	return r
}

func (r *Room) Has(id string) bool {
	for _, p := range r.players {
		if p.Id() == id {
			return true
		}
	}
	return false
}
