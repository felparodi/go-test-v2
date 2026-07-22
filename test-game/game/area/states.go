package area

import "juego-websocket/game/inter"

type AreaState struct {
	Items      []inter.Item
	Characters []inter.Character
}

func (ws *AreaState) GetItems() []inter.Item {
	return ws.Items
}

func (ws *AreaState) GetCharacters() []inter.Character {
	return ws.Characters
}
