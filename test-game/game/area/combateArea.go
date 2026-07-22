package area

import (
	"juego-websocket/game/board"
	"juego-websocket/game/inter"
)

type CombatArea struct {
	BasicArea
	timeline      []Turn
	actualTurn    Turn
	turnTiem      float64
	turnTimeStart float64
}

type Turn interface {
	getOwner() inter.Character
	getTimeStart() float64
}

type DrawTurn struct {
	ManyCards int
	Owner     inter.Character
	StartAt   float64
}

type PlayTurn struct {
	Card  board.Card
	Owner inter.Character
}
