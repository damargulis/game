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
	case "Alphabeta":
		p = player.AlphabetaPlayer{name, depth}
	default:
		fmt.Println("Player " + playerType + " not recognized")
		os.Exit(1)
	}
	return p
}

func printBoard(g game.Game, outputFile string) {
	s := g.BoardString()
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(s); err != nil {
		panic(err)
	}
}

func Play(g game.Game, gameFile string) int {
	var winner game.Player
	over := false
	for ; !over; over, winner = g.GameOver() {
		printBoard(g, gameFile)
		move := g.GetTurn(g.GetPlayerTurn())
		g = g.MakeMove(move)
	}
	printBoard(g, gameFile)
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
