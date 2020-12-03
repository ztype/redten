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
