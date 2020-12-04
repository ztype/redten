package game

type Room struct {
	players []*Player
}

func NewRoom(owner *Player) *Room {
	r := &Room{}
	r.players = append(r.players, owner)
	return r
}
