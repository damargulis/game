package main

import (
	"flag"
	"fmt"
	"github.com/damargulis/game/game"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	p1 := flag.String("p1", "Human", "Player 1")
	p2 := flag.String("p2", "Computer", "Player 2")
	flag.Parse()
	for i := 0; i < 10; i++ {
		fmt.Println("Game ", i)
		time.Sleep(500000000)
		fmt.Println(*p1, *p2)
		r := game.NewReversi(*p1, *p2)
		game.Play(r)
		t := game.NewTicTacToe(*p1, *p2)
		game.Play(t)
		c := game.NewCheckers(*p1, *p2)
		game.Play(c)
	}
}
