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
	//ch := make(chan int)

	g := game.NewCheckers(*p1, *p2, 10, 10)
	_ = game.Play(g, true)
}
