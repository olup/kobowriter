package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/asaskevich/EventBus"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Option struct {
	label  string
	action func()
}

func createMenu(title string, options []Option) func(screen *Screen, bus EventBus.Bus) func() {
	return func(screen *Screen, bus EventBus.Bus) func() {
		selected := 0
		onKey := func(event KeyEvent) {

			if event.keyValue == "KEY_UP" && selected > 0 {
				selected--
			}
			if event.keyValue == "KEY_DOWN" && selected < len(options)-1 {
				selected++
			}

			if event.keyValue == "KEY_ENTER" {
				options[selected].action()
			}

			line := 1

			matrix := screen.originalMatrix
			matrix = pasteMatrix(matrix, createMatrixFromText(title+"\n"+strings.Repeat("=", lenString(title)), lenString(title)), 4, line)

			line += 2

			for i, option := range options {
				optionMatrix := createMatrixFromText(option.label, lenString(option.label))
				if selected == i {
					optionMatrix = inverseMatrix(optionMatrix)
				}
				matrix = pasteMatrix(matrix, optionMatrix, 4, line+i)
			}

			screen.print(matrix)
		}

		bus.SubscribeAsync("KEY", onKey, false)

		// display
		bus.Publish("KEY", KeyEvent{})

		return func() {
			bus.Unsubscribe("KEY", onKey)
		}
	}
}

func menuMount(screen *Screen, bus EventBus.Bus) func() {
	options := []Option{
		{
			label: "Back",
			action: func() {
				bus.Publish("ROUTING", "document")
			},
		},
		{
			label: "Export as QR code",
			action: func() {
			},
		},
		{
			label: "Open Document",
			action: func() {
				bus.Publish("ROUTING", "file-menu")
			},
		},
		{
			label: "New Document",
			action: func() {
				id, _ := gonanoid.New()

				config := loadConfig()
				config.LastOpenedDocument = path.Join(saveLocation, id+".txt")
				saveConfig(config)

				bus.Publish("ROUTING", "document")
			},
		},
		{
			label: "Quit to XCSoar",
			action: func() {
				exec.Command("/opt/xcsoar/bin/KoboMenu").Start()
				bus.Publish("QUIT")
			},
		},
	}

	battreryCapacity, _ := os.ReadFile("/sys/class/power_supply/mc13892_bat/capacity")
	return createMenu("Menu / BAT "+string(battreryCapacity)+"%", options)(screen, bus)
}

func fileMenuMount(screen *Screen, bus EventBus.Bus) func() {
	files, _ := os.ReadDir(saveLocation)
	options := []Option{
		{
			label: "Back",
			action: func() {
				bus.Publish("ROUTING", "menu")
			},
		},
	}
	for _, file := range files {

		if strings.HasSuffix(file.Name(), ".txt") {
			filePath := path.Join(saveLocation, file.Name())
			content, _ := os.ReadFile(path.Join(saveLocation, file.Name()))

			label := strings.Split(string(content), "\n")[0]
			if lenString(label) > 30 {
				label = string([]rune(label)[0:30]) + "..."
			}
			options = append(options, Option{
				label: label,

				action: func() {
					config := loadConfig()
					config.LastOpenedDocument = filePath
					saveConfig(config)

					bus.Publish("ROUTING", "document")
				},
			})
		}

	}

	return createMenu("Open File", options)(screen, bus)
}
