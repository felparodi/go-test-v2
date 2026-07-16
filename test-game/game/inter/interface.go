package inter

import "net/http"

type Player interface {
	Start() error
	Send(message []byte) error
	ReadMessages()
	GetId() string
	GetCharacter() Character
	CloseConnection()
}

type Item interface {
	GetId() string
	GetPosition() Position
	SetPosition(Position)
	Update(float64, Size) []Event
	Collition(Item) []Event
	GetColitonArea() []ColitionaArea
}

type Coin interface {
	Item
}

type Bullet interface {
	Item
}

type Character interface {
	Item
	GetVelocity() Position
	GetPlayer() Player
	SetPlayer(Player)
	Move(float64, float64)
	AddScore(int)
	GetScore() int
	AddAction(Action)
}

type Event interface {
	GetEventName() string
	GetOwner() Item
	GetTragets() []Item
}

type ColitionaArea interface {
	GetPosition() Position
	GetSize() Size
}

type Size interface {
	GetHeight() float64
	GetWidth() float64
}

type Position interface {
	GetX() float64
	GetY() float64
	GetAngle() float64
	SetX(float64)
	SetY(float64)
	SetAngle(float64)
}

type World interface {
	GetSize() Size
	RemovePlayer(Player)
	RemovePlayerId(string)
	AddPlayer(Player)
	GetPlayer(string) (Player, bool)
	GetPlayers() []Player
	Update(float64)
	RLock()
	RUnlock()
	GetWorldState() WorldState
}

type Server interface {
	RemovePlayerId(string)
	GameLoop()
	HandleWebSocket(http.ResponseWriter, *http.Request)
}

type WorldState interface {
	GetCoins() []Coin
	GetCharacters() []Character
	GetPlayers() []Player
}

type Action interface {
	GetName() string
	GetData() interface{}
}
