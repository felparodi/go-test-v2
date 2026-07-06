package game

import (
	"encoding/json"
	"log"
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
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Score int     `json:"score"`
}

func NewServer() *Server {
	return &Server{
		world: NewWorld(),
		players: make(map[string]*Player),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // En producción, validar orígenes
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		gameLoopDone: make(chan bool),
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

	// Registrar nuevo jugador
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

	log.Printf("Jugador %s conectado", playerID)

	// Enviar estado inicial
	s.sendGameState(player)

	// Bucle de lectura de mensajes
	go s.readMessages(player)

	// Mantener conexión abierta
	<-s.gameLoopDone
}

// Leer mensajes del cliente
func (s *Server) readMessages(player *Player) {
	for {
		var msg Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al leer mensaje de %s: %v", player.ID, err)
			s.removePlayer(player.ID)
			return
		}

		s.handleMessage(player, msg)
	}
}

// Manejar mensajes del cliente
func (s *Server) handleMessage(player *Player, msg Message) {
	switch msg.Type {
	case "move":
		data := msg.Payload.(map[string]interface{})
		velocityX := data["velocityX"].(float64)
		velocityY := data["velocityY"].(float64)
		
		// Limitar velocidad máxima
		maxSpeed := 200.0
		speed := velocityX*velocityX + velocityY*velocityY
		if speed > maxSpeed*maxSpeed {
			scale := maxSpeed / speed
			velocityX *= scale
			velocityY *= scale
		}
		
		player.VelocityX = velocityX
		player.VelocityY = velocityY

	case "action":
		actionType := msg.Payload.(map[string]interface{})["action"].(string)
		log.Printf("Jugador %s realiza acción: %s", player.ID, actionType)
	}
}

// Enviar estado del juego a todos los jugadores
func (s *Server) sendGameStateToAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := s.getGameState()
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error al codificar estado: %v", err)
		return
	}

	// Enviar a todos los jugadores
	for _, player := range s.players {
		select {
		case <-time.After(5 * time.Millisecond):
			// Timeout para no bloquear
		default:
			err := player.Send(data)
			if err != nil {
				log.Printf("Error al enviar a %s: %v", player.ID, err)
			}
		}
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

// Obtener estado del juego
func (s *Server) getGameState() GameState {
	s.world.mu.RLock()
	defer s.world.mu.RUnlock()

	playersData := make(map[string]PlayerData)
	for id, player := range s.world.Players {
		playersData[id] = PlayerData{
			X:     player.X,
			Y:     player.Y,
			Score: player.Score,
		}
	}

	return GameState{
		Players: playersData,
		Items:   s.world.items,
	}
}

// Bucle principal del juego
func (s *Server) GameLoop() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	lastTime := time.Now()

	for {
		select {
		case <-s.gameLoopDone:
			return
		case <-ticker.C:
			now := time.Now()
			deltaTime := now.Sub(lastTime).Seconds()
			lastTime = now

			// Actualizar mundo
			s.world.Update(deltaTime)

			// Enviar estado a todos los jugadores
			s.sendGameStateToAll()
		}
	}
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