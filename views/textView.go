package views

import (
	"strings"
	"unicode/utf8"

	"github.com/olup/kobowriter/matrix"
	"github.com/olup/kobowriter/utils"
)

type TextView struct {
	content     string
	width       int
	height      int
	wrapContent []string
	cursorIndex int
	cursorPos   Position
	lineCount   []int
	scroll      int
}

type Position struct {
	x int
	y int
}

func (t *TextView) init(width int) {
	t.width = width
}

func (t *TextView) setContent(text string) {
	t.content = text
	t.wrapContent = strings.Split(utils.WrapText(text, t.width), "\n")

	lineCount := []int{}
	for _, line := range t.wrapContent {
		lineCount = append(lineCount, utf8.RuneCountInString(line)+1)
	}
	t.lineCount = lineCount
}

func (t *TextView) setCursorIndex(index int) {

	// Bounds
	if index < 0 {
		index = 0
	}
	if index > utils.LenString(t.content) {
		index = utils.LenString(t.content)
	}

	// Processing
	t.cursorIndex = index
	x := 0
	y := 0

	agg := 0

	for i, count := range t.lineCount {
		aggNext := count + agg
		if aggNext > t.cursorIndex {
			y = i
			x = t.cursorIndex - agg
			break
		}
		agg = aggNext
	}

	t.cursorPos = Position{
		x,
		y,
	}

	t.updateScroll()

}

func (t *TextView) setCursorPos(position Position) {
	// Bounds
	if position.y < 0 {
		position.y = 0
	}

	if position.x < 0 {
		position.x = 0
	}

	if position.y > len(t.lineCount)-1 {
		position.y = len(t.lineCount) - 1
	}

	if t.lineCount[position.y]-1 < position.x {
		position.x = t.lineCount[position.y] - 1
	}

	// Procesing

	agg := 0

	for i := 0; i < position.y; i++ {
		agg += t.lineCount[i]
	}

	agg += position.x

	t.cursorPos = position
	t.cursorIndex = agg
	t.updateScroll()

}

func (t *TextView) renderMatrix() matrix.Matrix {
	textMatrix := matrix.CreateMatrixFromText(t.content, t.width)
	if t.cursorPos.x >= 0 && t.cursorPos.y >= 0 && t.cursorPos.x < t.width {
		textMatrix[t.cursorPos.y][t.cursorPos.x].IsInverted = true
	}
	endBound := t.scroll + t.height
	if endBound > len(textMatrix) {
		endBound = len(textMatrix)
	}
	scrolledTextMatrix := textMatrix[t.scroll:endBound]
	return scrolledTextMatrix
}

func (t *TextView) updateScroll() {
	y := t.cursorPos.y

	if y > t.scroll+t.height-1 {
		t.scroll = y - 5
	}
	if y < t.scroll {
		t.scroll = y - t.height + 5
	}
	if t.scroll > len(t.wrapContent) {
		t.scroll = len(t.wrapContent) - 5
	}
	if t.scroll < 0 {
		t.scroll = 0
	}
}
