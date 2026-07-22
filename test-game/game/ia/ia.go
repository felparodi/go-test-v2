package ia

import (
	"juego-websocket/game/inter"
	"time"
)

type IA interface {
	inter.CharacterControler
	GetCharacter() inter.Character
	GetArea() inter.Area
}

type BasicIA struct {
	id            string
	character     inter.Character
	area          inter.Area
	server        inter.Server
	live          bool
	activeChannel chan bool
	strategy      IAStrategy
}

type IAStrategy func(IA) <-chan *Move

func newBasicIA(id string, c inter.Character, a inter.Area, str IAStrategy) BasicIA {
	return BasicIA{
		id:            id,
		character:     c,
		area:          a,
		strategy:      str,
		activeChannel: make(chan bool),
	}
}

func (ia *BasicIA) Start() (chan bool, error) {
	ia.live = true
	go ia.liveing()
	return ia.activeChannel, nil
}

func (ia *BasicIA) liveing() {
	for ia.live {
		strategy := ia.strategy(ia)
		for {
			time.Sleep(6000)
			move := <-strategy
			if move != nil {
				ia.character.Move(move.X*100, move.Y*100)
				if move.Actions != nil {
					for _, a := range move.Actions {
						ia.character.AddAction(a)
					}
				}
			} else {
				break
			}
		}
	}
}

// @TODO
func (ia *BasicIA) Stop() error {
	ia.live = false
	return nil
}

func (ia *BasicIA) GetId() string {
	return ia.id
}

func (ia *BasicIA) GetCharacter() inter.Character {
	return ia.character
}

func (ia *BasicIA) GetArea() inter.Area {
	return ia.area
}

type Action struct {
	name string
}

func (a *Action) GetName() string {
	return a.name
}

func (a *Action) GetData() interface{} {
	return nil
}

type Move struct {
	X       float64
	Y       float64
	Actions []*Action
}
