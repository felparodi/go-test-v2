package game

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
	case *Coin:
		c, _ := i.(*Coin)
		return coinToJson(c)
	case *Player:
		p, _ := i.(*Player)
		return playerToJson(p)
	default:
		break
	}
	return nil
}

func playerToJson(p *Player) PlayerData {
	return PlayerData{
		Id:    p.ID,
		X:     p.Character.Position.X,
		Y:     p.Character.Position.Y,
		Vx:    p.Character.Velocity.X,
		Vy:    p.Character.Velocity.Y,
		Score: p.Character.Score,
		Angle: p.Character.Position.Angle,
	}
}

func coinToJson(c *Coin) CoinData {
	return CoinData{
		Id: c.ID,
		X:  c.pos.X,
		Y:  c.pos.Y,
	}
}
