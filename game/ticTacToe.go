package game

import (
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type TicTacToe struct {
	board [3][3]string
	p1    game.Player
	p2    game.Player
	pTurn bool
	round int
}

type TicTacToeMove struct {
	row, col int
}

func (g TicTacToe) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g TicTacToe) GetHumanInput() game.Move {
	spot := readInts("Spot to place: ")
	return TicTacToeMove{row: spot[0], col: spot[1]}
}

func (g TicTacToe) GameOver() (bool, game.Player) {
	if g.board[0][0] == g.board[0][1] && g.board[0][0] == g.board[0][2] {
		if g.board[0][0] == "X" {
			return true, g.p1
		} else if g.board[0][0] == "O" {
			return true, g.p2
		}
	}
	if g.board[1][0] == g.board[1][1] && g.board[1][0] == g.board[1][2] {
		if g.board[1][0] == "X" {
			return true, g.p1
		} else if g.board[1][0] == "O" {
			return true, g.p2
		}
	}
	if g.board[2][0] == g.board[2][1] && g.board[2][0] == g.board[2][2] {
		if g.board[2][0] == "X" {
			return true, g.p1
		} else if g.board[2][0] == "O" {
			return true, g.p2
		}
	}
	if g.board[0][0] == g.board[1][0] && g.board[0][0] == g.board[2][0] {
		if g.board[0][0] == "X" {
			return true, g.p1
		} else if g.board[0][0] == "O" {
			return true, g.p2
		}
	}
	if g.board[0][1] == g.board[1][1] && g.board[0][1] == g.board[2][1] {
		if g.board[0][1] == "X" {
			return true, g.p1
		} else if g.board[0][1] == "O" {
			return true, g.p2
		}
	}
	if g.board[0][2] == g.board[1][2] && g.board[0][2] == g.board[2][2] {
		if g.board[0][2] == "X" {
			return true, g.p1
		} else if g.board[0][2] == "O" {
			return true, g.p2
		}
	}
	if g.board[0][0] == g.board[1][1] && g.board[0][0] == g.board[2][2] {
		if g.board[0][0] == "X" {
			return true, g.p1
		} else if g.board[0][0] == "O" {
			return true, g.p2
		}
	}
	if g.board[2][0] == g.board[1][1] && g.board[1][1] == g.board[0][2] {
		if g.board[2][0] == "X" {
			return true, g.p1
		} else if g.board[2][0] == "O" {
			return true, g.p2
		}
	}
	for _, row := range g.board {
		for _, p := range row {
			if p == "." {
				return false, player.ComputerPlayer{}
			}
		}
	}
	return true, player.HumanPlayer{"DRAW"}
}

func NewTicTacToe(p1 string, p2 string, depth1 int, depth2 int) *TicTacToe {
	g := new(TicTacToe)
	g.round = 0
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.board = [3][3]string{
		{".", ".", "."},
		{".", ".", "."},
		{".", ".", "."},
	}
	return g
}

func (g TicTacToe) MakeMove(m game.Move) game.Game {
	g.round++
	move := m.(TicTacToeMove)
	row := move.row
	col := move.col
	if g.pTurn {
		g.board[row][col] = "X"
	} else {
		g.board[row][col] = "O"
	}
	g.pTurn = !g.pTurn
	return g
}

func (g TicTacToe) isGoodMove(m TicTacToeMove) bool {
	row := m.row
	col := m.col
	return isInside(g, row, col) && g.board[row][col] == "."
}

func (g TicTacToe) GetPossibleMoves() []game.Move {
	var moves []game.Move
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				moves = append(moves, TicTacToeMove{row: i, col: j})
			}
		}
	}
	return moves
}

func (g TicTacToe) CurrentScore(p game.Player) int {
	return 0
}

func (g TicTacToe) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g TicTacToe) BoardString() string {
	s := "---\n"
	for _, row := range g.board {
		for _, p := range row {
			s += p
		}
		s += "\n"
	}
	s += "---"
	return s
}

func (g TicTacToe) GetRound() int {
	return g.round
}
