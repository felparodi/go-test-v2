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
		Players: make(map[string]*Player),
		Items:   make(map[string]Item),
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
	for _, player := range w.Players {
		events = append(events, player.update(deltaTime, w)...)
	}

	// Coliciones
	for _, p := range w.Players {
		for _, i := range w.Items {
			events = append(events, i.collition(p, w)...)
		}
	}
	// Event Loop

	//Coliciones
	for _, e := range events {
		log.Println("Event", e.getEventName(), e.getOwner(), e.getTragets())
		w.processEvent(e)
	}

}

func (w *World) getPlayer(playerId string) (*Player, bool) {
	player, exists := w.Players[playerId]
	return player, exists
}

func (w *World) addPlayer(player *Player) {
	w.mu.Lock()
	w.Players[player.ID] = player
	w.mu.Unlock()
}

func (w *World) removePlayer(player *Player) {
	w.mu.Lock()
	delete(w.Players, player.ID)
	w.mu.Unlock()
}

func (w *World) removePlayerId(playerId string) {
	w.mu.Lock()
	delete(w.Players, playerId)
	w.mu.Unlock()
}

func (w *World) getPlayers() []*Player {
	valores := make([]*Player, 0, len(w.Players))
	for _, v := range w.Players {
		valores = append(valores, v)
	}
	return valores
}

func (w *World) getItems() []Item {
	valores := make([]Item, 0, len(w.Items))
	for _, v := range w.Items {
		if v != nil {
			valores = append(valores, v) // Desreferenciar el puntero
		}
	}
	return valores
}

func (w *World) processEvent(e WorldEvent) {
	switch e.getEventName() {
	case "move-item-random-pose":
		e.getOwner().setPosition(getRandPosistion(w))
	case "add-point":
		for _, t := range e.getTragets() {
			switch t.(type) {
			case *Player:
				p, _ := t.(*Player)
				p.Score += 10
			default:
				break
			}
		}
	}
}
