package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type NineMensMorris struct {
	board                     [7][7]string
	p1                        game.Player
	p2                        game.Player
	pTurn, stage1, justMilled bool
	p1toPlace, p2toPlace      int
	round                     int
}

type NineMensMorrisMove struct {
	row1, col1, row2, col2 int
}

func (g NineMensMorris) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func NewNineMensMorris(p1 string, p2 string, depth1 int, depth2 int) *NineMensMorris {
	g := new(NineMensMorris)
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.stage1 = true
	g.justMilled = false
	g.p1toPlace = 9
	g.p2toPlace = 9
	g.round = 0
	g.board = [7][7]string{
		{".", "-", "-", ".", "-", "-", "."},
		{"|", ".", "-", ".", "-", ".", "|"},
		{"|", "|", ".", ".", ".", "|", "|"},
		{".", ".", ".", " ", ".", ".", "."},
		{"|", "|", ".", ".", ".", "|", "|"},
		{"|", ".", "-", ".", "-", ".", "|"},
		{".", "-", "-", ".", "-", "-", "."},
	}
	return g
}

func (g NineMensMorris) BoardString() string {
	s := "---------------\n"
	s += "  0 1 2 3 4 5 6\n"
	for i, row := range g.board {
		s += fmt.Sprintf("%v ", i)
		for _, p := range row {
			s += p
			s += " "
		}
		s += "\n"
	}
	s += "  0 1 2 3 4 5 6\n"
	s += "---------------"
	return s
}

func (g NineMensMorris) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func isIn(moves []game.Move, move game.Move) bool {
	for _, m := range moves {
		if m == move {
			return true
		}
	}
	return false
}

func (g NineMensMorris) GetHumanInput() game.Move {
	var move game.Move
	var possibleMoves = g.GetPossibleMoves()
	for !isIn(possibleMoves, move) {
		var spot []int
		if g.justMilled {
			spot = readInts("Peice to take: ")
		} else if g.stage1 {
			spot = readInts("Spot to place: ")
		} else {
			spot = readInts("Peice to move: ")
		}
		if g.stage1 || g.justMilled {
			move = NineMensMorrisMove{row1: spot[0], col1: spot[1]}
		} else {
			spot2 := readInts("Move to: ")
			move = NineMensMorrisMove{
				row1: spot[0],
				col1: spot[1],
				row2: spot2[0],
				col2: spot2[1],
			}
		}
	}
	return move
}

func (g NineMensMorris) GetPossibleMoves() []game.Move {
	if g.justMilled {
		var moves []game.Move
		var allMoves []game.Move
		var target string
		if g.pTurn {
			target = "O"
		} else {
			target = "X"
		}
		for i, row := range g.board {
			for j, spot := range row {
				if spot == target {
					if !g.isInMill(i, j) {
						moves = append(moves, NineMensMorrisMove{
							row1: i,
							col1: j,
						})
					}
					allMoves = append(allMoves, NineMensMorrisMove{
						row1: i,
						col1: j,
					})
				}
			}
		}
		if len(moves) > 0 {
			return moves
		} else {
			return allMoves
		}
	} else if g.stage1 {
		var moves []game.Move
		for i, row := range g.board {
			for j, spot := range row {
				if spot == "." {
					moves = append(moves, NineMensMorrisMove{
						row1: i,
						col1: j,
					})
				}
			}
		}
		return moves
	} else {
		var moves []game.Move
		var owns string
		if g.pTurn {
			owns = "X"
		} else {
			owns = "O"
		}
		for i, row := range g.board {
			for j, spot := range row {
				if spot == owns {
					checkRow := i
					checkCol := j + 1
					for isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "-" {
						checkCol++
					}
					if isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "." {
						moves = append(moves, NineMensMorrisMove{
							row1: i,
							col1: j,
							row2: checkRow,
							col2: checkCol,
						})
					}
					checkCol = j - 1
					for isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "-" {
						checkCol--
					}
					if isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "." {
						moves = append(moves, NineMensMorrisMove{
							row1: i,
							col1: j,
							row2: checkRow,
							col2: checkCol,
						})
					}
					checkCol = j
					checkRow = i + 1
					for isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "|" {
						checkRow++
					}
					if isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "." {
						moves = append(moves, NineMensMorrisMove{
							row1: i,
							col1: j,
							row2: checkRow,
							col2: checkCol,
						})
					}
					checkRow = i - 1
					for isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "|" {
						checkRow--
					}
					if isInside(g, checkRow, checkCol) && g.board[checkRow][checkCol] == "." {
						moves = append(moves, NineMensMorrisMove{
							row1: i,
							col1: j,
							row2: checkRow,
							col2: checkCol,
						})
					}
				}
			}
		}
		return moves
	}
}

func (g NineMensMorris) MakeMove(m game.Move) game.Game {
	g.round++
	move := m.(NineMensMorrisMove)
	if g.justMilled {
		g.board[move.row1][move.col1] = "."
		g.justMilled = false
		g.pTurn = !g.pTurn
		return g
	}
	var toRow, toCol int
	var own string
	if g.pTurn {
		own = "X"
	} else {
		own = "O"
	}
	if g.stage1 {
		if g.pTurn {
			g.board[move.row1][move.col1] = own
			g.p1toPlace--
		} else {
			g.board[move.row1][move.col1] = own
			g.p2toPlace--
		}
		if g.p1toPlace == 0 && g.p2toPlace == 0 {
			g.stage1 = false
		}
		toRow, toCol = move.row1, move.col1
	} else {
		g.board[move.row2][move.col2] = g.board[move.row1][move.col1]
		g.board[move.row1][move.col1] = "."
		toRow, toCol = move.row2, move.col2
	}
	if g.isInMill(toRow, toCol) {
		g.justMilled = true
	} else {
		g.pTurn = !g.pTurn
	}
	return g
}

func (g NineMensMorris) isInMill(row, col int) bool {
	own := g.board[row][col]
	horizontal := 0
	vertical := 0
	checkRow := row
	checkCol := col + 1
	for isInside(g, checkRow, checkCol) && (g.board[checkRow][checkCol] == "-" || g.board[checkRow][checkCol] == own) {
		if g.board[checkRow][checkCol] == own {
			horizontal += 1
		}
		checkCol++
	}
	checkCol = col - 1
	for isInside(g, checkRow, checkCol) && (g.board[checkRow][checkCol] == "-" || g.board[checkRow][checkCol] == own) {
		if g.board[checkRow][checkCol] == own {
			horizontal += 1
		}
		checkCol--
	}
	checkCol = col
	checkRow = row + 1
	for isInside(g, checkRow, checkCol) && (g.board[checkRow][checkCol] == "|" || g.board[checkRow][checkCol] == own) {
		if g.board[checkRow][checkCol] == own {
			vertical += 1
		}
		checkRow++
	}
	checkRow = row - 1
	for isInside(g, checkRow, checkCol) && (g.board[checkRow][checkCol] == "|" || g.board[checkRow][checkCol] == own) {
		if g.board[checkRow][checkCol] == own {
			vertical += 1
		}
		checkRow--
	}
	return horizontal >= 2 || vertical >= 2
}

func (g NineMensMorris) GameOver() (bool, game.Player) {
	if g.stage1 {
		return false, player.ComputerPlayer{}
	}
	if len(g.GetPossibleMoves()) == 0 {
		if g.pTurn {
			return true, g.p2
		} else {
			return true, g.p1
		}
	}
	p1 := 0
	p2 := 0
	for _, row := range g.board {
		for _, spot := range row {
			if spot == "X" {
				p1++
			} else if spot == "O" {
				p2++
			}
		}
	}
	if g.round > 500 {
		if p1 > p2 {
			return true, g.p1
		} else if p2 > p1 {
			return true, g.p2
		} else {
			return true, player.HumanPlayer{"DRAW"}
		}
	}
	if p1 < 3 {
		return true, g.p2
	} else if p2 < 3 {
		return true, g.p1
	} else {
		return false, player.ComputerPlayer{}
	}
}

func (g NineMensMorris) CurrentScore(p game.Player) int {
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

func (g NineMensMorris) GetRound() int {
	return g.round
}
