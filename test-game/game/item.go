package game

import "math/rand"

type Item interface {
	getId() string
	getPosition() Position
	getType() interface{}
	update(float64)
	collition(Item)
	getEvents() []ItemEvent
	cleanEvents()
}

type ItemEvent struct {
	Type string
}

type Position struct {
	X     float64
	Y     float64
	Angle float64
}

func getRandPosistion(w *World) Position {
	return Position{
		X: float64(rand.Intn(w.Width)),
		Y: float64(rand.Intn(w.Height)),
	}
}
