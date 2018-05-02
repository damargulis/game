package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type Checkers struct {
	board              [8][8]string
	p1                 game.Player
	p2                 game.Player
	pTurn, didJustJump bool
	jumpRow, jumpCol   int
	round              int
}

type CheckersMove struct {
	row1, col1, row2, col2 int
}

func (g Checkers) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g Checkers) GetHumanInput() game.Move {
	spot1 := readInts("Peice to move: ")
	spot2 := readInts("Move to: ")
	return CheckersMove{row1: spot1[0], col1: spot1[1], row2: spot2[0], col2: spot2[1]}
}

func (g Checkers) GameOver() (bool, game.Player) {
	p1Alive := false
	p2Alive := false
	for _, row := range g.board {
		for _, spot := range row {
			if spot == "X" || spot == "x" {
				p1Alive = true
			} else if spot == "O" || spot == "o" {
				p2Alive = true
			}
		}
	}
	if !p1Alive {
		return true, g.p2
	} else if !p2Alive {
		return true, g.p1
	} else {
		moves := g.GetPossibleMoves()
		if len(moves) == 0 || g.round > 500 {
			return true, player.HumanPlayer{"DRAW"}
		}
		return false, player.ComputerPlayer{}
	}
}

func NewCheckers(p1 string, p2 string, depth1 int, depth2 int) *Checkers {
	c := new(Checkers)
	c.p1 = getPlayer(p1, "Player 1", depth1)
	c.p2 = getPlayer(p2, "Player 2", depth2)
	c.pTurn = true
	c.board = [8][8]string{
		{".", "o", ".", "o", ".", "o", ".", "o"},
		{"o", ".", "o", ".", "o", ".", "o", "."},
		{".", "o", ".", "o", ".", "o", ".", "o"},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", ".", ".", "."},
		{"x", ".", "x", ".", "x", ".", "x", "."},
		{".", "x", ".", "x", ".", "x", ".", "x"},
		{"x", ".", "x", ".", "x", ".", "x", "."},
	}
	c.didJustJump = false
	c.round = 0
	return c
}

func (g Checkers) MakeMove(m game.Move) game.Game {
	g.round++
	move := m.(CheckersMove)
	g.board[move.row2][move.col2] = g.board[move.row1][move.col1]
	g.board[move.row1][move.col1] = "."
	if move.row1 == move.row2+2 || move.row1 == move.row2-2 {
		rowAvg := (move.row1 + move.row2) / 2
		colAvg := (move.col1 + move.col2) / 2
		g.board[rowAvg][colAvg] = "."
		g.didJustJump = true
		g.jumpRow = move.row2
		g.jumpCol = move.col2
	} else {
		g.didJustJump = false
	}
	if move.row2 == 0 && g.board[move.row2][move.col2] == "x" {
		g.board[move.row2][move.col2] = "X"
	} else if move.row2 == 7 && g.board[move.row2][move.col2] == "o" {
		g.board[move.row2][move.col2] = "O"
	}
	if g.didJustJump {
		moves := g.GetPossibleMoves()
		if len(moves) == 0 {
			g.pTurn = !g.pTurn
			g.didJustJump = false
			moves := g.GetPossibleMoves()
			if len(moves) == 0 {
				g.pTurn = !g.pTurn
			}
			return g
		} else {
			return g
		}
	} else {
		g.pTurn = !g.pTurn
		moves := g.GetPossibleMoves()
		if len(moves) == 0 {
			g.pTurn = !g.pTurn
		}
		return g
	}
}

func (g Checkers) isGoodMove(m CheckersMove) bool {
	possibleMoves := g.GetPossibleMoves()
	for _, move := range possibleMoves {
		if move == m {
			return true
		}
	}
	return false
}

func (g Checkers) CurrentScore(p game.Player) int {
	score := 0
	for _, row := range g.board {
		for _, spot := range row {
			if spot == "x" {
				score += 1
			} else if spot == "X" {
				score += 2
			} else if spot == "o" {
				score -= 1
			} else if spot == "O" {
				score -= 2
			}
		}
	}
	if p == g.p1 {
		return score
	} else {
		return -1 * score
	}
}

func (g Checkers) checkJump(i, j, rowDir, colDir int) []game.Move {
	var moves []game.Move
	peice := g.board[i][j]
	var target1, target2 string
	if peice == "x" || peice == "X" {
		target1 = "o"
		target2 = "O"
	} else {
		target1 = "x"
		target2 = "X"
	}
	if isInside(g, i+rowDir*2, j+colDir*2) && g.board[i+rowDir*2][j+colDir*2] == "." {
		if g.board[i+rowDir][j+colDir] == target1 || g.board[i+rowDir][j+colDir] == target2 {
			moves = append(moves, CheckersMove{
				row1: i,
				col1: j,
				row2: i + rowDir*2,
				col2: j + colDir*2,
			})
		}
	}
	return moves
}

func (g Checkers) checkMove(i, j, rowDir, colDir int) []game.Move {
	var moves []game.Move
	if isInside(g, i+rowDir, j+colDir) && g.board[i+rowDir][j+colDir] == "." {
		moves = append(moves, CheckersMove{
			row1: i,
			col1: j,
			row2: i + rowDir,
			col2: j + colDir,
		})
	}
	return moves
}

func (g Checkers) GetPossibleMoves() []game.Move {
	if g.didJustJump {
		var moves []game.Move
		peice := g.board[g.jumpRow][g.jumpCol]
		row := g.jumpRow
		col := g.jumpCol
		if peice == "x" || peice == "X" || peice == "O" {
			moves = append(moves, g.checkJump(row, col, -1, -1)...)
			moves = append(moves, g.checkJump(row, col, -1, 1)...)
		}
		if peice == "o" || peice == "O" || peice == "X" {
			moves = append(moves, g.checkJump(row, col, 1, -1)...)
			moves = append(moves, g.checkJump(row, col, 1, 1)...)
		}
		return moves
	}
	var moves []game.Move
	for i, row := range g.board {
		for j, spot := range row {
			if (g.pTurn && (spot == "X" || spot == "x")) || (!g.pTurn && spot == "O") {
				moves = append(moves, g.checkJump(i, j, -1, -1)...)
				moves = append(moves, g.checkJump(i, j, -1, 1)...)
			} else if (!g.pTurn && (spot == "O" || spot == "o")) || (g.pTurn && spot == "X") {
				moves = append(moves, g.checkJump(i, j, 1, -1)...)
				moves = append(moves, g.checkJump(i, j, 1, 1)...)
			}
		}
	}
	if len(moves) > 0 {
		return moves
	}
	for i, row := range g.board {
		for j, spot := range row {
			if (g.pTurn && (spot == "X" || spot == "x")) || (!g.pTurn && spot == "O") {
				moves = append(moves, g.checkMove(i, j, -1, -1)...)
				moves = append(moves, g.checkMove(i, j, -1, 1)...)
			}
			if (!g.pTurn && (spot == "O" || spot == "o")) || (g.pTurn && spot == "X") {
				moves = append(moves, g.checkMove(i, j, 1, -1)...)
				moves = append(moves, g.checkMove(i, j, 1, 1)...)
			}
		}
	}
	return moves
}

func (g Checkers) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Checkers) GetTurn(p game.Player) game.Move {
	m := p.GetTurn(g)
	//move := m.(CheckersMove)
	//for !g.isGoodMove(move) {
	//	m = p.GetTurn(g)
	//	move = m.(CheckersMove)
	//}
	return m
}

func (g Checkers) BoardString() string {
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

func (g Checkers) PrintBoard() {
	fmt.Println("-----------------")
	fmt.Println("  0 1 2 3 4 5 6 7")
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
