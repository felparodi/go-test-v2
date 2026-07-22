package game

import (
	"encoding/json"
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

	player := player.NewPlayer(playerID, conn, s, s.game)
	s.game.AddPlayer(player)
	log.Printf("Jugador %s conectado", playerID)
	s.sendGameState(player)
	activeplayer, _ := player.Start()
	for {
		select {
		case act := <-activeplayer:
			if !act {
				go s.RemovePlayerId(player.GetId())
			}
		case <-s.gameLoopDone:
			return
		}
	}
}

// Enviar estado del juego con optimización de broadcast // No es un broadcas puro ya que se los envia uno a uno
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
	for _, player := range s.game.GetPlayers() {
		wg.Add(1)
		go func(p inter.Player) {
			defer wg.Done()
			err := p.Send(data)
			if err != nil {
				log.Printf("Error al enviar a %s: %v", p.GetId(), err)
			}
			playerInfo := toJson(p)

			pData, err := json.Marshal(
				PlayerSate{
					Player: playerInfo,
					SendInfo: SendInfo{
						InfoType: "own",
					},
				},
			)
			if err == nil {
				err = p.Send(pData)
			}
			if err != nil {
				log.Printf("Error al enviar a %s: %v", p.GetId(), err)
			}
		}(player)
	}
	wg.Wait()
}

// Bucle principal del juego optimizado
func (s *Server) GameLoop() error {
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
			return nil

		case <-physicsTicker.C:
			s.game.Update(fixedDeltaTime)

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

	worldData := s.game.GetState()

	playersData := make(map[string]interface{})

	for _, player := range worldData.GetCharacters() {
		playersData[player.GetControler().GetId()] = toJson(player)
	}
	itemsData := []interface{}{}
	for _, item := range worldData.GetItems() {
		//log.Println(item)
		itemsData = append(itemsData, toJson(item))
	}

	r := GameState{
		CharacterData: playersData,
		Items:         itemsData,
		SendInfo: SendInfo{
			InfoType: "game-state",
		},
	}

	//log.Println(r)
	return r
}

// Enviar estado a un jugador específico
func (s *Server) sendGameState(player inter.Player) {
	state := s.getGameState()
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error al codificar estado: %v", err)
		return
	}
	player.Send(data)
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
