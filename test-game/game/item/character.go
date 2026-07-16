package item

import (
	"juego-websocket/game/inter"
	"math"
)

type Character struct {
	id       string
	position inter.Position
	velocity inter.Position
	oldPos   inter.Position
	score    int
	player   inter.Player
	acctions []inter.Action
}
type CharacterEvent struct {
	owner   *Character
	name    string
	targets []inter.Item
}

func (ce *CharacterEvent) GetOwner() inter.Item {
	return ce.owner
}

func (ce *CharacterEvent) GetEventName() string {
	return ce.name
}

func (ce *CharacterEvent) GetTragets() []inter.Item {
	return ce.targets
}

func NewCharacter(id string, s inter.Size) inter.Character {
	pos := GetRandPosistion(s)
	return &Character{
		id:       id,
		position: pos,
		oldPos:   pos,
		velocity: &Position{X: 0, Y: 0, Angle: 0},
		acctions: []inter.Action{},
		score:    0,
	}
}

func (c *Character) SetPlayer(p inter.Player) {
	c.player = p
}

func (c *Character) GetPlayer() inter.Player {
	return c.player
}

func (c *Character) Move(velocityX float64, velocityY float64) {
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
		c.velocity.SetAngle(math.Atan2(velocityY, velocityX))
		c.position.SetAngle(math.Atan2(velocityY, velocityX))
	}

	c.velocity.SetX(velocityX)
	c.velocity.SetY(velocityY)
}

func (c *Character) Update(deltaTime float64, s inter.Size) []inter.Event {
	//Copio la posicion anterior
	c.oldPos = c.position
	events := []inter.Event{}
	const friction = 0.92
	const minVelocity = 0.1
	// Aplicar fricción gradual
	if c.velocity.GetX() != 0 || c.velocity.GetY() != 0 {
		// Reducir velocidad gradualmente cuando no hay input
		c.velocity.SetX(c.velocity.GetX() * friction)
		c.velocity.SetY(c.velocity.GetY() * friction)

		// Si la velocidad es muy pequeña, detener
		if math.Abs(c.velocity.GetX()) < minVelocity {
			c.velocity.SetX(0)
		}
		if math.Abs(c.velocity.GetY()) < minVelocity {
			c.velocity.SetY(0)
		}
	}

	if c.velocity.GetX() != 0 || c.velocity.GetY() != 0 {
		// Ángulo de la velocidad (dirección del movimiento)
		velocityAngle := math.Atan2(c.velocity.GetY(), c.velocity.GetX())
		// Diferencia de ángulos (puedes devolver cualquiera de estos)
		c.velocity.SetAngle(velocityAngle) // Ángulo de la velocidad
	}

	// Mover con deltaTime para consistencia de velocidad
	c.position.SetX(c.position.GetX() + c.velocity.GetX()*deltaTime)
	c.position.SetY(c.position.GetY() + c.velocity.GetY()*deltaTime)
	if c.position.GetX() < 0 {
		c.position.SetX(0)
		c.velocity.SetX(0)
		events = append(events, &CharacterEvent{name: "limit-min-x", owner: c})
	}
	if c.position.GetX() > float64(s.GetWidth()) {
		c.position.SetX(s.GetWidth())
		c.velocity.SetX(0)
		events = append(events, &CharacterEvent{name: "limit-max-x", owner: c})
	}
	if c.position.GetY() < 0 {
		c.position.SetY(0)
		c.velocity.SetY(0)
		events = append(events, &CharacterEvent{name: "limit-min-y", owner: c})
	}
	if c.position.GetY() > float64(s.GetHeight()) {
		c.position.SetY(s.GetHeight())
		c.velocity.SetY(0)
		events = append(events, &CharacterEvent{name: "limit-max-y", owner: c})
	}
	for _, action := range c.acctions {
		switch action.GetName() {
		case "shoot":
			events = append(events, &CharacterEvent{name: "create-bullet", owner: c})
		}
	}
	clear(c.acctions)
	c.acctions = []inter.Action{}
	return events
}

func (c *Character) GetId() string {
	return c.id
}

func (c *Character) GetPosition() inter.Position {
	return c.position
}

func (c *Character) GetVelocity() inter.Position {
	return c.velocity
}

func (c *Character) SetPosition(pos inter.Position) {
	c.position = pos
}

func (c *Character) Collition(_ inter.Item) []inter.Event {
	return []inter.Event{}
}

func (c *Character) GetColitonArea() []inter.ColitionaArea {
	return []inter.ColitionaArea{}
}

func (c *Character) AddScore(score int) {
	c.score += score
}

func (c *Character) GetScore() int {
	return c.score
}

func (c *Character) AddAction(a inter.Action) {
	c.acctions = append(c.acctions, a)
}
