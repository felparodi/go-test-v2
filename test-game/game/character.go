package game

import (
	"math"
)

type Character struct {
	ID       string
	Position Position
	Velocity Position
	OldPos   Position
	Score    int
}

type CharacterEvent struct {
	owner   *Character
	name    string
	targets []Item
}

func (ce *CharacterEvent) getOwner() Item {
	return ce.owner
}

func (ce *CharacterEvent) getEventName() string {
	return ce.name
}

func (ce *CharacterEvent) getTragets() []Item {
	return ce.targets
}

func NewCharacter(id string, s *Server, w *World) *Character {
	pos := getRandPosistion(w)
	return &Character{
		ID:       id,
		Position: pos,
		OldPos:   pos,
		Velocity: Position{X: 0, Y: 0, Angle: 0},
		Score:    0,
	}
}

func (c *Character) move(velocityX float64, velocityY float64) {
	// Limitar velocidad máxima
	maxSpeed := 250.0
	speed := math.Sqrt(velocityX*velocityX + velocityY*velocityY)
	if speed > maxSpeed {
		scale := maxSpeed / speed
		velocityX *= scale
		velocityY *= scale
	}

	// Guardar última dirección si hay movimiento
	if math.Abs(velocityX) > 0.1 || math.Abs(velocityY) > 0.1 {
		c.Velocity.Angle = math.Atan2(velocityY, velocityX)
		c.Position.Angle = math.Atan2(velocityY, velocityX)
	}

	c.Velocity.X = velocityX
	c.Velocity.Y = velocityY
}

func (c *Character) update(deltaTime float64, w *World) []WorldEvent {
	//Copio la posicion anterior
	c.OldPos = c.Position
	events := []WorldEvent{}
	const friction = 0.92
	const minVelocity = 0.1
	// Aplicar fricción gradual
	if c.Velocity.X != 0 || c.Velocity.Y != 0 {
		// Reducir velocidad gradualmente cuando no hay input
		c.Velocity.X *= friction
		c.Velocity.Y *= friction

		// Si la velocidad es muy pequeña, detener
		if math.Abs(c.Velocity.X) < minVelocity {
			c.Velocity.X = 0
		}
		if math.Abs(c.Velocity.Y) < minVelocity {
			c.Velocity.Y = 0
		}
	}

	if c.Velocity.X != 0 || c.Velocity.Y != 0 {
		// Ángulo de la velocidad (dirección del movimiento)
		velocityAngle := math.Atan2(c.Velocity.Y, c.Velocity.X)
		// Diferencia de ángulos (puedes devolver cualquiera de estos)
		c.Velocity.Angle = velocityAngle // Ángulo de la velocidad
	}

	// Mover con deltaTime para consistencia de velocidad
	c.Position.X += c.Velocity.X * deltaTime
	c.Position.Y += c.Velocity.Y * deltaTime
	if c.Position.X < 0 {
		c.Position.X = 0
		c.Velocity.X = 0
		events = append(events, &CharacterEvent{name: "limit-min-x", owner: c})
	}
	if c.Position.X > float64(w.Width) {
		c.Position.X = float64(w.Width)
		c.Velocity.X = 0
		events = append(events, &CharacterEvent{name: "limit-max-x", owner: c})
	}
	if c.Position.Y < 0 {
		c.Position.Y = 0
		c.Velocity.Y = 0
		events = append(events, &CharacterEvent{name: "limit-min-y", owner: c})
	}
	if c.Position.Y > float64(w.Height) {
		c.Position.Y = float64(w.Height)
		c.Velocity.Y = 0
		events = append(events, &CharacterEvent{name: "limit-max-y", owner: c})
	}
	return events
}

func (c *Character) getId() string {
	return c.ID
}

func (c *Character) getPosition() Position {
	return c.Position
}

func (c *Character) setPosition(pos Position) {
	c.Position = pos
}

func (c *Character) collition(_ Item, _ *World) []WorldEvent {
	return []WorldEvent{}
}
