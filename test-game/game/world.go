package game

import (
	"log"
	"sync"
)

type WorldEvent interface {
	getEventName() string
	getOwner() Item
	getTragets() []Item
}

// Zona de juego con coordenadas
type World struct {
	Width   int
	Height  int
	Players map[string]*Player
	Items   map[string]Item
	Server  *Server
	mu      sync.RWMutex
}

type WorldState struct {
	Coins      []*Coin
	Characters []*Character
	Players    []*Player
}

func generateItems(cantItems int, w *World) []Item {
	items := []Item{}
	// Generar items aleatorios en el mapa
	for i := 0; i < cantItems; i++ {
		c := NewCoin(i, w)
		items = append(items, c)
	}
	return items
}

func NewWorld(s *Server) *World {
	world := &World{
		Width:   800,
		Height:  600,
		Items:   make(map[string]Item),
		Players: make(map[string]*Player),
		Server:  s,
	}
	for _, item := range generateItems(20, world) {
		world.Items[item.getId()] = item
	}
	return world
}

// Actualizar el mundo (simulación de física)
func (w *World) Update(deltaTime float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	events := []WorldEvent{}
	for _, item := range w.Items {
		events = append(events, item.update(deltaTime, w)...)
	}

	// Coliciones
	for _, it1 := range w.Items {
		for _, it2 := range w.Items {
			if it1.getId() != it2.getId() {
				events = append(events, it1.collition(it2, w)...)
			}
		}
	}
	// Event Loop

	//Coliciones
	for _, e := range events {
		log.Println("Event", e.getEventName(), e.getOwner(), e.getTragets())
		w.processEvent(e)
	}

}

func (w *World) addPlayer(p *Player) {
	w.mu.Lock()
	w.Players[p.getId()] = p
	w.Items[p.Character.getId()] = p.Character
	w.mu.Unlock()
}

func (w *World) removePlayer(p *Player) {
	w.mu.Lock()
	delete(w.Players, p.getId())
	delete(w.Items, p.Character.getId())
	w.mu.Unlock()
}

func (w *World) removePlayerId(playerId string) {
	w.mu.Lock()
	p, exits := w.Players[playerId]
	delete(w.Players, playerId)
	if exits {
		delete(w.Items, p.Character.getId())
	}
	w.mu.Unlock()
}

func (w *World) getWorldState() WorldState {
	ret := WorldState{
		Coins:      []*Coin{},
		Characters: []*Character{},
		Players:    []*Player{},
	}
	for _, item := range w.Items {
		switch item.(type) {
		case *Character:
			c, _ := (item).(*Character)
			ret.Characters = append(ret.Characters, c)
			if c.getPlayer() != nil {
				ret.Players = append(ret.Players, c.getPlayer())
			}
		case *Coin:
			c, _ := (item).(*Coin)
			ret.Coins = append(ret.Coins, c)
		}
	}
	return ret
}

func (w *World) processEvent(e WorldEvent) {
	switch e.getEventName() {
	case "move-item-random-pose":
		e.getOwner().setPosition(getRandPosistion(w))
	case "add-point":
		for _, t := range e.getTragets() {
			switch t.(type) {
			case *Character:
				c, _ := t.(*Character)
				c.Score += 10
			default:
				break
			}
		}
	}
}

func (w *World) getPlayer(id string) (*Player, bool) {
	player, exists := w.Players[id]
	return player, exists
}

func (w *World) getPlayers() []*Player {
	valores := make([]*Player, 0, len(w.Players))
	for _, v := range w.Players {
		valores = append(valores, v)
	}
	return valores
}
