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
	// clear && move to (0, 0)
	game.Clear()
	game.Cursor(game.Point{})

	// panel
	game.Draw(strings.Repeat(snakeWall, s.width+2))
	for i := 0; i < s.height; i++ {
		game.DrawAt(game.Point{X: i + 1}, snakeWall)
		for j := 0; j < s.width; j++ {
			game.Draw(snakeBlank)
		}
		game.DrawLine(snakeWall)
	}
	game.Draw(strings.Repeat(snakeWall, s.width+2))

	// snake
	for e := s.snake.Front(); e != nil; e = e.Next() {
		if p, ok := e.Value.(game.Point); ok {
			game.DrawAt(game.Point{X: p.X + 1, Y: p.Y + 1}, snakeBody)
		}
	}

	// snake head
	if head, ok := s.snake.Front().Value.(game.Point); ok {
		game.DrawAt(game.Point{X: head.X + 1, Y: head.Y + 1}, snakeHead)
	}

	// snake food
	game.DrawAt(game.Point{X: s.food.X + 1, Y: s.food.Y + 1}, snakeFood)

	// messages at right
	game.DrawAt(game.Point{X: 1, Y: s.width + 4}, fmt.Sprintf("Score: %d", s.snake.Len()))
	game.DrawAt(game.Point{X: 2, Y: s.width + 4}, "Tips:")
	game.DrawAt(game.Point{X: 3, Y: s.width + 7}, "q -> exit")
	game.DrawAt(game.Point{X: 4, Y: s.width + 7}, "a -> left")
	game.DrawAt(game.Point{X: 5, Y: s.width + 7}, "d -> right")
	game.DrawAt(game.Point{X: 6, Y: s.width + 7}, "w -> up")
	game.DrawAt(game.Point{X: 7, Y: s.width + 7}, "s -> down")
	game.DrawAt(game.Point{X: 8, Y: s.width + 4}, sub(s.msg))
	game.Cursor(game.Point{X: s.height + 1, Y: s.width + 3})
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
			// eat food, gen new
			s.genFood()
		} else {
			// eat nothing, remove tail
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
