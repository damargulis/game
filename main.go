package main

import (
	"flag"
	"fmt"
	"github.com/damargulis/game/game"
	"math/rand"
	"os"
	"time"
)

func writeOutput(fileName string, p1depth, p2depth, p1wins, p2wins, ties int) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v", p1depth, p2depth, p1wins, p2wins, ties)); err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	p1 := flag.String("p1", "Human", "Player 1")
	p2 := flag.String("p2", "Computer", "Player 2")
	flag.Parse()
	for depth1 := 0; depth1 <= 30; depth1++ {
		for depth2 := 0; depth2 <= 30; depth2++ {
			bwins := [3]int{0, 0, 0}
			twins := [3]int{0, 0, 0}
			c4wins := [3]int{0, 0, 0}
			rwins := [3]int{0, 0, 0}
			cwins := [3]int{0, 0, 0}
			for i := 0; i < 100; i++ {
				b := game.NewBoxes(*p1, *p2, depth1, depth2)
				bwin := game.Play(b)
				t := game.NewTicTacToe(*p1, *p2, depth1, depth2)
				twin := game.Play(t)
				c4 := game.NewConnect4(*p1, *p2, depth1, depth2)
				c4win := game.Play(c4)
				r := game.NewReversi(*p1, *p2, depth1, depth2)
				rwin := game.Play(r)
				c := game.NewCheckers(*p1, *p2, depth1, depth2)
				cwin := game.Play(c)
				twins[twin]++
				c4wins[c4win]++
				rwins[rwin]++
				cwins[cwin]++
				bwins[bwin]++
			}
			writeOutput("boxes.csv", depth1, depth2, bwins[1], bwins[2], bwins[0])
			writeOutput("checkers.csv", depth1, depth2, cwins[1], cwins[2], cwins[0])
			writeOutput("connect4.csv", depth1, depth2, c4wins[1], c4wins[2], c4wins[0])
			writeOutput("reversi.csv", depth1, depth2, rwins[1], rwins[2], rwins[0])
			writeOutput("tictactoe.csv", depth1, depth2, twins[1], twins[2], twins[0])
		}
	}
}
