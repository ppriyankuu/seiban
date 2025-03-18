package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	theme *tview.Theme
	pages *tview.Pages
)

func defaultTheme() *tview.Theme {
	return &tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorWheat, // Main background color for primitives.
		ContrastBackgroundColor:     tcell.ColorWheat, // Background color for contrasting elements.
		MoreContrastBackgroundColor: tcell.ColorWheat, // Background color for more contrasting elements.
		BorderColor:                 tcell.ColorWheat, // Box borders.
		TitleColor:                  tcell.ColorWheat, // Box titles.
		GraphicsColor:               tcell.ColorWheat, // Graphics.
		PrimaryTextColor:            tcell.ColorWheat, // Primary text.
		SecondaryTextColor:          tcell.ColorWheat, // Secondary text (e.g. labels).
		TertiaryTextColor:           tcell.ColorWheat, // Tertiary text (e.g. subtitles, notes).
		InverseTextColor:            tcell.ColorWheat, // Text on primary-colored backgrounds.
		ContrastSecondaryTextColor:  tcell.ColorWheat, // Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}

func Start(fileName string) error {
	app = tview.NewApplication()
	initiate(fileName)
	if err := app.Run(); err != nil {
		return fmt.Errorf("Error running the app: %s", err)
	}
	return nil
}

func initiate(fileName string) {
	theme = defaultTheme()
	boardPage := NewBoardPage(fileName)
	boardPageFrame := boardPage.Page()
	pages = tview.NewPages().AddPage("board", boardPageFrame, true, true)
	app.SetRoot(pages, true).SetFocus(boardPageFrame)
}
