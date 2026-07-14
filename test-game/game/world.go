package game

import (
	"math"
	"math/rand"
	"sync"
)

// Zona de juego con coordenadas
type World struct {
	Width   int
	Height  int
	Players map[string]*Player
	server  *Server
	mu      sync.RWMutex
	items   []Item
}

func NewWorld(s *Server) *World {
	return &World{
		Width:   800,
		Height:  600,
		Players: make(map[string]*Player),
		items:   generateItems(),
		server:  s,
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

	// Factor de fricción para movimiento más natural
	const friction = 0.92

	for _, player := range w.Players {
		// Aplicar fricción gradual
		if player.VelocityX != 0 || player.VelocityY != 0 {
			// Reducir velocidad gradualmente cuando no hay input
			player.VelocityX *= friction
			player.VelocityY *= friction

			// Si la velocidad es muy pequeña, detener
			if math.Abs(player.VelocityX) < 0.1 {
				player.VelocityX = 0
			}
			if math.Abs(player.VelocityY) < 0.1 {
				player.VelocityY = 0
			}
		}

		if player.VelocityX != 0 || player.VelocityY != 0 {
			// Ángulo de la velocidad (dirección del movimiento)
			velocityAngle := math.Atan2(player.VelocityY, player.VelocityX)
			// Diferencia de ángulos (puedes devolver cualquiera de estos)
			player.Angle = velocityAngle // Ángulo de la velocidad
		}

		// Mover con deltaTime para consistencia de velocidad
		player.X += player.VelocityX * deltaTime
		player.Y += player.VelocityY * deltaTime

		// Colisiones con bordes mejoradas
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

		// Recolectar items con mejor detección
		for i := len(w.items) - 1; i >= 0; i-- {
			item := w.items[i]
			dx := player.X - item.X
			dy := player.Y - item.Y
			distance := dx*dx + dy*dy
			if distance < 900 { // Radio de colección
				player.Score += item.Value
				w.items = append(w.items[:i], w.items[i+1:]...)

				// Regenerar item en nueva posición
				newItem := Item{
					ID:    "item_" + string(rune(i)),
					X:     float64(rand.Intn(w.Width)),
					Y:     float64(rand.Intn(w.Height)),
					Type:  "recurso",
					Value: 10 + rand.Intn(20),
				}
				w.items = append(w.items, newItem)
			}
		}
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
