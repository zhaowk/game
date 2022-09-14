package main

import (
	"container/list"
	"fmt"
	"github.com/zhaowk/game"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	snakeBlank = " "
	snakeBody  = "#"
	snakeHead  = "O"
	snakeWall  = "@"
	snakeFood  = "o"
)

type snake struct {
	width  int
	height int
	msg    string

	opCh      chan int
	snake     *list.List
	food      game.Point
	direction int
}

func (s *snake) Init(...interface{}) error {
	s.width = 10
	s.height = 10
	rand.Seed(time.Now().UnixNano())

	s.snake = list.New()
	pos := game.Point{X: s.height / 2, Y: s.width / 2}
	s.snake.PushFront(pos)
	s.direction = game.SysLeft
	s.genFood()

	s.opCh = make(chan int)

	go s.run()
	return nil
}

func (s *snake) Run(k int, _ string) {
	switch k {
	case 'w', 'W', game.SysUp:
		s.opCh <- game.SysUp
	case 's', 'S', game.SysDown:
		s.opCh <- game.SysDown
	case 'a', 'A', game.SysLeft:
		s.opCh <- game.SysLeft
	case 'd', 'D', game.SysRight:
		s.opCh <- game.SysRight
	case 'q', 'Q':
		os.Exit(0)
	}
}

func (s *snake) Next() bool {
	return true
}

func (s *snake) Finish() {
}

func (s *snake) run() {
	tick := time.Tick(10 * time.Millisecond)
	prev := time.Now()
	s.draw()
	for {
		select {
		case op := <-s.opCh:
			switch op {
			case game.SysUp: // switch
				s.direction = game.SysUp
			case game.SysDown:
				s.direction = game.SysDown
			case game.SysLeft:
				s.direction = game.SysLeft
			case game.SysRight:
				s.direction = game.SysRight
			}
		case <-tick:
			if time.Now().Add(-1 * time.Second).After(prev) {
				prev = time.Now()
				s.msg = time.Now().Format("2006-01-02 03:04:05")
				s.doMove()
				s.draw()
			}
		}
	}
}

func (s *snake) draw() {
	// clear screen
	fmt.Print("\x1B[1J")
	// move to left-top
	fmt.Print("\x1B[H")

	// panel
	fmt.Printf("\x1b[1H%s", strings.Repeat(snakeWall, s.width+2))
	for i := 0; i < s.height; i++ {
		fmt.Printf("\x1b[%dH%s", i+2, snakeWall)
		for j := 0; j < s.width; j++ {
			fmt.Print(snakeBlank)
		}
		fmt.Println(snakeWall)
	}
	fmt.Print(strings.Repeat(snakeWall, s.width+2))

	for e := s.snake.Front(); e != nil; e = e.Next() {
		if p, ok := e.Value.(game.Point); ok {
			fmt.Printf("\x1b[%d;%dH%s", p.X+2, p.Y+2, snakeBody)
		}
	}

	// snake head
	if head, ok := s.snake.Front().Value.(game.Point); ok {
		fmt.Printf("\x1b[%d;%dH%s", head.X+2, head.Y+2, snakeHead)
	}

	// snake food
	fmt.Printf("\x1b[%d;%dH%s", s.food.X+2, s.food.Y+2, snakeFood)

	// messages at right
	fmt.Printf("\x1b[2;%dH Score: %d", s.width+4, s.snake.Len())

	fmt.Printf("\x1b[3;%dH Tips:", s.width+4)
	fmt.Printf("\x1b[4;%dH    q -> exit", s.width+4)
	fmt.Printf("\x1b[5;%dH    a -> left", s.width+4)
	fmt.Printf("\x1b[6;%dH    d -> right", s.width+4)
	fmt.Printf("\x1b[7;%dH    w -> up", s.width+4)
	fmt.Printf("\x1b[8;%dH    s -> down", s.width+4)
	fmt.Printf("\x1b[9;%dH %s", s.width+4, sub(s.msg))
	fmt.Printf("\x1b[%d;%dH ", s.height+2, s.width+4)
}

func (s *snake) doMove() {
	if s.snake.Len() < 1 {
		return
	}

	if p, ok := s.snake.Front().Value.(game.Point); !ok {
		return
	} else {
		var q game.Point
		switch s.direction {
		case game.SysUp:
			q = p.Add(game.Point{X: -1})
		case game.SysDown:
			q = p.Add(game.Point{X: 1})
		case game.SysLeft:
			q = p.Add(game.Point{Y: -1})
		case game.SysRight:
			q = p.Add(game.Point{Y: 1})
		}

		if !s.check(q) { // game over
			s.msg = "Game over!"
			s.draw()
			time.Sleep(time.Second)
			os.Exit(0)
		}

		s.snake.PushFront(q)

		if s.food == q {
			// check
			s.doCheck()
			// eat, gen new
			s.genFood()
		} else {
			// not eat, remove tail
			s.snake.Remove(s.snake.Back())
		}
	}
}

func (s *snake) doCheck() {
	if s.snake.Len() == s.height*s.width {
		s.msg = "Win!"
		s.draw()
		time.Sleep(time.Second)
		os.Exit(0)
	}
}

func (s *snake) genFood() {
	n := rand.Intn(s.height*s.width - s.snake.Len())
	c := 0
	for i := 0; i < n && c < s.height*s.width; c++ {
		if s.isValid(c) {
			i++
		}
	}

	s.food = game.Point{X: c / s.width, Y: c % s.width}
}

func (s *snake) isValid(n int) bool {
	q := game.Point{X: n / s.width, Y: n % s.width}

	for e := s.snake.Front(); e != nil; e = e.Next() {
		if p, ok := e.Value.(game.Point); !ok || p == q {
			return false
		}
	}
	return true
}

func (s *snake) check(p game.Point) bool {
	if p.X < 0 || p.X >= s.height || p.Y < 0 || p.Y >= s.width { // out of range
		return false
	}

	for e := s.snake.Front(); e != nil; e = e.Next() { //
		if q, ok := e.Value.(game.Point); !ok || p == q {
			return false
		}
	}
	return true
}

func sub(s string) string {
	if len(s) < 20 {
		return s
	}
	return s[:20]
}
