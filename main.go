package main

import (
	"flag"
	"fmt"
	"github.com/damargulis/game/game"
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

func runExperiment(wrapper func(string, string, int, int) int, p1 string, p2 string, outputFile string, ch chan int) {
	for depth1 := 0; depth1 < 40; depth1++ {
		for depth2 := 0; depth2 < 40; depth2++ {
			wins := [3]int{0, 0, 0}
			start := time.Now()
			for i := 0; i < 100; i++ {
				win := wrapper(p1, p2, depth1, depth2)
				wins[win]++
			}
			t := time.Now()
			elapsed := t.Sub(start)
			writeOutput(outputFile, depth1, depth2, wins[1], wins[2], wins[0], elapsed)
		}
	}
	ch <- 0
}

func main() {
	rand.Seed(time.Now().Unix())
	p1 := flag.String("p1", "Human", "Player 1")
	p2 := flag.String("p2", "Computer", "Player 2")
	flag.Parse()
	boxesWrap := func(p1 string, p2 string, depth1 int, depth2 int) int {
		g := game.NewBoxes(p1, p2, depth1, depth2)
		return game.Play(g, "boxes.out")
	}
	ch := make(chan int)
	go runExperiment(boxesWrap, *p1, *p2, "boxes.csv", ch)
	checkersWrap := func(p1 string, p2 string, depth1 int, depth2 int) int {
		g := game.NewCheckers(p1, p2, depth1, depth2)
		return game.Play(g, "checkers.out")
	}
	go runExperiment(checkersWrap, *p1, *p2, "checkers.csv", ch)
	connect4Wrap := func(p1 string, p2 string, depth1 int, depth2 int) int {
		g := game.NewConnect4(p1, p2, depth1, depth2)
		return game.Play(g, "connect4.out")
	}
	go runExperiment(connect4Wrap, *p1, *p2, "connect4.csv", ch)
	reversiWrap := func(p1 string, p2 string, depth1 int, depth2 int) int {
		g := game.NewReversi(p1, p2, depth1, depth2)
		return game.Play(g, "reversi.out")
	}
	go runExperiment(reversiWrap, *p1, *p2, "reversi.csv", ch)
	tictactoeWrap := func(p1 string, p2 string, depth1 int, depth2 int) int {
		g := game.NewTicTacToe(p1, p2, depth1, depth2)
		return game.Play(g, "tictactoe.out")
	}
	go runExperiment(tictactoeWrap, *p1, *p2, "tictactoe.csv", ch)
	i := 0
	for i < 5 {
		<-ch
		i++
	}
	fmt.Println("Exiting")
}
