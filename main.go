package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const MAX_DEPTH = 5

type Player interface {
	getTurn(Game) Move
	getName() string
}

type HumanPlayer struct {
	name string
}

type Move interface {
}

type TicTacToeMove struct {
	row, col int
}

type CheckersMove struct {
	row1, col1, row2, col2 int
}

func (g TicTacToe) getHumanInput() Move {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	move := strings.Split(strings.TrimSpace(text), ",")
	row, col := move[0], move[1]
	rowI, _ := strconv.Atoi(row)
	colI, _ := strconv.Atoi(col)
	return TicTacToeMove{row: rowI, col: colI}
}

func (g Checkers) getHumanInput() Move {
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

func (p HumanPlayer) getTurn(g Game) Move {
	fmt.Println(p.name + " Take Your Turn: ")
	return g.getHumanInput()
}

func (p HumanPlayer) getName() string {
	return p.name
}

type ComputerPlayer struct{}

func (p ComputerPlayer) getName() string {
	return "Computer"
}

func (p ComputerPlayer) getTurn(g Game) Move {
	moves := g.getPossibleMoves()
	return moves[rand.Intn(len(moves))]
}

type MinimaxPlayer struct{}

func (p MinimaxPlayer) getName() string {
	return "Minimax Player"
}

type moveVal struct {
	move int
	val  int
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func (p MinimaxPlayer) getTurn(g Game) Move {
	moves := g.getPossibleMoves()
	ch := make(chan moveVal)
	scores := make([]int, len(moves))
	for i, move := range moves {
		go p.getMin(g, i, move, ch, 0)
	}
	i := 0
	for i < len(scores) {
		moveVal := <-ch
		i++
		scores[moveVal.move] = moveVal.val
	}
	bestScore := MinInt
	var bestMove Move
	for i, score := range scores {
		if score >= bestScore {
			bestScore = score
			bestMove = moves[i]
		}
	}
	return bestMove
}

func (p MinimaxPlayer) getMax(g Game, i int, m Move, ch chan moveVal, depth int) {
	if depth > MAX_DEPTH {
		ch <- moveVal{move: i, val: g.currentScore(p)}
		return
	}
	newG := g.makeMove(m)
	over, winner := newG.gameOver()
	if over {
		if winner == p {
			ch <- moveVal{move: i, val: MaxInt}
		} else if winner.getName() == "DRAW" {
			ch <- moveVal{move: i, val: 0}
		} else {
			ch <- moveVal{move: i, val: MinInt}
		}
	} else {
		moves := newG.getPossibleMoves()
		newCh := make(chan moveVal)
		scores := make([]int, len(moves))
		for j, move := range moves {
			go p.getMin(newG, j, move, newCh, depth+1)
		}
		j := 0
		for j < len(scores) {
			moveVal := <-newCh
			j++
			scores[moveVal.move] = moveVal.val
		}
		bestScore := MinInt
		for _, score := range scores {
			if score >= bestScore {
				bestScore = score
			}
		}
		ch <- moveVal{move: i, val: bestScore}
	}
}

func (p MinimaxPlayer) getMin(g Game, i int, m Move, ch chan moveVal, depth int) {
	if depth > MAX_DEPTH {
		ch <- moveVal{move: i, val: g.currentScore(p)}
		return
	}
	newG := g.makeMove(m)
	over, winner := newG.gameOver()
	if over {
		if winner == p {
			ch <- moveVal{move: i, val: MaxInt}
		} else if winner.getName() == "DRAW" {
			ch <- moveVal{move: i, val: 0}
		} else {
			ch <- moveVal{move: i, val: MinInt}
		}
	} else {
		moves := newG.getPossibleMoves()
		newCh := make(chan moveVal)
		scores := make([]int, len(moves))
		for j, move := range moves {
			go p.getMax(newG, j, move, newCh, depth+1)
		}
		j := 0
		for j < len(scores) {
			moveVal := <-newCh
			j++
			scores[moveVal.move] = moveVal.val
		}
		bestScore := MaxInt
		for _, score := range scores {
			if score <= bestScore {
				bestScore = score
			}
		}
		ch <- moveVal{move: i, val: bestScore}
	}
}

type Game interface {
	gameOver() (bool, Player)
	makeMove(Move) Game
	getPossibleMoves() []Move
	printBoard()
	getTurn(Player) Move
	getPlayerTurn() Player
	getHumanInput() Move
	currentScore(Player) int
}

type TicTacToe struct {
	board [3][3]string
	p1    Player
	p2    Player
	pTurn bool
}

type Checkers struct {
	board [8][8]string
	p1    Player
	p2    Player
	pTurn bool
}

func (g Checkers) gameOver() (bool, Player) {
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
		moves := g.getPossibleMoves()
		if len(moves) == 0 {
			return true, HumanPlayer{"DRAW"}
		}
		return false, ComputerPlayer{}
	}
}

func (g TicTacToe) gameOver() (bool, Player) {
	if g.board[0][0] == g.board[0][1] && g.board[0][0] == g.board[0][2] {
		if g.board[0][0] == "X" {
			return true, g.p1
		} else if g.board[0][0] == "O" {
			return true, g.p2
		}
	} else if g.board[1][0] == g.board[1][1] && g.board[1][0] == g.board[1][2] {
		if g.board[1][0] == "X" {
			return true, g.p1
		} else if g.board[1][0] == "O" {
			return true, g.p2
		}
	} else if g.board[2][0] == g.board[2][1] && g.board[2][0] == g.board[2][2] {
		if g.board[2][0] == "X" {
			return true, g.p1
		} else if g.board[2][0] == "O" {
			return true, g.p2
		}
	} else if g.board[0][0] == g.board[1][0] && g.board[0][0] == g.board[2][0] {
		if g.board[0][0] == "X" {
			return true, g.p1
		} else if g.board[0][0] == "O" {
			return true, g.p2
		}
	} else if g.board[0][1] == g.board[1][1] && g.board[0][1] == g.board[2][1] {
		if g.board[0][1] == "X" {
			return true, g.p1
		} else if g.board[0][1] == "O" {
			return true, g.p2
		}
	} else if g.board[0][2] == g.board[1][2] && g.board[0][2] == g.board[2][2] {
		if g.board[0][2] == "X" {
			return true, g.p1
		} else if g.board[0][2] == "O" {
			return true, g.p2
		}
	} else if g.board[0][0] == g.board[1][1] && g.board[0][0] == g.board[2][2] {
		if g.board[0][0] == "X" {
			return true, g.p1
		} else if g.board[0][0] == "O" {
			return true, g.p2
		}
	} else if g.board[2][0] == g.board[1][1] && g.board[1][1] == g.board[0][2] {
		if g.board[2][0] == "X" {
			return true, g.p1
		} else if g.board[2][0] == "O" {
			return true, g.p2
		}
	}
	for _, row := range g.board {
		for _, p := range row {
			if p == "." {
				return false, HumanPlayer{"DRAW"}
			}
		}
	}
	return true, HumanPlayer{"DRAW"}
}

func getPlayer(name string) Player {
	var p Player
	switch name {
	case "Human":
		p = HumanPlayer{"Player 1"}
	case "Computer":
		p = ComputerPlayer{}
	case "Minimax":
		p = MinimaxPlayer{}
	default:
		fmt.Println("Player " + name + " not recognized")
		os.Exit(1)
	}
	return p
}

func newTicTacToe(p1, p2 string) *TicTacToe {
	g := new(TicTacToe)
	g.p1 = getPlayer(p1)
	g.p2 = getPlayer(p2)
	g.pTurn = true
	g.board = [3][3]string{
		{".", ".", "."},
		{".", ".", "."},
		{".", ".", "."},
	}
	return g
}

func newCheckers(p1, p2 string) *Checkers {
	c := new(Checkers)
	c.p1 = getPlayer(p1)
	c.p2 = getPlayer(p2)
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
	return c
}

func (g TicTacToe) makeMove(m Move) Game {
	move := m.(TicTacToeMove)
	row := move.row
	col := move.col
	if g.pTurn {
		g.board[row][col] = "X"
	} else {
		g.board[row][col] = "O"
	}
	g.pTurn = !g.pTurn
	return g
}

func (g Checkers) makeMove(m Move) Game {
	move := m.(CheckersMove)
	g.board[move.row2][move.col2] = g.board[move.row1][move.col1]
	g.board[move.row1][move.col1] = "."
	if move.row1 == move.row2+2 || move.row1 == move.row2-2 {
		rowAvg := (move.row1 + move.row2) / 2
		colAvg := (move.col1 + move.col2) / 2
		g.board[rowAvg][colAvg] = "."
	}
	if move.row2 == 0 && g.board[move.row2][move.col2] == "x" {
		g.board[move.row2][move.col2] = "X"
	} else if move.row2 == 7 && g.board[move.row2][move.col2] == "o" {
		g.board[move.row2][move.col2] = "O"
	}
	g.pTurn = !g.pTurn
	moves := g.getPossibleMoves()
	if len(moves) == 0 {
		g.pTurn = !g.pTurn
	}
	return g
}

func (g Checkers) isGoodMove(m CheckersMove) bool {
	row1, row2, col1, col2 := m.row1, m.row2, m.col1, m.col2
	if row1 < 0 || row2 < 0 || col1 < 0 || col2 < 0 {
		return false
	} else if row1 >= 8 || row2 >= 8 || col1 >= 8 || col2 >= 8 {
		return false
	} else if g.board[row2][col2] != "." {
		return false
	} else if g.pTurn && g.board[row1][col1] == "x" {
		if row1 == row2+1 && (col1 == col2+1 || col1 == col2-1) {
			return true
		} else if row1 == row2+2 && col1 == col2+2 && (g.board[row2+1][col2+1] == "O" || g.board[row2+1][col2+1] == "o") {
			return true
		} else if row1 == row2+2 && col1 == col2-2 && (g.board[row2+1][col2-1] == "O" || g.board[row2+1][col2-1] == "o") {
			return true
		} else {
			return false
		}
	} else if g.pTurn && g.board[row1][col1] == "X" {
		if row1 == row2+1 && (col1 == col2+1 || col1 == col2-1) {
			return true
		} else if row1 == row2+2 && col1 == col2+2 && (g.board[row2+1][col2+1] == "O" || g.board[row2+1][col2+1] == "o") {
			return true
		} else if row1 == row2+2 && col1 == col2-2 && (g.board[row2+1][col2-1] == "O" || g.board[row2+1][col2-1] == "o") {
			return true
		} else if row1 == row2-1 && (col1 == col2+1 || col1 == col2-1) {
			return true
		} else if row1 == row2-2 && col1 == col2+2 && (g.board[row2-1][col2+1] == "O" || g.board[row2-1][col2+1] == "o") {
			return true
		} else if row1 == row2-2 && col1 == col2-2 && (g.board[row2-1][col2-1] == "O" || g.board[row2-1][col2-1] == "o") {
			return true
		} else {
			return false
		}
	} else if !g.pTurn && g.board[row1][col1] == "o" {
		if row1 == row2-1 && (col1 == col2+1 || col1 == col2-1) {
			return true
		} else if row1 == row2-2 && col1 == col2+2 && (g.board[row2-1][col2+1] == "X" || g.board[row2-1][col2+1] == "x") {
			return true
		} else if row1 == row2-2 && col1 == col2-2 && (g.board[row2-1][col2-1] == "X" || g.board[row2-1][col2-1] == "x") {
			return true
		} else {
			return false
		}
	} else if !g.pTurn && g.board[row1][col1] == "O" {
		if row1 == row2+1 && (col1 == col2+1 || col1 == col2-1) {
			return true
		} else if row1 == row2+2 && col1 == col2+2 && (g.board[row2+1][col2+1] == "X" || g.board[row2+1][col2+1] == "x") {
			return true
		} else if row1 == row2+2 && col1 == col2-2 && (g.board[row2+1][col2-1] == "X" || g.board[row2+1][col2-1] == "x") {
			return true
		} else if row1 == row2-1 && (col1 == col2+1 || col1 == col2-1) {
			return true
		} else if row1 == row2-2 && col1 == col2+2 && (g.board[row2-1][col2+1] == "X" || g.board[row2-1][col2+1] == "x") {
			return true
		} else if row1 == row2-2 && col1 == col2-2 && (g.board[row2-1][col2-1] == "X" || g.board[row2-1][col2-1] == "x") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (g TicTacToe) isGoodMove(m TicTacToeMove) bool {
	row := m.row
	col := m.col
	if row < 0 || col < 0 {
		return false
	} else if row > 2 || col > 2 {
		return false
	} else {
		return g.board[row][col] == "."
	}
}

func (g TicTacToe) getPossibleMoves() []Move {
	var moves []Move
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				moves = append(moves, TicTacToeMove{row: i, col: j})
			}
		}
	}
	return moves
}

func (g TicTacToe) currentScore(p Player) int {
	return 1
}

func (g Checkers) currentScore(p Player) int {
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

func (g Checkers) getPossibleMoves() []Move {
	var moves []Move
	for i, row := range g.board {
		for j, spot := range row {
			if g.pTurn && (spot == "X" || spot == "x") {
				if i-1 > 0 && j-1 > 0 && g.board[i-1][j-1] == "." {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 1,
						col2: j - 1,
					})
				}
				if i-1 > 0 && j+1 < 8 && g.board[i-1][j+1] == "." {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 1,
						col2: j + 1,
					})
				}
				if i-2 > 0 && j-2 > 0 && g.board[i-2][j-2] == "." && (g.board[i-1][j-1] == "O" || g.board[i-1][j-1] == "o") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 2,
						col2: j - 2,
					})
				}
				if i-2 > 0 && j+2 < 8 && g.board[i-2][j+2] == "." && (g.board[i-1][j+1] == "O" || g.board[i-1][j+1] == "o") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i - 2,
						col2: j + 2,
					})
				}
				if spot == "X" {
					if i+1 < 8 && j-1 > 0 && g.board[i+1][j-1] == "." {
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
					if i+2 < 8 && j-2 > 0 && g.board[i+2][j-2] == "." && (g.board[i+1][j-1] == "O" || g.board[i+1][j-1] == "o") {
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
				if i+1 < 8 && j-1 > 0 && g.board[i+1][j-1] == "." {
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
				if i+2 < 8 && j-2 > 0 && (g.board[i+1][j-1] == "X" || g.board[i+1][j-1] == "x") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i + 2,
						col2: j - 2,
					})
				}
				if i+2 < 8 && j+2 < 8 && (g.board[i+1][j+1] == "X" || g.board[i+1][j+1] == "x") {
					moves = append(moves, CheckersMove{
						row1: i,
						col1: j,
						row2: i + 2,
						col2: j + 2,
					})
				}
				if spot == "O" {
					if i-1 > 0 && j-1 > 0 && g.board[i-1][j-1] == "." {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 1,
							col2: j - 1,
						})
					}
					if i-1 > 0 && j+1 < 8 && g.board[i-1][j+1] == "." {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 1,
							col2: j + 1,
						})
					}
					if i-2 > 0 && j-2 > 0 && (g.board[i-1][j-1] == "x" || g.board[i-1][j-1] == "X") {
						moves = append(moves, CheckersMove{
							row1: i,
							col1: j,
							row2: i - 2,
							col2: j - 2,
						})
					}
					if i-2 > 0 && j+2 < 8 && (g.board[i-1][j+1] == "x" || g.board[i-1][j+1] == "X") {
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
	return moves
}

func (g Checkers) getPlayerTurn() Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g TicTacToe) getPlayerTurn() Player {
	if g.pTurn {
		return g.p1
	} else {
		return g.p2
	}
}

func (g TicTacToe) getTurn(p Player) Move {
	m := p.getTurn(g)
	move := m.(TicTacToeMove)
	for !g.isGoodMove(move) {
		m = p.getTurn(g)
		move = m.(TicTacToeMove)
	}
	return m
}

func (g Checkers) getTurn(p Player) Move {
	m := p.getTurn(g)
	move := m.(CheckersMove)
	for !g.isGoodMove(move) {
		m = p.getTurn(g)
		move = m.(CheckersMove)
	}
	return m
}

func play(g Game) {
	var winner Player
	over := false
	for ; !over; over, winner = g.gameOver() {
		g.printBoard()
		move := g.getTurn(g.getPlayerTurn())
		g = g.makeMove(move)
	}
	g.printBoard()
	name := winner.getName()
	if name == "DRAW" {
		fmt.Println("Its a draw!")
	} else {
		fmt.Println(name + " Wins!")
	}
}

func (g TicTacToe) printBoard() {
	fmt.Println("---")
	for _, row := range g.board {
		for _, p := range row {
			fmt.Print(p)
		}
		fmt.Println()
	}
	fmt.Println("---")
}

func (g Checkers) printBoard() {
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

func main() {
	rand.Seed(time.Now().Unix())
	p1 := flag.String("p1", "Human", "Player 1")
	p2 := flag.String("p2", "Computer", "Player 2")
	flag.Parse()
	for i := 0; i < 10; i++ {
		fmt.Println("Game ", i)
		time.Sleep(500000000)
		fmt.Println(*p1, *p2)
		t := newTicTacToe(*p1, *p2)
		play(t)
		c := newCheckers(*p1, *p2)
		play(c)
	}
}
