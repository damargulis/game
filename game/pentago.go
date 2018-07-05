package game

import (
	"bufio"
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
	"os"
	"strings"
)

type Pentago struct {
	board  [6][6]string
	p1     game.Player
	p2     game.Player
	pTurn  bool
	stage1 bool
	round  int
}

type PentagoMove struct {
	row, col, quad int
	clockwise      bool
}

func (g Pentago) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func NewPentago(p1 string, p2 string, depth1 int, depth2 int) *Pentago {
	g := new(Pentago)
	g.round = 0
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.stage1 = true
	g.board = [6][6]string{
		{".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", "."},
		{".", ".", ".", ".", ".", "."},
	}
	return g
}

func (g Pentago) BoardString() string {
	s := "--------------\n"
	for i, row := range g.board {
		s += fmt.Sprintf("%v ", i)
		for _, p := range row {
			s += p
			s += " "
		}
		s += "\n"
	}
	s += "  0 1 2 3 4 5 \n"
	s += "--------------"
	return s
}

func (g Pentago) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Pentago) GetHumanInput() game.Move {
	reader := bufio.NewReader(os.Stdin)
	if g.stage1 {
		spot1 := readInts("Spot to place: ")
		return PentagoMove{
			row: spot1[0],
			col: spot1[1],
		}
	} else {
		quad := readInts("Quadren to spin: ")
		quadI := quad[0]
		var dir string
		for dir != "CW" && dir != "CCW" {
			fmt.Println("Direction (CW/CCW): ")
			dir, _ = reader.ReadString('\n')
			dir = strings.TrimSpace(dir)
		}
		return PentagoMove{
			quad:      quadI,
			clockwise: dir == "CW",
		}
	}
}

var cRows = [4]int{1, 1, 4, 4}
var cCols = [4]int{1, 4, 1, 4}

func (g Pentago) GetPossibleMoves() []game.Move {
	var moves []game.Move
	if g.stage1 {
		for i, row := range g.board {
			for j, spot := range row {
				if spot == "." {
					moves = append(moves, PentagoMove{
						row: i,
						col: j,
					})
				}
			}
		}
	} else {
		skippable := false
		for i := range cRows {
			cRow := cRows[i]
			cCol := cCols[i]
			if g.board[cRow-1][cCol-1] == g.board[cRow-1][cCol+1] && g.board[cRow-1][cCol-1] == g.board[cRow+1][cCol-1] && g.board[cRow-1][cCol-1] == g.board[cRow+1][cCol+1] {
				if g.board[cRow-1][cCol] == g.board[cRow+1][cCol] && g.board[cRow-1][cCol] == g.board[cRow][cCol-1] && g.board[cRow-1][cCol] == g.board[cRow][cCol+1] {
					if !skippable {
						moves = append(moves, PentagoMove{
							quad:      i,
							clockwise: true,
						})
					}
					skippable = true
					continue
				}
			}
			moves = append(moves, PentagoMove{
				quad:      i,
				clockwise: true,
			}, PentagoMove{
				quad:      i,
				clockwise: false,
			})
		}
	}
	return moves
}

func (g Pentago) MakeMove(m game.Move) game.Game {
	g.round++
	move := m.(PentagoMove)
	if g.stage1 {
		if g.pTurn {
			g.board[move.row][move.col] = "X"
		} else {
			g.board[move.row][move.col] = "O"
		}
		g.stage1 = false
	} else {
		quad := move.quad
		dir := move.clockwise
		cRow := cRows[quad]
		cCol := cCols[quad]
		if dir {
			tmp := g.board[cRow-1][cCol-1]
			g.board[cRow-1][cCol-1] = g.board[cRow+1][cCol-1]
			g.board[cRow+1][cCol-1] = g.board[cRow+1][cCol+1]
			g.board[cRow+1][cCol+1] = g.board[cRow-1][cCol+1]
			g.board[cRow-1][cCol+1] = tmp
			tmp = g.board[cRow-1][cCol]
			g.board[cRow-1][cCol] = g.board[cRow][cCol-1]
			g.board[cRow][cCol-1] = g.board[cRow+1][cCol]
			g.board[cRow+1][cCol] = g.board[cRow][cCol+1]
			g.board[cRow][cCol+1] = tmp
		} else {
			tmp := g.board[cRow-1][cCol-1]
			g.board[cRow-1][cCol-1] = g.board[cRow-1][cCol+1]
			g.board[cRow-1][cCol+1] = g.board[cRow+1][cCol+1]
			g.board[cRow+1][cCol+1] = g.board[cRow+1][cCol-1]
			g.board[cRow+1][cCol-1] = tmp
			tmp = g.board[cRow-1][cCol]
			g.board[cRow-1][cCol] = g.board[cRow][cCol+1]
			g.board[cRow][cCol+1] = g.board[cRow+1][cCol]
			g.board[cRow+1][cCol] = g.board[cRow][cCol-1]
			g.board[cRow][cCol-1] = tmp
		}
		g.pTurn = !g.pTurn
		g.stage1 = true
	}
	return g
}

func (g Pentago) GameOver() (bool, game.Player) {
	hasSpace := false
	p1win := false
	p2win := false
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				hasSpace = true
				continue
			}
			if isInside(g, i+4, j) {
				if spot == g.board[i+1][j] && spot == g.board[i+2][j] && spot == g.board[i+3][j] && spot == g.board[i+4][j] {
					if spot == "X" {
						p1win = true
					} else {
						p2win = true
					}
				}
			}
			if isInside(g, i, j+4) {
				if spot == g.board[i][j+1] && spot == g.board[i][j+2] && spot == g.board[i][j+3] && spot == g.board[i][j+4] {
					if spot == "X" {
						p1win = true
					} else {
						p2win = true
					}
				}
			}
			if isInside(g, i+4, j+4) {
				if spot == g.board[i+1][j+1] && spot == g.board[i+2][j+2] && spot == g.board[i+3][j+3] && spot == g.board[i+4][j+4] {
					if spot == "X" {
						p1win = true
					} else {
						p2win = true
					}
				}
			}
			if isInside(g, i-4, j+4) {
				if spot == g.board[i-1][j+1] && spot == g.board[i-2][j+2] && spot == g.board[i-3][j+3] && spot == g.board[i-4][j+4] {
					if spot == "X" {
						p1win = true
					} else {
						p2win = true
					}
				}
			}
		}
	}
	if p1win && p2win {
		return true, player.HumanPlayer{"DRAW"}
	} else if p1win {
		return true, g.p1
	} else if p2win {
		return true, g.p2
	} else if g.stage1 && !hasSpace {
		return true, player.HumanPlayer{"DRAW"}
	} else {
		return false, player.ComputerPlayer{}
	}
}

func (g Pentago) CurrentScore(p game.Player) int {
	return 0
}

func (g Pentago) GetRound() int {
	return g.round
}
