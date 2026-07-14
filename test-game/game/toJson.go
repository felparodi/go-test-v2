package game

type PlayerData struct {
	Id    string  `json:"playerId"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Score int     `json:"score"`
	Angle float64 `json:"angle"`
}

type CoinData struct {
	Id string  `json:"playerId"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

func itemToJson(i Item) interface{} {
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
		X:     p.Position.X,
		Y:     p.Position.Y,
		Score: p.Score,
		Angle: p.Position.Angle,
	}
}

func coinToJson(c *Coin) CoinData {
	return CoinData{
		Id: c.ID,
		X:  c.pos.X,
		Y:  c.pos.Y,
	}
}
