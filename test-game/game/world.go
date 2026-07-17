package game

import (
	"juego-websocket/game/ia"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"log"
	"math/rand"
	"sync"
)

// Zona de juego con coordenadas
type World struct {
	size    inter.Size
	players map[string]inter.Player
	items   map[string]inter.Item
	server  inter.Server
	mu      sync.RWMutex
}

type Size struct {
	Height float64
	Width  float64
}

func (s *Size) GetHeight() float64 {
	return s.Height
}

func (s *Size) GetWidth() float64 {
	return s.Width
}

func (s *Size) Copy() inter.Size {
	return &Size{
		Height: s.Height,
		Width:  s.Width,
	}
}

type WorldState struct {
	Items      []inter.Item
	Characters []inter.Character
	Players    []inter.Player
}

func (ws *WorldState) GetItems() []inter.Item {
	return ws.Items
}

func (ws *WorldState) GetCharacters() []inter.Character {
	return ws.Characters
}

func (ws *WorldState) GetPlayers() []inter.Player {
	return ws.Players
}

func generateCoins(cantItems int, w inter.World) []inter.Item {
	items := []inter.Item{}
	// Generar items aleatorios en el mapa
	for i := 0; i < cantItems; i++ {
		c := item.NewCoin(i, w.GetSize())
		items = append(items, c)
	}
	return items
}

func generateNPC(cantItems int, s inter.Server, w inter.World) []inter.Item {
	items := []inter.Item{}
	// Generar items aleatorios en el mapa
	basic := rand.Intn(cantItems + 1)
	for i := 0; i < basic; i++ {
		ia := ia.NewBasicIA(i, s, w)
		go ia.Start()
		items = append(items, ia.GetCharacter())
	}
	greed := rand.Intn(cantItems - basic + 1)
	for i := 0; i < greed; i++ {
		ia := ia.NewGreedIA(1, s, w)
		go ia.Start()
		items = append(items, ia.GetCharacter())
	}
	return items
}

func NewWorld(s inter.Server) inter.World {
	world := &World{
		size:    &Size{Width: 800, Height: 600},
		items:   make(map[string]inter.Item),
		players: make(map[string]inter.Player),
		server:  s,
	}
	for _, item := range generateCoins(0, world) {
		world.items[item.GetId()] = item
	}
	for _, item := range generateNPC(0, world.server, world) {
		world.items[item.GetId()] = item
	}
	return world
}

func (w *World) GetSize() inter.Size {
	return w.size
}

// Actualizar el mundo (simulación de física)
func (w *World) Update(deltaTime float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	events := []inter.Event{}
	// Actualizar items
	for _, item := range w.items {
		events = append(events, item.Update(deltaTime, w.size)...)
	}
	// Coliciones
	// Mejorar con mayas se colizion
	for _, it1 := range w.items {
		for _, it2 := range w.items {
			if it1.GetId() != it2.GetId() {
				events = append(events, it1.Collition(it2)...)
			}
		}
	}
	// Event Loop
	for _, e := range events {
		log.Println("Event", e.GetEventName(), e.GetOwner(), e.GetTragets())
		w.processEvent(e)
	}
}

func (w *World) AddPlayer(p inter.Player) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.players[p.GetId()] = p
	w.items[p.GetCharacter().GetId()] = p.GetCharacter()
}

func (w *World) RemovePlayer(p inter.Player) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.players, p.GetId())
	delete(w.items, p.GetCharacter().GetId())
}

func (w *World) RemovePlayerId(playerId string) {
	w.mu.Lock()
	p, exits := w.players[playerId]
	delete(w.players, playerId)
	if exits {
		delete(w.items, p.GetCharacter().GetId())
	}
	w.mu.Unlock()
}

func (w *World) RenamePlayer(oldName string, newName string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	player, exist := w.players[oldName]
	if exist {
		player.SetId(newName)
		delete(w.players, oldName)
		w.players[newName] = player
	}

}

func (w *World) GetCoins() []inter.Coin {
	ret := []inter.Coin{}
	for _, item := range w.items {
		switch item.(type) {
		case inter.Coin:
			c, _ := (item).(inter.Coin)
			ret = append(ret, c)
		}
	}
	return ret
}

func (w *World) GetWorldState() inter.WorldState {
	ret := &WorldState{
		Items:      []inter.Item{},
		Characters: []inter.Character{},
		Players:    []inter.Player{},
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
	for _, player := range w.players {
		ret.Players = append(ret.Players, player)
	}
	return ret
}

func (w *World) processEvent(e inter.Event) {
	switch e.GetEventName() {
	case "move-item-random-pose":
		e.GetOwner().SetPosition(item.GetRandPosistion(w.size))
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

func (w *World) GetPlayer(id string) (inter.Player, bool) {
	player, exists := w.players[id]
	return player, exists
}

func (w *World) GetPlayers() []inter.Player {
	valores := make([]inter.Player, 0, len(w.players))
	for _, v := range w.players {
		valores = append(valores, v)
	}
	return valores
}

func (w *World) RLock() {
	w.mu.RLock()
}

func (w *World) RUnlock() {
	w.mu.RUnlock()
}
