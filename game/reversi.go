package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type Reversi struct {
	board [8][8]string
	p1    game.Player
	p2    game.Player
	pTurn bool
}

type ReversiMove struct {
	row, col int
}

func (g Reversi) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g Reversi) BoardString() string {
	s := "-----------------\n"
	s += "  0 1 2 3 4 5 6 7\n"
	for i, row := range g.board {
		s += fmt.Sprintf("%v ", i)
		for _, p := range row {
			s += p
			s += " "
		}
		s += "\n"
	}
	s += "  0 1 2 3 4 5 6 7\n"
	s += "-----------------\n"
	return s
}

func (g Reversi) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Reversi) GetHumanInput() game.Move {
	spot := readInts("Spot to place: ")
	return ReversiMove{row: spot[0], col: spot[1]}
}

func (g Reversi) checkMove(i, j, rowDir, colDir int) bool {
	var target, match string
	if g.pTurn {
		target, match = "O", "X"
	} else {
		target, match = "X", "O"
	}
	if isInside(g, i+rowDir, j+colDir) && g.board[i+rowDir][j+colDir] == target {
		rowCheck := i + rowDir
		colCheck := j + colDir
		for isInside(g, rowCheck, colCheck) && g.board[rowCheck][colCheck] == target {
			rowCheck += rowDir
			colCheck += colDir
		}
		if isInside(g, rowCheck, colCheck) && g.board[rowCheck][colCheck] == match {
			return true
		}
	}
	return false
}

func (g Reversi) GetPossibleMoves() []game.Move {
	var moves []game.Move
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				if g.checkMove(i, j, 0, 1) ||
					g.checkMove(i, j, 1, 1) ||
					g.checkMove(i, j, 1, 0) ||
					g.checkMove(i, j, 1, 0) ||
					g.checkMove(i, j, 1, -1) ||
					g.checkMove(i, j, 0, -1) ||
					g.checkMove(i, j, -1, -1) ||
					g.checkMove(i, j, -1, 0) ||
					g.checkMove(i, j, -1, 1) {
					moves = append(moves, ReversiMove{
						row: i,
						col: j,
					})
				}
			}
		}
	}
	return moves
}

func (g Reversi) checkAndFill(i, j, rowDir, colDir int) Reversi {
	var target, match string
	if g.pTurn {
		target, match = "O", "X"
	} else {
		target, match = "X", "O"
	}
	if isInside(g, i+rowDir, j+colDir) && g.board[i+rowDir][j+colDir] == target {
		rowCheck := i + rowDir
		colCheck := j + colDir
		for isInside(g, rowCheck, colCheck) && g.board[rowCheck][colCheck] == target {
			rowCheck += rowDir
			colCheck += colDir
		}
		if isInside(g, rowCheck, colCheck) && g.board[rowCheck][colCheck] == match {
			for r, c := i, j; r != rowCheck || c != colCheck; r, c = r+rowDir, c+colDir {
				g.board[r][c] = match
			}
		}
	}
	return g
}

func (g Reversi) MakeMove(m game.Move) game.Game {
	move := m.(ReversiMove)
	var match string
	if g.pTurn {
		match = "X"
	} else {
		match = "O"
	}
	g.board[move.row][move.col] = match
	g = g.checkAndFill(move.row, move.col, 0, 1)
	g = g.checkAndFill(move.row, move.col, 1, 1)
	g = g.checkAndFill(move.row, move.col, 1, 0)
	g = g.checkAndFill(move.row, move.col, 1, -1)
	g = g.checkAndFill(move.row, move.col, 0, -1)
	g = g.checkAndFill(move.row, move.col, -1, -1)
	g = g.checkAndFill(move.row, move.col, -1, 0)
	g = g.checkAndFill(move.row, move.col, -1, 1)
	g.pTurn = !g.pTurn
	possibleMoves := g.GetPossibleMoves()
	if len(possibleMoves) == 0 {
		g.pTurn = !g.pTurn
	}
	return g
}

func (g Reversi) GameOver() (bool, game.Player) {
	possibleMoves := g.GetPossibleMoves()
	if len(possibleMoves) == 0 {
		score := g.CurrentScore(g.p1)
		if score > 0 {
			return true, g.p1
		} else if score < 0 {
			return true, g.p2
		} else {
			return true, player.HumanPlayer{"DRAW"}
		}
	} else {
		return false, player.ComputerPlayer{}
	}
}

func (g Reversi) CurrentScore(p game.Player) int {
	score := 0
	for _, row := range g.board {
		for _, spot := range row {
			if spot == "X" {
				score++
			} else if spot == "O" {
				score--
			}
		}
	}
	if p == g.p1 {
		return score
	} else {
		return -1 * score
	}
}

func NewReversi(p1 string, p2 string, depth1 int, depth2 int) *Reversi {
	r := new(Reversi)
	r.p1 = getPlayer(p1, "Player 1", depth1)
	r.p2 = getPlayer(p2, "Player 2", depth2)
	r.pTurn = true
	r.board = [8][8]string{
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", "X", "O", ".", ".", "."},
		{".", ".", ".", "O", "X", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
	}
	return r
}
