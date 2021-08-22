package main

import (
	"strings"
)

type MatrixElement struct {
	content    rune
	isInverted bool
}

type Matrix [][]MatrixElement

func createNewMatrix(width int, height int) Matrix {
	a := make([][]MatrixElement, height)
	for i := range a {
		a[i] = make([]MatrixElement, width)
		for j := range a[i] {
			a[i][j] = MatrixElement{
				content:    ' ',
				isInverted: false,
			}
		}
	}
	return a
}

func createMatrixFromText(text string, width int) Matrix {

	wrapped := wrapText(text, int(width))

	wrapedArray := strings.Split(wrapped, "\n")
	result := createNewMatrix(width, len(wrapedArray))

	for i := range result {
		for j := range result[i] {
			if j < lenString(wrapedArray[i]) {
				result[i][j].content = []rune(wrapedArray[i])[j]
			}
		}
	}

	return result
}

func pasteMatrix(baseMatrix Matrix, topMatrix Matrix, offsetX int, offsetY int) Matrix {
	resultMatrix := copyMatrix(baseMatrix)
	for i := range resultMatrix {
		localI := i - offsetY
		if localI < 0 || localI >= len(topMatrix) {
			continue
		}

		for j := range resultMatrix[i] {
			localJ := j - offsetX
			if localJ < 0 || localJ >= len(topMatrix[localI]) {
				continue
			}

			resultMatrix[i][j] = topMatrix[localI][localJ]

		}
	}

	return resultMatrix
}

func matrixToText(matrix Matrix) string {
	stringz := make([]string, len(matrix))
	for i := range matrix {
		for _, elem := range matrix[i] {
			stringz[i] = stringz[i] + string(elem.content)
		}
	}
	return strings.Join(stringz, "")
}

func inverseMatrix(in Matrix) (out Matrix) {
	out = in
	for i := range out {
		for j, elem := range out[i] {
			out[i][j].isInverted = !elem.isInverted
		}
	}
	return
}

func copyMatrix(in Matrix) (out Matrix) {
	out = createNewMatrix(len(in[0]), len(in))
	for i := range out {
		for j := range out[i] {
			out[i][j].content = in[i][j].content
			out[i][j].isInverted = in[i][j].isInverted
		}
	}
	return
}
