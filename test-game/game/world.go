package game

import (
	"log"
	"sync"
)

// Zona de juego con coordenadas
type World struct {
	Width   int
	Height  int
	Players map[string]*Player
	Items   map[string]Item
	server  *Server
	mu      sync.RWMutex
}

func generateItems(cantItems int, w *World) []Item {
	items := []Item{}
	// Generar items aleatorios en el mapa
	for i := 0; i < cantItems; i++ {
		c := NewCoin(i, w)
		log.Println("generateItems.c", i, c)
		items = append(items, c)
	}
	log.Println("generateItems", items)
	return items
}

func NewWorld(s *Server) *World {
	world := &World{
		Width:   800,
		Height:  600,
		Players: make(map[string]*Player),
		Items:   make(map[string]Item),
		server:  s,
	}
	for _, item := range generateItems(20, world) {
		log.Println("NewWorld.for", item.getId(), item)
		world.Items[item.getId()] = item
	}
	log.Println("NewWorld", world)
	return world
}

// Actualizar el mundo (simulación de física)
func (w *World) Update(deltaTime float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	//items := []Item(append(w.Items))
	for _, player := range w.Players {
		player.update(deltaTime)

		// Recolectar items con mejor detección
		/*
			for _, item := range w.Items {
				if item.isCollition(player) {
					player.Score += 10
					w.items = append(w.items[:i], w.items[i+1:]...)

					// Regenerar item en nueva posición
					newItem := NewCoin(i, w)
					w.items = append(w.items, newItem)
				}
			}
		*/
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
