package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var (
	defaultMaps = []boxMap{
		{
			"########",
			"# .. p #",
			"# oo   #",
			"#      #",
			"########",
		},
		{
			"#############",
			"# ...       #",
			"# ooo  p    #",
			"#           #",
			"#############",
		},
		{
			"###########",
			"# ..   #. #",
			"# oo   #o #",
			"#    p    #",
			"###########",
		},
	}
)

type boxMap []string

// String map stringer
func (b boxMap) String() (r string) {
	for _, s := range b {
		r += s + "\n"
	}
	return
}

// Width map width(cols)
func (b boxMap) Width() int {
	width := 0
	for _, s := range b {
		if len(s) > width {
			width = len(s)
		}
	}
	return width
}

// Height map height(rows)
func (b boxMap) Height() int {
	return len(b)
}

// loadMap : load maps from `path`
// A map is a text file with suffix `.txt`
func loadMap(path string) (maps []boxMap, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	maps = make([]boxMap, 0)
	var m []string

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			m, err = readMap(path + "/" + entry.Name())
			if err != nil {
				return
			}
			maps = append(maps, m)
		}
	}

	return
}

// readMap: read map from `file`
// A valid text has only such symbols:
//
//	`.`: a target
//	`o`: a box
//	`O`: a box on a target
//	`p`: the player
//	`P`: then player on the target
//	`#`: the wall
//	` `: blank
//	`\n`: new line
func readMap(file string) (s boxMap, e error) {
	bs, e := os.ReadFile(file)
	if e != nil {
		return
	}

	bytes.Replace(bs, []byte{'\r', '\n'}, []byte{'\n'}, -1) // windows
	bytes.Replace(bs, []byte{'\r'}, []byte{'\n'}, -1)       // mac

	for i, b := range bs {
		if !isValid(b) {
			return nil, fmt.Errorf("file format error, %d", i)
		}
	}

	s = strings.Split(string(bs), "\n")
	return
}

func isValid(b byte) bool {
	switch b {
	case '.', 'o', 'O', 'p', 'P', '#', ' ', '\n':
		return true
	default:
		return false
	}
}
