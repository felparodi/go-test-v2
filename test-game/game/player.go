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
	Position    Position
	Velocity    Position
	OldPos      Position
	Score       int
	Conn        *websocket.Conn
	mu          sync.Mutex
	world       *World
	server      *Server
	rateLimiter *RateLimiter
	events      []ItemEvent
}

type RateLimiter struct {
	lastUpdate time.Time
	count      int
	mu         sync.Mutex
}

func NewPlayer(id string, conn *websocket.Conn, s *Server, w *World) *Player {
	pos := getRandPosistion(w)
	return &Player{
		ID:          id,
		Position:    pos,
		OldPos:      pos,
		Velocity:    Position{X: 0, Y: 0, Angle: 0},
		Score:       0,
		Conn:        conn,
		server:      s,
		world:       w,
		rateLimiter: &RateLimiter{},
		events:      []ItemEvent{},
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
	data := msg.Payload.(map[string]interface{})
	if id, ok := data["playerId"].(string); ok && id != p.ID {
		p.mu.Lock()
		p.world.removePlayer(p)
		p.ID = id
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
		player.Velocity.Angle = math.Atan2(velocityY, velocityX)
		player.Position.Angle = math.Atan2(velocityY, velocityX)
	}

	player.Velocity.X = velocityX
	player.Velocity.Y = velocityY
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

func (player *Player) update(deltaTime float64) {
	//Copio la posicion anterior
	player.OldPos = player.Position
	const friction = 0.92
	const minVelocity = 0.1
	// Aplicar fricción gradual
	if player.Velocity.X != 0 || player.Velocity.Y != 0 {
		// Reducir velocidad gradualmente cuando no hay input
		player.Velocity.X *= friction
		player.Velocity.Y *= friction

		// Si la velocidad es muy pequeña, detener
		if math.Abs(player.Velocity.X) < minVelocity {
			player.Velocity.X = 0
		}
		if math.Abs(player.Velocity.Y) < minVelocity {
			player.Velocity.Y = 0
		}
	}

	if player.Velocity.X != 0 || player.Velocity.Y != 0 {
		// Ángulo de la velocidad (dirección del movimiento)
		velocityAngle := math.Atan2(player.Velocity.Y, player.Velocity.X)
		// Diferencia de ángulos (puedes devolver cualquiera de estos)
		player.Velocity.Angle = velocityAngle // Ángulo de la velocidad
	}

	// Mover con deltaTime para consistencia de velocidad
	player.Position.X += player.Velocity.X * deltaTime
	player.Position.Y += player.Velocity.Y * deltaTime
	if player.Position.X < 0 {
		player.Position.X = 0
		player.Velocity.X = 0
		player.events = append(player.events, ItemEvent{Type: "limit-min-x"})
	}
	if player.Position.X > float64(player.world.Width) {
		player.Position.X = float64(player.world.Width)
		player.Velocity.X = 0
		player.events = append(player.events, ItemEvent{Type: "limit-max-x"})
	}
	if player.Position.Y < 0 {
		player.Position.Y = 0
		player.Velocity.Y = 0
		player.events = append(player.events, ItemEvent{Type: "limit-min-y"})
	}
	if player.Position.Y > float64(player.world.Height) {
		player.Position.Y = float64(player.world.Height)
		player.Velocity.Y = 0
		player.events = append(player.events, ItemEvent{Type: "limit-max-y"})
	}
}

func (p *Player) isCollition(_ Item) bool {
	return false
}

func (p *Player) getId() string {
	return p.ID
}

func (p *Player) getType() interface{} {
	return p
}

func (p *Player) getPosition() Position {
	return p.Position
}

func (p *Player) getEvents() []ItemEvent {
	return p.events
}

func (p *Player) cleanEvents() {

}

func (p *Player) collition(_ Item) {

}
