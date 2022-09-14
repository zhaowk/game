package game

// Game interface
type Game interface {
	Init(...interface{}) error
	Run(int, string)
	Next() bool
	Finish()
}

// RunGame run Game with args
func RunGame(g Game, args ...interface{}) {
	if g == nil {
		panic("empty game")
	}

	if err := g.Init(args...); err != nil {
		panic(err)
	}

	for g.Next() {
		g.Run(GetCh())
	}

	g.Finish()
}
