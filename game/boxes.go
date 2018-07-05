package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type Boxes struct {
	board [9][9]string
	p1    game.Player
	p2    game.Player
	pTurn bool
	round int
}

type BoxesMove struct {
	row, col int
}

func (g Boxes) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g Boxes) BoardString() string {
	s := "*******************\n"
	s += "  0 1 2 3 4 5 6 7 8\n"
	for i, row := range g.board {
		s += fmt.Sprintf("%v ", i)
		for _, spot := range row {
			s += spot
			s += " "
		}
		s += "\n"
	}
	s += "  0 1 2 3 4 5 6 7 8\n"
	s += "*******************"
	return s
}

func (g Boxes) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Boxes) GetHumanInput() game.Move {
	spot := readInts("Place a line at: ")
	rowI, colI := spot[0], spot[1]
	return BoxesMove{
		row: rowI,
		col: colI,
	}
}

func (g Boxes) GetPossibleMoves() []game.Move {
	var moves []game.Move
	for i, row := range g.board {
		for j, spot := range row {
			if (i+j)%2 != 0 && spot == " " {
				moves = append(moves, BoxesMove{
					row: i,
					col: j,
				})
			}
		}
	}
	return moves
}

func (g Boxes) checkSpot(i, j int) bool {
	return isInside(g, i+1, j) &&
		isInside(g, i-1, j) &&
		isInside(g, i, j+1) &&
		isInside(g, i, j-1) &&
		g.board[i+1][j] != " " &&
		g.board[i-1][j] != " " &&
		g.board[i][j+1] != " " &&
		g.board[i][j-1] != " "
}

func (g Boxes) MakeMove(move game.Move) game.Game {
	g.round++
	didClaim := false
	m := move.(BoxesMove)
	var own string
	if g.pTurn {
		own = "X"
	} else {
		own = "O"
	}
	if m.row%2 == 0 {
		g.board[m.row][m.col] = "-"
		if g.checkSpot(m.row-1, m.col) {
			g.board[m.row-1][m.col] = own
			didClaim = true
		}
		if g.checkSpot(m.row+1, m.col) {
			g.board[m.row+1][m.col] = own
			didClaim = true
		}
	} else {
		g.board[m.row][m.col] = "|"
		if g.checkSpot(m.row, m.col-1) {
			g.board[m.row][m.col-1] = own
			didClaim = true
		}
		if g.checkSpot(m.row, m.col+1) {
			g.board[m.row][m.col+1] = own
			didClaim = true
		}
	}
	if !didClaim {
		g.pTurn = !g.pTurn
	}
	return g
}

func (g Boxes) GameOver() (bool, game.Player) {
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

func (g Boxes) CurrentScore(p game.Player) int {
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

func NewBoxes(p1 string, p2 string, depth1 int, depth2 int) *Boxes {
	g := new(Boxes)
	g.round = 0
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.board = [9][9]string{
		{".", " ", ".", " ", ".", " ", ".", " ", "."},
		{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		{".", " ", ".", " ", ".", " ", ".", " ", "."},
		{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		{".", " ", ".", " ", ".", " ", ".", " ", "."},
		{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		{".", " ", ".", " ", ".", " ", ".", " ", "."},
		{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		{".", " ", ".", " ", ".", " ", ".", " ", "."},
	}
	return g
}

func (g Boxes) GetRound() int {
	return g.round
}
