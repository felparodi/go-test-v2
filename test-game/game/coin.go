package game

import (
	"fmt"
)

type Coin struct {
	ID     string
	pos    Position
	events []ItemEvent
}

func (c *Coin) getId() string {
	return c.ID
}

func (c *Coin) getType() interface{} {
	return c
}

func (c *Coin) getPosition() Position {
	return c.pos
}

func NewCoin(i int, w *World) *Coin {

	return &Coin{
		ID:     fmt.Sprintf("coin_%d", i),
		pos:    getRandPosistion(w),
		events: []ItemEvent{},
	}
}

func (c *Coin) collition(i Item) {
	if i.getType() == "player" {
		if c.isCollition(i) {
			c.events = []ItemEvent{
				{Type: "add-point"},
				{Type: "remove-item"},
				{Type: "create-item"},
			}
		}
	}
}

func (c *Coin) isCollition(i Item) bool {
	dx := i.getPosition().X - c.getPosition().X
	dy := i.getPosition().Y - c.getPosition().Y
	distance := dx*dx + dy*dy
	return distance < 900
}

func (c *Coin) update(deltaTime float64) {

}

func (c *Coin) getEvents() []ItemEvent {
	return c.events
}

func (c *Coin) cleanEvents() {
	c.events = []ItemEvent{}
}
