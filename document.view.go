package main

import (
	"os"
	"path"
	"time"

	"github.com/asaskevich/EventBus"
)

func documentMount(screen *Screen, bus EventBus.Bus, documentPath string) func() {
	docContent := []byte("")
	if documentPath != "" {
		docContent, _ = os.ReadFile(documentPath)
	}
	text := &TextView{
		width:       int(screen.state.MaxCols) - 4,
		height:      int(screen.state.MaxRows) - 2,
		content:     "",
		scroll:      0,
		cursorIndex: 0,
	}

	text.setContent(string(docContent))
	text.setCursorIndex(lenString(string(docContent)))

	onEvent := func(event KeyEvent) {
		linesToMove := 1
		if event.isCtrl {
			linesToMove = text.height
		}

		// if date combo
		if event.isChar {
			text.setContent(insertAt(text.content, event.keyChar, text.cursorIndex))
			text.setCursorIndex(text.cursorIndex + 1)
		} else {
			// if is modifier key
			switch event.keyValue {
			case "KEY_BACKSPACE":
				text.setContent(deleteAt(text.content, text.cursorIndex))
				text.setCursorIndex(text.cursorIndex - 1)
			case "KEY_DEL":
				if text.cursorIndex < lenString(text.content) {
					text.setContent(deleteAt(text.content, text.cursorIndex+1))
				}
			case "KEY_SPACE":
				text.setContent(insertAt(text.content, " ", text.cursorIndex))
				text.setCursorIndex(text.cursorIndex + 1)
			case "KEY_ENTER":
				text.setContent(insertAt(text.content, "\n", text.cursorIndex))
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
				text.setContent(insertAt(text.content, time.Now().Format("02/01/2006"), text.cursorIndex))
				text.setCursorIndex(text.cursorIndex + 10)
			}
		}

		compiledMatrix := pasteMatrix(screen.originalMatrix, text.renderMatrix(), 2, 1)
		screen.print(compiledMatrix)

		if documentPath != "" {
			os.WriteFile(path.Join(documentPath), []byte(text.content), 0644)
		}
	}

	bus.SubscribeAsync("KEY", onEvent, false)

	// display
	bus.Publish("KEY", KeyEvent{})

	return func() {
		bus.Unsubscribe("KEY", onEvent)
	}
}
