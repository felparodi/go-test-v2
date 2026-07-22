package position

import (
	"juego-websocket/game/inter"
)

type Position struct {
	X     float64
	Y     float64
	Angle float64
}

func NewPosition(x, y, angle float64) inter.Position {
	return &Position{
		X:     x,
		Y:     y,
		Angle: angle,
	}
}

func (p *Position) GetX() float64 {
	return p.X
}

func (p *Position) GetY() float64 {
	return p.Y
}

func (p *Position) GetAngle() float64 {
	return p.Angle
}

func (p *Position) SetX(x float64) {
	p.X = x
}

func (p *Position) SetY(y float64) {
	p.Y = y
}

func (p *Position) SetAngle(angle float64) {
	p.Angle = angle
}

func (p *Position) Copy() inter.Position {
	return &Position{
		X:     p.X,
		Y:     p.Y,
		Angle: p.Angle,
	}
}
