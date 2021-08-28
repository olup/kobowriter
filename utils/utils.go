package utils

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Config struct {
	LastOpenedDocument string `json:"lastOpenDocument"`
}

func LoadConfig(saveLocation string) Config {
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

func SaveConfig(config Config, saveLocation string) {

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	content, _ := json.Marshal(config)
	os.WriteFile(path.Join(saveLocation, "config.json"), []byte(content), 777)
}

func IsLetter(s string) bool {
	return !strings.Contains(s, "KEY")
}

func InsertAt(text string, insert string, index int) string {
	if index == LenString(text) {
		return text + insert
	}
	runeText := []rune(text)
	return string(append(runeText[:index], append([]rune(insert), runeText[index:]...)...))
}

func DeleteAt(text string, index int) string {
	runeText := []rune(text)
	return string(append(runeText[:index-1], runeText[index:]...))
}

func LenString(s string) int {
	return utf8.RuneCountInString(s)
}
