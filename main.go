package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func generateFooterColorBars(theme string) string {
	colors := getThemeColors(theme)

	footerText := theme + "\n"
	for i := 0; i < 16; i++ {
		footerText += fmt.Sprintf("[%s]████", colors[i])
		if i == 7 {
			footerText += "\n"
		}
	}

	return footerText
}

func setSearchedItem(s string, list *tview.List) {
	if s == "" {
		return
	}

	for i := 0; i < list.GetItemCount(); i++ {
		if text, _ := list.GetItemText(i); strings.Contains(text, s) {
			list.SetCurrentItem(i)
			return
		}
	}
	return
}

func main() {

	app := tview.NewApplication()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("[yellow]CRITTY-THEME")

	themes := getThemesList()

	list := getListItem(themes, app)

	input := getInputField(app)

	selected_theme, _ := list.GetItemText(list.GetCurrentItem())

	footer_left := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("Current: " + generateFooterColorBars(getCurrentThemeName()))

	footer_right := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("Selected: " + generateFooterColorBars(selected_theme))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			list.Blur()
			input.Focus(nil)
			return nil
		} else if event.Key() == tcell.KeyEnter {
			t, _ := list.GetItemText(list.GetCurrentItem())
			ok := changeTheme(t)
			if !ok {
				panic("!ok")
			}
			footer_left.SetText("Current: " + generateFooterColorBars(t))
		}

		return event
	})

	list.SetChangedFunc(func(index int, text, secText string, shortcut rune) {
		footer_right.SetText("Selected: " + generateFooterColorBars(text))
	})

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			input.Blur()
			list.Focus(nil)
			setSearchedItem(input.GetText(), list)
			return nil
		}
		return event
	})

	outer_grid := tview.NewGrid().
		SetRows(6, 1, 0, 1, 4, 7).
		SetColumns(50, 0, 0, 50).
		SetBorders(true).
		AddItem(title, 1, 1, 1, 2, 0, 0, false).
		AddItem(list, 2, 1, 1, 2, 0, 0, false).
		AddItem(input, 3, 1, 1, 2, 0, 0, true).
		AddItem(footer_left, 4, 1, 1, 1, 0, 0, false).
		AddItem(footer_right, 4, 2, 1, 1, 0, 0, false)

	app.SetRoot(outer_grid, true).EnableMouse(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
