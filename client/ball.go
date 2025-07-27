package main

import (
	"sync"

	"github.com/faiface/pixel"
)

type Ball struct {
	Pos    pixel.Vec
	Radius float64
	Vel    pixel.Vec
	mu     sync.Mutex `json:"-"`
}

func NewBall(x, y, radius float64) Ball {
	return Ball{
		Pos:    pixel.V(x, y),
		Radius: radius,
		Vel:    pixel.V(0, 0),
	}
}
func (b *Ball) GetPos() pixel.Vec {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Pos
}
func (b *Ball) setVelX(x float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Vel.X = x
}
func (b *Ball) setVelY(y float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Vel.Y = y
}
func (b *Ball) GetVel() pixel.Vec {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Vel
}

func (ball *Ball) Move() {
	ball.mu.Lock()
	defer ball.mu.Unlock()

	ball.Pos.X += ball.Vel.X
	ball.Pos.Y += ball.Vel.Y
}
