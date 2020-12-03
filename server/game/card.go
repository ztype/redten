package game

import (
	"fmt"
	"math/rand"
	"time"
)

var cardvalue = map[string]int{
	"A": 12,
	"2": 13,
	// "3":  1,
	// "4":  2,
	// "5":  3,
	"6":  4,
	"7":  5,
	"8":  6,
	"9":  7,
	"10": 8,
	"J":  9,
	"Q":  10,
	"K":  11,
}

var cardshape = map[string]int{
	"S": 1,
	"H": 2,
	"C": 3,
	"D": 4,
}

var Kings = map[string]int{
	"kx": 14, //小王
	"kd": 15, //大王
}

type Card struct {
	Shape string
	Value string
}

func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.Shape, c.Value)
}

// 40 cards for redten
func NewCardSet() []*Card {
	var set []*Card
	for v := range cardvalue {
		for s := range cardshape {
			c := Card{Shape: s, Value: v}
			set = append(set, &c)
		}
	}
	return set
}

func shuffleCards(cards []*Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}
