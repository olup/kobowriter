package views

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/asaskevich/EventBus"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/olup/kobowriter/event"
	"github.com/olup/kobowriter/matrix"
	"github.com/olup/kobowriter/screener"
	"github.com/olup/kobowriter/utils"
)

type Option struct {
	label  string
	action func()
}

func createMenu(title string, options []Option) func(screen *screener.Screen, bus EventBus.Bus) func() {
	return func(screen *screener.Screen, bus EventBus.Bus) func() {
		selected := 0
		onKey := func(e event.KeyEvent) {

			if e.KeyValue == "KEY_UP" && selected > 0 {
				selected--
			}
			if e.KeyValue == "KEY_DOWN" && selected < len(options)-1 {
				selected++
			}

			if e.KeyValue == "KEY_ENTER" {
				options[selected].action()
			}

			line := 1

			matrixx := screen.GetOriginalMatrix()
			matrixx = matrix.PasteMatrix(matrixx, matrix.CreateMatrixFromText(title+"\n"+strings.Repeat("=", utils.LenString(title)), utils.LenString(title)), 4, line)

			line += 2

			for i, option := range options {
				optionMatrix := matrix.CreateMatrixFromText(option.label, utils.LenString(option.label))
				if selected == i {
					optionMatrix = matrix.InverseMatrix(optionMatrix)
				}
				matrixx = matrix.PasteMatrix(matrixx, optionMatrix, 4, line+i)
			}

			screen.Print(matrixx)
		}

		bus.SubscribeAsync("KEY", onKey, false)

		// display
		bus.Publish("KEY", event.KeyEvent{})

		return func() {
			bus.Unsubscribe("KEY", onKey)
		}
	}
}

func MainMenu(screen *screener.Screen, bus EventBus.Bus, saveLocation string) func() {
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
				bus.Publish("ROUTING", "qr")
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

				config := utils.LoadConfig(saveLocation)
				config.LastOpenedDocument = path.Join(saveLocation, id+".txt")
				utils.SaveConfig(config, saveLocation)

				bus.Publish("ROUTING", "document")
			},
		},
		{
			label: "Settings",
			action: func() {
				bus.Publish("ROUTING", "settings-menu")
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

	return createMenu("Menu", options)(screen, bus)
}

func FileMenu(screen *screener.Screen, bus EventBus.Bus, saveLocation string) func() {
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
			if utils.LenString(label) > 30 {
				label = string([]rune(label)[0:30]) + "..."
			}
			options = append(options, Option{
				label: label,

				action: func() {
					config := utils.LoadConfig(saveLocation)
					config.LastOpenedDocument = filePath
					utils.SaveConfig(config, saveLocation)

					bus.Publish("ROUTING", "document")
				},
			})
		}

	}

	return createMenu("Open File", options)(screen, bus)
}

func SettingsMenu(screen *screener.Screen, bus EventBus.Bus, saveLocation string) func() {
	options := []Option{
		{
			label: "Back",
			action: func() {
				bus.Publish("ROUTING", "menu")
			},
		},
		{
			label: "Toggle light",
			action: func() {
				lightPath := "/sys/class/backlight/mxc_msp430_fl.0/brightness"
				light := "0"
				presentLightRaw, _ := os.ReadFile(lightPath)
				presentLight := strings.TrimSuffix(string(presentLightRaw), "\n")

				if presentLight == "0" {
					light = "10"
				} else {
					light = "0"
				}

				os.WriteFile(lightPath, []byte(light), os.ModePerm)
			},
		},
	}

	return createMenu("Open File", options)(screen, bus)
}
