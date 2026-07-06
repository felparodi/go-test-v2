package game

import (
	"sync"
)

// Zona de juego con coordenadas
type World struct {
	Width   int
	Height  int
	Players map[string]*Player
	mu      sync.RWMutex
	items   []Item
}

type Item struct {
	ID    string
	X     float64
	Y     float64
	Type  string
	Value int
}

func NewWorld() *World {
	return &World{
		Width:   800,
		Height:  600,
		Players: make(map[string]*Player),
		items:   generateItems(),
	}
}

func generateItems() []Item {
	items := []Item{}
	// Generar items aleatorios en el mapa
	for i := 0; i < 20; i++ {
		items = append(items, Item{
			ID:    "item_" + string(rune(i)),
			X:     float64(i * 30 % 800),
			Y:     float64(i * 40 % 600),
			Type:  "recurso",
			Value: 10,
		})
	}
	return items
}

// Actualizar el mundo (simulación de física)
func (w *World) Update(deltaTime float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Actualizar cada jugador
	for _, player := range w.Players {
		// Movimiento con inercia
		player.X += player.VelocityX * deltaTime
		player.Y += player.VelocityY * deltaTime

		// Colisiones con bordes
		if player.X < 0 {
			player.X = 0
			player.VelocityX = 0
		}
		if player.X > float64(w.Width) {
			player.X = float64(w.Width)
			player.VelocityX = 0
		}
		if player.Y < 0 {
			player.Y = 0
			player.VelocityY = 0
		}
		if player.Y > float64(w.Height) {
			player.Y = float64(w.Height)
			player.VelocityY = 0
		}

		// Recolectar items cercanos
		for i := len(w.items) - 1; i >= 0; i-- {
			item := w.items[i]
			dx := player.X - item.X
			dy := player.Y - item.Y
			if dx*dx+dy*dy < 2500 { // Radio de 50 unidades
				player.Score += item.Value
				// Remover el item
				w.items = append(w.items[:i], w.items[i+1:]...)
			}
		}
	}
}