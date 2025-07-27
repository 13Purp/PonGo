package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type UdpHandler struct {
	CmdChan     chan uint8
	localAddr   *net.UDPAddr
	serverAddr  *net.UDPAddr
	isPlayerOne bool
}

func NewUdpHandler(portNum string, gameState *GameState, isPlayerOne bool) *UdpHandler {
	localAddr, _ := net.ResolveUDPAddr("udp", portNum) // fixed port
	serverAddr, _ := net.ResolveUDPAddr("udp", "192.168.0.26:9999")
	handler := &UdpHandler{
		serverAddr:  serverAddr,
		localAddr:   localAddr,
		CmdChan:     make(chan uint8, 1),
		isPlayerOne: isPlayerOne,
	}

	//go handler.SyncWithServer(gameState)
	return handler
}

func (handler *UdpHandler) SyncWithServer(gameState *GameState) {

	conn, err := net.DialUDP("udp", handler.localAddr, handler.serverAddr)
	if err != nil {
		fmt.Println("UDP dial error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		input := <-handler.CmdChan

		gameState.Input = input
		gameState.fromOp = false
		data, err := json.Marshal(gameState)
		if err != nil {
			fmt.Println("Error marshaling GameState:", err)
			continue
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("UDP write error:", err)
			continue
		}

		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("UDP read error:", err)
			continue
		}

		newGameState := *gameState

		err = json.Unmarshal(buf[:n], &newGameState)
		if err != nil {
			fmt.Println("Failed to parse GameState:", err)
			continue
		}

		fmt.Println("newseqA", newGameState.Seq, "/", gameState.Seq)
		if newGameState.Seq > gameState.Seq {
			gameState.Ball = newGameState.Ball
			gameState.Input = 0
			fmt.Print("newseqT", newGameState.Seq)
			if (!newGameState.fromOp && handler.isPlayerOne) || (newGameState.fromOp && !handler.isPlayerOne) {
				//fmt.Println("P1 MOV", gameState.Seq)
				gameState.PaddleP1 = newGameState.PaddleP1
			}
			if (!newGameState.fromOp && !handler.isPlayerOne) || (newGameState.fromOp && handler.isPlayerOne) {
				fmt.Println("P2 MOV", gameState.Seq)
				gameState.PaddleP2 = newGameState.PaddleP2
			}
			gameState.Score = newGameState.Score
			gameState.Seq = newGameState.Seq
			gameState.fromOp = false
		}

	}

}
func (handler *UdpHandler) SyncWithServerSeq(gameState *GameState) {

	conn, err := net.DialUDP("udp", handler.localAddr, handler.serverAddr)
	if err != nil {
		fmt.Println("UDP dial error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 4096)

	input := <-handler.CmdChan

	gameState.Input = input
	gameState.fromOp = false
	if handler.isPlayerOne {
		gameState.PaddleP2 = Paddle{}
	} else {
		gameState.PaddleP1 = Paddle{}
	}
	data, err := json.Marshal(gameState)
	if err != nil {
		fmt.Println("Error marshaling GameState:", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("UDP write error:", err)

	}

	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("UDP read error:", err)

	}

	newGameState := *gameState

	err = json.Unmarshal(buf[:n], &newGameState)
	if err != nil {
		fmt.Println("Failed to parse GameState:", err)

	}

	fmt.Println("newseqA", newGameState.Seq, "/", gameState.Seq)
	if newGameState.Seq > gameState.Seq {
		gameState.Ball = newGameState.Ball
		gameState.Input = 0
		fmt.Print("newseqT", newGameState.Seq)
		gameState.PaddleP1 = newGameState.PaddleP1
		gameState.PaddleP2 = newGameState.PaddleP2
		gameState.Score = newGameState.Score
		gameState.Seq = newGameState.Seq
		gameState.fromOp = false
	}

}
