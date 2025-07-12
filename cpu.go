package main

import (
	"math/rand"
	"time"
)

type CPU struct {
	paddle   *Paddle
	ball     *Ball
	cmdChan  chan CPUCommand
	delayMin time.Duration
	delayMax time.Duration
}

func NewCPU(paddle *Paddle, ball *Ball) *CPU {
	cpu := &CPU{
		paddle:   paddle,
		ball:     ball,
		cmdChan:  make(chan CPUCommand, 10),
		delayMin: 30 * time.Millisecond,
		delayMax: 60 * time.Millisecond,
	}

	go cpu.loop()
	return cpu
}

func (cpu *CPU) Move() {
	select {
	case cmd := <-cpu.cmdChan:
		switch cmd {
		case MoveUp:
			cpu.paddle.MoveUp(5, 600)
		case MoveDown:
			cpu.paddle.MoveDown(5, 0)
		}
	default:
	}
}

type CPUCommand int

const (
	NoOp CPUCommand = iota
	MoveUp
	MoveDown
)

func (cpu *CPU) loop() {
	for {
		delay := time.Duration(rand.Intn(int(cpu.delayMax-cpu.delayMin))) + cpu.delayMin
		time.Sleep(delay)

		vel := cpu.ball.GetVel()
		if vel.X < 0 {
			cpu.cmdChan <- NoOp
			continue
		}

		pos := cpu.ball.GetPos()
		numberOfCommands := rand.Intn(5) + 1
		if pos.Y < cpu.paddle.Pos.Y {
			for range numberOfCommands {
				cpu.cmdChan <- MoveDown
			}
		} else if pos.Y > cpu.paddle.Pos.Y {
			for range numberOfCommands {
				cpu.cmdChan <- MoveUp
			}
		} else {
			for range numberOfCommands {
				cpu.cmdChan <- NoOp
			}
		}
	}
}
