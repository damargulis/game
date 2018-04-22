package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
	"os"
)

func getPlayer(playerType string, name string, depth int) game.Player {
	var p game.Player
	switch playerType {
	case "Human":
		p = player.HumanPlayer{name}
	case "Computer":
		p = player.ComputerPlayer{name}
	case "Minimax":
		p = player.MinimaxPlayer{name, depth}
	default:
		fmt.Println("Player " + playerType + " not recognized")
		os.Exit(1)
	}
	return p
}

func Play(g game.Game) int {
	var winner game.Player
	over := false
	for ; !over; over, winner = g.GameOver() {
		g.PrintBoard()
		move := g.GetTurn(g.GetPlayerTurn())
		g = g.MakeMove(move)
	}
	g.PrintBoard()
	name := winner.GetName()
	if name == "DRAW" {
		fmt.Println("Its a draw!")
		return 0
	} else {
		fmt.Println(name + " Wins!")
		if name == "Player 1" {
			return 1
		} else {
			return 2
		}
	}
}
