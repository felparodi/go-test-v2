package world

import (
	"juego-websocket/game/item"
	"juego-websocket/game/player"
)

type World interface {
	getSize() item.Size
	getPlayers() []player.Player
	removePlayer(player.Player)
	addPlayer(player.Player)
	removePlayerId(string)
	getWorldState()
}

type WorldState struct {
	Coins      []item.Item
	Characters []item.Item
	Players    []player.Player
}
