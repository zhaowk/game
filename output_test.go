package game

import (
	"testing"
)

func TestCursor(t *testing.T) {
	Clear()
	Cursor(Point{})
	DrawAt(Point{0, 3}, "111")
	Draw("123")
	DrawLine("456")
	DrawColor8(Foreground, ColorRed, "123")
	DrawColor256(Foreground, 0x1f, "123")
	DrawColorRGB(Foreground, RGB{0xcc, 0xcc, 0xcc}, "123")
	DrawSgr("hello world")

	DrawLine("1")
	Draw("1234567890")
	CursorUp(1)
	Draw(ClearLinePrev)
	CursorDown(1)
	CursorPos(1, 2)
	Draw(ClearLineNext)
	DrawSgr("hello world", SgrBold, SgrFaint, SgrItalic, SgrUnderLine, SgrReverse)
}
