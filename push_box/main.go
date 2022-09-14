package main

import (
	"os"
)

func main() {
	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	runGame(path)
}
