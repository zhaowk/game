package main

import (
	"github.com/zhaowk/game"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		game.RunGame(&pushBoxMul{}, os.Args[1])
	} else {
		game.RunGame(&pushBoxMul{})
	}
}
