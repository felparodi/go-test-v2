package game

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID          string
	X           float64
	Y           float64
	VelocityX   float64
	VelocityY   float64
	Angle       float64
	Score       int
	Conn        *websocket.Conn
	mu          sync.Mutex
	world       *World
	server      *Server
	rateLimiter *RateLimiter
}

type PlayerData struct {
	Id    string  `json:"playerId"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Score int     `json:"score"`
	Vx    float64 `json:"vx,omitempty"`
	Vy    float64 `json:"vy,omitempty"`
	Angle float64 `json:"angle"`
}

type RateLimiter struct {
	lastUpdate time.Time
	count      int
	mu         sync.Mutex
}

func NewPlayer(id string, conn *websocket.Conn, s *Server, w *World) *Player {
	return &Player{
		ID:          id,
		X:           400,
		Y:           300,
		Score:       0,
		Angle:       0,
		Conn:        conn,
		server:      s,
		world:       w,
		rateLimiter: &RateLimiter{},
	}
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
	p.server.initMessage(p, msg)
}

// Setear el vector de movimiento
func (player *Player) moveMessage(msg Message) {
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
		player.Angle = math.Atan2(velocityY, velocityX)
	}

	player.VelocityX = velocityX
	player.VelocityY = velocityY
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

func (p *Player) toData() PlayerData {
	return PlayerData{
		Id:    p.ID,
		X:     p.X,
		Y:     p.Y,
		Score: p.Score,
		Vx:    p.VelocityX,
		Vy:    p.VelocityY,
		Angle: p.Angle,
	}
}
