package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/asaskevich/EventBus"
	"github.com/shermp/go-fbink-v2/gofbink"
)

func findKeyboard(screen *Screen, bus EventBus.Bus) {
	// get key logger
	keyboard := keylogger.FindKeyboardDevice()

	buttonLogger, _ := keylogger.New("/dev/input/event0")
	buttonChannel := buttonLogger.Read()

	screen.fb.ClearScreen(&gofbink.FBInkConfig{
		IsFlashing: true,
	})

	for len(keyboard) <= 0 {
		printAlert("No keyboard found.\n\nPlug your keyboard or clic main button to quit.\n\nNote that [USB OTG MODE] must be turned on in order to detect the keyboard.", 30, screen)
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

	screen.fb.ClearScreen(&gofbink.FBInkConfig{})
	fmt.Println("Found a keyboard at", keyboard)

	k, _ := keylogger.New(keyboard)
	go getKeyEvent(k, bus)
	bus.Publish("ROUTING", "document")
	return
}
