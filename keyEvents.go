package main

import (
	"github.com/MarinX/keylogger"
	"github.com/asaskevich/EventBus"
)

type KeyEvent struct {
	isCtrl      bool
	isAlt       bool
	isAltGr     bool
	isShift     bool
	isShiftLock bool
	keyCode     int
	isChar      bool
	keyChar     string
	keyValue    string
}

func getKeyEvent(k *keylogger.KeyLogger, b EventBus.Bus) {
	event := KeyEvent{
		isShift:     false,
		isShiftLock: false,
		isAltGr:     false,
		isAlt:       false,
		isCtrl:      false,
	}

	events := k.Read()
	for e := range events {
		if e.Type == keylogger.EvKey {

			keyValue := KeyCode[int(e.Code)]
			if keyValue == "" {
				continue
			}

			event.keyChar = ""
			event.isChar = false
			event.keyCode = int(e.Code)
			event.keyValue = keyValue

			if e.KeyPress() {
				switch keyValue {
				case "KEY_L_SHIFT", "KEY_R_SHIFT":
					event.isShift = true
				case "KEY_CAPSLOCK":
					event.isShiftLock = !event.isShiftLock
				case "KEY_ALT_GR":
					event.isAltGr = true
				case "KEY_L_ALT":
					event.isAlt = true
				case "KEY_L_CTRL", "KEY_R_CTRL":
					event.isCtrl = true

				}
			}

			if e.KeyRelease() {
				switch keyValue {
				case "KEY_L_SHIFT", "KEY_R_SHIFT":
					event.isShift = false
				case "KEY_ALT_GR":
					event.isAltGr = false
				case "KEY_L_GR":
					event.isAlt = false
				case "KEY_L_CTRL", "KEY_R_CTRL":
					event.isCtrl = false
				}
			} else {

				// letters
				if isLetter(keyValue) {
					event.isChar = true
					if event.isShift || event.isShiftLock {
						event.keyChar = KeyCodeMaj[int(e.Code)]
					} else if event.isAltGr {
						event.keyChar = KeyCodeAltGr[int(e.Code)]
					} else {
						event.keyChar = KeyCode[int(e.Code)]
					}
				}

				b.Publish("KEY", event)
			}

		}
	}
	println("lost keyboadr")
	b.Publish("REQUIRE_KEYBOARD")
}
