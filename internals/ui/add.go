package ui

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	command "github.com/ppriyankuu/seiban/pkg/commands"
	"github.com/rivo/tview"
)

func NewAddPage(p *BoardPage, pos int) tview.Primitive {
	width, height := GetSize()
	form := tview.NewForm().
		AddInputField("Task", "", width/4, nil, nil).
		AddInputField("Task Description", "", width/4, nil, nil)
	form.SetFieldBackgroundColor(tcell.ColorWheat)  // Set background color of input fields
	form.SetFieldTextColor(tcell.ColorBlack)        // Set text color for better contrast
	form.SetButtonBackgroundColor(tcell.ColorWheat) // Set button background color
	form.SetButtonTextColor(tcell.ColorBlack)       // Set button text color
	form.SetBorderColor(theme.BorderColor)

	form = form.AddButton("Save", func() {
		taskName := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
		taskName = strings.TrimSpace(taskName)
		if len(taskName) <= 0 {
			return
		}
		taskDesc := form.GetFormItemByLabel("Task Description").(*tview.InputField).GetText()
		taskDesc = strings.TrimSpace(taskDesc)
		addTaskCommand := command.CreateAddTaskCommand(p.activeListIdx, taskName, taskDesc, pos)
		if err := p.command.Execute(addTaskCommand); err != nil {
			app.Stop()
			log.Fatal(err)
		}
		p.redraw(p.activeListIdx)
		pages.SwitchToPage("board")
	}).AddButton("Cancel", func() {
		closeAddPage()
	})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			closeAddPage()
		}
		return event
	})
	form.SetBorder(true).SetTitle("Create Task").SetTitleAlign(tview.AlignCenter)
	return GetCenteredModal(form, width/2, height/2)
}

func closeAddPage() {
	pages.RemovePage("add")
	pages.SwitchToPage("board")
}
