package main

import (
	"github.com/shermp/go-fbink-v2/gofbink"
)

type Screen struct {
	originalMatrix Matrix
	presentMatrix  Matrix
	fb             *gofbink.FBInk
	state          gofbink.FBInkState
}

func (s *Screen) init() {
	s.state = gofbink.FBInkState{}

	fbinkOpts := gofbink.FBInkConfig{}
	rOpts := gofbink.RestrictedConfig{
		Fontmult: 3,
		Fontname: gofbink.Ctrld,
	}
	s.fb = gofbink.New(&fbinkOpts, &rOpts)

	s.fb.Open()
	s.fb.Init(&fbinkOpts)
	s.fb.AddOTfont("/mnt/onboard/.adds/kobowriter/inc.otf", gofbink.FntRegular)

	s.fb.GetState(&fbinkOpts, &s.state)

	// clear screen on initialisation
	s.fb.ClearScreen(&gofbink.FBInkConfig{
		IsFlashing: true,
	})

	s.presentMatrix = createNewMatrix(int(s.state.MaxCols), int(s.state.MaxRows))
	s.originalMatrix = createNewMatrix(int(s.state.MaxCols), int(s.state.MaxRows))

	println("Screen struct inited")

}

func (s *Screen) clean() {
	s.fb.Close()
}

func (s *Screen) print(matrix Matrix) {
	printDiff(s.presentMatrix, matrix, s.fb)
	s.presentMatrix = matrix
}

func same(a MatrixElement, b MatrixElement) bool {
	return a.content == b.content && a.isInverted == b.isInverted
}

func printDiff(previous Matrix, next Matrix, fb *gofbink.FBInk) {
	for i := range previous {
		for j := range previous[i] {
			if !same(previous[i][j], next[i][j]) {
				fb.FBprint(string(next[i][j].content), &gofbink.FBInkConfig{
					Row:        int16(i),
					Col:        int16(j),
					NoRefresh:  true,
					IsInverted: next[i][j].isInverted,
				})
			}

		}
	}

	fb.Refresh(0, 0, 0, 0, gofbink.HWDither(gofbink.WfmAUTO), &gofbink.FBInkConfig{})
}
