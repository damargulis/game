package main

import (
	"flag"
	"fmt"
	"github.com/damargulis/game/game"
	"github.com/damargulis/graph"
	"math/rand"
	"os"
	"time"
)

func writeOutput(fileName string, p1depth, p2depth, p1wins, p2wins, ties int, time time.Duration) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", p1depth, p2depth, p1wins, p2wins, ties, time)); err != nil {
		panic(err)
	}
}

func runExperiment(wrapper func(string, string, int, int) int, p1 string, p2 string, outputFile string, ch chan int, title string, graphFile string) {
	for i := 0; i < 500; i++ {
		for curMax := 0; curMax < 40; curMax++ {
			for curMin := 0; curMin < curMax; curMin++ {
				wins := [3]int{0, 0, 0}
				start := time.Now()
				win := wrapper(p1, p2, curMax, curMin)
				wins[win]++
				t := time.Now()
				elapsed := t.Sub(start)
				writeOutput(outputFile, curMax, curMin, wins[1], wins[2], wins[0], elapsed)

				wins = [3]int{0, 0, 0}
				start = time.Now()
				win = wrapper(p1, p2, curMin, curMax)
				wins[win]++
				t = time.Now()
				elapsed = t.Sub(start)
				writeOutput(outputFile, curMin, curMax, wins[1], wins[2], wins[0], elapsed)
			}
			wins := [3]int{0, 0, 0}
			start := time.Now()
			win := wrapper(p1, p2, curMax, curMax)
			wins[win]++
			t := time.Now()
			elapsed := t.Sub(start)
			writeOutput(outputFile, curMax, curMax, wins[1], wins[2], wins[0], elapsed)
			fmt.Printf("%v finished depth %v on iteration %v\n", title, curMax, i)
			if i == 0 {
				graph.CreateGraph(title, outputFile, graphFile)
			}
		}
		graph.CreateGraph(title, outputFile, graphFile)
	}
	ch <- 0
}

func main() {
	rand.Seed(time.Now().Unix())
	p1 := flag.String("p1", "Alphabeta", "Player 1")
	p2 := flag.String("p2", "Alphabeta", "Player 2")
	flag.Parse()
	ch := make(chan int)

	abaloneWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewAbalone(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	boxesWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewBoxes(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	checkersWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewCheckers(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	connect4Wrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewConnect4(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	mancalaWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewMancala(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	martianChessWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewMartianChess(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	nineMensMorrisWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewNineMensMorris(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	pentagoWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewPentago(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	reversiWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewReversi(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	tictactoeWrap := func(p1 string, p2 string, p1depth int, p2depth int) int {
		g := game.NewTicTacToe(p1, p2, p1depth, p2depth)
		return game.Play(g, false)
	}
	go runExperiment(abaloneWrap, *p1, *p2, "abalone.csv", ch, "Abalone", "abalone.png")
	go runExperiment(boxesWrap, *p1, *p2, "boxes.csv", ch, "Boxes", "boxes.png")
	go runExperiment(checkersWrap, *p1, *p2, "checkers.csv", ch, "Checkers", "checkers.png")
	go runExperiment(connect4Wrap, *p1, *p2, "connect4.csv", ch, "Connect4", "connect4.png")
	go runExperiment(mancalaWrap, *p1, *p2, "mancala.csv", ch, "Mancala", "mancala.png")
	go runExperiment(martianChessWrap, *p1, *p2, "martianchess.csv", ch, "Martian Chess", "martianchess.png")
	go runExperiment(nineMensMorrisWrap, *p1, *p2, "ninemensmorris.csv", ch, "Nine Men's Morris", "ninemensmorris.png")
	go runExperiment(pentagoWrap, *p1, *p2, "pentago.csv", ch, "Pentago", "pentago.png")
	go runExperiment(reversiWrap, *p1, *p2, "reversi.csv", ch, "Reversi", "reversi.png")
	go runExperiment(tictactoeWrap, *p1, *p2, "tictactoe.csv", ch, "Tic Tac Toe", "tictactoe.png")
	for i := 0; i < 10; i++ {
		<-ch
	}
}
