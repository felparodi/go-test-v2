package ia

import (
	"juego-websocket/game/inter"
	"math"
	"math/rand"
)

func RandomDirection() (float64, float64) {
	angle := float64(rand.Intn(8) * 45)
	x := math.Cos(angle)
	y := math.Sin(angle)
	return x, y
}

func NormalizeMove(x, y float64, pos inter.Position, area inter.Area) (float64, float64) {
	size := area.GetSize()
	if pos.GetX()+x < area.GetSize().GetMinWidth() {
		x = 0
	}
	if pos.GetY()+y < area.GetSize().GetMinHeight() {
		y = 0
	}
	if pos.GetX()+x > size.GetMaxWidth() {
		x = 0
	}
	if pos.GetY()+y > size.GetMaxHeight() {
		y = 0
	}
	return x, y
}
