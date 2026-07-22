package board

import (
	"math/rand"
)

type Board interface {
	GetHand() []Card //Devuelve la lista de la mano
	DrawCard()       //Roba una Carta
	Shuffle()        //Mescla el deck
	ReShuffle()      //Mescla el descarte en el deck
}

type BasicBoard struct {
	Deck    []Card
	Hand    []Card
	Discard []Card
	Banish  []Card
}

func (b *BasicBoard) GetHand() []Card {
	return b.Hand
}

func (b *BasicBoard) DrawCard() {
	if len(b.Deck) < 1 {
		b.ReShuffle()
	}
	if len(b.Deck) > 0 {
		b.Hand = append(b.Hand, b.Deck[0])
		b.Deck = b.Deck[1:]
	}
}

func (b *BasicBoard) Shuffle() {
	for i := len(b.Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		b.Deck[i], b.Deck[j] = b.Deck[j], b.Deck[i]
	}
}

func (b *BasicBoard) ReShuffle() {
	b.Deck = append(b.Deck, b.Discard...)
	b.Discard = []Card{}
	for i := len(b.Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		b.Deck[i], b.Deck[j] = b.Deck[j], b.Deck[i]
	}
}

func NewPlayerBoard() Board {
	return &BasicBoard{
		Deck: []Card{
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("defense"),
			NewNormalCard("defense"),
			NewNormalCard("defense"),
		},
	}
}

func NewIABoard() Board {
	return &BasicBoard{
		Deck: []Card{
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("shoot"),
			NewNormalCard("defense"),
		},
	}
}
