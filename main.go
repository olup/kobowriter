package main

import (
	"fmt"
	"os/exec"

	"github.com/asaskevich/EventBus"

	_ "embed"

	"github.com/olup/kobowriter/screener"
	"github.com/olup/kobowriter/utils"
	"github.com/olup/kobowriter/views"
)

var saveLocation = "/mnt/onboard/.adds/kobowriter"
var filename = "autosave.txt"

func main() {
	fmt.Println("Program started")

	// kill all nickel related stuff. Will need a reboot to find back the usual
	fmt.Println("Killing kobo programs ...")
	// exec.Command(`killall`, `-q`, `-TERM`, `nickel`, `hindenburg`, `sickel`, `fickel`, `adobehost`, `foxitpdf`, `iink`, `fmon`).Run()
	exec.Command("killall", "-s", "SIGKILL", "KoboMenu").Run()

	// rotate screen
	fmt.Println("Rotate screen ...")
	exec.Command(`fbdepth`, `--rota`, `2`).Run()

	// initialise fbink
	fmt.Println("Init FBInk ...")

	screen := screener.InitScreen()
	defer screen.Clean()

	bus := EventBus.New()

	c := make(chan bool)
	defer close(c)

	bus.SubscribeAsync("REQUIRE_KEYBOARD", func() {
		findKeyboard(screen, bus)
	}, false)

	bus.SubscribeAsync("QUIT", func() {
		screen.PrintAlert("Good Bye !", 500)

		// quitting
		c <- true
		return
	}, false)

	var unmount func()
	bus.SubscribeAsync("ROUTING", func(routeName string) {
		if unmount != nil {
			unmount()
		}

		switch routeName {
		case "document":
			config := utils.LoadConfig(saveLocation)
			unmount = views.Document(screen, bus, config.LastOpenedDocument)
		case "menu":
			unmount = views.MainMenu(screen, bus, saveLocation)
		case "file-menu":
			unmount = views.FileMenu(screen, bus, saveLocation)
		case "settings-menu":
			unmount = views.SettingsMenu(screen, bus, saveLocation)
		case "qr":
			unmount = views.Qr(screen, bus, saveLocation)

		default:
			unmount = views.Document(screen, bus, "")
		}

	}, false)

	// init
	bus.Publish("REQUIRE_KEYBOARD")

	for quit := range c {
		if quit {
			break
		}
	}

	println("yo")

}
