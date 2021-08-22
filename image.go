package main

import (
	"image"
	"image/png"
	"io"
)

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) ([]byte, error) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels []byte

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels = append(pixels, byte(r/257), byte(g/257), byte(b/257), byte(a/257))
		}
	}

	return pixels, nil
}

// Get the bi-dimensional pixel array
func getPixelsFromImage(img image.Image) ([]byte, error) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels []byte

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels = append(pixels, byte(r/257), byte(g/257), byte(b/257), byte(a/257))
		}
	}

	return pixels, nil
}
