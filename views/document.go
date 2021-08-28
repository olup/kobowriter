package views

import (
	"os"
	"path"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/olup/kobowriter/event"
	"github.com/olup/kobowriter/matrix"
	"github.com/olup/kobowriter/screener"
	"github.com/olup/kobowriter/utils"
)

func Document(screen *screener.Screen, bus EventBus.Bus, documentPath string) func() {
	docContent := []byte("")
	if documentPath != "" {
		docContent, _ = os.ReadFile(documentPath)
	}
	text := &TextView{
		width:       int(screen.Width) - 4,
		height:      int(screen.Height) - 2,
		content:     "",
		scroll:      0,
		cursorIndex: 0,
	}

	text.setContent(string(docContent))
	text.setCursorIndex(utils.LenString(string(docContent)))

	onEvent := func(e event.KeyEvent) {
		linesToMove := 1
		if e.IsCtrl {
			linesToMove = text.height
		}

		// if date combo
		if e.IsChar {
			text.setContent(utils.InsertAt(text.content, e.KeyChar, text.cursorIndex))
			text.setCursorIndex(text.cursorIndex + 1)
		} else {
			// if is modifier key
			switch e.KeyValue {
			case "KEY_BACKSPACE":
				text.setContent(utils.DeleteAt(text.content, text.cursorIndex))
				text.setCursorIndex(text.cursorIndex - 1)
			case "KEY_DEL":
				if text.cursorIndex < utils.LenString(text.content) {
					text.setContent(utils.DeleteAt(text.content, text.cursorIndex+1))
				}
			case "KEY_SPACE":
				text.setContent(utils.InsertAt(text.content, " ", text.cursorIndex))
				text.setCursorIndex(text.cursorIndex + 1)
			case "KEY_ENTER":
				text.setContent(utils.InsertAt(text.content, "\n", text.cursorIndex))
				text.setCursorIndex(text.cursorIndex + 1)
			case "KEY_RIGHT":
				text.setCursorIndex(text.cursorIndex + 1)
			case "KEY_LEFT":
				text.setCursorIndex(text.cursorIndex - 1)
			case "KEY_DOWN":
				text.setCursorPos(Position{
					x: text.cursorPos.x,
					y: text.cursorPos.y + linesToMove,
				})
			case "KEY_UP":
				text.setCursorPos(Position{
					x: text.cursorPos.x,
					y: text.cursorPos.y - linesToMove,
				})
			case "KEY_ESC":
				bus.Publish("ROUTING", "menu")
			case "KEY_F1":
				text.setContent(utils.InsertAt(text.content, time.Now().Format("02/01/2006"), text.cursorIndex))
				text.setCursorIndex(text.cursorIndex + 10)
			}
		}

		compiledMatrix := matrix.PasteMatrix(screen.GetOriginalMatrix(), text.renderMatrix(), 2, 1)
		screen.Print(compiledMatrix)

		if documentPath != "" {
			os.WriteFile(path.Join(documentPath), []byte(text.content), 0644)
		}
	}

	bus.SubscribeAsync("KEY", onEvent, false)

	// display
	bus.Publish("KEY", event.KeyEvent{})

	return func() {
		bus.Unsubscribe("KEY", onEvent)
	}
}
