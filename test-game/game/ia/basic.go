package ia

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"juego-websocket/game/position"
	"math/rand"
)

func NewBasicIA(id int, a inter.Area) IA {
	idName := fmt.Sprintf("IA_BASIC_%d", id)
	pos := position.GetRandPosistion(a.GetSize())
	ret := &IAData{
		id:        fmt.Sprintf(idName, id),
		character: item.NewCharacter(pos),
		area:      a,
		strategy:  basicStragey,
	}
	ret.character.SetControler(ret)
	return ret
}

func basicStragey(b *IAData) <-chan *Move {
	canal := make(chan *Move)
	go func() {
		moveTime := rand.Intn(150) * 10
		x, y := RandomDirection()
		actions := []*Action{}
		for t := 0; t < moveTime; t++ {
			x, y := NormalizeMove(x, y, b.character.GetPosition(), b.area)
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
