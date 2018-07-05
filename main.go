package main

import (
	"fmt"
	"github.com/damargulis/game/game"
	//	interfaces "github.com/damargulis/game/interfaces"
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

func runExperiment(wrapper func(string, string, int, int) int, p1 string, p2 string, outputFile string) {
	fmt.Println("Running Experiment", outputFile, "with players", p1, p2)
	outputFile = outputFile + ".csv"
	for curMax := 0; curMax < 40; curMax += 2 {
		for curMin := 0; curMin < curMax; curMin += 2 {
			wins := [3]int{0, 0, 0}
			start := time.Now()
			for i := 0; i < 100; i++ {
				win := wrapper(p1, p2, curMax, curMin)
				wins[win]++
			}
			t := time.Now()
			elapsed := t.Sub(start)
			writeOutput(outputFile, curMax, curMin, wins[1], wins[2], wins[0], elapsed)
			wins = [3]int{0, 0, 0}
			start = time.Now()
			for i := 0; i < 100; i++ {
				win := wrapper(p1, p2, curMin, curMax)
				wins[win]++
			}
			t = time.Now()
			elapsed = t.Sub(start)
			writeOutput(outputFile, curMin, curMax, wins[1], wins[2], wins[0], elapsed)
		}
		wins := [3]int{0, 0, 0}
		start := time.Now()
		for i := 0; i < 100; i++ {
			win := wrapper(p1, p2, curMax, curMax)
			wins[win]++
		}
		t := time.Now()
		elapsed := t.Sub(start)
		writeOutput(outputFile, curMax, curMax, wins[1], wins[2], wins[0], elapsed)
		fmt.Println("Finished depth", curMax)
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	nineWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewNineMensMorris(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	pentagoWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewPentago(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	reversiWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewReversi(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	ticWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewTicTacToe(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	martianWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewMartianChess(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	mancalaWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewMancala(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	connectWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewConnect4(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	checkersWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewCheckers(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	boxesWrap := func(p1, p2 string, depth1, depth2 int) int {
		g := game.NewBoxes(p1, p2, depth1, depth2)
		return game.Play(g, true)
	}
	//abaloneWrap := func(p1, p2 string, depth1, depth2 int) int {
	//	g := game.NewAbalone(p1, p2, depth1, depth2)
	//	return game.Play(g, true)
	//}
	wraps := []func(string, string, int, int) int{
		ticWrap,
		pentagoWrap,
		mancalaWrap,
		boxesWrap,
		connectWrap,
		reversiWrap,
		checkersWrap,
		martianWrap,
		nineWrap,
		//	abaloneWrap,
	}
	names := []string{
		"tictactoe",
		"pentago",
		"mancala",
		"boxes",
		"connect4",
		"reversi",
		"checkers",
		"martianchess",
		"ninemensmorris",
		//	"abalone",
	}
	for i, wrap := range wraps {
		runExperiment(wrap, "Montecarlo", "Montecarlo", names[i])
	}
}
