package event

import "juego-websocket/game/inter"

type BasicEvent struct {
	Owner   inter.Item
	Name    string
	Targets []inter.Item
}

func NewEvent(name string, owner inter.Item, targets []inter.Item) inter.Event {
	return &BasicEvent{
		Name:    name,
		Owner:   owner,
		Targets: targets,
	}
}

func (ce *BasicEvent) GetOwner() inter.Item {
	return ce.Owner
}

func (ce *BasicEvent) GetEventName() string {
	return ce.Name
}

func (ce *BasicEvent) GetTragets() []inter.Item {
	if ce.Targets != nil {
		return ce.Targets
	}
	return []inter.Item{}
}
