package game

type Player interface {
	GetTurn(Game) Move
	GetName() string
}

type Move interface {
}

type Game interface {
	PrintBoard()
	GetPlayerTurn() Player
	GetHumanInput() Move
	GetPossibleMoves() []Move
	GetTurn(Player) Move
	MakeMove(Move) Game
	GameOver() (bool, Player)
	CurrentScore(Player) int
}
