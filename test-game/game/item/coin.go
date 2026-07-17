package item

import (
	"fmt"
	"juego-websocket/game/inter"
)

type Coin struct {
	id  string
	pos inter.Position
}

type CoinEvent struct {
	owner   *Coin
	name    string
	targets []inter.Item
}

func (ce *CoinEvent) GetOwner() inter.Item {
	return ce.owner
}

func (ce *CoinEvent) GetEventName() string {
	return ce.name
}

func (ce *CoinEvent) GetTragets() []inter.Item {
	return ce.targets
}

func (c *Coin) GetId() string {
	return c.id
}

func (c *Coin) GetPosition() inter.Position {
	return c.pos
}

func (c *Coin) SetPosition(pos inter.Position) {
	c.pos = pos
}

func NewCoin(i int, s inter.Size) inter.Coin {
	return &Coin{
		id:  fmt.Sprintf("coin_%d", i),
		pos: GetRandPosistion(s),
	}
}

func (c *Coin) Collition(i inter.Item) []inter.Event {
	events := []inter.Event{}
	if c.isCollition(i) {
		switch i.(type) {
		case *Character:
			events = append(events,
				&CoinEvent{name: "add-points", owner: c, targets: []inter.Item{i}},
				&CoinEvent{name: "move-item-random-pose", owner: c},
			)
		}
	}
	return events
}

func (c *Coin) isCollition(i inter.Item) bool {
	dx := i.GetPosition().GetX() - c.GetPosition().GetX()
	dy := i.GetPosition().GetY() - c.GetPosition().GetY()
	distance := dx*dx + dy*dy
	return distance < 900
}

func (c *Coin) Update(_ float64, _ inter.Size) []inter.Event {
	return []inter.Event{}
}

func (c *Coin) GetColitonArea() []inter.ColitionaArea {
	return []inter.ColitionaArea{}
}

func (c *Coin) ProcessEvent(e inter.Event) {

}

func (c *Coin) GetPoint() int {
	return 10
}

func (c *Coin) SetPoint(int) {

}
