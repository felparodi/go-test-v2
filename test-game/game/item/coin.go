package item

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/position"
)

type Coin interface {
	inter.Item
	SetPoint(int)
	GetPoint() int
}

type BasicCoin struct {
	id  string
	pos inter.Position
}

type CoinEvent struct {
	owner   Coin
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

func (c *BasicCoin) GetId() string {
	return c.id
}

func (c *BasicCoin) GetPosition() inter.Position {
	return c.pos
}

func (c *BasicCoin) SetPosition(pos inter.Position) {
	c.pos = pos
}

func NewCoin(i int, s inter.Size) Coin {
	return &BasicCoin{
		id:  fmt.Sprintf("coin_%d", i),
		pos: position.GetRandPosistion(s),
	}
}

func (c *BasicCoin) Collition(i inter.Item) []inter.Event {
	events := []inter.Event{}
	if c.isCollition(i) {
		switch i.(type) {
		case inter.Character:
			events = append(events,
				&CoinEvent{name: "add-points", owner: c, targets: []inter.Item{i}},
				&CoinEvent{name: "move-item-random-pose", owner: c},
			)
		}
	}
	return events
}

func (c *BasicCoin) isCollition(i inter.Item) bool {
	dx := i.GetPosition().GetX() - c.GetPosition().GetX()
	dy := i.GetPosition().GetY() - c.GetPosition().GetY()
	distance := dx*dx + dy*dy
	return distance < 900
}

func (c *BasicCoin) Update(_ float64, _ inter.Size) []inter.Event {
	return []inter.Event{}
}

func (c *BasicCoin) GetColitonArea() []inter.ColitionaArea {
	return []inter.ColitionaArea{}
}

func (c *BasicCoin) ProcessEvent(e inter.Event) {

}

func (c *BasicCoin) GetPoint() int {
	return 10
}

func (c *BasicCoin) SetPoint(int) {

}

func (c *BasicCoin) GetType() string {
	return "coin"
}
