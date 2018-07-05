package player

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type ComboTimePlayer struct {
	Name    string
	MaxTime int
}

func (p ComboTimePlayer) GetName() string {
	return p.Name
}

func (p ComboTimePlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	if len(moves) == 1 {
		return moves[0]
	}
	moves = p.alphaStage(g, moves, float64(p.MaxTime)/float64(2)*1000)
	if len(moves) == 1 {
		return moves[0]
	}
	moves = p.montecarloStage(g, moves, float64(p.MaxTime)/float64(2)*1000)
	return moves[rand.Intn(len(moves))]
}

func (p ComboTimePlayer) montecarloStage(g game.Game, moves []game.Move, time float64) []game.Move {
	wins := make([]float64, len(moves))
	attempts := make([]int, len(moves))

	timer := make(chan int)
	result := make(chan float64)

	move := rand.Intn(len(moves))
	attempts[move]++
	newG := g.MakeMove(moves[move])
	go sleep(time, timer)
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
			fmt.Printf("Number of iterations: %v\n", iters)
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
			return bestMoves
		}
	}
}

func (p ComboTimePlayer) runSimulation(g game.Game, result chan float64, depth int) {
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

func (p ComboTimePlayer) alphaStage(g game.Game, moves []game.Move, time float64) []game.Move {
	scores := make([]int, len(moves))
	timer := make(chan int)
	result := make(chan int)

	go sleep(time, timer)
	maxDepth := 0
	move := 0
	go p.checkMove(g, moves[move], maxDepth, result)
	move++
	if move >= len(moves) {
		move = 0
		maxDepth++
	}
	for {
		select {
		case r := <-result:
			if r == MaxInt {
				return []game.Move{moves[move]}
			} else if r == MinInt {
				moves = append(moves[:move], moves[move+1:]...)
				scores = append(scores[:move], scores[move+1:]...)
				if len(moves) == 1 {
					return moves
				}
				if move >= len(moves) {
					move = 0
					maxDepth++
				}
				go p.checkMove(g, moves[move], maxDepth, result)
			} else {
				scores[move] = r
				move++
				if move >= len(moves) {
					move = 0
					maxDepth++
				}
				go p.checkMove(g, moves[move], maxDepth, result)
			}
		case <-timer:
			fmt.Printf("Max depth: %v\n", maxDepth)
			bestScore := MinInt
			for _, score := range scores {
				if score > bestScore {
					bestScore = score
				}
			}
			var bestMoves []game.Move
			for i, score := range scores {
				if score == bestScore {
					bestMoves = append(bestMoves, moves[i])
				}
			}
			return bestMoves
		}
	}
}

func (p ComboTimePlayer) checkMove(g game.Game, m game.Move, maxDepth int, r chan int) {
	r <- p.getScore(g, m, 0, MinInt, MaxInt, maxDepth)
}

func (p ComboTimePlayer) getScore(g game.Game, m game.Move, depth, alpha, beta, maxDepth int) int {
	if depth > maxDepth {
		return g.CurrentScore(p)
	}
	newG := g.MakeMove(m)
	over, winner := newG.GameOver()
	if over {
		if winner == p {
			return MaxInt
		} else if winner.GetName() == "DRAW" {
			return 0
		} else {
			return MinInt
		}
	} else {
		player := newG.GetPlayerTurn()
		if p == player {
			return p.getMax(newG, depth+1, alpha, beta, maxDepth)
		} else {
			return p.getMin(newG, depth+1, alpha, beta, maxDepth)
		}
	}
}

func (p ComboTimePlayer) getMax(g game.Game, depth int, alpha int, beta int, maxDepth int) int {
	moves := g.GetPossibleMoves()
	v := MinInt
	for _, move := range moves {
		score := p.getScore(g, move, depth, alpha, beta, maxDepth)
		if score > v {
			v = score
		}
		if v > alpha {
			alpha = v
		}
		if beta <= alpha {
			break
		}
	}
	return v
}

func (p ComboTimePlayer) getMin(g game.Game, depth int, alpha int, beta int, maxDepth int) int {
	moves := g.GetPossibleMoves()
	v := MaxInt
	for _, move := range moves {
		score := p.getScore(g, move, depth, alpha, beta, maxDepth)
		if score < v {
			v = score
		}
		if beta < v {
			beta = v
		}
		if beta <= alpha {
			break
		}
	}
	return v
}
