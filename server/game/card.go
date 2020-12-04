package game

import (
	"fmt"
	"math/rand"
	"sort"
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
	"S": 1, //Spade,black ♠
	"H": 2, //Heart,red ♥
	"C": 3, //Club,black ♣
	"D": 4, //Diamond,red ♦
}

var Kings = map[string]int{
	"kx": 14, //小王
	"kd": 15, //大王
}

//single,double,sisterdouble(2 or more double),chain(3 or more)
//boom(three)
//doubleten(2 red ten)

type Card struct {
	Shape string
	Value string
	face  string
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

func sortCards(cards []*Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cardvalue[cards[i].Value] < cardvalue[cards[j].Value]
	})
}

//s :fmt.Sprintf("%s%s",c.Shape,c.Value)
func NewCard(s string) *Card {
	if len(s) < 2 {
		return nil
	}
	return &Card{Shape: string(s[0]), Value: s[0:]}
}

func (c *Card) String() string {
	if c.face == "" {
		c.face = fmt.Sprintf("%s%s", c.Shape, c.Value)
	}
	return c.face
}

//
func (c *Card) SmallerThan(b *Card) bool {
	if b.String() == "H10" {
		return false
	}
	if c.String() == "H10" {
		return true
	}
	if b.String() == "D10" {
		return false
	}
	return cardvalue[c.Value] < cardvalue[b.Value]
}

func IsChain(cards []*Card) bool {
	if len(cards) < 3 {
		return false
	}
	sortCards(cards)
	last := cards[0]
	for i := 0; i < len(cards); i++ {
		if cards[i].String() == "H10" || cards[i].String() == "D10" {
			return false
		}
		if cards[i].Value == "2" {
			return false
		}
		if i > 0 && cardvalue[cards[i].Value] != cardvalue[last.Value]+1 {
			return false
		}
		last = cards[i]
	}
	return true
}
