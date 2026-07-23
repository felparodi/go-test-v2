package ia

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"juego-websocket/game/size"
	"math/rand"
)

func NewDummyIA(id int, a inter.Area) IA {
	idName := fmt.Sprintf("IA_BASIC_%d", id)
	pos := a.GetSize().GetRandPosistion()
	c := item.NewCharacter(pos)
	ia := newBasicIA(idName, c, a, dummyStragey)
	ret := &ia
	ret.character.SetControler(ret)
	return ret
}

// @TDOO se pega a los bordes
func dummyStragey(ia IA) <-chan *Move {
	canal := make(chan *Move)
	go func() {
		moveTime := rand.Intn(15) * 10
		x, y := RandomDirection()
		actions := []*Action{}
		for t := 0; t < moveTime; t++ {
			position := ia.GetCharacter().GetPosition()
			area := ia.GetArea()
			x, y, _ := size.NormalizeMove(x*5, y*5, position, area.GetSize())
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
