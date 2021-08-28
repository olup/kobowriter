package views

import (
	"os"

	"github.com/asaskevich/EventBus"
	"github.com/olup/kobowriter/event"
	"github.com/olup/kobowriter/screener"
	"github.com/olup/kobowriter/utils"
	"github.com/skip2/go-qrcode"
)

func Qr(screen *screener.Screen, bus EventBus.Bus, saveLocation string) func() {
	onKey := func(event event.KeyEvent) {
		screen.Clear()
		bus.Publish("ROUTING", "menu")
	}

	bus.SubscribeAsync("KEY", onKey, false)

	// Display QR on mount
	screen.Clear()
	config := utils.LoadConfig(saveLocation)
	content, err := os.ReadFile(config.LastOpenedDocument)
	if err != nil {
		bus.Publish("ROUTING", "menu")
	}

	image, _ := qrcode.Encode(string(content), qrcode.High, 800)

	screen.PrintPng(image, 800, 800, 100, 100)

	return func() {
		bus.Unsubscribe("KEY", onKey)
	}
}
