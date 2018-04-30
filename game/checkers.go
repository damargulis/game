package game

import (
	"bufio"
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
	"os"
	"strconv"
	"strings"
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

func (g Checkers) GetHumanInput() game.Move {
	fmt.Println("Peice to move: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	spot1 := strings.Split(strings.TrimSpace(text), ",")
	row1, col1 := spot1[0], spot1[1]
	row1I, _ := strconv.Atoi(row1)
	col1I, _ := strconv.Atoi(col1)
	fmt.Println("Move to: ")
	text2, _ := reader.ReadString('\n')
	spot2 := strings.Split(strings.TrimSpace(text2), ",")
	row2, col2 := spot2[0], spot2[1]
	row2I, _ := strconv.Atoi(row2)
	col2I, _ := strconv.Atoi(col2)
	return CheckersMove{row1: row1I, col1: col1I, row2: row2I, col2: col2I}
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

func (g Checkers) GetPossibleMoves() []game.Move {
	if g.didJustJump {
		var moves []game.Move
		peice := g.board[g.jumpRow][g.jumpCol]
		row := g.jumpRow
		col := g.jumpCol
		if peice == "x" || peice == "X" {
			if row-2 >= 0 && col-2 >= 0 && g.board[row-2][col-2] == "." && (g.board[row-1][col-1] == "o" || g.board[row-1][col-1] == "O") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row - 2,
					col2: col - 2,
				})
			}
			if row-2 >= 0 && col+2 < 8 && g.board[row-2][col+2] == "." && (g.board[row-1][col+1] == "o" || g.board[row-1][col+1] == "O") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row - 2,
					col2: col + 2,
				})
			}
		}
		if peice == "X" {
			if row+2 < 8 && col-2 >= 0 && g.board[row+2][col-2] == "." && (g.board[row+1][col-1] == "o" || g.board[row+1][col-1] == "O") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row + 2,
					col2: col - 2,
				})
			}
			if row+2 < 8 && col+2 < 8 && g.board[row+2][col+2] == "." && (g.board[row+1][col+1] == "o" || g.board[row+1][col+1] == "O") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row + 2,
					col2: col + 2,
				})
			}
		}
		if peice == "o" || peice == "O" {
			if row+2 < 8 && col-2 >= 0 && g.board[row+2][col-2] == "." && (g.board[row+1][col-1] == "x" || g.board[row+1][col-1] == "X") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row + 2,
					col2: col - 2,
				})
			}
			if row+2 < 8 && col+2 < 8 && g.board[row+2][col+2] == "." && (g.board[row+1][col+1] == "x" || g.board[row+1][col+1] == "X") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row + 2,
					col2: col + 2,
				})
			}
		}
		if peice == "O" {
			if row-2 >= 0 && col-2 >= 0 && g.board[row-2][col-2] == "." && (g.board[row-1][col-1] == "x" || g.board[row-1][col-1] == "X") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row - 2,
					col2: col - 2,
				})
			}
			if row-2 >= 0 && col+2 < 8 && g.board[row-2][col+2] == "." && (g.board[row-1][col+1] == "x" || g.board[row-1][col+1] == "X") {
				moves = append(moves, CheckersMove{
					row1: row,
					col1: col,
					row2: row - 2,
					col2: col + 2,
				})
			}
		}
		return moves
	}
	var moves []game.Move
	for i, row := range g.board {
		for j, spot := range row {
			if g.pTurn && (spot == "X" || spot == "x") {
				if i-2 >= 0 && j-2 >= 0 && g.board[i-2][j-2] == "." && (g.board[i-1][j-1] == "O" || g.board[i-1][j-1] == "o") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 2,
						col2: j - 2,
					})
				}
				if i-2 >= 0 && j+2 < 8 && g.board[i-2][j+2] == "." && (g.board[i-1][j+1] == "O" || g.board[i-1][j+1] == "o") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 2,
						col2: j + 2,
					})
				}
				if spot == "X" {
					if i+2 < 8 && j-2 >= 0 && g.board[i+2][j-2] == "." && (g.board[i+1][j-1] == "O" || g.board[i+1][j-1] == "o") {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i + 2,
							col2: j - 2,
						})
					}
					if i+2 < 8 && j+2 < 8 && g.board[i+2][j+2] == "." && (g.board[i+1][j+1] == "O" || g.board[i+1][j+1] == "o") {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i + 2,
							col2: j + 2,
						})
					}
				}
			} else if !g.pTurn && (spot == "O" || spot == "o") {
				if i+2 < 8 && j-2 >= 0 && g.board[i+2][j-2] == "." && (g.board[i+1][j-1] == "X" || g.board[i+1][j-1] == "x") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i + 2,
						col2: j - 2,
					})
				}
				if i+2 < 8 && j+2 < 8 && g.board[i+2][j+2] == "." && (g.board[i+1][j+1] == "X" || g.board[i+1][j+1] == "x") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i + 2,
						col2: j + 2,
					})
				}
				if spot == "O" {
					if i-2 >= 0 && j-2 >= 0 && g.board[i-2][j-2] == "." && (g.board[i-1][j-1] == "x" || g.board[i-1][j-1] == "X") {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 2,
							col2: j - 2,
						})
					}
					if i-2 >= 0 && j+2 < 8 && g.board[i-2][j+2] == "." && (g.board[i-1][j+1] == "x" || g.board[i-1][j+1] == "X") {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 2,
							col2: j + 2,
						})
					}
				}
			}
		}
	}
	if len(moves) > 0 {
		return moves
	}
	for i, row := range g.board {
		for j, spot := range row {
			if g.pTurn && (spot == "X" || spot == "x") {
				if i-1 >= 0 && j-1 >= 0 && g.board[i-1][j-1] == "." {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 1,
						col2: j - 1,
					})
				}
				if i-1 >= 0 && j+1 < 8 && g.board[i-1][j+1] == "." {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 1,
						col2: j + 1,
					})
				}
				if spot == "X" {
					if i+1 < 8 && j-1 >= 0 && g.board[i+1][j-1] == "." {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i + 1,
							col2: j - 1,
						})
					}
					if i+1 < 8 && j+1 < 8 && g.board[i+1][j+1] == "." {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i + 1,
							col2: j + 1,
						})
					}
				}
			} else if !g.pTurn && (spot == "O" || spot == "o") {
				if i+1 < 8 && j-1 >= 0 && g.board[i+1][j-1] == "." {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i + 1,
						col2: j - 1,
					})
				}
				if i+1 < 8 && j+1 < 8 && g.board[i+1][j+1] == "." {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i + 1,
						col2: j + 1,
					})
				}
				if spot == "O" {
					if i-1 >= 0 && j-1 >= 0 && g.board[i-1][j-1] == "." {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 1,
							col2: j - 1,
						})
					}
					if i-1 >= 0 && j+1 < 8 && g.board[i-1][j+1] == "." {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 1,
							col2: j + 1,
						})
					}
				}
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
