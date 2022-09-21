//go:build linux || darwin

package game

import (
	"fmt"
	"strconv"
)

// Clear screen
func Clear() {
	fmt.Print(ClearPrev)
	fmt.Print(ClearAll)
}

// DrawAt draw string `s` at Point `p`. Point(x, y) starts with (0, 0) from left-top (x => row, y => col)
func DrawAt(p Point, s string) {
	fmt.Print(cursorPos(p.X+1, p.Y+1) + s)
}

// Draw string `s` at current Position
func Draw(s string) {
	fmt.Print(s)
}

// DrawLine Draw string `s` at current Position with new line
func DrawLine(s string) {
	fmt.Println(s)
}

// DrawLineAt Draw string `s` at Point `p` with new line
func DrawLineAt(p Point, s string) {
	fmt.Println(cursorPos(p.X+1, p.Y+1) + s)
}

// DrawSgr Draw content `s` with Sgr definitions
func DrawSgr(content string, sgr ...string) {
	fmt.Printf("%s%s%s", SgrSet(sgr...), content, SgrReset())
}

// DrawColor8 draw `s` with color(ColorType, Color8)
func DrawColor8(t ColorType, c Color8, s string) {
	DrawSgr(s, SgrColor(t, c))
}

// DrawColor256 draw `s` with color256(ColorType, uint8)
func DrawColor256(t ColorType, c uint8, s string) {
	DrawSgr(s, SgrColor8bit(t, c))
}

// DrawColorRGB draw `s` with colorRgb(ColorType, RGB)
func DrawColorRGB(t ColorType, c RGB, s string) {
	DrawSgr(s, SgrColorRGB(t, c))
}

// csi 转义序列定义， 参考：
// https://en.wikipedia.org/wiki/ANSI_escape_code
// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97

const (
	_CSI          = "\x1b["
	ClearAfter    = _CSI + "0J"
	ClearPrev     = _CSI + "1J"
	ClearAll      = _CSI + "2J"
	ClearLineNext = _CSI + "0K"
	ClearLinePrev = _CSI + "1K"
	ClearLineFull = _CSI + "2K"
	SCP           = _CSI + "s" // 保存光标位置
	RCP           = _CSI + "u" // 恢复光标位置
)

// Cursor move cursor to Point(p) starts with (0,0) from left-top
func Cursor(p Point) {
	fmt.Print(cursorPos(p.X+1, p.Y+1))
}

// CursorUp move cursor up
func CursorUp(n int) {
	fmt.Printf("%s%dA", _CSI, n)
}

// CursorDown move cursor down
func CursorDown(n int) {
	fmt.Printf("%s%dB", _CSI, n)
}

// CursorForward move cursor forward
func CursorForward(n int) {
	fmt.Printf("%s%dC", _CSI, n)
}

// CursorBack move cursor back
func CursorBack(n int) {
	fmt.Printf("%s%dD", _CSI, n)
}

// CursorNextLine move cursor next line
func CursorNextLine(n int) {
	fmt.Printf("%s%dE", _CSI, n)
}

// CursorPrevLine move cursor prev line
func CursorPrevLine(n int) {
	fmt.Printf("%s%dF", _CSI, n)
}

// CursorPos move cursor to (x, y) starts with (1, 1)
func CursorPos(x, y int) {
	fmt.Print(cursorPos(x, y))
}

// CursorSave save current cursor position
func CursorSave() {
	fmt.Print(SCP)
}

// CursorRestore restore saved position
func CursorRestore() {
	fmt.Print(RCP)
}

func cursorPos(x, y int) string {
	return fmt.Sprintf("%s%d;%dH", _CSI, x, y)
}

func scrollUp(n int) string {
	return fmt.Sprintf("%s%dS", _CSI, n)
}

func scrollDown(n int) string {
	return fmt.Sprintf("%s%dT", _CSI, n)
}

// Sgr get sgr string
func Sgr(m string) string {
	return _CSI + m + "m"
}

// SgrSet set multi sgr attributes and return the string
func SgrSet(s ...string) string {
	if len(s) == 0 {
		return ""
	}

	sgr := s[0]
	for i := 1; i < len(s); i++ {
		sgr += ";" + s[i]
	}
	return Sgr(sgr)
}

// SgrReset reset sgr
func SgrReset() string {
	return SgrNormal
}

// ColorType color type of Foreground, Background, BrightForeground, BrightBackground
type ColorType uint8

// String color type string
func (c ColorType) String() string {
	switch c {
	case Foreground:
		return SgrForeground
	case Background:
		return SgrBackground
	case BrightForeground:
		return SgrBrightForeground
	case BrightBackground:
		return SgrBrightBackground
	}
	return ""
}

// Color8 color of ColorBlack, ColorRed, ..., ColorWhite
type Color8 uint8

// String Color8 symbol
func (c Color8) String() string {
	if c < ColorMax {
		return strconv.Itoa(int(c))
	}
	return ""
}

// RGB color of RGB(R, G, B)
type RGB [3]uint8

// String RGB symbol
func (r RGB) String() string {
	return fmt.Sprintf("%d;%d;%d", r[0], r[1], r[2])
}

const (
	Foreground ColorType = iota
	Background
	BrightForeground
	BrightBackground
)

const (
	ColorBlack Color8 = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorMax
)

// SgrColor sgr color 8 string
func SgrColor(t ColorType, color Color8) string {
	if t > BrightBackground || color >= ColorMax {
		return ""
	}
	return fmt.Sprintf("%s%d", t.String(), color)
}

// SgrColor8bit sgr color 8bit string
func SgrColor8bit(t ColorType, color uint8) string {
	if t > BrightBackground {
		return ""
	}
	return fmt.Sprintf("%s%s%d", t.String(), SgrColor256, color)
}

// SgrColorRGB sgr color rgb string
func SgrColorRGB(t ColorType, color RGB) string {
	if t > BrightBackground {
		return ""
	}
	return fmt.Sprintf("%s%s%s", t.String(), SgrColorRgb, color.String())
}

const (
	SgrNormal = _CSI + "m"

	SgrBold       = "1"
	SgrFaint      = "2"
	SgrItalic     = "3"
	SgrUnderLine  = "4"
	SgrSlowBlink  = "5"
	SgrRapidBlink = "6"
	SgrReverse    = "7"
	SgrHide       = "8"
	SgrStrike     = "9"

	SgrPriFont = "10"
	SgrFont1   = "11"
	SgrFont2   = "12"
	SgrFont3   = "13"
	SgrFont4   = "14"
	SgrFont5   = "15"
	SgrFont6   = "16"
	SgrFont7   = "17"
	SgrFont8   = "18"
	SgrFont9   = "19"

	SgrNotBold       = "21"
	SgrNotFaint      = "22"
	SgrNotItalic     = "23"
	SgrNotUnderLine  = "24"
	SgrNotSlowBlink  = "25"
	SgrNotRapidBlink = "26"
	SgrNotReverse    = "27"
	SgrNotHide       = "28"
	SgrNotStrike     = "29"

	SgrOverLine    = "53"
	SgrNotOverLine = "55"

	SgrIdeogramUnderLine       = "60"
	SgrIdeogramDoubleUnderLine = "61"
	SgrIdeogramOverLine        = "62"
	SgrIdeogramDoubleOverLine  = "63"
	SgrIdeogramStressMark      = "64"
	SgrIdeogramClear           = "65"

	SgrForeground       = "3"
	SgrBackground       = "4"
	SgrBrightForeground = "9"
	SgrBrightBackground = "10"

	SgrColorBlack   = "0"
	SgrColorRed     = "1"
	SgrColorGreen   = "2"
	SgrColorYellow  = "3"
	SgrColorBlue    = "4"
	SgrColorMagenta = "5"
	SgrColorCyan    = "6"
	SgrColorWhite   = "7"
	SgrColor256     = "8;5;"
	SgrColorRgb     = "8;2;"
	SgrColorDft     = "9"
)
