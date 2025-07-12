package main

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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
func MovePlayerOne(win *pixelgl.Window, paddle *Paddle) {
	if win.Pressed(pixelgl.KeyW) {
		paddle.MoveUp(5, 600)
	}
	if win.Pressed(pixelgl.KeyS) {
		paddle.MoveDown(5, 0)
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

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pong",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	const paddleWidth = 10.0
	const paddleHeight = 100.0
	score := Score{LeftScore: 0, RightScore: 0}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreText := text.New(pixel.V(350, 560), atlas)
	leftPaddle := NewPaddle(20, 300, paddleWidth, paddleHeight)
	rightPaddle := NewPaddle(770, 300, paddleWidth, paddleHeight)
	ball := NewBall(400, 300, 10)
	CPU := NewCPU(&rightPaddle, &ball)
	ball.Vel = pixel.V(-2.5, rand.Float64()*4-2)

	if err != nil {
		panic(err)
	}

	for !window.Closed() {
		window.Clear(colornames.Black)

		CalcScoreAndRespawn(&ball, &score)

		MovePlayerOne(window, &leftPaddle)

		CPU.Move()
		display := imdraw.New(nil)
		DrawPaddle(window, display, leftPaddle)
		DrawPaddle(window, display, rightPaddle)

		HandleWallCollision(&ball)
		HandlePaddleCollision(&ball, &leftPaddle)
		HandlePaddleCollision(&ball, &rightPaddle)

		ball.Move()
		DrawBall(window, display, &ball)

		scoreText.Clear()
		fmt.Fprintf(scoreText, "%d : %d", score.LeftScore, score.RightScore)
		scoreText.Draw(window, pixel.IM.Scaled(scoreText.Orig, 2))

		display.Draw(window)
		window.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
