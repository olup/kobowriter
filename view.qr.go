package main

import (
	"os"

	"github.com/asaskevich/EventBus"
	"github.com/shermp/go-fbink-v2/gofbink"
	"github.com/skip2/go-qrcode"
)

func mountQr(screen *Screen, bus EventBus.Bus) func() {
	onKey := func(event KeyEvent) {
		screen.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})
		bus.Publish("ROUTING", "menu")
	}

	bus.SubscribeAsync("KEY", onKey, false)

	// Display QR on mount
	screen.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})
	config := loadConfig()
	content, err := os.ReadFile(config.LastOpenedDocument)
	if err != nil {
		bus.Publish("ROUTING", "menu")
	}

	image, _ := qrcode.Encode(string(content), qrcode.High, 800)

	screen.printPng(image, 800, 800, 100, 100)

	return func() {
		bus.Unsubscribe("KEY", onKey)
	}
}
