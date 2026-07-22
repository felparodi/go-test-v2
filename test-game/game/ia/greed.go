package ia

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"juego-websocket/game/position"
	"juego-websocket/game/size"
	"math"
	"math/rand"
)

func NewGreedIA(id int, a inter.Area) IA {
	idName := fmt.Sprintf("IA_GREED_%d", id)
	pos := position.GetRandPosistion(a.GetSize())
	c := item.NewCharacter(pos)
	ia := newBasicIA(idName, c, a, greadyStrategy)
	ret := &ia
	ret.character.SetControler(ret)
	return ret
}

func greadyStrategy(ia IA) <-chan *Move {
	canal := make(chan *Move)
	position := ia.GetCharacter().GetPosition()
	area := ia.GetArea()
	coin := getClosedCoin(position, nil, area)
	go func() {
		moveTime := rand.Intn(150) * 10
		position := ia.GetCharacter().GetPosition()
		area := ia.GetArea()
		coin := getClosedCoin(position, coin, area)
		if coin == nil {
			canal <- &Move{
				X: 0,
				Y: 0,
			}
		}
		for t := 0; t < moveTime; t++ {
			position := ia.GetCharacter().GetPosition()
			area := ia.GetArea()
			angle := angleToNearestMultipleOf45(position, coin.GetPosition())
			x := math.Cos(angle)
			y := math.Sin(angle)
			x, y = size.NormalizeMove(x, y, position, area.GetSize())
			canal <- &Move{
				X: x,
				Y: y,
			}
		}
		canal <- nil
	}()
	return canal
}

func getClosedCoin(p1 inter.Position, last item.Coin, area inter.Area) item.Coin {
	closed := last
	closedDist := float64(999999999999999)
	if last != nil {
		cdx := p1.GetX() - closed.GetPosition().GetX()
		cdy := p1.GetY() - closed.GetPosition().GetY()
		closedDist = cdx * cdy
	}
	for _, i := range area.SearchItems(func(i inter.Item) bool {
		_, e := i.(item.Coin)
		return !e
	}) {
		c, _ := i.(item.Coin)
		dx := p1.GetX() - c.GetPosition().GetX()
		dy := p1.GetY() - c.GetPosition().GetY()
		nDist := dx*dx + dy*dy + 1
		if nDist < closedDist {
			closed = c
			closedDist = nDist
		}
	}
	return closed
}

func angleToNearestMultipleOf45(p1, p2 inter.Position) float64 {
	// 1. Calcular el vector
	dx := p2.GetX() - p1.GetX()
	dy := p2.GetY() - p1.GetY()

	// 2. Calcular el ángulo en radianes con Atan2 (devuelve [-π, π])
	angRad := math.Atan2(dy, dx)

	// 3. Convertir a grados
	angDeg := angRad * (180.0 / math.Pi)

	// 4. Normalizar a [0, 360) para trabajar cómodamente
	if angDeg < 0 {
		angDeg += 360.0
	}

	// 5. Redondear al múltiplo de 45° más cercano
	//    Dividimos entre 45, redondeamos al entero más cercano, y multiplicamos de vuelta
	multiples := math.Round(angDeg / 45.0)
	nearestAngle := multiples * 45.0

	// 6. Asegurar que si da 360 (por redondeo desde 337.5+) lo dejemos en 0
	if nearestAngle >= 360.0 {
		nearestAngle = 0.0
	}

	return nearestAngle
}
