package main

import (
	"github.com/zhaowk/game"
	"sort"
	"strings"
)

type block interface {
	Switch() block
	Points() []game.Point
	GetLine(n int) string
}

type block2 struct {
	base []game.Point
}

func (b block2) Points() []game.Point {
	return b.base
}

func (b block2) Switch() block {
	return &block2{base: b.base}
}

func (b block2) GetLine(n int) string {
	if n < 2 {
		return strings.Repeat(russiaBlockBlk, 2)
	}
	return ""
}

type block3 struct {
	base []game.Point
}

func (b block3) Points() []game.Point {
	return b.base
}

func (b block3) Switch() block {
	points := make([]game.Point, len(b.base))
	for i, p := range b.base {
		switch p.X {
		case 0:
			points[i] = game.Point{X: p.Y, Y: 2}
		case 1:
			points[i] = game.Point{X: p.Y, Y: p.X}
		case 2:
			points[i] = game.Point{X: p.Y}
		}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Less(points[j])
	})
	return &block3{base: points}
}

func (b block3) GetLine(n int) string {
	s := make([]byte, 3)
	s[0] = russiaBlockEmpty[0]
	s[1] = russiaBlockEmpty[0]
	s[2] = russiaBlockEmpty[0]
	if n < 3 {
		for _, p := range b.base {
			if p.X == n {
				s[p.Y] = russiaBlockBlk[0]
			}
		}
	}
	return string(s)
}

type block4 struct {
	base []game.Point
}

func (b block4) Points() []game.Point {
	return b.base
}

func (b block4) Switch() block {
	var points []game.Point
	if b.base[1].X == 0 { // 横 -> 竖
		points = []game.Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
	} else { // 竖 -> 横
		points = []game.Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
	}
	return &block4{base: points}
}

func (b block4) GetLine(n int) string {
	s := make([]byte, 4)
	s[0] = russiaBlockEmpty[0]
	s[1] = russiaBlockEmpty[0]
	s[2] = russiaBlockEmpty[0]
	s[3] = russiaBlockEmpty[0]
	if n < 4 {
		for _, p := range b.base {
			if p.X == n {
				s[p.Y] = russiaBlockBlk[0]
			}
		}
	}
	return string(s)
}
