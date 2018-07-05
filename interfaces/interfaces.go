package game

type Player interface {
	GetTurn(Game) Move
	GetName() string
}

type Move interface {
}

type Game interface {
	BoardString() string
	GetPlayerTurn() Player
	GetHumanInput() Move
	GetPossibleMoves() []Move
	MakeMove(Move) Game
	GameOver() (bool, Player)
	CurrentScore(Player) int
	GetBoardDimensions() (int, int)
	GetRound() int
}
