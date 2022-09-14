package main

import (
	"fmt"
	"os"
	"time"

	"github.com/zhaowk/game"
)

const (
	PushBoxBlank        = 0
	PushBoxTarget       = 1
	PushBoxBox          = 2
	PushBoxTargetBox    = PushBoxTarget | PushBoxBox // 3
	PushBoxPerson       = 4
	PushBoxPersonTarget = PushBoxTarget | PushBoxPerson // 5
	PushBoxWall         = 8
)

type gameItem byte

func (g gameItem) toByte() byte {
	switch g {
	case PushBoxBlank:
		return ' '
	case PushBoxTarget:
		return '.'
	case PushBoxBox:
		return 'o'
	case PushBoxTargetBox:
		return 'O'
	case PushBoxPerson:
		return 'p'
	case PushBoxPersonTarget:
		return 'P'
	case PushBoxWall:
		return '#'
	default:
		return ' '
	}
}

func (g gameItem) fromByte(b byte) gameItem {
	switch b {
	case ' ':
		return PushBoxBlank
	case '.':
		return PushBoxTarget
	case 'o':
		return PushBoxBox
	case 'O':
		return PushBoxTargetBox
	case 'p':
		return PushBoxPerson
	case 'P':
		return PushBoxPersonTarget
	case '#':
		return PushBoxWall
	default:
		return PushBoxBlank
	}
}

type gamePane [][]gameItem

// String gamePane stringer
func (g gamePane) String() (r string) {
	for _, s := range g {
		for _, c := range s {
			r += fmt.Sprintf("%c", c.toByte())
		}
		r += "\n"
	}
	return
}

type pushBox struct {
	original   boxMap
	origPerson game.Point
	runtime    gamePane
	runPerson  game.Point
	width      int
	height     int
	msg        string
}

func (g *pushBox) Init(args ...interface{}) error {
	if len(args) == 0 {
		return fmt.Errorf("empty game")
	}

	if m, ok := args[0].(boxMap); !ok {
		return fmt.Errorf("unknown game")
	} else if err := g.init(m); err != nil {
		return err
	}
	return nil
}

func (g *pushBox) Run(k int, _ string) {
	switch k {
	case 'w', 'W', game.SysUp:
		g.move(-1, 0)
	case 's', 'S', game.SysDown:
		g.move(1, 0)
	case 'a', 'A', game.SysLeft:
		g.move(0, -1)
	case 'd', 'D', game.SysRight:
		g.move(0, 1)
	case 'r', 'R':
		_ = g.init(g.original)
	case 'q', 'Q':
		os.Exit(0)
	}
}

func (g *pushBox) Next() bool {
	g.draw()
	for i := range g.runtime {
		for j := range g.runtime[i] {
			switch g.runtime[i][j] {
			case PushBoxPersonTarget, PushBoxTarget, PushBoxBox:
				return true
			}
		}
	}

	return false
}

func (g *pushBox) Finish() {
	g.msg = "congratulations!"
	g.draw()
	time.Sleep(300 * time.Millisecond)
}

func (g *pushBox) validMap(pane boxMap) (game.Point, bool) {
	target, box, person, p := 0, 0, 0, game.Point{}
	for i := range pane {
		for j := range pane[i] {
			switch pane[i][j] {
			case 'P':
				target++
				person++
				p = game.Point{X: i, Y: j}
			case 'p':
				person++
				p = game.Point{X: i, Y: j}
			case '.':
				target++
			case 'o':
				box++
			}
		}
	}

	if target != box || person != 1 {
		return p, false
	}
	return p, true
}

func (g *pushBox) init(pane boxMap) error {
	var valid bool
	if g.origPerson, valid = g.validMap(pane); !valid {
		return fmt.Errorf("not a valid map")
	}

	g.original = pane
	g.height = pane.Height()
	g.width = pane.Width()
	if pane == nil {
		return nil
	}

	g.runtime = make(gamePane, g.height)

	for i := range g.runtime {
		g.runtime[i] = make([]gameItem, g.width)
		for j := range pane[i] {
			g.runtime[i][j] = g.runtime[i][j].fromByte(pane[i][j])
		}
	}

	g.runPerson = g.origPerson
	return nil
}

func (g *pushBox) move(x, y int) {
	a, b := g.runPerson.X+x, g.runPerson.Y+y
	// range check
	if a < 0 || b < 0 || a >= g.height || b >= g.width {
		g.msg = "out of range!" + fmt.Sprint(a, b)
		return
	}

	if g.runtime[a][b] == PushBoxWall { // wall
		g.msg = "wall!"
	} else if g.runtime[a][b]&PushBoxBox > 0 {
		// do push
		c, d := a+x, b+y
		if c < 0 || d < 0 || c >= g.height || d >= g.width || // out of range
			g.runtime[c][d]&PushBoxWall > 0 || // wall
			g.runtime[c][d]&PushBoxBox > 0 { // box
			g.msg = "can not push box!"
		} else {
			// do push
			g.runtime[c][d] |= PushBoxBox
			g.runtime[a][b] &= ^PushBoxBox & 0xff
			g.doMove(a, b)
		}
	} else { // 无墙无箱
		g.doMove(a, b)
	}
}

func (g *pushBox) doMove(x, y int) {
	// do move
	g.runtime[g.runPerson.X][g.runPerson.Y] &= ^PushBoxPerson & 0xff
	g.runtime[x][y] |= PushBoxPerson
	g.runPerson = game.Point{X: x, Y: y}
	g.msg = ""
}

func (g *pushBox) draw() {
	if g.width <= 0 || g.height <= 0 {
		return
	}

	// clear screen
	fmt.Print("\x1B[1J")
	// move to left-top
	fmt.Print("\x1B[H")

	// panel
	for _, s := range g.runtime {
		for _, c := range s {
			fmt.Printf("%c", c.toByte())
		}
		fmt.Println()
	}

	// messages at right
	fmt.Printf("\x1b[1;%dH Tips: push all `o` to `.`", g.width+4)
	fmt.Printf("\x1b[2;%dH press w,s,a,d to move `p`", g.width+4)
	fmt.Printf("\x1b[3;%dH press r to reset, q to exit", g.width+4)
	fmt.Printf("\x1b[4;%dH %s", g.width+4, g.msg)
	fmt.Printf("\x1b[%d;1H ", g.height+1)
	fmt.Printf("height:%d, width:%d", g.height, g.width)
}

func runGame(path string) {
	maps, err := loadMap(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
	if len(maps) == 0 {
		maps = defaultMaps
	}

	for _, m := range maps {
		game.RunGame(&pushBox{}, m)
	}
}
