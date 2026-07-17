package ia

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"math"
	"math/rand"
	"time"
)

type IA interface {
	inter.CharacterControler
}

type Basic struct {
	id        string
	character inter.Character
	world     inter.World
	server    inter.Server
	live      bool
}

func NewBasicIA(id int, s inter.Server, w inter.World) IA {
	idName := fmt.Sprintf("IA_Basic_%d", id)
	ret := &Basic{
		id:        fmt.Sprintf(idName, id),
		character: item.NewCharacter(idName, w.GetSize()),
		server:    s,
		world:     w,
	}
	ret.character.SetControler(ret)
	return ret
}

/*
*	Que empiece a escucar al usuario
**/
func (b *Basic) Start() error {
	b.live = true
	moveTime := 100
	for b.live {
		x, y := RandomDirection()

		for t := 0; t < moveTime; t++ {
			pos := b.character.GetPosition()
			size := b.world.GetSize()
			if pos.GetX()+x < 0 {
				x = -x
			}
			if pos.GetY()+y < 0 {
				y = -y
			}
			if pos.GetX()+x > size.GetWidth() {
				x = -x
			}
			if pos.GetY()+y > size.GetHeight() {
				y = -y
			}
			b.character.Move(x*100, y*100)
			time.Sleep(6000)
			if rand.Intn(10) > 7 {
				b.character.AddAction(&Action{name: "shoot"})
			}
		}
		moveTime = rand.Intn(150) * 10
	}
	return nil
}

func (b *Basic) End() error {
	b.live = false
	return nil
}

func (b *Basic) GetId() string {
	return b.id
}

func (b *Basic) GetCharacter() inter.Character {
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

func RandomDirection() (float64, float64) {
	angle := float64(rand.Intn(8) * 45)
	x := math.Cos(angle)
	y := math.Sin(angle)
	return x, y
}
