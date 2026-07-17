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

func NormalizeMove(x, y float64, pos inter.Position, world inter.World) (float64, float64) {
	size := world.GetSize()
	if pos.GetX()+x < 0 {
		x = 0
	}
	if pos.GetY()+y < 0 {
		y = 0
	}
	if pos.GetX()+x > size.GetWidth() {
		x = 0
	}
	if pos.GetY()+y > size.GetHeight() {
		y = 0
	}
	return x, y
}
