package main

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Score struct {
	LeftScore  int
	RightScore int
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
func MovePlayerOne(input uint8, paddle *Paddle) {
	if input == 1 {
		paddle.MoveUp(5, 600)
		fmt.Println("P1 UP")
	}
	if input == 2 {
		paddle.MoveDown(5, 0)
		fmt.Println("P1 DWN")

	}
}

func MovePlayerTwo(input uint8, paddle *Paddle) {
	if input == 1 {
		paddle.MoveUp(5, 600)
		fmt.Println("P2 UP")

	}
	if input == 2 {
		paddle.MoveDown(5, 0)
		fmt.Println("P2 DWN")

	}
}

func CalcScoreAndRespawn(ball *Ball, score *Score) {
	result := respawnBallIfOut(ball)
	if result == 1 {
		score.RightScore++
	}
	if result == -1 {
		score.LeftScore++
	}
}
func HandlePaddleCollision(ball *Ball, paddle *Paddle) {
	if !hasColided(ball, paddle) {
		return
	}
	var ballPos = ball.GetPos()
	ball.setVelX(ball.GetVel().X * -1)
	impactY := (ballPos.Y - (paddle.Pos.Y + paddle.Height/2)) / (paddle.Height / 2)
	ball.setVelY(impactY * 3)
}
func HandleWallCollision(ball *Ball) {
	if ball.Pos.Y-ball.Radius < 0 || ball.Pos.Y+ball.Radius > 600 {
		ball.Vel.Y *= -1
	}
}
func DrawPaddle(win *pixelgl.Window, imd *imdraw.IMDraw, paddle Paddle) {
	imd.Color = colornames.White
	p1Start, p1End := paddle.Rect()
	imd.Push(p1Start, p1End)
	imd.Rectangle(0)
}
func DrawBall(win *pixelgl.Window, imd *imdraw.IMDraw, ball *Ball) {
	imd.Color = colornames.White
	imd.Push(ball.GetPos())
	imd.Circle(ball.Radius, 0)
	imd.Draw(win)
}

func hasColided(ball *Ball, paddle *Paddle) bool {
	pMin, pMax := paddle.Rect()

	closestX := clamp(ball.GetPos().X, pMin.X, pMax.X)
	closestY := clamp(ball.GetPos().Y, pMin.Y, pMax.Y)

	dx := ball.Pos.X - closestX
	dy := ball.Pos.Y - closestY

	return dx*dx+dy*dy < ball.Radius*ball.Radius
}

func respawnBallIfOut(ball *Ball) int {
	if ball.GetPos().X < 0 {
		*ball = NewBall(400, 300, 10)
		ball.setVelX(2.5)
		ball.setVelY(rand.Float64()*4 - 2)
		return 1
	}
	if ball.GetPos().X > 800 {
		*ball = NewBall(400, 300, 10)
		ball.setVelX(-2.5)
		ball.setVelY(rand.Float64()*4 - 2)
		return -1
	}
	return 0
}
