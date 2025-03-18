package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const helpText = `
  Controls
  
	Navigation
	────────────────────────────────
    j → Move down
    k → Move up
    h → Move left
    l → Move right
    g → Focus first
    G → Focus last
	
	Task Management
	────────────────────────────────
    a → Add under cursor
    A → Add at end
    d → Mark as done
    D → Delete
    e → Edit task
	
	Movement
	────────────────────────────────
    L → Move right
    H → Move left
    J → Move down
    K → Move up
	
	Actions
	────────────────────────────────
    Enter → View info  
    u → Undo
    Ctrl+R → Redo
    q → Quit

	────────────────────────────────
`

// displays the help page that contains all the keybinds of the application
func NewHelpPage(p *BoardPage) tview.Primitive {
	help := tview.NewModal().
		SetText(helpText).
		SetBackgroundColor(theme.PrimitiveBackgroundColor).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "OK" {
				closeHelpPage()
			}
		})

	help.SetTextColor(tcell.ColorWheat)
	help.SetBorderColor(theme.BorderColor)
	help.SetButtonTextColor(tcell.ColorBlack)
	help.SetButtonBackgroundColor(tcell.ColorWheat)
	help.SetBackgroundColor(tcell.ColorBlack) // Set modal background

	help.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			closeHelpPage()
		}
		switch event.Rune() {
		case 'q':
			closeHelpPage()
		}
		return event
	})
	width, height := GetSize()
	return GetCenteredModal(help, width/2, height/2)
}

func closeHelpPage() {
	pages.HidePage("help")
	pages.SwitchToPage("board")
}
