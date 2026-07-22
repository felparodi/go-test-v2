package area

import (
	"juego-websocket/game/ia"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"math/rand"
)

func GenerateCoins(cantItems int, s inter.Size) []inter.Item {
	items := []inter.Item{}
	// Generar items aleatorios en el mapa
	for i := 0; i < cantItems; i++ {
		c := item.NewCoin(i, s)
		items = append(items, c)
	}
	return items
}

func GenerateNPC(cantItems int, a inter.Area) []inter.Item {
	items := []inter.Item{}
	// Generar items aleatorios en el mapa
	basic := rand.Intn(cantItems + 1)
	for i := 0; i < basic; i++ {
		ia := ia.NewBasicIA(i, a)
		go ia.Start()
		items = append(items, ia.GetCharacter())
	}
	greed := rand.Intn(cantItems - basic + 1)
	for i := 0; i < greed; i++ {
		ia := ia.NewGreedIA(1, a)
		go ia.Start()
		items = append(items, ia.GetCharacter())
	}
	return items
}
