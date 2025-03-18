package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// displays the information about a task.
func NewInfoPage(p *BoardPage, listIdx, taskIdx int) tview.Primitive {
	task, err := p.data.GetTask(listIdx, taskIdx)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	info := tview.NewModal().
		SetText(fmt.Sprintf("Task: %v\n Task Description: %v", task.ItemName, task.ItemDescription)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				closeInfoPage()
			}
		})
	info.SetBorderColor(theme.BorderColor)          // Set inner border color
	info.SetTextColor(tcell.ColorBlack)             // Set text color
	info.SetButtonBackgroundColor(tcell.ColorWheat) // Set button background color
	info.SetButtonTextColor(tcell.ColorBlack)       // Set button text color
	info.SetBorderPadding(0, 0, 0, 0)
	info.SetBackgroundColor(theme.PrimitiveBackgroundColor) // Set modal background
	info.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			closeInfoPage()
		}
		switch event.Rune() {
		case 'q':
			closeInfoPage()
		}
		return event
	})
	width, height := GetSize()
	return GetCenteredModal(info, width/2, height/2)
}

func closeInfoPage() {
	pages.RemovePage("info")
	pages.SwitchToPage("board")
}
