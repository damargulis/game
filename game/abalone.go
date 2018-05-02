package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
	"math"
)

type Abalone struct {
	board [9][9]string
	p1    game.Player
	p2    game.Player
	pTurn bool
	round int
}

type AbaloneMove struct {
	startRow, startCol, endRow, endCol, moveRow, moveCol int
}

func NewAbalone(p1 string, p2 string, depth1 int, depth2 int) *Abalone {
	g := new(Abalone)
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.board = [9][9]string{
		{" ", " ", " ", " ", "O", "O", ".", ".", "."},
		{" ", " ", " ", "O", "O", ".", ".", ".", "."},
		{" ", " ", "O", "O", "O", ".", ".", ".", "."},
		{" ", "O", "O", "O", ".", ".", ".", ".", "X"},
		{"O", "O", "O", ".", ".", ".", "X", "X", "X"},
		{"O", ".", ".", ".", ".", "X", "X", "X", " "},
		{".", ".", ".", ".", "X", "X", "X", " ", " "},
		{".", ".", ".", ".", "X", "X", " ", " ", " "},
		{".", ".", ".", "X", "X", " ", " ", " ", " "},
	}
	g.round = 0
	return g
}

func (g Abalone) GetBoardDimensions() (int, int) {
	return len(g.board), len(g.board[0])
}

func (g Abalone) BoardString() string {
	s := "----------------------\n"
	rowLabel := 0
	startRow := 0
	startCol := 0
	buffer := 0
	for g.board[startRow][startCol] == " " {
		startRow++
		buffer++
	}

	for isInside(g, startRow, startCol) {
		i := 0
		for i < buffer {
			s += " "
			i++
		}
		buffer--
		s += fmt.Sprintf("%v ", rowLabel)
		rowLabel++
		curRow := startRow
		curCol := startCol
		for isInside(g, curRow, curCol) {
			s += fmt.Sprintf("%v ", g.board[curRow][curCol])
			curRow--
			curCol++
		}
		startRow++
		s += "\n"
	}
	startRow--
	startCol++
	buffer = 1
	for g.board[startRow][startCol] != " " {
		i := 0
		for i < buffer {
			s += " "
			i++
		}
		buffer++
		s += fmt.Sprintf("%v ", rowLabel)
		rowLabel++
		curRow := startRow
		curCol := startCol
		for isInside(g, curRow, curCol) {
			s += fmt.Sprintf("%v ", g.board[curRow][curCol])
			curRow--
			curCol++
		}
		s += fmt.Sprintf("%v", len(g.board)-startCol)
		startCol++
		s += "\n"
	}

	s += "       0 1 2 3 4\n"
	s += "--------------------\n"
	return s
}

func (g Abalone) PrintBoard() {
	fmt.Println(g.BoardString())
}

func (g Abalone) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Abalone) humanToGrid(row int, col int) (int, int) {
	rRow := len(g.board) - col - 1
	rCol := row + ((len(g.board) / 2) - rRow)
	return rRow, rCol
}

func (g Abalone) GetHumanInput() game.Move {
	spot1 := readInts("Select start marble: ")
	spot2 := readInts("Select end marble: ")
	spot3 := readInts("Move 1st marble to: ")
	row1I, col1I := g.humanToGrid(spot1[0], spot1[1])
	row2I, col2I := g.humanToGrid(spot2[0], spot2[1])
	row3I, col3I := g.humanToGrid(spot3[0], spot3[1])
	return AbaloneMove{
		startRow: row1I,
		startCol: col1I,
		endRow:   row2I,
		endCol:   col2I,
		moveRow:  row3I,
		moveCol:  col3I,
	}
}

func (g Abalone) getBroadMoves(i, j, rowDir, colDir int, own, target string) []game.Move {
	var moves []game.Move
	if isInside(g, i+rowDir, j+colDir) && g.board[i+rowDir][j+colDir] == "." {
		broadRowDir := rowDir + 1
		broadColDir := colDir + 1
		if broadRowDir >= 2 {
			broadRowDir = -1
		}
		if broadColDir >= 2 {
			broadColDir = -1
		}
		if isInside(g, i+broadRowDir, j+broadColDir) &&
			g.board[i+broadRowDir][j+broadColDir] == own &&
			isInside(g, i+broadRowDir+rowDir, j+broadColDir+colDir) &&
			g.board[i+broadRowDir+rowDir][j+broadColDir+colDir] == "." {
			moves = append(moves, AbaloneMove{
				startRow: i,
				startCol: j,
				endRow:   i + broadRowDir,
				endCol:   j + broadColDir,
				moveRow:  i + rowDir,
				moveCol:  j + colDir,
			})
			if isInside(g, i+broadRowDir*2, j+broadColDir*2) &&
				g.board[i+broadRowDir*2][j+broadColDir*2] == own &&
				isInside(g, i+broadRowDir*2+rowDir, j+broadColDir*2) &&
				g.board[i+broadRowDir*2+rowDir][j+broadColDir*2+colDir] == "." {
				moves = append(moves, AbaloneMove{
					startRow: i,
					startCol: j,
					endRow:   i + broadRowDir*2,
					endCol:   j + broadColDir*2,
					moveRow:  i + rowDir,
					moveCol:  j + colDir,
				})
			}
		}
		broadRowDir++
		broadColDir++
		if broadRowDir >= 2 {
			broadRowDir = -1
		}
		if broadColDir >= 2 {
			broadColDir = -1
		}
		if isInside(g, i+broadRowDir, j+broadColDir) &&
			g.board[i+broadRowDir][j+broadColDir] == own &&
			isInside(g, i+broadRowDir+rowDir, j+broadColDir+colDir) &&
			g.board[i+broadRowDir+rowDir][j+broadColDir+colDir] == "." {
			moves = append(moves, AbaloneMove{
				startRow: i,
				startCol: j,
				endRow:   i + broadRowDir,
				endCol:   j + broadColDir,
				moveRow:  i + rowDir,
				moveCol:  j + colDir,
			})
			if isInside(g, i+broadRowDir*2, j+broadColDir*2) &&
				g.board[i+broadRowDir*2][j+broadColDir*2] == own &&
				isInside(g, i+broadRowDir*2+rowDir, j+broadColDir*2+colDir) &&
				g.board[i+broadRowDir*2+rowDir][j+broadColDir*2+colDir] == "." {
				moves = append(moves, AbaloneMove{
					startRow: i,
					startCol: j,
					endRow:   i + broadRowDir*2,
					endCol:   j + broadColDir*2,
					moveRow:  i + rowDir,
					moveCol:  j + colDir,
				})
			}
		}
	}
	return moves
}

func (g Abalone) getArrowMoves(i, j, rowDir, colDir int, owns, target string) []game.Move {
	var moves []game.Move
	if isInside(g, i+rowDir, j+colDir) &&
		(g.board[i+rowDir][j+colDir] == "." ||
			g.board[i+rowDir][j+colDir] == target) {
		backing := 0
		attacking := 0
		curRow := i - rowDir
		curCol := j - colDir
		for isInside(g, curRow, curCol) &&
			g.board[curRow][curCol] == owns {
			backing++
			curRow -= rowDir
			curCol -= colDir
		}
		curRow = i + rowDir
		curCol = j + colDir
		for isInside(g, curRow, curCol) &&
			g.board[curRow][curCol] == target {
			attacking++
			curRow += rowDir
			curCol += colDir
		}
		if isInside(g, curRow, curCol) &&
			g.board[curRow][curCol] == owns {
			return moves
		}
		var p int
		if backing > 2 {
			p = 2
		} else {
			p = backing
		}
		for ; p >= attacking; p-- {
			moves = append(moves, AbaloneMove{
				startRow: i,
				startCol: j,
				endRow:   i - (p * rowDir),
				endCol:   j - (p * colDir),
				moveRow:  i + rowDir,
				moveCol:  j + colDir,
			})
		}
	}
	return moves
}

func (g Abalone) GetPossibleMoves() []game.Move {
	var moves []game.Move
	var owns, target string
	if g.pTurn {
		owns = "X"
		target = "O"
	} else {
		owns = "O"
		target = "X"
	}
	for i, row := range g.board {
		for j, spot := range row {
			if spot == owns {
				moves = append(moves, g.getArrowMoves(i, j, 1, 0, owns, target)...)
				moves = append(moves, g.getArrowMoves(i, j, 1, -1, owns, target)...)
				moves = append(moves, g.getArrowMoves(i, j, 0, -1, owns, target)...)
				moves = append(moves, g.getArrowMoves(i, j, -1, 0, owns, target)...)
				moves = append(moves, g.getArrowMoves(i, j, -1, 1, owns, target)...)
				moves = append(moves, g.getArrowMoves(i, j, 0, 1, owns, target)...)

				moves = append(moves, g.getBroadMoves(i, j, 1, 0, owns, target)...)
				moves = append(moves, g.getBroadMoves(i, j, 1, -1, owns, target)...)
				moves = append(moves, g.getBroadMoves(i, j, 0, -1, owns, target)...)
				moves = append(moves, g.getBroadMoves(i, j, -1, 0, owns, target)...)
				moves = append(moves, g.getBroadMoves(i, j, -1, 1, owns, target)...)
				moves = append(moves, g.getBroadMoves(i, j, 0, 1, owns, target)...)
			}
		}
	}
	return moves
}

func (g Abalone) GetTurn(p game.Player) game.Move {
	m := p.GetTurn(g)
	return m
}

func (g Abalone) MakeMove(m game.Move) game.Game {
	g.round++
	move := m.(AbaloneMove)
	var marbleRowDir, marbleColDir int
	var moveRowDir, moveColDir int
	if move.endRow-move.startRow == 0 {
		marbleRowDir = 0
	} else {
		marbleRowDir = (move.endRow - move.startRow) / int(math.Abs(float64(move.endRow-move.startRow)))
	}
	if move.endCol-move.startCol == 0 {
		marbleColDir = 0
	} else {
		marbleColDir = (move.endCol - move.startCol) / int(math.Abs(float64(move.endCol-move.startCol)))
	}
	if move.moveRow-move.startRow == 0 {
		moveRowDir = 0
	} else {
		moveRowDir = (move.moveRow - move.startRow) / int(math.Abs(float64(move.moveRow-move.startRow)))
	}
	if move.moveCol-move.startCol == 0 {
		moveColDir = 0
	} else {
		moveColDir = (move.moveCol - move.startCol) / int(math.Abs(float64(move.moveCol-move.startCol)))
	}

	var movingMarbles []string
	curRow := move.startRow
	curCol := move.startCol
	for curRow != move.endRow+marbleRowDir || curCol != move.endCol+marbleColDir {
		movingMarbles = append(movingMarbles, g.board[curRow][curCol])
		g.board[curRow][curCol] = "."
		curRow += marbleRowDir
		curCol += marbleColDir

	}
	if marbleRowDir == 0 && marbleColDir == 0 {
		movingMarbles = append(movingMarbles, g.board[curRow][curCol])
		g.board[curRow][curCol] = "."
	}
	curRow = move.moveRow
	curCol = move.moveCol
	var replacing string
	if isInside(g, curRow, curCol) {
		replacing = g.board[curRow][curCol]
	} else {
		replacing = " "
	}
	for _, marble := range movingMarbles {
		g.board[curRow][curCol] = marble
		curRow += marbleRowDir
		curCol += marbleColDir
	}
	g.pTurn = !g.pTurn
	if replacing == "X" || replacing == "O" {
		curRow = move.moveRow + moveRowDir
		curCol = move.moveCol + moveColDir
		var newSpot string
		if isInside(g, curRow, curCol) {
			newSpot = g.board[curRow][curCol]
		} else {
			newSpot = " "
		}
		for newSpot == "X" || newSpot == "O" {
			curRow += moveRowDir
			curCol += moveColDir
			if isInside(g, curRow, curCol) {
				newSpot = g.board[curRow][curCol]
			} else {
				newSpot = " "
			}
		}
		if newSpot == "." {
			g.board[curRow][curCol] = replacing
		}
	}
	return g
}

func (g Abalone) GameOver() (bool, game.Player) {
	if g.round > 500 {
		return true, player.HumanPlayer{"DRAW"}
	}
	if len(g.GetPossibleMoves()) == 0 {
		return true, player.HumanPlayer{"DRAW"}
	}
	p1left := 0
	p2left := 0
	for _, row := range g.board {
		for _, spot := range row {
			if spot == "X" {
				p1left++
			} else if spot == "O" {
				p2left++
			}
		}
	}
	if p1left <= 8 {
		return true, g.p2
	} else if p2left <= 8 {
		return true, g.p1
	}
	return false, player.ComputerPlayer{}
}

func (g Abalone) _distToEdge(row, col, rowDir, colDir int) int {
	curRow := row
	curCol := col
	dist := 0
	for isInside(g, curRow, curCol) && g.board[curRow][curCol] != " " {
		curRow += rowDir
		curCol += colDir
		dist++
	}
	return dist
}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func (g Abalone) distToEdge(row, col int) int {
	minDist := MaxInt
	dist := g._distToEdge(row, col, 0, 1)
	if dist < minDist {
		minDist = dist
	}
	dist = g._distToEdge(row, col, 0, -1)
	if dist < minDist {
		minDist = dist
	}
	dist = g._distToEdge(row, col, 1, 0)
	if dist < minDist {
		minDist = dist
	}
	dist = g._distToEdge(row, col, -1, 0)
	if dist < minDist {
		minDist = dist
	}
	dist = g._distToEdge(row, col, -1, 1)
	if dist < minDist {
		minDist = dist
	}
	dist = g._distToEdge(row, col, 1, -1)
	if dist < minDist {
		minDist = dist
	}
	return minDist
}

func (g Abalone) CurrentScore(p game.Player) int {
	score := 0
	var target, owns string
	if p == g.p1 {
		target = "O"
		owns = "X"
	} else {
		target = "X"
		owns = "O"
	}
	for _, row := range g.board {
		for _, spot := range row {
			if spot == owns {
				score += 1
			} else if spot == target {
				score -= 1
			}
		}
	}
	return score
}
