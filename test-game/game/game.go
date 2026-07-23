package game

import (
	"encoding/json"
	"fmt"
	"juego-websocket/game/area"
	"juego-websocket/game/inter"
	"log"
	"math/rand"
	"sync"
	"time"
)

const MAX_PLAYER = 5

type BasicGame struct {
	mu            sync.RWMutex
	players       map[string]inter.Player
	worlds        map[string]inter.Area
	playerOnWorld map[string]string
	worldPlayers  map[string]map[string]struct{}
	worldsMu      map[string]*sync.RWMutex
	activeChannel chan bool
	threadChannel chan bool
}

func NewGame(s inter.Server) inter.Game {
	return &BasicGame{
		players:       map[string]inter.Player{},
		worlds:        map[string]inter.Area{},
		playerOnWorld: make(map[string]string),
		worldPlayers:  make(map[string]map[string]struct{}),
		worldsMu:      make(map[string]*sync.RWMutex),
		activeChannel: make(chan bool),
		threadChannel: make(chan bool),
	}
}

// @TODO implementar
func (g *BasicGame) Start() (chan bool, error) {
	go g.runUpdateThread()
	return g.activeChannel, nil
}

// @TODO implementar
func (g *BasicGame) Stop() error {
	defer func() { g.activeChannel <- false }()
	g.threadChannel <- false
	return nil
}

// @TODO ver como parar si no hay player o si no hay que hacerlo
func (g *BasicGame) runUpdateThread() error {
	physicsTicker := time.NewTicker(20 * time.Millisecond)
	defer physicsTicker.Stop()
	const fixedDeltaTime = 1.0 / 60.0
	lastTime := time.Now().UnixMilli()
	for {
		actualTimer := time.Now().UnixMilli()
		delta := float64(int(actualTimer-lastTime)) / 600.0
		//log.Println("RT", lastTime, actualTimer, delta)
		select {
		case <-physicsTicker.C:
			g.update(delta)
		case <-g.threadChannel:
			return nil
		}
		lastTime = actualTimer
	}
}

func (g *BasicGame) AddPlayer(p inter.Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.addPlayer(p)
}

// @TODO cuando deberia tirar error
func (g *BasicGame) addPlayer(p inter.Player) error {
	g.players[p.GetId()] = p
	worldName := ""
	for world, players := range g.worldPlayers {
		if len(players) < MAX_PLAYER {
			worldName = world
			break
		}
	}
	log.Println("addPlayer.1", worldName)
	var world inter.Area
	if worldName != "" {
		world = g.worlds[worldName]
	} else {
		worldName = fmt.Sprintf("World_%d", rand.Intn(1000))
		//@Por si ser repite nombre
		/*
			for _, exist := g.worldPlayers[worldName]; exist; {
				worldName = fmt.Sprintf("World_%d", rand.Intn(1000))
			}
		*/
		world = area.NewWorldArea(worldName)
		world.Start()
		g.worlds[worldName] = world
		g.worldsMu[worldName] = &sync.RWMutex{}
		g.worldPlayers[worldName] = make(map[string]struct{})
	}
	log.Println("addPlayer.2", world)
	g.playerOnWorld[p.GetId()] = worldName
	g.worldPlayers[worldName][p.GetId()] = struct{}{}
	c := p.GetCharacter()
	c.SetPosition(world.GetSize().GetRandPosistion())
	world.AddItem(c)
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
		worldName := g.playerOnWorld[id]
		g.worldsMu[worldName].Lock()
		defer g.worldsMu[worldName].Unlock()
		delete(g.players, id)
		delete(g.playerOnWorld, id)
		delete(g.worldPlayers[worldName], id)
		c := player.GetCharacter()
		g.worlds[worldName].RemoveItem(c)
		if len(g.worldPlayers[worldName]) == 0 {
			world := g.worlds[worldName]
			world.Stop()
			delete(g.worlds, worldName)
			delete(g.worldPlayers, worldName)
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
		worldName := g.playerOnWorld[oldId]
		g.worldsMu[worldName].Lock()
		defer g.worldsMu[worldName].Unlock()
		player.SetId(newId)
		delete(g.players, oldId)
		g.players[newId] = player
		g.playerOnWorld[newId] = worldName
		delete(g.playerOnWorld, oldId)
		delete(g.worldPlayers[worldName], oldId)
		g.worldPlayers[worldName][newId] = struct{}{}
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

// @TODO Ver si funciona la idea
func (g *BasicGame) update(delta float64) error {
	worldsSet := make(map[string]struct{})
	for _, world := range g.playerOnWorld {
		worldsSet[world] = struct{}{}
	}
	for worldName := range worldsSet {
		go g.updateWorld(worldName, delta)
	}
	return nil
}

func (g *BasicGame) updateWorld(worldName string, delta float64) error {
	g.worldsMu[worldName].Lock()
	defer g.worldsMu[worldName].Unlock()
	g.worlds[worldName].Update(delta)
	worldState := g.worlds[worldName].GetState()
	playersData := make(map[string]interface{})
	for _, player := range worldState.GetCharacters() {
		playersData[player.GetControler().GetId()] = toJson(player)
	}
	itemsData := []interface{}{}
	for _, item := range worldState.GetItems() {
		//log.Println(item)
		itemsData = append(itemsData, toJson(item))
	}
	state := GameState{
		CharacterData: playersData,
		Items:         itemsData,
		SendInfo: SendInfo{
			InfoType: "game-state",
		},
	}
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("Error al codificar estado: %v", err)
		return err
	}
	for playerId := range g.worldPlayers[worldName] {
		go g.players[playerId].Send(data)
		//@TODO
	}
	return err
}
