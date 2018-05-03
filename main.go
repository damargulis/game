package main

import (
	"flag"
	"fmt"
	"github.com/damargulis/game/game"
	//	interfaces "github.com/damargulis/game/interfaces"
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

func runExperiment(wrapper func(string, string, int, int) int, p1 string, p2 string, outputFile string, title string, graphFile string) {
	for i := 0; i < 500; i++ {
		for curMax := 0; curMax < 40; curMax += 1 {
			for curMin := 0; curMin < curMax; curMin += 1 {
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
}

func main() {
	rand.Seed(time.Now().Unix())
	p1 := flag.String("p1", "Montecarlo", "Player 1")
	p2 := flag.String("p2", "Montecarlo", "Player 2")
	flag.Parse()

	ch := make(chan int)
	nineWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewNineMensMorris(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	pentagoWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewPentago(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	reversiWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewReversi(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	ticWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewTicTacToe(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	martianWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewMartianChess(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	mancalaWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewMancala(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	connectWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewConnect4(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	checkersWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewCheckers(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	boxesWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewBoxes(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	abaloneWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewAbalone(p1, p2, depth1, depth2)
		return game.Play(g, false)
	}
	go runExperiment(connectWrap, *p1, *p2, "connect4_mc.csv", "Connect 4", "connect_mc.png")
	go runExperiment(mancalaWrap, *p1, *p2, "mancala_mc.csv", "Mancala", "mancala_mc.png")
	go runExperiment(martianWrap, *p1, *p2, "martianchess_mc.csv", "Martian Chess", "martianchess_mc.png")
	go runExperiment(reversiWrap, *p1, *p2, "reversi_mc.csv", "Reversi", "reversi_mc.png")
	go runExperiment(pentagoWrap, *p1, *p2, "pentago_mc.csv", "Pentago", "pentago_mc.png")
	go runExperiment(nineWrap, *p1, *p2, "ninemensmorris_mc.csv", "Nine Men's Morris", "ninemensmorris_mc.png")
	go runExperiment(boxesWrap, *p1, *p2, "boxes_mc.csv", "Boxes", "boxes_mc.png")
	go runExperiment(checkersWrap, *p1, *p2, "checkers_mc.csv", "Checkers", "checkers_mc.png")
	go runExperiment(ticWrap, *p1, *p2, "tictactoe_mc.csv", "Tic Tac Toe", "tictactoe_mc.png")
	go runExperiment(abaloneWrap, *p1, *p2, "abalone_mc.csv", "Abalone", "abalone_mc.png")
	for i := 0; i < 10; i++ {
		<-ch
	}
	//wins := 0
	//games := [10]interfaces.Game{
	//	game.NewTicTacToe(*p1, *p2, 40, 0),
	//	game.NewReversi(*p1, *p2, 40, 25),
	//	game.NewPentago(*p1, *p2, 40, 50),
	//	game.NewNineMensMorris(*p1, *p2, 40, 35),
	//	game.NewMartianChess(*p1, *p2, 40, 55),
	//	game.NewMancala(*p1, *p2, 40, 40),
	//	game.NewConnect4(*p1, *p2, 40, 30),
	//	game.NewCheckers(*p1, *p2, 40, 30),
	//	game.NewBoxes(*p1, *p2, 40, 50),
	//	game.NewAbalone(*p1, *p2, 40, 60),
	//}
	//for _, g := range games {
	//	_ = game.Play(g, true)
	//}
	//fmt.Println("Total wins: ", wins)
}
