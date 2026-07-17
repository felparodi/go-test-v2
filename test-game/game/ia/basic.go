package ia

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"math/rand"
)

func NewBasicIA(id int, s inter.Server, w inter.World) IA {
	idName := fmt.Sprintf("IA_BASIC_%d", id)
	ret := &IAData{
		id:        fmt.Sprintf(idName, id),
		character: item.NewCharacter(w.GetSize()),
		server:    s,
		world:     w,
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
			x, y := NormalizeMove(x, y, b.character.GetPosition(), b.world)
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
