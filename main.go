package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Player interface {
	getTurn(Game) Move
	getName() string
}

type HumanPlayer struct {
	name string
}

type Move struct {
	row, col int
}

func (p HumanPlayer) getTurn(g Game) Move {
	fmt.Println(p.name + " Take Your Turn: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	move := strings.Split(strings.TrimSpace(text), ",")
	row, col := move[0], move[1]
	rowI, _ := strconv.Atoi(row)
	colI, _ := strconv.Atoi(col)
	return Move{row: rowI, col: colI}
}

func (p HumanPlayer) getName() string {
	return p.name
}

type ComputerPlayer struct{}

func (p ComputerPlayer) getName() string {
	return "Computer"
}

func (p ComputerPlayer) getTurn(g Game) Move {
	row := int(math.Floor(rand.Float64() * 3))
	col := int(math.Floor(rand.Float64() * 3))
	return Move{row: row, col: col}
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
		go p.getMin(g, i, move, ch)
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

func (p MinimaxPlayer) getMax(g Game, i int, m Move, ch chan moveVal) {
	newG := g.makeNewMove(m)
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
			go p.getMin(newG, j, move, newCh)
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

func (p MinimaxPlayer) getMin(g Game, i int, m Move, ch chan moveVal) {
	newG := g.makeNewMove(m)
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
			go p.getMax(newG, j, move, newCh)
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

type Game struct {
	board [3][3]string
	p1    Player
	p2    Player
	pTurn bool
}

func (g Game) gameOver() (bool, Player) {
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

func newGame(p1, p2 string) *Game {
	g := new(Game)
	switch p1 {
	case "Human":
		g.p1 = HumanPlayer{"Player 1"}
	case "Computer":
		g.p1 = ComputerPlayer{}
	case "Minimax":
		g.p1 = MinimaxPlayer{}
	default:
		fmt.Println("Player 1 not recognized")
		os.Exit(1)
	}
	switch p2 {
	case "Human":
		g.p2 = HumanPlayer{"Player 1"}
	case "Computer":
		g.p2 = ComputerPlayer{}
	case "Minimax":
		g.p2 = MinimaxPlayer{}
	default:
		fmt.Println("Player 2 not recognized")
		os.Exit(1)
	}
	g.pTurn = true
	g.board = [3][3]string{
		{".", ".", "."},
		{".", ".", "."},
		{".", ".", "."},
	}
	return g
}

func (g Game) makeNewMove(m Move) Game {
	row := m.row
	col := m.col
	if g.pTurn {
		g.board[row][col] = "X"
	} else {
		g.board[row][col] = "O"
	}
	g.pTurn = !g.pTurn
	return g
}

func (g *Game) makeMove(m Move) {
	row := m.row
	col := m.col
	if g.pTurn {
		g.board[row][col] = "X"
	} else {
		g.board[row][col] = "O"
	}
}

func (g Game) isGoodMove(m Move) bool {
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

func (g Game) getPossibleMoves() []Move {
	var moves []Move
	for i, row := range g.board {
		for j, spot := range row {
			if spot == "." {
				moves = append(moves, Move{row: i, col: j})
			}
		}
	}
	return moves
}

func (g Game) getTurn(p Player) Move {
	m := p.getTurn(g)
	for !g.isGoodMove(m) {
		m = p.getTurn(g)
	}
	return m
}

func (g Game) play() {
	var winner Player
	over := false
	for ; !over; over, winner = g.gameOver() {
		g.printBoard()
		var m Move
		if g.pTurn {
			m = g.getTurn(g.p1)
		} else {
			m = g.getTurn(g.p2)
		}
		(&g).makeMove(m)
		g.pTurn = !g.pTurn
	}
	g.printBoard()
	name := winner.getName()
	if name == "DRAW" {
		fmt.Println("Its a draw!")
	} else {
		fmt.Println(name + " Wins!")
	}
}

func (g Game) printBoard() {
	fmt.Println("---")
	for _, row := range g.board {
		for _, p := range row {
			fmt.Print(p)
		}
		fmt.Println()
	}
	fmt.Println("---")
}

func main() {
	p1 := flag.String("p1", "Human", "Player 1")
	p2 := flag.String("p2", "Computer", "Player 2")
	flag.Parse()
	g := newGame(*p1, *p2)
	g.play()
}
