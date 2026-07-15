package game

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID          string
	Character   *Character
	Conn        *websocket.Conn
	mu          sync.Mutex
	world       *World
	server      *Server
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	lastUpdate time.Time
	count      int
	mu         sync.Mutex
}

func NewPlayer(id string, conn *websocket.Conn, s *Server, w *World) *Player {
	ret := &Player{
		ID:          id,
		Character:   NewCharacter(id, w),
		Conn:        conn,
		server:      s,
		world:       w,
		rateLimiter: &RateLimiter{},
	}
	ret.Character.setPlayer(ret)
	return ret
}

func (p *Player) Send(message []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Conn.WriteMessage(websocket.TextMessage, message)
}

/*
*	Que empiece a escucar al usuario
**/
func (p *Player) start() error {
	return nil
}

// Rate limiting mejorado
func (p *Player) checkRateLimit() bool {
	rl := p.rateLimiter
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

func (p *Player) readMessages() {
	for {
		var msg Message
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al leer mensaje de %s: %v", p.ID, err)
			p.server.removePlayer(p.ID)
			return
		}

		if !p.checkRateLimit() {
			log.Printf("Rate limit excedido para %s", p.ID)
			continue
		}

		p.handleMessage(msg)
	}
}

func (p *Player) initMessage(msg Message) {
	data := msg.Payload.(map[string]interface{})
	if id, ok := data["playerId"].(string); ok && id != p.ID {
		p.mu.Lock()
		p.world.removePlayer(p)
		p.ID = id
		p.Character.ID = id
		p.world.addPlayer(p)
		p.mu.Unlock()
		log.Printf("Jugador renombrado a %s", id)
	}
}

// Setear el vector de movimiento
func (player *Player) moveMessage(msg Message) {
	data := msg.Payload.(map[string]interface{})
	velocityX := data["velocityX"].(float64)
	velocityY := data["velocityY"].(float64)
	player.Character.move(velocityX, velocityY)
}

func (player *Player) actionMessage(msg Message) {
	actionType := msg.Payload.(map[string]interface{})["action"].(string)
	log.Printf("Jugador %s realiza acción: %s", player.ID, actionType)
}

func (p *Player) handleMessage(msg Message) {
	switch msg.Type {
	case "init":
		p.initMessage(msg)
	case "move":
		p.moveMessage(msg)
	case "action":
		p.actionMessage(msg)
	default:
		log.Printf("Mensaje desconocido de %s: %s", p.ID, msg.Type)
	}
}

func (p *Player) getId() string {
	return p.ID
}

func (p *Player) cleanEvents() {

}

func (p *Player) collition(_ Item, _ *World) []WorldEvent {
	return []WorldEvent{}
}
