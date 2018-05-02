package game

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
	"github.com/damargulis/game/player"
)

type Boxes struct {
	vertical   [4][5]bool
	horizontal [5][4]bool
	owners     [4][4]string
	p1         game.Player
	p2         game.Player
	pTurn      bool
}

type BoxesMove struct {
	horizontal bool
	row, col   int
}

func (g Boxes) GetBoardDimensions() (int, int) {
	return len(g.owners), len(g.owners)
}

func (g Boxes) BoardString() string {
	s := "***********\n"
	s += "  012345678\n"
	rowNum := 0
	for i, row := range g.horizontal {
		s += fmt.Sprintf("%v ", rowNum)
		rowNum++
		for _, spot := range row {
			s += "."
			if spot {
				s += "-"
			} else {
				s += " "
			}
		}
		s += ".\n"
		if i >= len(g.vertical) {
			continue
		}
		s += fmt.Sprintf("%v ", rowNum)
		rowNum++
		for j, spot := range g.vertical[i] {
			if spot {
				s += "|"
			} else {
				s += " "
			}
			if i < len(g.owners) && j < len(g.owners[i]) {
				s += g.owners[i][j]
			}
		}
		s += "\n"
	}
	s += "  012345678\n"
	s += "***********\n"
	return s
}

func (g Boxes) PrintBoard() {
	fmt.Println("***********")
	fmt.Println("  012345678")
	rowNum := 0
	for i, row := range g.horizontal {
		fmt.Print(rowNum)
		fmt.Print(" ")
		rowNum++
		for _, spot := range row {
			fmt.Print(".")
			if spot {
				fmt.Print("-")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print(".")
		fmt.Print("\n")
		if i >= len(g.vertical) {
			continue
		}
		fmt.Print(rowNum)
		fmt.Print(" ")
		rowNum++
		for j, spot := range g.vertical[i] {
			if spot {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
			if i < len(g.owners) && j < len(g.owners[i]) {
				fmt.Print(g.owners[i][j])
			}
		}
		fmt.Print("\n")
	}
	fmt.Println("  012345678")
	fmt.Println("***********")
}

func (g Boxes) GetPlayerTurn() game.Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g Boxes) GetHumanInput() game.Move {
	colI, rowI := 0, 0
	for (colI+rowI)%2 == 0 {
		spot := readInts("Place a line at: ")
		rowI, colI = spot[0], spot[1]
	}
	var horizontal bool
	if rowI%2 == 0 {
		horizontal = true
	} else {
		horizontal = false
	}
	return BoxesMove{
		row:        rowI / 2,
		col:        colI / 2,
		horizontal: horizontal,
	}
}

func (g Boxes) GetPossibleMoves() []game.Move {
	var moves []game.Move
	for i, row := range g.horizontal {
		for j, spot := range row {
			if !spot {
				moves = append(moves, BoxesMove{
					row:        i,
					col:        j,
					horizontal: true,
				})
			}
		}
	}
	for i, row := range g.vertical {
		for j, spot := range row {
			if !spot {
				moves = append(moves, BoxesMove{
					row:        i,
					col:        j,
					horizontal: false,
				})
			}
		}
	}
	return moves
}

func (g Boxes) isGoodMove(m BoxesMove) bool {
	row := m.row
	col := m.col
	if m.horizontal {
		if row < len(g.horizontal) && col < len(g.horizontal[row]) {
			return !g.horizontal[row][col]
		}
	} else {
		if row < len(g.vertical) && col < len(g.vertical[row]) {
			return !g.vertical[row][col]
		}
	}
	return false
}

func (g Boxes) GetTurn(p game.Player) game.Move {
	m := p.GetTurn(g)
	move := m.(BoxesMove)
	for !g.isGoodMove(move) {
		m = p.GetTurn(g)
		move = m.(BoxesMove)
	}
	return m
}

func (g Boxes) MakeMove(move game.Move) game.Game {
	didClaim := false
	m := move.(BoxesMove)
	if m.horizontal {
		g.horizontal[m.row][m.col] = true
		if m.row-1 >= 0 && g.horizontal[m.row-1][m.col] {
			if g.vertical[m.row-1][m.col] && g.vertical[m.row-1][m.col+1] {
				if g.pTurn {
					g.owners[m.row-1][m.col] = "X"
				} else {
					g.owners[m.row-1][m.col] = "O"
				}
				didClaim = true
			}
		}
		if m.row+1 < len(g.horizontal) && g.horizontal[m.row+1][m.col] {
			if g.vertical[m.row][m.col] && g.vertical[m.row][m.col+1] {
				if g.pTurn {
					g.owners[m.row][m.col] = "X"
				} else {
					g.owners[m.row][m.col] = "O"
				}
				didClaim = true
			}
		}
	} else {
		g.vertical[m.row][m.col] = true
		if m.col-1 >= 0 && g.vertical[m.row][m.col-1] {
			if g.horizontal[m.row][m.col-1] && g.horizontal[m.row][m.col-1] {
				if g.pTurn {
					g.owners[m.row][m.col-1] = "X"
				} else {
					g.owners[m.row][m.col-1] = "O"
				}
				didClaim = true
			}
		}
		if m.col+1 < len(g.vertical[m.row]) && g.vertical[m.row][m.col+1] {
			if g.horizontal[m.row][m.col] && g.horizontal[m.row+1][m.col] {
				if g.pTurn {
					g.owners[m.row][m.col] = "X"
				} else {
					g.owners[m.row][m.col] = "O"
				}
				didClaim = true
			}
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
	for _, row := range g.owners {
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
	g.p1 = getPlayer(p1, "Player 1", depth1)
	g.p2 = getPlayer(p2, "Player 2", depth2)
	g.pTurn = true
	g.owners = [4][4]string{
		{" ", " ", " ", " "},
		{" ", " ", " ", " "},
		{" ", " ", " ", " "},
		{" ", " ", " ", " "},
	}
	g.vertical = [4][5]bool{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}
	g.horizontal = [5][4]bool{
		{false, false, false, false},
		{false, false, false, false},
		{false, false, false, false},
		{false, false, false, false},
		{false, false, false, false},
	}
	return g
}
