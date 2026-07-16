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

type CoinData struct {
	Id string  `json:"playerId"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

func toJson(i interface{}) interface{} {
	switch i.(type) {
	case inter.Coin:
		c, _ := i.(inter.Coin)
		return coinToJson(c)
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

func coinToJson(c inter.Coin) CoinData {
	return CoinData{
		Id: c.GetId(),
		X:  c.GetPosition().GetX(),
		Y:  c.GetPosition().GetY(),
	}
}
