package game

import (
	"juego-websocket/game/inter"
)

type PlayerData struct {
	Id    string  `json:"playerId"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Vx    float64 `json:"vx"`
	Vy    float64 `json:"vy"`
	Score int     `json:"score"`
	Angle float64 `json:"angle"`
}

type ItemData struct {
	Id   string  `json:"id"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Type string  `json:"type"`
}

func toJson(i interface{}) interface{} {
	switch i.(type) {
	case inter.Item:
		it, _ := i.(inter.Item)
		switch it.(type) {
		case inter.Bullet:
			return itemToJson(it, "bullet")
		case inter.Coin:
			return itemToJson(it, "coin")
		}
	case inter.Player:
		p, _ := i.(inter.Player)
		return playerToJson(p)
	default:
		break
	}
	return nil
}

func playerToJson(p inter.Player) PlayerData {
	return PlayerData{
		Id:    p.GetId(),
		X:     p.GetCharacter().GetPosition().GetX(),
		Y:     p.GetCharacter().GetPosition().GetY(),
		Vx:    p.GetCharacter().GetVelocity().GetX(),
		Vy:    p.GetCharacter().GetVelocity().GetY(),
		Score: p.GetCharacter().GetScore(),
		Angle: p.GetCharacter().GetPosition().GetAngle(),
	}
}

func itemToJson(c inter.Item, t string) ItemData {
	return ItemData{
		Id:   c.GetId(),
		X:    c.GetPosition().GetX(),
		Y:    c.GetPosition().GetY(),
		Type: t,
	}
}
