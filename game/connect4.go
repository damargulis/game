package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type Connect4 struct {
	board [8][8]string
	p1    game.Player
	p2    game.Player
	pTurn bool
}

type Connect4Move struct {
	col int
}

func (g Connect4) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g Connect4) BoardString() string {
	s := "-----------------\n"
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

func (g Connect4) PrintBoard() {
	fmt.Println("-----------------")
	for i, row := range g.board {
		fmt.Print(i)
		fmt.Print(" ")
		for _, p := range row {
			fmt.Print(p)
			fmt.Print(" ")
		}
		fmt.Println()
	}
	fmt.Println("  0 1 2 3 4 5 6 7")
	fmt.Println("-----------------")
}

func (g Connect4) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Connect4) GetHumanInput() game.Move {
	col := readInts("Column to move in: ")
	return Connect4Move{col: col[0]}
}

func (g Connect4) GetPossibleMoves() []game.Move {
	var moves []game.Move
	topRow := g.board[0]
	for i, spot := range topRow {
		if spot == "." {
			moves = append(moves, Connect4Move{col: i})
		}
	}
	return moves
}

func (g Connect4) GetTurn(p game.Player) game.Move {
	m := p.GetTurn(g)
	return m
}

func (g Connect4) MakeMove(m game.Move) game.Game {
	move := m.(Connect4Move)
	col := move.col
	i := 0
	for isInside(g, i, col) && g.board[i][col] == "." {
		i++
	}
	if g.pTurn {
		g.board[i-1][col] = "X"
	} else {
		g.board[i-1][col] = "O"
	}
	g.pTurn = !g.pTurn
	return g
}

func (g Connect4) checkMatch(i, j, rowDir, colDir int) bool {
	return isInside(g, i+rowDir*3, j+colDir+3) &&
		g.board[i][j] == g.board[i+rowDir][j+colDir] &&
		g.board[i][j] == g.board[i+rowDir*2][j+colDir*2] &&
		g.board[i][j] == g.board[i+rowDir*3][j+colDir*3]
}

func (g Connect4) GameOver() (bool, game.Player) {
	hasSpace := false
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				hasSpace = true
				continue
			}
			if g.checkMatch(i, j, 1, 0) {
				if g.board[i][j] == "X" {
					return true, g.p1
				} else {
					return true, g.p2
				}
			}
			if g.checkMatch(i, j, 0, 1) {
				if g.board[i][j] == "X" {
					return true, g.p1
				} else {
					return true, g.p2
				}
			}
			if g.checkMatch(i, j, 1, 1) {
				if g.board[i][j] == "X" {
					return true, g.p1
				} else {
					return true, g.p2
				}
			}
			if g.checkMatch(i, j, -1, 1) {
				if g.board[i][j] == "X" {
					return true, g.p1
				} else {
					return true, g.p2
				}
			}
		}
	}
	if hasSpace {
		return false, player.ComputerPlayer{}
	} else {
		return true, player.HumanPlayer{"DRAW"}
	}
}

func (g Connect4) CurrentScore(p game.Player) int {
	return 0
}

func NewConnect4(p1 string, p2 string, depth1 int, depth2 int) *Connect4 {
	c := new(Connect4)
	c.p1 = getPlayer(p1, "Player 1", depth1)
	c.p2 = getPlayer(p2, "Player 2", depth2)
	c.pTurn = true
	c.board = [8][8]string{
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
	}
	return c
}
