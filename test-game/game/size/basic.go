package size

import (
	"juego-websocket/game/inter"
	"juego-websocket/game/position"
	"log"
	"math/rand"
)

type BasicSize struct {
	Height  float64
	Width   float64
	Padding float64
}

func NewSize(h, w float64) inter.Size {
	return &BasicSize{
		Height:  h,
		Width:   w,
		Padding: 0,
	}
}

func NewSizePadding(h, w, p float64) inter.Size {
	return &BasicSize{
		Height:  h,
		Width:   w,
		Padding: p,
	}
}

func (s *BasicSize) GetHeight() float64 {
	return s.Height
}

func (s *BasicSize) GetWidth() float64 {
	return s.Width
}

func (s *BasicSize) GetMaxWidth() float64 {
	return s.Width - s.Padding
}

func (s *BasicSize) GetMinWidth() float64 {
	return s.Padding
}

func (s *BasicSize) GetMaxHeight() float64 {
	return s.Height - s.Padding
}

func (s *BasicSize) GetMinHeight() float64 {
	return s.Padding
}

func (s *BasicSize) Copy() inter.Size {
	return &BasicSize{
		Height:  s.Height,
		Width:   s.Width,
		Padding: s.Padding,
	}
}

func NormalizeMove(x, y float64, pos inter.Position, size inter.Size) (float64, float64, []string) {
	limits := []string{}
	if pos.GetX()+x < size.GetMinWidth() {
		limits = append(limits, "limit-min-x")
		x = 0
	}
	if pos.GetY()+y < size.GetMinHeight() {
		limits = append(limits, "limit-min-y")
		y = 0
	}
	if pos.GetX()+x > size.GetMaxWidth() {
		limits = append(limits, "limit-max-x")
		x = 0
	}
	if pos.GetY()+y > size.GetMaxHeight() {
		limits = append(limits, "limit-max-y")
		y = 0
	}
	return x, y, limits
}

func (s *BasicSize) GetRandPosistion() inter.Position {
	x := float64(rand.Intn(int(s.GetMaxWidth()-s.GetMinWidth())) + int(s.GetMinWidth()))
	y := float64(rand.Intn(int(s.GetMaxHeight()-s.GetMinHeight())) + int(s.GetMinHeight()))
	angle := float64(rand.Intn(8) * 45)
	log.Println("RP", s, x, y, angle)
	return position.NewPosition(x, y, angle)
}
