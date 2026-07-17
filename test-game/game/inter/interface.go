package inter

import "net/http"

type Player interface {
	CharacterControler
	SetId(string)
	Send(message []byte) error
}

type Item interface {
	GetId() string
	GetPosition() Position
	SetPosition(Position)
	Update(float64, Size) []Event
	Collition(Item) []Event
	GetColitonArea() []ColitionaArea
	ProcessEvent(Event)
}

type Coin interface {
	Item
	SetPoint(int)
	GetPoint() int
}

type Bullet interface {
	Item
	GetOwner() Item
}

type CharacterControler interface {
	GetId() string
	GetCharacter() Character
	Start() error
	End() error
}

type Character interface {
	Item
	GetVelocity() Position
	Move(float64, float64)
	SetScore(int)
	GetScore() int
	AddAction(Action)
	SetControler(CharacterControler)
	GetControler() CharacterControler
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
	Copy() Size
}

type Position interface {
	GetX() float64
	GetY() float64
	GetAngle() float64
	SetX(float64)
	SetY(float64)
	SetAngle(float64)
	Copy() Position
}

type World interface {
	GetSize() Size
	RemovePlayer(Player)
	RemovePlayerId(string)
	RenamePlayer(string, string)
	AddPlayer(Player)
	GetPlayer(string) (Player, bool)
	GetPlayers() []Player
	Update(float64)
	RLock()
	RUnlock()
	GetWorldState() WorldState
	GetCoins() []Coin
}

type Server interface {
	RemovePlayerId(string)
	GameLoop()
	HandleWebSocket(http.ResponseWriter, *http.Request)
}

type WorldState interface {
	GetItems() []Item
	GetCharacters() []Character
	GetPlayers() []Player
}

type Action interface {
	GetName() string
	GetData() interface{}
}
