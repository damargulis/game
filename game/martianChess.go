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

type MartianChess struct {
	board              [8][4]string
	p1                 game.Player
	p2                 game.Player
	p1points, p2points int
	pTurn              bool
	lastMove           MartianChessMove
	round              int
}

type MartianChessMove struct {
	startRow, startCol, endRow, endCol int
}

func NewMartianChess(p1 string, p2 string, depth1 int, depth2 int) *MartianChess {
	g := new(MartianChess)
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.board = [8][4]string{
		{"Q", "Q", "D", "."},
		{"Q", "D", "P", "."},
		{"D", "P", "P", "."},
		{".", ".", ".", "."},
		{".", ".", ".", "."},
		{".", "P", "P", "D"},
		{".", "P", "D", "Q"},
		{".", "D", "Q", "Q"},
	}
	g.round = 0
	return g
}

func (g MartianChess) BoardString() string {
	s := "---------\n"
	s += "  0 1 2 3\n"
	for i, row := range g.board {
		s += fmt.Sprintf("%v ", i)
		for _, p := range row {
			s += p
			s += " "
		}
		s += "\n"
	}
	s += "  0 1 2 3\n"
	s += fmt.Sprintf("P1: %v P2: %v\n", g.p1points, g.p2points)
	s += "---------\n"
	return s
}

func (g MartianChess) PrintBoard() {
	fmt.Println(g.BoardString())
}

func (g MartianChess) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g MartianChess) GetHumanInput() game.Move {
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
	return MartianChessMove{startRow: row1I, startCol: col1I, endRow: row2I, endCol: col2I}
}

func in(arr []int, check int) bool {
	for _, i := range arr {
		if i == check {
			return true
		}
	}
	return false
}

func (g MartianChess) getDroneMoves(i, j, rowDir, colDir int, rows []int) []MartianChessMove {
	var moves []MartianChessMove
	if i+rowDir >= 0 && i+rowDir < len(g.board) && j+colDir >= 0 && j+colDir < len(g.board[i+rowDir]) {
		if g.board[i+rowDir][j+colDir] == "." || !in(rows, i+rowDir) {
			moves = append(moves, MartianChessMove{
				startRow: i,
				startCol: j,
				endRow:   i + rowDir,
				endCol:   j + colDir,
			})
		}
		if g.board[i+rowDir][j+colDir] == "." {
			endRow := i + rowDir*2
			endCol := j + colDir*2
			if endRow >= 0 && endRow < len(g.board) && endCol >= 0 && endCol < len(g.board[endRow]) {
				if g.board[endRow][endCol] == "." || !in(rows, endRow) {
					moves = append(moves, MartianChessMove{
						startRow: i,
						startCol: j,
						endRow:   endRow,
						endCol:   endCol,
					})
				}
			}
		}
	}
	return moves
}

func (g MartianChess) getPawnMoves(i, j, rowDir, colDir int, rows []int) []MartianChessMove {
	var moves []MartianChessMove
	endRow := i + rowDir
	endCol := j + colDir
	if endRow >= 0 && endRow < len(g.board) && endCol >= 0 && endCol < len(g.board[endRow]) {
		if g.board[endRow][endCol] == "." || !in(rows, endRow) {
			moves = append(moves, MartianChessMove{
				startRow: i,
				startCol: j,
				endRow:   endRow,
				endCol:   endCol,
			})
		}
	}
	return moves
}

func (g MartianChess) getQueenMoves(i, j, rowDir, colDir int, rows []int) []MartianChessMove {
	var moves []MartianChessMove
	endRow := i + rowDir
	endCol := j + colDir
	for endRow >= 0 && endRow < len(g.board) && endCol >= 0 && endCol < len(g.board[endRow]) && g.board[endRow][endCol] == "." {
		moves = append(moves, MartianChessMove{
			startRow: i,
			startCol: j,
			endRow:   endRow,
			endCol:   endCol,
		})
		endRow += rowDir
		endCol += colDir
	}
	if endRow >= 0 && endRow < len(g.board) && endCol >= 0 && endCol < len(g.board[endRow]) && !in(rows, endRow) {
		moves = append(moves, MartianChessMove{
			startRow: i,
			startCol: j,
			endRow:   endRow,
			endCol:   endCol,
		})
	}
	return moves
}

func (g MartianChess) GetPossibleMoves() []game.Move {
	var rows []int
	if g.pTurn {
		rows = []int{7, 6, 5, 4}
	} else {
		rows = []int{0, 1, 2, 3}
	}
	var moves []MartianChessMove
	for _, i := range rows {
		row := g.board[i]
		for j, spot := range row {
			if spot == "Q" {
				moves = append(moves, g.getQueenMoves(i, j, -1, -1, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, -1, 0, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, -1, 1, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, 0, -1, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, 0, 1, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, 1, -1, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, 1, 0, rows)...)
				moves = append(moves, g.getQueenMoves(i, j, 1, 1, rows)...)
			} else if spot == "D" {
				moves = append(moves, g.getDroneMoves(i, j, -1, 0, rows)...)
				moves = append(moves, g.getDroneMoves(i, j, 1, 0, rows)...)
				moves = append(moves, g.getDroneMoves(i, j, 0, -1, rows)...)
				moves = append(moves, g.getDroneMoves(i, j, 0, 1, rows)...)
			} else if spot == "P" {
				moves = append(moves, g.getPawnMoves(i, j, -1, -1, rows)...)
				moves = append(moves, g.getPawnMoves(i, j, -1, 1, rows)...)
				moves = append(moves, g.getPawnMoves(i, j, 1, -1, rows)...)
				moves = append(moves, g.getPawnMoves(i, j, 1, 1, rows)...)
			}
		}
	}
	var rMoves []game.Move
	for _, move := range moves {
		if move.startRow == g.lastMove.endRow && move.startCol == g.lastMove.endCol && move.endRow == g.lastMove.startRow && move.endCol == g.lastMove.startCol {
			continue
		}
		rMoves = append(rMoves, move)
	}
	return rMoves
}

func (g MartianChess) GetTurn(p game.Player) game.Move {
	m := p.GetTurn(g)
	return m
}

func (g MartianChess) MakeMove(m game.Move) game.Game {
	g.round++
	move := m.(MartianChessMove)
	startRow := move.startRow
	startCol := move.startCol
	endRow := move.endRow
	endCol := move.endCol
	if g.board[endRow][endCol] == "Q" {
		if g.pTurn {
			g.p1points += 3
		} else {
			g.p2points += 3
		}
	} else if g.board[endRow][endCol] == "D" {
		if g.pTurn {
			g.p1points += 2
		} else {
			g.p2points += 2
		}
	} else if g.board[endRow][endCol] == "P" {
		if g.pTurn {
			g.p1points += 1
		} else {
			g.p2points += 1
		}
	}
	g.board[endRow][endCol] = g.board[startRow][startCol]
	g.board[startRow][startCol] = "."
	g.pTurn = !g.pTurn
	g.lastMove = move
	return g
}

func (g MartianChess) GameOver() (bool, game.Player) {
	if g.round > 500 {
		return true, player.HumanPlayer{"DRAW"}
	}
	rows1 := []int{0, 1, 2, 3}
	rows2 := []int{4, 5, 6, 7}
	p1Alive := false
	p2Alive := false
	for _, i := range rows1 {
		row := g.board[i]
		for _, spot := range row {
			if spot != "." {
				p1Alive = true
				break
			}
		}
		if p1Alive {
			break
		}
	}
	for _, i := range rows2 {
		row := g.board[i]
		for _, spot := range row {
			if spot != "." {
				p2Alive = true
				break
			}
		}
		if p2Alive {
			break
		}
	}
	if p1Alive && p2Alive {
		return false, player.ComputerPlayer{}
	} else {
		if g.p1points == g.p2points {
			return true, player.HumanPlayer{"DRAW"}
		} else if g.p1points > g.p2points {
			return true, g.p1
		} else {
			return true, g.p2
		}
	}
}

func (g MartianChess) CurrentScore(p game.Player) int {
	if p == g.p1 {
		return g.p1points - g.p2points
	} else {
		return g.p2points - g.p1points
	}
}
