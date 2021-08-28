package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/asaskevich/EventBus"

	"github.com/olup/kobowriter/event"
	"github.com/olup/kobowriter/screener"
)

func findKeyboard(screen *screener.Screen, bus EventBus.Bus) {
	// get key logger
	keyboard := keylogger.FindKeyboardDevice()

	buttonLogger, _ := keylogger.New("/dev/input/event0")
	buttonChannel := buttonLogger.Read()

	screen.Clear()

	for len(keyboard) <= 0 {
		screen.PrintAlert("No keyboard found.\n\nPlug your keyboard or clic main button to quit.\n\nNote that [USB OTG MODE] must be turned on in order to detect the keyboard.", 30)
		time.Sleep(1 * time.Second)
		select {
		case _ = <-buttonChannel:
			println("Quitting program")
			bus.Publish("QUIT")
			exec.Command("/opt/xcsoar/bin/KoboMenu").Start()
			return
		default:
		}
		keyboard = keylogger.FindKeyboardDevice()
	}

	screen.Clear()
	fmt.Println("Found a keyboard at", keyboard)

	k, _ := keylogger.New(keyboard)
	go event.BindKeyEvent(k, bus)
	bus.Publish("ROUTING", "document")
	return
}
