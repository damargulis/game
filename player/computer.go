package player

import (
	"github.com/damargulis/game/interfaces"
	"math/rand"
)

type ComputerPlayer struct {
	Name string
}

func (p ComputerPlayer) GetName() string {
	return p.Name
}

func (p ComputerPlayer) GetTurn(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	return moves[rand.Intn(len(moves))]
}
