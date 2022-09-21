package main

import (
	"fmt"
	"github.com/zhaowk/game"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	//russiaBlockWall  = "\x1b[40m \x1b[0m"
	russiaBlockWall  = "#"
	russiaBlockBlk   = "@"
	russiaBlockEmpty = " "
)

var (
	blocks = []block{
		&block2{base: []game.Point{{0, 0}, {0, 1}, {1, 0}, {1, 1}}}, // ::
		&block3{base: []game.Point{{0, 2}, {1, 0}, {1, 1}, {1, 2}}}, // ..:
		&block3{base: []game.Point{{0, 0}, {1, 0}, {1, 1}, {1, 2}}}, // :..
		&block3{base: []game.Point{{0, 1}, {0, 2}, {1, 0}, {1, 1}}}, // .:'
		&block3{base: []game.Point{{0, 0}, {0, 1}, {1, 1}, {1, 2}}}, // ':.
		&block3{base: []game.Point{{0, 1}, {1, 0}, {1, 1}, {1, 2}}}, // .:.
		&block4{base: []game.Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}}}, // ....
	}
)

type russiaBlock struct {
	width   int
	height  int
	runtime [][]byte
	msg     string
	score   int
	stop    bool

	opCh chan int
	curr block
	next block
	pos  game.Point
}

func (b *russiaBlock) Init(...interface{}) error {
	b.width = 10
	b.height = 15
	rand.Seed(time.Now().UnixNano())
	b.runtime = make([][]byte, b.height)
	for i := 0; i < b.height; i++ {
		b.runtime[i] = make([]byte, b.width)
	}

	b.genNext()
	b.curr = b.next
	b.pos = game.Point{Y: b.width / 2}
	b.genNext()

	b.opCh = make(chan int)

	go b.run()
	return nil
}

func (b *russiaBlock) Run(k int, _ string) {
	switch k {
	case 'w', 'W', game.SysUp:
		b.opCh <- game.SysUp
	case 's', 'S', game.SysDown:
		b.opCh <- game.SysDown
	case 'a', 'A', game.SysLeft:
		b.opCh <- game.SysLeft
	case 'd', 'D', game.SysRight:
		b.opCh <- game.SysRight
	case 'q', 'Q':
		b.stop = true
	}
}

func (b *russiaBlock) Next() bool {
	return !b.stop
}

func (b *russiaBlock) Finish() {
}

func (b *russiaBlock) run() {
	tick := time.Tick(10 * time.Millisecond)
	prev := time.Now()
	b.draw()
	for !b.stop {
		select {
		case op := <-b.opCh:
			switch op {
			case game.SysUp: // switch
				b.doSwitch()
			case game.SysDown:
				b.doRapidDown()
			case game.SysLeft:
				b.doLeft()
			case game.SysRight:
				b.doRight()
			}
			b.draw()
		case <-tick:
			if time.Now().Add(-1 * time.Second).After(prev) {
				prev = time.Now()
				b.msg = time.Now().Format("2006-01-02 03:04:05")
				b.doDown()
				b.draw()
			}
		}
	}
}

func (b *russiaBlock) genNext() {
	idx, roll := rand.Intn(len(blocks)), rand.Intn(4)
	bl := blocks[idx]
	for i := 0; i < roll; i++ {
		bl = bl.Switch()
	}
	b.next = bl
}

func (b *russiaBlock) draw() {
	game.Clear()
	game.Cursor(game.Point{})

	// panel
	game.Draw(strings.Repeat(russiaBlockWall, b.width+2))

	for i, s := range b.runtime {
		game.DrawAt(game.Point{X: i + 1}, russiaBlockWall)
		for _, c := range s {
			if c == 0 || c == ' ' {
				game.Draw(russiaBlockEmpty)
			} else {
				game.Draw(russiaBlockBlk)
			}
		}
		game.DrawLine(russiaBlockWall)
	}
	game.Draw(strings.Repeat(russiaBlockWall, b.width+2))

	for _, p := range b.curr.Points() {
		game.DrawAt(game.Point{X: b.pos.X + p.X + 1, Y: b.pos.Y + p.Y + 1}, russiaBlockBlk)
	}

	// messages at right
	game.DrawAt(game.Point{X: 1, Y: b.width + 4}, fmt.Sprintf("Score: %d", b.score))
	game.DrawAt(game.Point{X: 2, Y: b.width + 4}, "Next:")
	game.DrawAt(game.Point{X: 3, Y: b.width + 6}, b.next.GetLine(0))
	game.DrawAt(game.Point{X: 4, Y: b.width + 6}, b.next.GetLine(1))
	game.DrawAt(game.Point{X: 5, Y: b.width + 6}, b.next.GetLine(2))
	game.DrawAt(game.Point{X: 6, Y: b.width + 6}, b.next.GetLine(3))
	game.DrawAt(game.Point{X: 7, Y: b.width + 4}, "Tips:")
	game.DrawAt(game.Point{X: 8, Y: b.width + 7}, "q -> exit")
	game.DrawAt(game.Point{X: 9, Y: b.width + 7}, "a -> left")
	game.DrawAt(game.Point{X: 10, Y: b.width + 7}, "d -> right")
	game.DrawAt(game.Point{X: 11, Y: b.width + 7}, "w -> switch")
	game.DrawAt(game.Point{X: 12, Y: b.width + 7}, "s -> down")
	game.DrawAt(game.Point{X: 13, Y: b.width + 4}, sub(b.msg))
	game.Cursor(game.Point{X: b.height + 1, Y: b.width + 3})
}

func (b *russiaBlock) doSwitch() {
	s := b.curr.Switch()
	if b.isValid(b.pos, s) {
		b.curr = s
	}
}

func (b *russiaBlock) doLeft() {
	b.doMove(0, -1)
}

func (b *russiaBlock) doRight() {
	b.doMove(0, 1)
}

func (b *russiaBlock) doDown() bool {
	b.doMove(1, 0)
	return b.doCheck()
}

func (b *russiaBlock) doMove(x, y int) {
	target := b.pos.Add(game.Point{X: x, Y: y})

	if b.isValid(target, b.curr) {
		b.pos = target
	}
}

func (b *russiaBlock) doRapidDown() {
	// just do down
	for i := 0; i < b.height; i++ {
		if b.doDown() {
			return
		}
	}
}

func (b *russiaBlock) doCheck() (merged bool) {
	down := b.pos.Add(game.Point{X: 1})

	if !b.isValid(down, b.curr) {
		merged = true
		lines := make([]int, 0)
		for _, p := range b.curr.Points() {
			// do merge
			b.runtime[b.pos.X+p.X][b.pos.Y+p.Y] = 1
			filled := true
			for _, r := range b.runtime[b.pos.X+p.X] {
				if r == 0 {
					filled = false
					break
				}
			}
			if filled {
				lines = append(lines, b.pos.X+p.X)
			}
		}
		// do score check
		if len(lines) > 0 {
			b.score += 1 << (len(lines) - 1)
			b.movePane(lines)
		}

		// generate new
		b.curr = b.next
		b.pos = game.Point{Y: b.width / 2}
		b.genNext()

		// game over
		if !b.isValid(b.pos, b.curr) {
			b.msg = "Game over!"
			b.draw()
			time.Sleep(time.Second)
			os.Exit(0)
		}
	}
	return
}

func (b *russiaBlock) isValid(pos game.Point, blk block) bool {
	for _, bl := range blk.Points() {
		if p := pos.Add(bl); !b.check(p) {
			return false
		}
	}
	return true
}

func (b *russiaBlock) check(p game.Point) bool {
	return p.X >= 0 && p.X < b.height && p.Y >= 0 && p.Y < b.width && b.runtime[p.X][p.Y] == 0
}

func (b *russiaBlock) movePane(lines []int) {
	for _, line := range lines {
		for i := line; i > 0; i-- {
			copy(b.runtime[i], b.runtime[i-1]) // copy prev line -> next line
		}

		for i := 0; i < len(b.runtime[0]); i++ { // the first line, fill 0
			b.runtime[0][i] = 0
		}
	}
}

func sub(s string) string {
	if len(s) < 20 {
		return s
	}
	return s[:20]
}
