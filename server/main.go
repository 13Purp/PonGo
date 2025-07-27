package main

//server
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"

	"github.com/faiface/pixel"
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
type ClientInfo struct {
	Addr       *net.UDPAddr
	IsPlayer1  bool
	PairedWith string
	Seq        uint64
}

var SEQ uint64

func getOpAdress(clientPairs map[string]ClientInfo, clientAddr *net.UDPAddr) *net.UDPAddr {
	return clientPairs[clientPairs[clientAddr.String()].PairedWith].Addr
}

func calcNextState(gameState *GameState, IsPlayer1 bool) {

	if gameState.Ball.GetVel().X == 0 && gameState.Ball.GetVel().Y == 0 {
		gameState.Ball.Vel = pixel.V(-2.5, rand.Float64()*4-2)

	}

	CalcScoreAndRespawn(&gameState.Ball, &gameState.Score)

	if IsPlayer1 {
		MovePlayerOne(gameState.Input, &gameState.PaddleP1)
		fmt.Print(" ", gameState.Seq)

	} else {
		MovePlayerTwo(gameState.Input, &gameState.PaddleP2)
		fmt.Print(" ", gameState.Seq)

	}

	HandleWallCollision(&gameState.Ball)
	HandlePaddleCollision(&gameState.Ball, &gameState.PaddleP1)
	HandlePaddleCollision(&gameState.Ball, &gameState.PaddleP2)

	gameState.Ball.Move()
	SEQ++
	gameState.Seq = SEQ

}

func pairIncomingClients(clientPairs map[string]ClientInfo, clientAddr *net.UDPAddr) bool {
	addrStr := clientAddr.String()
	if _, ok := clientPairs[addrStr]; !ok {
		paired := false
		for otherAddrStr, info := range clientPairs {
			if info.PairedWith == "" {
				clientPairs[addrStr] = ClientInfo{
					Addr:       clientAddr,
					IsPlayer1:  false,
					PairedWith: otherAddrStr,
				}

				info.PairedWith = addrStr
				clientPairs[otherAddrStr] = info

				fmt.Println("Paired", addrStr, "with", otherAddrStr)
				paired = true
				break
			}
		}
		if !paired {
			clientPairs[addrStr] = ClientInfo{
				Addr:       clientAddr,
				IsPlayer1:  true,
				PairedWith: "",
			}
			fmt.Println("Waiting to pair", addrStr)
		}

		return paired
	}

	return true

}

func main() {

	SEQ = 0

	addr, _ := net.ResolveUDPAddr("udp", ":9999")
	conn, _ := net.ListenUDP("udp", addr)
	buf := make([]byte, 4096)
	clientPairs := make(map[string]ClientInfo)

	for {

		n, clientAddr, _ := conn.ReadFromUDP(buf)

		paired := pairIncomingClients(clientPairs, clientAddr)

		var gameState GameState

		err := json.Unmarshal(buf[:n], &gameState)
		if err != nil {
			fmt.Println("Failed to parse GameState:", err)
		}

		if paired {
			calcNextState(&gameState, clientPairs[clientAddr.String()].IsPlayer1)
		}

		gameState.fromOp = false
		data, err := json.Marshal(gameState)
		if err != nil {
			fmt.Println("Error marshaling GameState:", err)
			return
		}

		conn.WriteToUDP(data, clientAddr)
		gameState.fromOp = true
		SEQ++
		gameState.Seq = SEQ
		data, err = json.Marshal(gameState)
		if err != nil {
			fmt.Println("Error marshaling GameState:", err)
			return
		}
		if paired {
			opAdress := getOpAdress(clientPairs, clientAddr)
			conn.WriteToUDP(data, opAdress)
		}
	}
}
