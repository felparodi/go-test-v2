package game

import (
	"sync"
)

type Game struct {
	mu     sync.RWMutex
	worlds map[string]*World
}
