package player

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type MonteCarloTimePlayer struct {
	Name    string
	MaxTime int
}

func (p MonteCarloTimePlayer) GetName() string {
	return p.Name
}

func (p MonteCarloTimePlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	if len(moves) == 1 {
		return moves[0]
	}
	wins := make([]float64, len(moves))
	attempts := make([]int, len(moves))

	timer := make(chan int)
	result := make(chan float64)

	move := rand.Intn(len(moves))
	attempts[move]++
	newG := g.MakeMove(moves[move])
	go sleep(float64(p.MaxTime)*1000, timer)
	go p.runSimulation(newG, result, 0)
	iters := 0
	for {
		select {
		case r := <-result:
			iters++
			wins[move] += r
			attempts[move]++
			move = rand.Intn(len(moves))
			newG = g.MakeMove(moves[move])
			go p.runSimulation(newG, result, 0)
		case <-timer:
			fmt.Printf("Number iterations: %v\n", iters)
			bestScore := float64(MinInt)
			scores := make([]float64, len(moves))
			for i := range moves {
				if attempts[i] > 0 {
					scores[i] = wins[i] / float64(attempts[i])
				} else {
					scores[i] = 0
				}
				if scores[i] >= bestScore {
					bestScore = scores[i]
				}
			}
			var bestMoves []game.Move
			for i, score := range scores {
				if score == bestScore {
					bestMoves = append(bestMoves, moves[i])
				}
			}
			return bestMoves[rand.Intn(len(bestMoves))]
		}
	}
}

func (p MonteCarloTimePlayer) runSimulation(g game.Game, result chan float64, depth int) {
	over, winner := g.GameOver()
	if over {
		if winner == p {
			result <- float64(MaxInt-depth) / float64(MaxInt)
		} else if winner.GetName() == "DRAW" {
			result <- 0
		} else {
			result <- float64(MinInt+depth) / float64(MaxInt)
		}
	} else {
		moves := g.GetPossibleMoves()
		move := moves[rand.Intn(len(moves))]
		newG := g.MakeMove(move)
		p.runSimulation(newG, result, depth+1)
	}
}
