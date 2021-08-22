package main

import (
	"encoding/json"
	"math"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Config struct {
	LastOpenedDocument string `json:"lastOpenDocument"`
}

func loadConfig() Config {
	content, err := os.ReadFile(path.Join(saveLocation, "config.json"))

	if err != nil {
		id, _ := gonanoid.New()
		return Config{
			LastOpenedDocument: id + ".txt",
		}
	}

	var config Config

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(content, &config)
	return config
}

func saveConfig(config Config) {

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	content, _ := json.Marshal(config)
	os.WriteFile(path.Join(saveLocation, "config.json"), []byte(content), 777)
}

func isLetter(s string) bool {
	return !strings.Contains(s, "KEY")
}

func insertAt(text string, insert string, index int) string {
	if index == lenString(text) {
		return text + insert
	}
	runeText := []rune(text)
	return string(append(runeText[:index], append([]rune(insert), runeText[index:]...)...))
}

func deleteAt(text string, index int) string {
	runeText := []rune(text)
	return string(append(runeText[:index-1], runeText[index:]...))
}

func lenString(s string) int {
	return utf8.RuneCountInString(s)
}

func printAlert(message string, width int, s *Screen) {
	matrix := createMatrixFromText(message, width)
	x := math.Floor((float64(s.state.MaxCols)/2)-float64(width)/2) - 1
	y := math.Floor((float64(s.state.MaxRows)/2)-float64(len(matrix))/2) - 1
	outerMatrix := createNewMatrix(width+2, len(matrix)+2)
	matrix = pasteMatrix(outerMatrix, matrix, 1, 1)
	matrix = inverseMatrix(matrix)
	s.print(pasteMatrix(s.originalMatrix, matrix, int(x), int(y)))
}
