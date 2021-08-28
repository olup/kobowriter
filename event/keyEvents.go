package event

import (
	"github.com/MarinX/keylogger"
	"github.com/asaskevich/EventBus"
	"github.com/olup/kobowriter/utils"
)

type KeyEvent struct {
	IsCtrl      bool
	IsAlt       bool
	IsAltGr     bool
	IsShift     bool
	IsShiftLock bool
	KeyCode     int
	IsChar      bool
	KeyChar     string
	KeyValue    string
}

func BindKeyEvent(k *keylogger.KeyLogger, b EventBus.Bus) {
	event := KeyEvent{
		IsShift:     false,
		IsShiftLock: false,
		IsAltGr:     false,
		IsAlt:       false,
		IsCtrl:      false,
	}

	events := k.Read()
	for e := range events {
		if e.Type == keylogger.EvKey {

			keyValue := KeyCode[int(e.Code)]
			if keyValue == "" {
				continue
			}

			event.KeyChar = ""
			event.IsChar = false
			event.KeyCode = int(e.Code)
			event.KeyValue = keyValue

			if e.KeyPress() {
				switch keyValue {
				case "KEY_L_SHIFT", "KEY_R_SHIFT":
					event.IsShift = true
				case "KEY_CAPSLOCK":
					event.IsShiftLock = !event.IsShiftLock
				case "KEY_ALT_GR":
					event.IsAltGr = true
				case "KEY_L_ALT":
					event.IsAlt = true
				case "KEY_L_CTRL", "KEY_R_CTRL":
					event.IsCtrl = true

				}
			}

			if e.KeyRelease() {
				switch keyValue {
				case "KEY_L_SHIFT", "KEY_R_SHIFT":
					event.IsShift = false
				case "KEY_ALT_GR":
					event.IsAltGr = false
				case "KEY_L_GR":
					event.IsAlt = false
				case "KEY_L_CTRL", "KEY_R_CTRL":
					event.IsCtrl = false
				}
			} else {

				// letters
				if utils.IsLetter(keyValue) {
					event.IsChar = true
					if event.IsShift || event.IsShiftLock {
						event.KeyChar = KeyCodeMaj[int(e.Code)]
					} else if event.IsAltGr {
						event.KeyChar = KeyCodeAltGr[int(e.Code)]
					} else {
						event.KeyChar = KeyCode[int(e.Code)]
					}
				}

				b.Publish("KEY", event)
			}

		}
	}
	println("lost keyboadr")
	b.Publish("REQUIRE_KEYBOARD")
}
