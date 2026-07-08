package game

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Player struct {
	ID        string
	X         float64
	Y         float64
	VelocityX float64
	VelocityY float64
	Angle	  float64
	Score     int
	Conn      *websocket.Conn
	mu        sync.Mutex
}

func NewPlayer(id string, conn *websocket.Conn) *Player {
	return &Player{
		ID:        id,
		X:         400,
		Y:         300,
		Score:     0,
		Angle:     0,
		Conn:      conn,
	}
}

func (p *Player) Send(message []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Conn.WriteMessage(websocket.TextMessage, message)
}