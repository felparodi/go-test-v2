package game

import (
	"juego-websocket/game/area"
	"juego-websocket/game/inter"
	"juego-websocket/game/position"
	"sync"
)

type BasicGame struct {
	mu            sync.RWMutex
	players       map[string]inter.Player
	worlds        map[string]inter.Area
	activeChannel chan bool
}

func NewGame(s inter.Server) inter.Game {
	return &BasicGame{
		players: map[string]inter.Player{},
		worlds: map[string]inter.Area{
			"0": area.NewWorldArea(s),
		},
		activeChannel: make(chan bool),
	}
}

// @TODO implementar
func (g *BasicGame) Start() (chan bool, error) {
	return g.activeChannel, nil
}

// @TODO implementar
func (g *BasicGame) Stop() error {
	return nil
}

func (g *BasicGame) AddPlayer(p inter.Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.addPlayer(p)
}

// @TODO cuando deberia tirar error
func (g *BasicGame) addPlayer(p inter.Player) error {
	g.players[p.GetId()] = p
	c := p.GetCharacter()
	c.SetPosition(position.GetRandPosistion(g.worlds["0"].GetSize()))
	g.worlds["0"].AddItem(c)
	return nil
}

func (g *BasicGame) RemovePlayer(p inter.Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.removePlayer(p)
}

// @TODO cuando debolver error
// @TODO ver si hacer algo con los areas
func (g *BasicGame) removePlayer(p inter.Player) error {
	return g.removePlayerId(p.GetId())
}

func (g *BasicGame) RemovePlayerId(id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.removePlayerId(id)
}

// @TODO cuando tirar error
func (g *BasicGame) removePlayerId(id string) error {
	player, exist := g.players[id]
	if exist {
		delete(g.players, id)
		c := player.GetCharacter()
		for _, world := range g.worlds {
			world.RemoveItem(c)
		}
	}
	return nil
}

func (g *BasicGame) RenamePlayer(oldId, newId string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.renamePlayer(oldId, newId)
}

// @TODO Ver que hacer con los punteros hacia abajo
// @TODO Ver cuando tirar error
func (g *BasicGame) renamePlayer(oldId, newId string) error {
	player, exist := g.players[oldId]
	if exist {
		player.SetId(newId)
		delete(g.players, oldId)
		g.players[newId] = player
	}
	return nil
}

func (g *BasicGame) GetPlayer(newId string) (inter.Player, bool) {
	player, exist := g.players[newId]
	return player, exist
}

func (g *BasicGame) GetPlayers() []inter.Player {
	ret := []inter.Player{}
	for _, player := range g.players {
		ret = append(ret, player)
	}
	return ret
}

func (g *BasicGame) GetState() inter.AreaState {
	return g.worlds["0"].GetState()
}

func (g *BasicGame) Update(delta float64) error {
	return g.worlds["0"].Update(delta)
}
