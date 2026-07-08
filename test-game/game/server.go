package game

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	world        *World
	players      map[string]*Player
	upgrader     websocket.Upgrader
	mu           sync.RWMutex
	gameLoopDone chan bool
	// Nuevo: para optimizar broadcasts
	broadcastChan chan []byte
	// Nuevo: rate limiting
	rateLimiter map[string]*RateLimiter
}

type RateLimiter struct {
	lastUpdate time.Time
	count      int
	mu         sync.Mutex
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type GameState struct {
	Players map[string]PlayerData `json:"players"`
	Items   []Item                `json:"items"`
}

type PlayerData struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Score     int     `json:"score"`
	Vx        float64 `json:"vx,omitempty"`
	Vy        float64 `json:"vy,omitempty"`
	LastAngle float64 `json:"lastAngle,omitempty"` // Nueva: última dirección
}

func NewServer() *Server {
	return &Server{
		world: NewWorld(),
		players: make(map[string]*Player),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		gameLoopDone:  make(chan bool),
		broadcastChan: make(chan []byte, 100),
		rateLimiter:   make(map[string]*RateLimiter),
	}
}

// Manejar conexión WebSocket
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar conexión: %v", err)
		return
	}
	defer conn.Close()

	playerID := r.URL.Query().Get("id")
	if playerID == "" {
		playerID = conn.RemoteAddr().String()
	}

	player := NewPlayer(playerID, conn)

	s.mu.Lock()
	s.players[playerID] = player
	s.world.mu.Lock()
	s.world.Players[playerID] = player
	s.world.mu.Unlock()
	s.mu.Unlock()

	s.rateLimiter[playerID] = &RateLimiter{}

	log.Printf("Jugador %s conectado", playerID)

	s.sendGameState(player)
	go s.readMessages(player)

	<-s.gameLoopDone
}

// Leer mensajes del cliente mejorado
func (s *Server) readMessages(player *Player) {
	for {
		var msg Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al leer mensaje de %s: %v", player.ID, err)
			s.removePlayer(player.ID)
			return
		}

		if !s.checkRateLimit(player.ID) {
			log.Printf("Rate limit excedido para %s", player.ID)
			continue
		}

		s.handleMessage(player, msg)
	}
}

// Rate limiting mejorado
func (s *Server) checkRateLimit(playerID string) bool {
	rl, exists := s.rateLimiter[playerID]
	if !exists {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastUpdate) > time.Second {
		rl.count = 0
		rl.lastUpdate = now
	}

	rl.count++
	return rl.count <= 60
}

func (s *Server) initMessage(player *Player, msg Message) {
	data := msg.Payload.(map[string]interface{})
	if id, ok := data["playerId"].(string); ok && id != player.ID {
		s.mu.Lock()
		delete(s.players, player.ID)
		s.world.mu.Lock()
		delete(s.world.Players, player.ID)
		
		player.ID = id
		s.players[id] = player
		s.world.Players[id] = player
		s.world.mu.Unlock()
		s.mu.Unlock()
		log.Printf("Jugador renombrado a %s", id)
	}
}

func (s *Server) moveMessage(player *Player, msg Message) {
	data := msg.Payload.(map[string]interface{})
	velocityX := data["velocityX"].(float64)
	velocityY := data["velocityY"].(float64)

	// Limitar velocidad máxima
	maxSpeed := 250.0
	speed := math.Sqrt(velocityX*velocityX + velocityY*velocityY)
	if speed > maxSpeed {
		scale := maxSpeed / speed
		velocityX *= scale
		velocityY *= scale
	}

	// Guardar última dirección si hay movimiento
	if math.Abs(velocityX) > 0.1 || math.Abs(velocityY) > 0.1 {
		player.LastAngle = math.Atan2(velocityY, velocityX)
	}

	player.VelocityX = velocityX
	player.VelocityY = velocityY
}

func (s *Server) actionMessage(player *Player, msg Message) {
	actionType := msg.Payload.(map[string]interface{})["action"].(string)
	log.Printf("Jugador %s realiza acción: %s", player.ID, actionType)
}


// Manejar mensajes mejorado
func (s *Server) handleMessage(player *Player, msg Message) {
	switch msg.Type {
	case "init": s.initMessage(player, msg)
	case "move": s.moveMessage(player, msg)
	case "action": s.actionMessage(player, msg)
	default:
		log.Printf("Mensaje desconocido de %s: %s", player.ID, msg.Type)
	}
}

// Enviar estado del juego con optimización de broadcast
func (s *Server) sendGameStateToAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := s.getGameState()
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error al codificar estado: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, player := range s.players {
		wg.Add(1)
		go func(p *Player) {
			defer wg.Done()
			err := p.Send(data)
			if err != nil {
				log.Printf("Error al enviar a %s: %v", p.ID, err)
			}
		}(player)
	}
	wg.Wait()
}

// Bucle principal del juego optimizado
func (s *Server) GameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	physicsTicker := time.NewTicker(10 * time.Millisecond)
	defer physicsTicker.Stop()

	lastTime := time.Now()
	accumulator := 0.0
	const fixedDeltaTime = 1.0 / 60.0

	for {
		select {
		case <-s.gameLoopDone:
			return

		case <-physicsTicker.C:
			s.world.Update(fixedDeltaTime)

		case <-ticker.C:
			now := time.Now()
			deltaTime := now.Sub(lastTime).Seconds()
			lastTime = now

			accumulator += deltaTime
			if accumulator >= fixedDeltaTime {
				s.sendGameStateToAll()
				accumulator = 0
			}
		}
	}
}

// Obtener estado del juego incluyendo velocidades
func (s *Server) getGameState() GameState {
	s.world.mu.RLock()
	defer s.world.mu.RUnlock()

	playersData := make(map[string]PlayerData)
	for id, player := range s.world.Players {
		playersData[id] = PlayerData{
			X:         player.X,
			Y:         player.Y,
			Score:     player.Score,
			Vx:        player.VelocityX,
			Vy:        player.VelocityY,
			LastAngle: player.LastAngle,
		}
	}

	return GameState{
		Players: playersData,
		Items:   s.world.items,
	}
}

// Enviar estado a un jugador específico
func (s *Server) sendGameState(player *Player) {
	state := s.getGameState()
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error al codificar estado: %v", err)
		return
	}
	player.Send(data)
}

// Eliminar jugador
func (s *Server) removePlayer(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if player, exists := s.players[id]; exists {
		player.Conn.Close()
		delete(s.players, id)
		
		s.world.mu.Lock()
		delete(s.world.Players, id)
		s.world.mu.Unlock()
		
		log.Printf("Jugador %s desconectado", id)
	}
}