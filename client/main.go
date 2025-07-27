package main

//client
import (
	"flag"
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type GameState struct {
	PaddleP1 Paddle
	PaddleP2 Paddle
	Ball     Ball
	Score    Score
	Input    uint8
	Seq      uint64
	fromOp   bool
}

func run() {
	port := flag.Int("port", 30001, "local UDP port to bind to")
	player := flag.Int("player", 1, "Player number: 1 or 2")

	flag.Parse()
	isPlayerOne := *player == 1

	portNum := fmt.Sprintf(":%d", *port)
	fmt.Println(isPlayerOne)

	cfg := pixelgl.WindowConfig{
		Title:  "Pong",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	const paddleWidth = 10.0
	const paddleHeight = 100.0
	score := Score{LeftScore: 0, RightScore: 0}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreText := text.New(pixel.V(350, 560), atlas)
	leftPaddle := NewPaddle(20, 300, paddleWidth, paddleHeight)
	rightPaddle := NewPaddle(770, 300, paddleWidth, paddleHeight)
	ball := NewBall(400, 300, 10)

	gameState := GameState{
		PaddleP1: leftPaddle,
		PaddleP2: rightPaddle,
		Ball:     ball,
		Score:    score,
	}

	udpHandler := NewUdpHandler(portNum, &gameState, isPlayerOne)

	for !window.Closed() {
		window.Clear(colornames.Black)

		MovePlayer(window, udpHandler)

		udpHandler.SyncWithServerSeq(&gameState)

		display := imdraw.New(nil)
		DrawPaddle(window, display, gameState.PaddleP1)
		DrawPaddle(window, display, gameState.PaddleP2)
		DrawBall(window, display, &gameState.Ball)

		scoreText.Clear()
		fmt.Fprintf(scoreText, "%d : %d", gameState.Score.LeftScore, gameState.Score.RightScore)
		scoreText.Draw(window, pixel.IM.Scaled(scoreText.Orig, 2))

		display.Draw(window)
		window.Update()
	}
}

func main() {

	pixelgl.Run(run)

}
