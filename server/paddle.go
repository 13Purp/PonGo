package main

import "github.com/faiface/pixel"

type Paddle struct {
	Pos    pixel.Vec
	Width  float64
	Height float64
}

func NewPaddle(x, y, width, height float64) Paddle {
	return Paddle{
		Pos:    pixel.V(x, y),
		Width:  width,
		Height: height,
	}
}

func (p *Paddle) MoveUp(dy, maxY float64) {
	newY := p.Pos.Y + dy
	if newY+p.Height > maxY {
		newY = maxY - p.Height
	}
	p.Pos.Y = newY
}

func (p *Paddle) MoveDown(dy, minY float64) {
	newY := p.Pos.Y - dy
	if newY < minY {
		newY = minY
	}
	p.Pos.Y = newY
}

// Rect returns the rectangle coordinates of the paddle for drawing/collision.
func (p *Paddle) Rect() (pixel.Vec, pixel.Vec) {
	return p.Pos, pixel.V(p.Pos.X+p.Width, p.Pos.Y+p.Height)
}
