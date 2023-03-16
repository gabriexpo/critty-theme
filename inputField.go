package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getInputField(app *tview.Application) *tview.InputField {
	input := tview.NewInputField().
		SetLabel("Search theme: ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEsc {
				app.Stop()
			}
			return
		})

	return input
}
