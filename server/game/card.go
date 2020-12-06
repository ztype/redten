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
	"H": 1, //Heart,red ♥
	"D": 2, //Diamond,red ♦
	"S": 3, //Spade,black ♠
	"C": 4, //Club,black ♣
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

func ShuffleCards(cards []*Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func CutCards(cs []*Card, n int) [][]*Card {
	set := make([][]*Card, 0)
	x := len(cs) / n
	for i := 0; i < len(cs); i += x {
		sub := make([]*Card, x)
		copy(sub, cs[i:i+x])
		set = append(set, sub)
	}
	return set
}

func SortCards(cards []*Card) {
	sort.Slice(cards, func(i, j int) bool {
		b := cards[i].CompareTo(cards[j])
		if b == 0 {
			return cardshape[cards[i].Shape] < cardshape[cards[j].Shape]
		}
		return b > 0
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

//c<b:1, c==b:0, c>b:-1
func (c *Card) CompareTo(b *Card) int {
	if b.String() == "H10" {
		return 1
	}
	if c.String() == "H10" {
		return -1
	}
	//c or b != H10
	if b.String() == "D10" {
		return 1
	}
	if c.String() == "D10" {
		return -1
	}
	if cardvalue[c.Value] == cardvalue[b.Value] {
		return 0
	}
	if cardvalue[c.Value] < cardvalue[b.Value] {
		return 1
	}
	return -1
}

func IsChain(cards []*Card) bool {
	if len(cards) < 3 {
		return false
	}
	SortCards(cards)
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
