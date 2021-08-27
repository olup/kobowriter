package main

import (
	"fmt"
	"os/exec"

	"github.com/asaskevich/EventBus"
	"github.com/shermp/go-fbink-v2/gofbink"

	_ "embed"
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
	screen := Screen{}
	screen.init()
	defer screen.clean()

	bus := EventBus.New()

	c := make(chan bool)
	defer close(c)

	bus.SubscribeAsync("REQUIRE_KEYBOARD", func() {
		findKeyboard(&screen, bus)
	}, false)

	bus.SubscribeAsync("QUIT", func() {
		screen.fb.FBprint("Good bye", &gofbink.FBInkConfig{
			IsInverted: true,
			IsCleared:  true,
		})

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
			config := loadConfig()
			unmount = documentMount(&screen, bus, config.LastOpenedDocument)
		case "menu":
			unmount = menuMount(&screen, bus)
		case "file-menu":
			unmount = fileMenuMount(&screen, bus)
		case "settings-menu":
			unmount = settingsMenuMount(&screen, bus)
		case "qr":
			unmount = mountQr(&screen, bus)

		default:
			unmount = documentMount(&screen, bus, "")
		}

	}, false)

	// init
	bus.Publish("REQUIRE_KEYBOARD")

	for quit := range c {
		if quit {
			break
		}
	}

}
