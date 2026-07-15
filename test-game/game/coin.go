package game

import (
	"fmt"
)

type Coin struct {
	ID     string
	pos    Position
	events []ItemEvent
}

type CoinEvent struct {
	owner   *Coin
	name    string
	targets []Item
}

func (ce *CoinEvent) getOwner() Item {
	return ce.owner
}

func (ce *CoinEvent) getEventName() string {
	return ce.name
}

func (ce *CoinEvent) getTragets() []Item {
	return ce.targets
}

func (c *Coin) getId() string {
	return c.ID
}

func (c *Coin) getPosition() Position {
	return c.pos
}

func (c *Coin) setPosition(pos Position) {
	c.pos = pos
}

func NewCoin(i int, w *World) *Coin {

	return &Coin{
		ID:     fmt.Sprintf("coin_%d", i),
		pos:    getRandPosistion(w),
		events: []ItemEvent{},
	}
}

func (c *Coin) collition(i Item, w *World) []WorldEvent {
	events := []WorldEvent{}
	if c.isCollition(i) {
		switch i.(type) {
		case *Player:
			events = append(events,
				&CoinEvent{
					name:    "add-point",
					owner:   c,
					targets: []Item{i},
				},
				&CoinEvent{
					name:  "move-item-random-pose",
					owner: c,
				},
			)
		}
	}
	return events
}

func (c *Coin) isCollition(i Item) bool {
	dx := i.getPosition().X - c.getPosition().X
	dy := i.getPosition().Y - c.getPosition().Y
	distance := dx*dx + dy*dy
	return distance < 900
}

func (c *Coin) update(_ float64, _ *World) []WorldEvent {
	return []WorldEvent{}
}
