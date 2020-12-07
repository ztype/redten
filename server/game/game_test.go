package game

import (
	"encoding/json"
	"log"
	"testing"
)

func Test_card(t *testing.T) {
	cs := NewCardSet()
	log.Println("before")
	log.Println(cs)

	ShuffleCards(cs)
	log.Println("after shuffle")
	log.Println(cs)

	cc := make([]*Card, len(cs))
	copy(cc, cs)

	SortCards(cs)
	log.Println("after sort")
	log.Println(cs)

	sets := CutCards(cc, 2)
	log.Println("cut&sort")
	for i, set := range sets {
		SortCards(set)
		log.Println(i, set)
	}
}

func Test_roomId(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := NewRoomId()
		log.Println(id)
	}
}

func Test_msg(t *testing.T) {
	str := `{"cmd":"join_room","data":{"room_id":"1F7B"}}`
	var m Msg
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(m)
}
