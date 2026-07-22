package ia

import (
	"math"
	"math/rand"
)

func RandomDirection() (float64, float64) {
	angle := float64(rand.Intn(8) * 45)
	x := math.Cos(angle)
	y := math.Sin(angle)
	return x, y
}
