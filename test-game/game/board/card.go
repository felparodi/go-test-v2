package board

import (
	"fmt"
	"juego-websocket/game/inter"
	"math/rand"
)

// @Temporal
type Card interface {
	GetEvents() []inter.Event
}

type Modify struct {
}

type NormalCard struct {
	Id       string //Codigi
	CodeCard string
}

func NewNormalCard(code string) Card {
	return &NormalCard{
		Id:       fmt.Sprintf("card_%s_%d", code, rand.Int63()),
		CodeCard: code,
	}
}

func (n *NormalCard) GetEvents() []inter.Event {
	return []inter.Event{}
}

// @TODO Ver A mucho Futuro
type ModifyCard struct {
	NormalCard
	Modifys []Modify
}
