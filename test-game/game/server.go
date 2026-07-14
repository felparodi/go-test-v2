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
	worlds       map[string]*World
	upgrader     websocket.Upgrader
	mu           sync.RWMutex
	gameLoopDone chan bool
	// Nuevo: para optimizar broadcasts
	broadcastChan chan []byte
}

type GameState struct {
	Players map[string]interface{} `json:"players"`
	Items   []interface{}          `json:"items"`
}

func NewServer() *Server {
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
	s.worlds = map[string]*World{"0": NewWorld(s)}
	return s
}

func (s *Server) getWorldTo() *World {
	return s.worlds["0"]
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

	world := s.getWorldTo()
	player := NewPlayer(playerID, conn, s, world)

	s.mu.Lock()
	world.addPlayer(player)
	s.mu.Unlock()

	log.Printf("Jugador %s conectado", playerID)

	s.sendGameState(player)
	go player.readMessages()

	<-s.gameLoopDone
}

// @TODO Player word
func (s *Server) getPlayerWorld(p *Player) *World {
	return s.worlds["0"]
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
	for _, player := range s.getPlayers() {
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
			s.worlds["0"].Update(fixedDeltaTime)

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
	s.worlds["0"].mu.RLock()
	defer s.worlds["0"].mu.RUnlock()

	playersData := make(map[string]interface{})
	for id, player := range s.worlds["0"].Players {
		playersData[id] = itemToJson(player)
	}

	itemsData := []interface{}{}
	for _, item := range s.worlds["0"].getItems() {
		//log.Println(item)
		itemsData = append(itemsData, itemToJson(item))
	}

	r := GameState{
		Players: playersData,
		Items:   itemsData,
	}

	//log.Println(r)
	return r
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

	if player, exists := s.getPlayer(id); exists {
		player.Conn.Close()

		s.worlds["0"].removePlayerId(id)

		log.Printf("Jugador %s desconectado", id)
	}
}

/*
*
 */
func (s *Server) getPlayer(playerId string) (*Player, bool) {
	return s.worlds["0"].getPlayer(playerId)
}

func (s *Server) getPlayers() []*Player {
	players := make([]*Player, 0)
	for _, word := range s.worlds {
		players = append(players, word.getPlayers()...)
	}
	return players
}
