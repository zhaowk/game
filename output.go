//go:build linux || darwin

package game

import (
	"fmt"
	"strconv"
)

func Clear() {
	fmt.Print(ClearAll)
}

func DrawAt(p Point, s string) {
	fmt.Print(cursorPos(p.X+1, p.Y+1) + s)
}

func Draw(s string) {
	fmt.Print(s)
}

func DrawLine(s string) {
	fmt.Println(s)
}

func DrawLineAt(p Point, s string) {
	fmt.Println(cursorPos(p.X+1, p.Y+1) + s)
}

func DrawSgr(content string, sgr ...string) {
	fmt.Printf("%s%s%s", SgrSet(sgr...), content, SgrReset())
}

func DrawColor8(t ColorType, c Color8, s string) {
	DrawSgr(s, SgrColor(t, c))
}

func DrawColor256(t ColorType, c uint8, s string) {
	DrawSgr(s, SgrColor8bit(t, c))
}

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

func Cursor(p Point) {
	fmt.Print(cursorPos(p.X+1, p.Y+1))
}

func CursorUp(n int) {
	fmt.Printf("%s%dA", _CSI, n)
}

func CursorDown(n int) {
	fmt.Printf("%s%dB", _CSI, n)
}

func CursorForward(n int) {
	fmt.Printf("%s%dC", _CSI, n)
}

func CursorBack(n int) {
	fmt.Printf("%s%dD", _CSI, n)
}

func CursorNextLine(n int) {
	fmt.Printf("%s%dE", _CSI, n)
}

func CursorPrevLine(n int) {
	fmt.Printf("%s%dF", _CSI, n)
}

func CursorPos(x, y int) {
	fmt.Print(cursorPos(x, y))
}

func CursorSave() {
	fmt.Print(SCP)
}

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

func Sgr(m string) string {
	return _CSI + m + "m"
}

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

func SgrReset() string {
	return SgrNormal
}

type ColorType uint8

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

type Color8 uint8

func (c Color8) String() string {
	if c < ColorMax {
		return strconv.Itoa(int(c))
	}
	return ""
}

type RGB [3]uint8

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

func SgrColor(t ColorType, color Color8) string {
	if t > BrightBackground || color >= ColorMax {
		return ""
	}
	return fmt.Sprintf("%s%d", t.String(), color)
}

func SgrColor8bit(t ColorType, color uint8) string {
	if t > BrightBackground {
		return ""
	}
	return fmt.Sprintf("%s%s%d", t.String(), SgrColor256, color)
}

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
