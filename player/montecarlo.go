package player

import (
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type MonteCarloPlayer struct {
	Name    string
	MaxSims int
}

func (p MonteCarloPlayer) GetName() string {
	return p.Name
}

func (p MonteCarloPlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	wins := make([]int, len(moves))
	attempts := make([]int, len(moves))
	for i := 0; i < p.MaxSims; i++ {
		move := rand.Intn(len(moves))
		attempts[move]++
		newG := g.MakeMove(moves[move])
		winner := p.runSimulation(newG)
		if winner == p {
			wins[move] += 2
		} else if winner.GetName() == "DRAW" {
			wins[move] += 1
		}
	}
	bestScore := 0.0
	scores := make([]float64, len(moves))
	for i := range moves {
		if attempts[i] > 0 {
			scores[i] = float64(wins[i]) / float64(attempts[i])
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

func (p MonteCarloPlayer) runSimulation(g game.Game) game.Player {
	over, winner := g.GameOver()
	if over {
		return winner
	} else {
		moves := g.GetPossibleMoves()
		move := moves[rand.Intn(len(moves))]
		newG := g.MakeMove(move)
		return p.runSimulation(newG)
	}
}
