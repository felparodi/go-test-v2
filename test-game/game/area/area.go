package area

import (
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"sync"
)

// Zona de juego con coordenadas
type BasicArea struct {
	size            inter.Size
	items           map[string]inter.Item
	itemsDictionary map[string][]string
	server          inter.Server
	activeChannel   chan bool
	mu              sync.RWMutex
}

func newBasicArea(server inter.Server, size inter.Size) BasicArea {
	return BasicArea{
		server:          server,
		size:            size,
		items:           map[string]inter.Item{},
		itemsDictionary: map[string][]string{},
		activeChannel:   make(chan bool),
	}
}

// @TODO implement
func (a *BasicArea) Start() (chan bool, error) {
	return a.activeChannel, nil
}

// @TODO implement
func (a *BasicArea) Stop() error {
	return nil
}

func (w *BasicArea) GetSize() inter.Size {
	return w.size
}

// Actualizar el mundo (simulación de física)
func (a *BasicArea) Update(deltaTime float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.update(deltaTime)
}

func (a *BasicArea) update(deltaTime float64) error {
	events := []inter.Event{}
	// Actualizar items
	for _, item := range a.items {
		events = append(events, item.Update(deltaTime, a.size)...)
	}
	// Coliciones
	// Mejorar con mayas se colizion
	for _, it1 := range a.items {
		for _, it2 := range a.items {
			if it1.GetId() != it2.GetId() {
				events = append(events, it1.Collition(it2)...)
			}
		}
	}
	// Event Loop
	for _, e := range events {
		//log.Println("Event", e.GetEventName(), e.GetOwner(), e.GetTragets())
		a.processEvent(e)
	}
	return nil
}

func (w *BasicArea) SearchItems(find func(inter.Item) bool) []inter.Item {
	ret := []inter.Item{}
	for _, item := range w.items {
		if find(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

func (w *BasicArea) GetState() inter.AreaState {
	ret := &AreaState{
		Items:      []inter.Item{},
		Characters: []inter.Character{},
	}
	for _, item := range w.items {
		switch item.(type) {
		case inter.Character:
			c, _ := (item).(inter.Character)
			ret.Characters = append(ret.Characters, c)
		case inter.Item:
			c, _ := (item).(inter.Item)
			ret.Items = append(ret.Items, c)
		}
	}
	return ret
}

func (w *BasicArea) processEvent(e inter.Event) {
	switch e.GetEventName() {
	case "move-item-random-pose":
		e.GetOwner().SetPosition(w.size.GetRandPosistion())
	case "remove":
		delete(w.items, e.GetOwner().GetId())
	case "create-bullet":
		bullet := item.NewBullet(e.GetOwner())
		w.items[bullet.GetId()] = bullet
	}
	if e.GetTragets() != nil {
		for _, t := range e.GetTragets() {
			t.ProcessEvent(e)
		}
	}
}

func (w *BasicArea) AddItem(item inter.Item) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.addItem(item)
}

// @TODO Tirar error si ya se uso el id?
func (a *BasicArea) addItem(item inter.Item) error {
	a.items[item.GetId()] = item
	list, exist := a.itemsDictionary[item.GetType()]
	if exist {
		a.itemsDictionary[item.GetType()] = append(list, item.GetId())
	} else {
		a.itemsDictionary[item.GetType()] = []string{item.GetId()}
	}
	return nil
}

func (a *BasicArea) RemoveItem(item inter.Item) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.removeItem(item)
}

// @TODO cuando tiraria error?
func (a *BasicArea) removeItem(item inter.Item) error {
	delete(a.items, item.GetId())
	a.deleteItemDictionary(item.GetType(), item.GetId())
	return nil
}

func (a *BasicArea) deleteItemDictionary(itemType, id string) {
	skip := -1
	itemIds, exist := a.itemsDictionary[itemType]
	if exist {
		for index, itemId := range itemIds {
			if itemId == id {
				skip = index
				break
			}
		}
	}
	if skip >= 0 {
		a.itemsDictionary[itemType] = append(itemIds[:skip], itemIds[skip+1:]...)
	}
}

func (a *BasicArea) RemoveItemId(id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.removeItemId(id)
}

// @TODO cuando deberia debolver error
func (a *BasicArea) removeItemId(id string) error {
	item, exist := a.items[id]
	if exist {
		delete(a.items, id)
		a.deleteItemDictionary(item.GetType(), id)
	}
	return nil
}

func (a *BasicArea) SmallSearchItems(itemType string, filter func(inter.Item) bool) []inter.Item {
	itemIds, exist := a.itemsDictionary[itemType]
	ret := []inter.Item{}
	if exist {
		for _, itemId := range itemIds {
			item, existItem := a.items[itemId]
			if existItem && filter(item) {
				ret = append(ret, item)
			}
		}
	}
	return ret
}
