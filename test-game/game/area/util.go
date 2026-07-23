package area

import (
	"juego-websocket/game/ia"
	"juego-websocket/game/inter"
	"juego-websocket/game/item"
	"log"
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
	dummy := rand.Intn(cantItems + 1)
	log.Printf("Se crean %d Dummy NPC", dummy)
	for i := 0; i < dummy; i++ {
		ia := ia.NewDummyIA(i, a)
		ia.Start()
		items = append(items, ia.GetCharacter())
	}
	greed := cantItems - dummy
	log.Printf("Se crean %d Greedy NPC", greed)
	for i := 0; i < greed; i++ {
		ia := ia.NewGreedIA(i, a)
		ia.Start()
		items = append(items, ia.GetCharacter())
	}
	return items
}
