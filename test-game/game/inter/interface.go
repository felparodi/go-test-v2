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
	GetType() string
}

type CharacterControler interface {
	Thread
	GetId() string
	GetCharacter() Character
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
	GetMaxHeight() float64
	GetMaxWidth() float64
	GetMinHeight() float64
	GetMinWidth() float64
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

type Thread interface {
	Start() (chan bool, error)
	Stop() error
}

type Game interface {
	Thread
	AddPlayer(Player) error
	RemovePlayer(Player) error
	RemovePlayerId(string) error
	RenamePlayer(string, string) error
	GetPlayer(string) (Player, bool)
	GetPlayers() []Player
	Update(float64) error //@TODO temporal
	GetState() AreaState  //@TODO temporal
}

type Area interface {
	Thread
	GetSize() Size
	AddItem(Item) error
	RemoveItem(Item) error
	RemoveItemId(string) error
	Update(float64) error
	GetState() AreaState
	SearchItems(func(Item) bool) []Item
	SmallSearchItems(string, func(Item) bool) []Item
}

type Server interface {
	RemovePlayerId(string) error
	GameLoop() error
	HandleWebSocket(http.ResponseWriter, *http.Request)
}

type AreaState interface {
	GetItems() []Item
	GetCharacters() []Character
}

type Action interface {
	GetName() string
	GetData() interface{}
}
