package screener

import (
	"bytes"
	"image"
	"math"

	"github.com/fogleman/gg"
	"github.com/olup/kobowriter/matrix"
	"github.com/shermp/go-fbink-v2/gofbink"
)

type Screen struct {
	originalMatrix matrix.Matrix
	presentMatrix  matrix.Matrix
	fb             *gofbink.FBInk
	state          gofbink.FBInkState
	Width          int
	Height         int
}

var dc = gg.NewContext(25, 40)
var charCache = map[string][]byte{}

func InitScreen() (s *Screen) {
	s = &Screen{}
	dc.LoadFontFace("inc.ttf", 48)
	s.state = gofbink.FBInkState{}

	fbinkOpts := gofbink.FBInkConfig{}
	rOpts := gofbink.RestrictedConfig{
		Fontmult: 3,
		Fontname: gofbink.Ctrld,
	}
	s.fb = gofbink.New(&fbinkOpts, &rOpts)

	s.fb.Open()
	s.fb.Init(&fbinkOpts)
	s.fb.AddOTfont("/mnt/onboard/.adds/kobowriter/inc.ttf", gofbink.FntRegular)

	s.fb.GetState(&fbinkOpts, &s.state)

	// clear screen on initialisation
	s.fb.ClearScreen(&gofbink.FBInkConfig{
		IsFlashing: true,
	})

	s.Width = int(s.state.MaxCols)
	s.Height = int(s.state.MaxRows)

	s.presentMatrix = matrix.CreateNewMatrix(s.Width, s.Height)
	s.originalMatrix = matrix.CreateNewMatrix(s.Width, s.Height)

	println("Screen struct inited")

	return

}

func (s *Screen) Clean() {
	s.fb.Close()
}

func (s *Screen) Print(matrix matrix.Matrix) {
	printDiff(s.presentMatrix, matrix, s.fb)
	s.presentMatrix = matrix
}

func same(a matrix.MatrixElement, b matrix.MatrixElement) bool {
	return a.Content == b.Content && a.IsInverted == b.IsInverted
}

func printDiff(previous matrix.Matrix, next matrix.Matrix, fb *gofbink.FBInk) {
	for i := range previous {
		for j := range previous[i] {
			if !same(previous[i][j], next[i][j]) {
				fb.FBprint(string(next[i][j].Content), &gofbink.FBInkConfig{
					Row:        int16(i),
					Col:        int16(j),
					NoRefresh:  true,
					IsInverted: next[i][j].IsInverted,
				})

			}

		}
	}

	fb.Refresh(0, 0, 0, 0, gofbink.HWDither(gofbink.WfmAUTO), &gofbink.FBInkConfig{})
}

func (s *Screen) PrintPng(imgBytes []byte, w int, h int, x int, y int) {
	img, _, _ := image.Decode(bytes.NewReader(imgBytes))
	buffer, _ := getPixelsFromImage(img)
	s.fb.PrintRawData(buffer, w, h, uint16(x), uint16(y), &gofbink.FBInkConfig{})
}

func getCharImage(s string) []byte {
	if char, ok := charCache[s]; ok {
		return char
	} else {
		dc.SetRGB(1, 1, 1)
		dc.Clear()

		dc.SetRGB(0, 0, 0)
		dc.DrawString(s, 0, 35)
		img := dc.Image()
		buffer, _ := getPixelsFromImage(img)
		charCache[s] = buffer
		return buffer
	}
}

func (s *Screen) PrintAlert(message string, width int) {
	matrixx := matrix.CreateMatrixFromText(message, width)
	x := math.Floor((float64(s.state.MaxCols)/2)-float64(width)/2) - 1
	y := math.Floor((float64(s.state.MaxRows)/2)-float64(len(matrixx))/2) - 1
	outerMatrix := matrix.CreateNewMatrix(width+2, len(matrixx)+2)
	matrixx = matrix.PasteMatrix(outerMatrix, matrixx, 1, 1)
	matrixx = matrix.InverseMatrix(matrixx)
	s.Print(matrix.PasteMatrix(s.originalMatrix, matrixx, int(x), int(y)))
}

func (s *Screen) Clear() {
	s.fb.ClearScreen(&gofbink.FBInkConfig{})
	s.presentMatrix = matrix.FillMatrix(s.presentMatrix, ' ')
}

func (s *Screen) GetOriginalMatrix() matrix.Matrix {
	return matrix.CopyMatrix(s.originalMatrix)
}
