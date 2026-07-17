package game

import (
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	id          string
	character   inter.Character
	conn        *websocket.Conn
	mu          sync.Mutex
	world       inter.World
	server      inter.Server
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	lastUpdate time.Time
	count      int
	mu         sync.Mutex
}

func NewPlayer(id string, conn *websocket.Conn, s inter.Server, w inter.World) inter.Player {
	ret := &Player{
		id:          id,
		character:   item.NewCharacter(id, w.GetSize()),
		conn:        conn,
		server:      s,
		world:       w,
		rateLimiter: &RateLimiter{},
	}
	ret.character.SetControler(ret)
	return ret
}

func (p *Player) Send(message []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conn.WriteMessage(websocket.TextMessage, message)
}

/*
*	Que empiece a escucar al usuario
**/
func (p *Player) Start() error {
	p.readMessages()
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
		err := p.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al leer mensaje de %s: %v", p.id, err)
			p.server.RemovePlayerId(p.id)
			return
		}

		if !p.checkRateLimit() {
			log.Printf("Rate limit excedido para %s", p.id)
			continue
		}

		p.handleMessage(msg)
	}
}

func (p *Player) initMessage(msg Message) {
	data := msg.Payload.(map[string]interface{})
	if id, ok := data["playerId"].(string); ok && id != p.id {
		p.mu.Lock()
		p.world.RemovePlayer(p)
		p.id = id
		p.world.AddPlayer(p)
		p.mu.Unlock()
		log.Printf("Jugador renombrado a %s", id)
	}
}

// Setear el vector de movimiento
func (player *Player) moveMessage(msg Message) {
	data := msg.Payload.(map[string]interface{})
	velocityX := data["velocityX"].(float64)
	velocityY := data["velocityY"].(float64)
	player.character.Move(velocityX, velocityY)
}

func (player *Player) actionMessage(msg Message) {
	actionType := msg.Payload.(map[string]interface{})["action"].(string)
	log.Printf("Jugador %s realiza acción: %s", player.id, actionType)
	switch actionType {
	case "z", "Z":
		player.character.AddAction(&Action{name: "shoot"})
	}
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
		log.Printf("Mensaje desconocido de %s: %s", p.id, msg.Type)
	}
}

func (p *Player) GetId() string {
	return p.id
}

func (p *Player) GetCharacter() inter.Character {
	return p.character
}

func (p *Player) End() error {
	return p.conn.Close()
}

type Action struct {
	name string
}

func (a *Action) GetName() string {
	return a.name
}

func (a *Action) GetData() interface{} {
	return nil
}
