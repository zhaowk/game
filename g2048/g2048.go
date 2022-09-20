package main

import (
	"fmt"
	"github.com/zhaowk/game"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	gWall    = "#"
	gBall    = "@"
	gSpace   = " "
	gGameMax = 2048
)

type block uint16

func (b block) String() string {
	switch {
	case b == 0:
		return strings.Repeat(" ", 4)
	case b < 10:
		return fmt.Sprintf("   %d", b)
	case b < 100:
		return fmt.Sprintf("  %d", b)
	case b < 1000:
		return fmt.Sprintf(" %d", b)
	case b < 10000:
		return fmt.Sprintf("%d", b)
	default:
		return strings.Repeat(" ", 4)
	}
}

type g2048 struct {
	size int
	pane [][]block
}

func (g *g2048) Init(...interface{}) error {
	rand.Seed(time.Now().UnixNano())
	g.size = 4
	g.pane = make([][]block, g.size)
	for i := 0; i < g.size; i++ {
		g.pane[i] = make([]block, g.size)
	}

	m, n := rand.Intn(g.size*g.size), rand.Intn(g.size*g.size)
	if m == n {
		n = (n + g.size + 3) % (g.size * g.size)
	}

	g.pane[m/g.size][m%g.size] = block(2)
	g.pane[n/g.size][n%g.size] = block(2)

	g.draw()

	return nil
}

func (g *g2048) Run(k int, _ string) {
	switch k {
	case 'w', 'W', game.SysUp:
		g.moveUp()
	case 's', 'S', game.SysDown:
		g.moveDown()
	case 'a', 'A', game.SysLeft:
		g.moveLeft()
	case 'd', 'D', game.SysRight:
		g.moveRight()
	case 'q', 'Q':
		game.DrawLine("Exiting...")
		os.Exit(0)
	}
}

func (g *g2048) Next() bool {
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.pane[i][j] >= gGameMax {
				game.DrawLine("Congratulations!")
				return false
			}
		}
	}
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.pane[i][j] == 0 {
				return true
			}

			if j > 0 && g.pane[i][j] == g.pane[i][j-1] {
				return true
			}

			if i > 0 && g.pane[i][j] == g.pane[i-1][j] {
				return true
			}
		}
	}
	return false
}

func (g *g2048) Finish() {}

func (g *g2048) draw() {
	game.Clear()
	game.DrawLineAt(game.Point{}, strings.Repeat(gWall, g.size*7+1))

	for i := 0; i < g.size; i++ {
		for k := 0; k < 3; k++ {
			game.Draw(gWall)
			for j := 0; j < g.size; j++ {
				game.Draw(fmt.Sprintf("%s%s", g.getContent(k, g.pane[i][j]), gSpace))
			}
			game.CursorBack(1)
			game.DrawLine(gWall)
		}
		game.DrawLine(fmt.Sprintf("%s%s%s", gWall, strings.Repeat(gSpace, g.size*7-1), gWall))
	}
	game.CursorUp(1)
	game.DrawLine(strings.Repeat(gWall, g.size*7+1))
}

func (g *g2048) getContent(level int, b block) string {
	if b == 0 {
		return strings.Repeat(gSpace, 6)
	}
	switch level {
	case 0, 2:
		return strings.Repeat(gBall, 6)
	case 1:
		return fmt.Sprintf("%s%s%s", gBall, b.String(), gBall)
	default:
		return strings.Repeat(gSpace, 6)
	}
}

func (g *g2048) move(rev bool, genFunc func(i int) []block, setFunc func(x int, line []block)) {
	for i := 0; i < g.size; i++ {
		setFunc(i, g.doMerge(rev, genFunc(i)))
	}

	g.genNext()
	g.draw()
}

func (g *g2048) moveLeft() {
	g.move(false, g.genRow, g.setRow)
}

func (g *g2048) moveRight() {
	g.move(true, g.genRow, g.setRow)
}

func (g *g2048) moveUp() {
	g.move(false, g.genCol, g.setCol)
}

func (g *g2048) moveDown() {
	g.move(true, g.genCol, g.setCol)
}

func (g *g2048) genRow(i int) []block {
	tmp := make([]block, len(g.pane[i]))
	copy(tmp, g.pane[i])
	return g.pane[i]
}

func (g *g2048) genCol(i int) []block {
	var tmp []block
	for j := 0; j < g.size; j++ {
		tmp = append(tmp, g.pane[j][i])
	}
	return tmp
}

func (g *g2048) setRow(x int, line []block) {
	for i := 0; i < len(line); i++ {
		if x >= 0 && x < g.size {
			g.pane[x][i] = line[i]
		}
	}
}

func (g *g2048) setCol(y int, line []block) {
	for i := 0; i < len(line); i++ {
		if y >= 0 && y < g.size {
			g.pane[i][y] = line[i]
		}
	}
}

func (g *g2048) doMerge(reverse bool, item []block) []block {
	sort.Slice(item, func(i, j int) bool {
		return (reverse && i > j) || i < j
	})
	target := make([]block, 0)

	for _, b := range item {
		if b != 0 {
			target = append(target, b)
		}
	}

	var complete bool

	for !complete {
		complete = true
		for i := 1; i < len(target); i++ {
			if target[i] == target[i-1] {
				complete = false
				target[i-1] <<= 1

				var tmp []block
				if i < len(target)-1 {
					tmp = target[i+1:]
				}
				target = append(target[:i], tmp...)
			}
		}
	}

	sort.Slice(target, func(i, j int) bool {
		return (reverse && i > j) || i < j
	})

	tmp := make([]block, len(item)-len(target))
	if reverse {
		return append(tmp, target...)
	} else {
		return append(target, tmp...)
	}
}

func (g *g2048) genNext() {
	empty := make([][2]int, 0)
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if g.pane[i][j] == 0 {
				empty = append(empty, [2]int{i, j})
			}
		}
	}

	if len(empty) == 0 {
		game.DrawLine("Game over!")
		os.Exit(0)
	}

	pos := empty[rand.Intn(len(empty))]
	g.pane[pos[0]][pos[1]] = 2
}
