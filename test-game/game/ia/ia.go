package ia

import (
	"juego-websocket/game/inter"
	"time"
)

type IA interface {
	inter.CharacterControler
}

type IAData struct {
	id        string
	character inter.Character
	world     inter.World
	server    inter.Server
	live      bool
	strategy  func(*IAData) <-chan *Move
}

/*
*	Que empiece a escucar al usuario
**/
func (b *IAData) Start() error {
	b.live = true
	for b.live {
		strategy := b.strategy(b)
		for {
			time.Sleep(6000)
			move := <-strategy
			if move != nil {
				b.character.Move(move.X*100, move.Y*100)
				if move.Actions != nil {
					for _, a := range move.Actions {
						b.character.AddAction(a)
					}
				}
			} else {
				break
			}
		}

	}
	return nil
}

func (b *IAData) End() error {
	b.live = false
	return nil
}

func (b *IAData) GetId() string {
	return b.id
}

func (b *IAData) GetCharacter() inter.Character {
	return b.character
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
