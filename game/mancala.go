package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type Mancala struct {
	board     [2][6]int
	p1        game.Player
	p2        game.Player
	p1Capture int
	p2Capture int
	pTurn     bool
}

type MancalaMove struct {
	row, col int
}

func (g Mancala) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func NewMancala(p1 string, p2 string, depth1 int, depth2 int) *Mancala {
	g := new(Mancala)
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.board = [2][6]int{
		{4, 4, 4, 4, 4, 4},
		{4, 4, 4, 4, 4, 4},
	}
	g.p1Capture = 0
	g.p2Capture = 0
	return g
}

func (g Mancala) BoardString() string {
	s := "-------------------\n"
	s += "   0  1  2  3  4  5\n"
	s += "0 "
	for _, p := range g.board[0] {
		if p < 10 {
			s += " "
		}
		s += fmt.Sprintf("%v", p)
		s += " "
	}
	s += "\n"
	s += fmt.Sprintf(" %v                 %v", g.p1Capture, g.p2Capture)
	s += "\n"
	s += "1 "
	for _, p := range g.board[1] {
		if p < 10 {
			s += " "
		}
		s += fmt.Sprintf("%v", p)
		s += " "
	}
	s += "\n"
	s += "   0  1  2  3  4  5\n"
	s += "-------------------\n"
	return s
}

func (g Mancala) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Mancala) GetHumanInput() game.Move {
	colA := readInts("Col to move: ")
	col := colA[0]
	if g.pTurn {
		return MancalaMove{
			row: 0,
			col: col,
		}
	} else {
		return MancalaMove{
			row: 1,
			col: col,
		}
	}
}

func (g Mancala) GetPossibleMoves() []game.Move {
	var row int
	var moves []game.Move
	if g.pTurn {
		row = 0
	} else {
		row = 1
	}
	for j, col := range g.board[row] {
		if col > 0 {
			moves = append(moves, MancalaMove{
				row: row,
				col: j,
			})
		}
	}
	return moves
}

func (g Mancala) MakeMove(m game.Move) game.Game {
	move := m.(MancalaMove)
	amtInHand := g.board[move.row][move.col]
	g.board[move.row][move.col] = 0
	curRow := move.row
	curCol := move.col
	for amtInHand > 0 {
		if curRow == 0 {
			if curCol == 0 {
				if g.pTurn {
					g.p1Capture++
					amtInHand--
					if amtInHand == 0 {
						return g
					} else {
						curRow = 1
					}
				} else {
					curRow = 1
				}
			} else {
				curCol--
			}
		} else if curRow == 1 {
			if curCol == 5 {
				if !g.pTurn {
					g.p2Capture++
					amtInHand--
					if amtInHand == 0 {
						return g
					} else {
						curRow = 0
					}
				} else {
					curRow = 0
				}
			} else {
				curCol++
			}
		} else {
			panic("Unexpected row")
		}
		g.board[curRow][curCol]++
		amtInHand--
	}
	if g.board[curRow][curCol] == 1 {
		if g.pTurn && curRow == 0 {
			g.p1Capture += 1 + g.board[1][curCol]
			g.board[0][curCol] = 0
			g.board[1][curCol] = 0
		} else if !g.pTurn && curRow == 1 {
			g.p2Capture += 1 + g.board[0][curCol]
			g.board[0][curCol] = 0
			g.board[1][curCol] = 0
		}
	}
	g.pTurn = !g.pTurn
	return g
}

func (g Mancala) GameOver() (bool, game.Player) {
	p1Alive := false
	p2Alive := false
	for i := range g.board[0] {
		if g.board[0][i] > 0 {
			p1Alive = true
		}
		if g.board[1][i] > 0 {
			p2Alive = true
		}
	}
	if p1Alive && p2Alive {
		return false, player.ComputerPlayer{}
	} else {
		for i := range g.board[0] {
			g.p1Capture += g.board[0][i]
			g.p2Capture += g.board[1][i]
			g.board[0][i] = 0
			g.board[1][i] = 0
		}
		if g.p1Capture > g.p2Capture {
			return true, g.p1
		} else if g.p2Capture > g.p1Capture {
			return true, g.p2
		} else {
			return true, player.HumanPlayer{"DRAW"}
		}
	}
}

func (g Mancala) CurrentScore(p game.Player) int {
	if p == g.p1 {
		return g.p1Capture - g.p2Capture
	} else {
		return g.p2Capture - g.p1Capture
	}
}
