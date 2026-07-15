package game

import "math/rand"

type Item interface {
	getId() string
	getPosition() Position
	setPosition(Position)
	update(float64, *World) []Event
	collition(Item, *World) []Event
}

type Position struct {
	X     float64
	Y     float64
	Angle float64
}

type Event interface {
	getEventName() string
	getOwner() Item
	getTragets() []Item
}

func getRandPosistion(w *World) Position {
	return Position{
		X: float64(rand.Intn(w.Width)),
		Y: float64(rand.Intn(w.Height)),
	}
}
