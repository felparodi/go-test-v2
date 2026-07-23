package area

import (
	"juego-websocket/game/inter"
	"juego-websocket/game/size"
)

type WorldArea struct {
	BasicArea
}

func NewWorldArea(s inter.Server) inter.Area {
	size := size.NewSizePadding(600, 800, 100)
	world := &WorldArea{
		BasicArea: newBasicArea(s, size),
	}
	for _, item := range GenerateCoins(15, size) {
		world.items[item.GetId()] = item
	}
	for _, item := range GenerateNPC(0, world) {
		world.items[item.GetId()] = item
	}
	return world
}
