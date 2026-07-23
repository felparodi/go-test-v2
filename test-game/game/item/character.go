package item

import (
	"fmt"
	"juego-websocket/game/event"
	"juego-websocket/game/inter"
	"juego-websocket/game/position"
	"juego-websocket/game/size"
	"math"
	"math/rand"
)

type Character struct {
	id              string
	position        inter.Position
	velocity        inter.Position
	oldPos          inter.Position
	score           int
	controler       inter.CharacterControler
	acctions        []inter.Action
	acctionColdDown float64
}

func NewCharacter(pos inter.Position) inter.Character {
	return &Character{
		id:              fmt.Sprintf("Character_%d", rand.Intn(9999999)),
		position:        pos,
		oldPos:          pos,
		velocity:        position.NewPosition(0, 0, 0),
		acctions:        []inter.Action{},
		acctionColdDown: 0,
		score:           0,
	}
}

func (c *Character) SetControler(p inter.CharacterControler) {
	c.controler = p
}

func (c *Character) GetControler() inter.CharacterControler {
	return c.controler
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
	c.oldPos = c.position.Copy()
	c.acctionColdDown -= deltaTime
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
	vx, vy, limits := size.NormalizeMove(
		c.velocity.GetX()*deltaTime,
		c.velocity.GetY()*deltaTime,
		c.position,
		s,
	)

	// Mover con deltaTime para consistencia de velocidad
	c.position.SetX(c.position.GetX() + vx)
	c.position.SetY(c.position.GetY() + vy)
	for _, limit := range limits {
		switch limit {
		case "limit-min-x", "limit-max-x":
			c.velocity.SetX(0)
		case "limit-min-y", "limit-max-y":
			c.velocity.SetY(0)
		}
		events = append(events, event.NewEvent(limit, c, nil))
	}

	for _, action := range c.acctions {
		switch action.GetName() {
		case "shoot":
			events = append(events, event.NewEvent("create-bullet", c, nil))
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

func (c *Character) SetScore(score int) {
	c.score = score
}

func (c *Character) GetScore() int {
	return c.score
}

func (c *Character) AddAction(a inter.Action) {
	if c.acctionColdDown <= 0 {
		c.acctions = append(c.acctions, a)
		c.acctionColdDown = 1
	}
}

func (c *Character) ProcessEvent(e inter.Event) {
	switch e.GetEventName() {
	case "add-points":
		c.score += 10
	case "remove-points":
		c.score -= 3
		if c.score < 0 {
			c.score = 0
		}
	}
}

func (c *Character) GetType() string {
	return "character"
}
