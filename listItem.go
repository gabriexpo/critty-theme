package main

import "github.com/rivo/tview"

func getListItem(themes []string, app *tview.Application) *tview.List {
	list := tview.NewList().
		SetDoneFunc(func() {
			app.Stop()
			return
		})

	// Set the themes adas items
	for i, t := range themes {
		list.AddItem(t, "", rune(i+int('a')), nil)
	}

	return list
}
