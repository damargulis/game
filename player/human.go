package player

import (
	"fmt"
	"github.com/damargulis/game/interfaces"
)

type HumanPlayer struct {
	Name string
}

func (p HumanPlayer) GetTurn(g game.Game) game.Move {
	fmt.Println(p.Name + " Take Your Turn: ")
	return g.GetHumanInput()
}

func (p HumanPlayer) GetName() string {
	return p.Name
}
