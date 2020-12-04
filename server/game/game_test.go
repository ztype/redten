package game

import (
	"log"
	"testing"
)

func Test_card(t *testing.T) {
	cs := NewCardSet()
	log.Println("before")
	log.Println(cs)

	shuffleCards(cs)
	log.Println("after")
	log.Println(cs)
}

func Test_roomId(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := NewRoomId()
		log.Println(id)
	}
}
