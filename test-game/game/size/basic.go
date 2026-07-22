package size

import "juego-websocket/game/inter"

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

func NormalizeMove(x, y float64, pos inter.Position, size inter.Size) (float64, float64) {
	if pos.GetX()+x < size.GetMinWidth() {
		x = 0
	}
	if pos.GetY()+y < size.GetMinHeight() {
		y = 0
	}
	if pos.GetX()+x > size.GetMaxWidth() {
		x = 0
	}
	if pos.GetY()+y > size.GetMaxHeight() {
		y = 0
	}
	return x, y
}
