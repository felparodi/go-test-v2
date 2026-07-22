package item

import (
	"fmt"
	"juego-websocket/game/inter"
	"juego-websocket/game/position"
	"math"
	"math/rand"
)

type Bullet interface {
	inter.Item
	GetOwner() inter.Item
}

type BasicBullet struct {
	id     string
	pos    inter.Position
	oldPos inter.Position
	owner  inter.Item
	vector inter.Position
}

type BulletEvent struct {
	owner   Bullet
	name    string
	targets []inter.Item
}

func (ce *BulletEvent) GetOwner() inter.Item {
	return ce.owner
}

func (ce *BulletEvent) GetEventName() string {
	return ce.name
}

func (ce *BulletEvent) GetTragets() []inter.Item {
	return ce.targets
}

func (b *BasicBullet) GetId() string {
	return b.id
}

func (b *BasicBullet) GetPosition() inter.Position {
	return b.pos
}

func (b *BasicBullet) SetPosition(pos inter.Position) {
	b.pos = pos
}

func NewBullet(owner inter.Item) Bullet {
	angle := owner.GetPosition().GetAngle()
	return NewBulletAngle(owner, angle)
}

func NewBulletAngle(owner inter.Item, angle float64) Bullet {
	//log.Println("New Bullet", owner.GetId())
	module := float64(150)
	vector := position.NewPosition(
		math.Cos(angle)*module,
		math.Sin(angle)*module,
		angle,
	)
	return &BasicBullet{
		id:     fmt.Sprintf("Bullet_%d", rand.Intn(9999999)),
		pos:    owner.GetPosition().Copy(),
		owner:  owner,
		vector: vector,
	}
}

func (b *BasicBullet) Collition(i inter.Item) []inter.Event {
	events := []inter.Event{}
	if b.isCollition(i) {
		switch i.(type) {
		case inter.Character:
			if i.GetId() != b.owner.GetId() {
				events = append(events,
					&BulletEvent{name: "remove-points", owner: b, targets: []inter.Item{i}},
					&BulletEvent{name: "remove", owner: b},
				)
			}
		}
	}
	return events
}

func (b *BasicBullet) isCollition(i inter.Item) bool {
	dx := i.GetPosition().GetX() - b.GetPosition().GetX()
	dy := i.GetPosition().GetY() - b.GetPosition().GetY()
	distance := dx*dx + dy*dy
	return distance < 300
}

func (b *BasicBullet) Update(deltaTime float64, s inter.Size) []inter.Event {
	b.oldPos = b.pos
	events := []inter.Event{}
	// Mover con deltaTime para consistencia de velocidad
	b.pos.SetX(b.pos.GetX() + b.vector.GetX()*deltaTime)
	b.pos.SetY(b.pos.GetY() + b.vector.GetY()*deltaTime)
	if b.pos.GetX() < 0 || b.pos.GetX() > float64(s.GetWidth()) || b.pos.GetY() < 0 || b.pos.GetY() > float64(s.GetHeight()) {
		events = append(events, &BulletEvent{name: "remove", owner: b})
	}
	return events
}

func (b *BasicBullet) GetColitonArea() []inter.ColitionaArea {
	return []inter.ColitionaArea{}
}

func (b *BasicBullet) GetOwner() inter.Item {
	return b.owner
}

func (b *BasicBullet) ProcessEvent(e inter.Event) {

}

func (c *BasicBullet) GetType() string {
	return "bullet"
}
