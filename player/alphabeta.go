package player

import (
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type AlphabetaPlayer struct {
	Name     string
	MaxDepth int
}

func (p AlphabetaPlayer) GetName() string {
	return p.Name
}

func (p AlphabetaPlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	scores := make([]int, len(moves))
	v := MinInt
	alpha := MinInt
	beta := MaxInt

	if len(moves) > p.MaxDepth {
		return moves[rand.Intn(len(moves))]
	}

	for i, move := range moves {
		score := p.getScore(g, move, len(moves)+1, alpha, beta)
		scores[i] = score
		if score > v {
			v = score
		}
		if v > alpha {
			alpha = v
		}
		if beta <= alpha {
			return move
		}
	}
	bestScore := MinInt
	for _, score := range scores {
		if score >= bestScore {
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

func (p AlphabetaPlayer) getScore(g game.Game, m game.Move, depth int, alpha int, beta int) int {
	if depth > p.MaxDepth {
		return g.CurrentScore(p)
	}
	newG := g.MakeMove(m)
	over, winner := newG.GameOver()
	if over {
		if winner == p {
			return MaxInt - depth
		} else if winner.GetName() == "DRAW" {
			return 0
		} else {
			return MinInt + depth
		}
	} else {
		player := newG.GetPlayerTurn()
		if p == player {
			return p.getMax(newG, depth+1, alpha, beta)
		} else {
			return p.getMin(newG, depth+1, alpha, beta)
		}
	}
}

func (p AlphabetaPlayer) getMax(g game.Game, depth int, alpha int, beta int) int {
	moves := g.GetPossibleMoves()
	v := MinInt
	for _, move := range moves {
		score := p.getScore(g, move, depth+len(moves), alpha, beta)
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

func (p AlphabetaPlayer) getMin(g game.Game, depth int, alpha int, beta int) int {
	moves := g.GetPossibleMoves()
	v := MaxInt
	for _, move := range moves {
		score := p.getScore(g, move, depth+len(moves), alpha, beta)
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
