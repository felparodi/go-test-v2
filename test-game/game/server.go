package game

import (
	"juego-websocket/game/inter"
	"juego-websocket/game/player"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	game         inter.Game
	upgrader     websocket.Upgrader
	mu           sync.RWMutex
	gameLoopDone chan bool
	// Nuevo: para optimizar broadcasts
	broadcastChan chan []byte
}

type SendInfo struct {
	InfoType string `json:"type"`
}
type GameState struct {
	SendInfo
	CharacterData map[string]interface{} `json:"players"`
	Items         []interface{}          `json:"items"`
}

type PlayerSate struct {
	SendInfo
	Player interface{} `json:"player"`
}

func NewServer() inter.Server {
	s := &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		gameLoopDone:  make(chan bool),
		broadcastChan: make(chan []byte, 100),
	}
	s.game = NewGame(s)
	return s
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

	player := player.NewPlayer(playerID, conn, s.game)
	s.game.AddPlayer(player)
	log.Printf("Jugador %s conectado", playerID)
	active, _ := player.Start()
	for {
		select {
		case act := <-active:
			if !act {
				go s.RemovePlayerId(player.GetId())
			}
		case <-s.gameLoopDone:
			return
		}
	}
}

// Bucle principal del juego optimizado
func (s *Server) GameLoop() error {
	s.game.Start()
	defer s.game.Stop()
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	const fixedDeltaTime = 1.0 / 60.0

	for {
		select {
		case <-s.gameLoopDone:
			return nil
		}
	}
}

// Eliminar jugador
func (s *Server) RemovePlayerId(id string) error {
	log.Printf("Jugador %s desconectado", id)
	s.mu.Lock()
	defer s.mu.Unlock()
	if player, exists := s.game.GetPlayer(id); exists {
		log.Printf("Jugador %s desconectado", id)
		player.Stop()
		log.Printf("Jugador %s desconectado", id)
		s.game.RemovePlayerId(id)
		log.Printf("Jugador %s desconectado", id)
	}
	log.Printf("Jugador %s desconectado", id)
	return nil
}
