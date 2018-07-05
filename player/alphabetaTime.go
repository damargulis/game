package player

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type AlphabetaTimePlayer struct {
	Name    string
	MaxTime int
}

func (p AlphabetaTimePlayer) GetName() string {
	return p.Name
}

func (p AlphabetaTimePlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	if len(moves) == 1 {
		return moves[0]
	}
	scores := make([]int, len(moves))
	timer := make(chan int)
	result := make(chan int)

	go sleep(float64(p.MaxTime)*1000, timer)
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
				return moves[move]
			} else if r == MinInt {
				moves = append(moves[:move], moves[move+1:]...)
				scores = append(scores[:move], scores[move+1:]...)
				if len(moves) == 1 {
					return moves[0]
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
			return bestMoves[rand.Intn(len(bestMoves))]
		}
	}
}

func (p AlphabetaTimePlayer) checkMove(g game.Game, m game.Move, maxDepth int, r chan int) {
	r <- p.getScore(g, m, 0, MinInt, MaxInt, maxDepth)
}

func (p AlphabetaTimePlayer) getScore(g game.Game, m game.Move, depth, alpha, beta, maxDepth int) int {
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

func (p AlphabetaTimePlayer) getMax(g game.Game, depth int, alpha int, beta int, maxDepth int) int {
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

func (p AlphabetaTimePlayer) getMin(g game.Game, depth int, alpha int, beta int, maxDepth int) int {
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
