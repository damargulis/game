package player

import (
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type MinimaxPlayer struct {
	Name     string
	MaxDepth int
}

func (p MinimaxPlayer) GetName() string {
	return p.Name
}

type moveVal struct {
	move int
	val  int
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func (p MinimaxPlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	ch := make(chan moveVal)
	scores := make([]int, len(moves))
	for i, move := range moves {
		go p.getScore(g, i, move, ch, 0)
	}
	i := 0
	for i < len(scores) {
		moveVal := <-ch
		i++
		scores[moveVal.move] = moveVal.val
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

func (p MinimaxPlayer) getScore(g game.Game, i int, m game.Move, ch chan moveVal, depth int) {
	if depth > p.MaxDepth {
		ch <- moveVal{move: i, val: g.CurrentScore(p)}
		return
	}
	newG := g.MakeMove(m)
	over, winner := newG.GameOver()
	if over {
		if winner == p {
			ch <- moveVal{move: i, val: MaxInt - depth}
		} else if winner.GetName() == "DRAW" {
			ch <- moveVal{move: i, val: 0}
		} else {
			ch <- moveVal{move: i, val: MinInt + depth}
		}
	} else {
		player := newG.GetPlayerTurn()
		if p == player {
			ch <- moveVal{move: i, val: p.getMax(newG, depth+1)}
		} else {
			ch <- moveVal{move: i, val: p.getMin(newG, depth+1)}
		}
	}
}

func (p MinimaxPlayer) getMax(g game.Game, depth int) int {
	moves := g.GetPossibleMoves()
	ch := make(chan moveVal)
	scores := make([]int, len(moves))
	for i, move := range moves {
		go p.getScore(g, i, move, ch, depth+len(moves))
	}
	i := 0
	for i < len(scores) {
		moveVal := <-ch
		i++
		scores[moveVal.move] = moveVal.val
	}
	bestScore := MinInt
	for _, score := range scores {
		if score >= bestScore {
			bestScore = score
		}
	}
	return bestScore
}

func (p MinimaxPlayer) getMin(g game.Game, depth int) int {
	moves := g.GetPossibleMoves()
	ch := make(chan moveVal)
	scores := make([]int, len(moves))
	for i, move := range moves {
		go p.getScore(g, i, move, ch, depth+len(moves))
	}
	i := 0
	for i < len(scores) {
		moveVal := <-ch
		i++
		scores[moveVal.move] = moveVal.val
	}
	bestScore := MaxInt
	for _, score := range scores {
		if score <= bestScore {
			bestScore = score
		}
	}
	return bestScore
}
