package ia

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"juego-websocket/game/position"
	"juego-websocket/game/size"
	"math/rand"
)

func NewDummyIA(id int, a inter.Area) IA {
	idName := fmt.Sprintf("IA_BASIC_%d", id)
	pos := position.GetRandPosistion(a.GetSize())
	c := item.NewCharacter(pos)
	ia := newBasicIA(idName, c, a, dummyStragey)
	ret := &ia
	ret.character.SetControler(ret)
	return ret
}

func dummyStragey(ia IA) <-chan *Move {
	canal := make(chan *Move)
	go func() {
		moveTime := rand.Intn(150) * 10
		x, y := RandomDirection()
		actions := []*Action{}
		for t := 0; t < moveTime; t++ {
			position := ia.GetCharacter().GetPosition()
			area := ia.GetArea()
			x, y := size.NormalizeMove(x, y, position, area.GetSize())
			if rand.Intn(10) > 7 {
				actions = append(actions, &Action{name: "shoot"})
			}
			canal <- &Move{
				X:       x,
				Y:       y,
				Actions: actions,
			}
		}
		canal <- nil
	}()
	return canal
}
