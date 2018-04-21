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

type Reversi struct {
	board [8][8]string
	p1    game.Player
	p2    game.Player
	pTurn bool
}

type ReversiMove struct {
	row, col int
}

func (g Reversi) PrintBoard() {
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

func (g Reversi) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Reversi) GetHumanInput() game.Move {
	fmt.Println("Spot to place: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	spot := strings.Split(strings.TrimSpace(text), ",")
	row, col := spot[0], spot[1]
	rowI, _ := strconv.Atoi(row)
	colI, _ := strconv.Atoi(col)
	return ReversiMove{row: rowI, col: colI}
}

func (g Reversi) GetPossibleMoves() []game.Move {
	var moves []game.Move
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				var target, match string
				if g.pTurn {
					target, match = "O", "X"
				} else {
					target, match = "X", "O"
				}
				if j+1 < 8 && g.board[i][j+1] == target {
					rowCheck := i
					colCheck := j + 1
					for colCheck < 8 && g.board[rowCheck][colCheck] == target {
						colCheck++
					}
					if colCheck < 8 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
						continue
					}
				}
				if i+1 < 8 && j+1 < 8 && g.board[i+1][j+1] == target {
					rowCheck := i + 1
					colCheck := j + 1
					for rowCheck < 8 && colCheck < 8 && g.board[rowCheck][colCheck] == target {
						colCheck++
						rowCheck++
					}
					if colCheck < 8 && rowCheck < 8 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
						continue
					}
				}
				if i+1 < 8 && g.board[i+1][j] == target {
					rowCheck := i + 1
					colCheck := j
					for rowCheck < 8 && g.board[rowCheck][colCheck] == target {
						rowCheck++
					}
					if rowCheck < 8 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
						continue
					}
				}
				if i+1 < 8 && j-1 >= 0 && g.board[i+1][j-1] == target {
					rowCheck := i + 1
					colCheck := j - 1
					for rowCheck < 8 && colCheck >= 0 && g.board[rowCheck][colCheck] == target {
						rowCheck++
						colCheck--
					}
					if rowCheck < 8 && colCheck >= 0 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
						continue
					}
				}
				if j-1 >= 0 && g.board[i][j-1] == target {
					rowCheck := j
					colCheck := j - 1
					for colCheck >= 0 && g.board[rowCheck][colCheck] == target {
						colCheck--
					}
					if colCheck >= 0 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
						continue
					}
				}
				if i-1 >= 0 && j-1 >= 0 && g.board[i-1][j-1] == target {
					rowCheck := i - 1
					colCheck := j - 1
					for rowCheck >= 0 && colCheck >= 0 && g.board[rowCheck][colCheck] == target {
						rowCheck--
						colCheck--
					}
					if rowCheck >= 0 && colCheck >= 0 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
						continue
					}
				}
				if i-1 >= 0 && g.board[i-1][j] == target {
					rowCheck := i - 1
					colCheck := j
					for rowCheck >= 0 && g.board[rowCheck][colCheck] == target {
						rowCheck--
					}
					if rowCheck >= 0 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
					}
				}
				if i-1 >= 0 && j+1 < 8 && g.board[i-1][j+1] == target {
					rowCheck := i - 1
					colCheck := j + 1
					for rowCheck >= 0 && colCheck < 8 && g.board[rowCheck][colCheck] == target {
						rowCheck--
						colCheck++
					}
					if rowCheck >= 0 && colCheck < 8 && g.board[rowCheck][colCheck] == match {
						moves = append(moves, ReversiMove{
							row: i,
							col: j,
						})
					}
				}
			}
		}
	}
	return moves
}

func (g Reversi) GetTurn(p game.Player) game.Move {
	m := p.GetTurn(g)
	return m
}

func (g Reversi) MakeMove(m game.Move) game.Game {
	move := m.(ReversiMove)
	var target, match string
	if g.pTurn {
		target, match = "O", "X"
	} else {
		target, match = "X", "O"
	}
	g.board[move.row][move.col] = match
	if move.col+1 < 8 && g.board[move.row][move.col+1] == target {
		rowCheck := move.row
		colCheck := move.col + 1
		for colCheck < 8 && g.board[rowCheck][colCheck] == target {
			colCheck++
		}
		if colCheck < 8 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; j++ {
				g.board[i][j] = match
			}
		}
	}
	if move.row+1 < 8 && move.col+1 < 8 && g.board[move.row+1][move.col+1] == target {
		rowCheck := move.row + 1
		colCheck := move.col + 1
		for rowCheck < 8 && colCheck < 8 && g.board[rowCheck][colCheck] == target {
			rowCheck++
			colCheck++
		}
		if rowCheck < 8 && colCheck < 8 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; i, j = i+1, j+1 {
				g.board[i][j] = match
			}
		}
	}
	if move.row+1 < 8 && g.board[move.row+1][move.col] == target {
		rowCheck := move.row + 1
		colCheck := move.col
		for rowCheck < 8 && g.board[rowCheck][colCheck] == target {
			rowCheck++
		}
		if rowCheck < 8 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; i++ {
				g.board[i][j] = match
			}
		}
	}
	if move.row+1 < 8 && move.col-1 >= 0 && g.board[move.row+1][move.col-1] == target {
		rowCheck := move.row + 1
		colCheck := move.col - 1
		for rowCheck < 8 && colCheck >= 0 && g.board[rowCheck][colCheck] == target {
			rowCheck++
			colCheck--
		}
		if rowCheck < 8 && colCheck >= 0 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; i, j = i+1, j-1 {
				g.board[i][j] = match
			}
		}
	}
	if move.col-1 >= 0 && g.board[move.row][move.col-1] == target {
		rowCheck := move.row
		colCheck := move.col - 1
		for colCheck >= 0 && g.board[rowCheck][colCheck] == target {
			colCheck--
		}
		if colCheck >= 0 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; j-- {
				g.board[i][j] = match
			}
		}
	}
	if move.row-1 >= 0 && move.col-1 >= 0 && g.board[move.row-1][move.col-1] == target {
		rowCheck := move.row - 1
		colCheck := move.col - 1
		for rowCheck >= 0 && colCheck >= 0 && g.board[rowCheck][colCheck] == target {
			rowCheck--
			colCheck--
		}
		if rowCheck >= 0 && colCheck >= 0 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; i, j = i-1, j-1 {
				g.board[i][j] = match
			}
		}
	}
	if move.row-1 >= 0 && g.board[move.row-1][move.col] == target {
		rowCheck := move.row - 1
		colCheck := move.col
		for rowCheck >= 0 && g.board[rowCheck][colCheck] == target {
			rowCheck--
		}
		if rowCheck >= 0 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; i-- {
				g.board[i][j] = match
			}
		}
	}
	if move.row-1 >= 0 && move.col+1 < 8 && g.board[move.row-1][move.col+1] == target {
		rowCheck := move.row - 1
		colCheck := move.col + 1
		for rowCheck >= 0 && colCheck < 8 && g.board[rowCheck][colCheck] == target {
			rowCheck--
			colCheck++
		}
		if rowCheck >= 0 && colCheck < 8 && g.board[rowCheck][colCheck] == match {
			for i, j := move.row, move.col; i != rowCheck || j != colCheck; i, j = i-1, j+1 {
				g.board[i][j] = match
			}
		}
	}
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

func NewReversi(p1, p2 string) *Reversi {
	r := new(Reversi)
	r.p1 = getPlayer(p1, "Player 1")
	r.p2 = getPlayer(p2, "Player 2")
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
