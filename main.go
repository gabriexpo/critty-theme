package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gabriexpo/go-functional"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func generateFooterColorBars(theme ColorScheme, name string) string {

	footerText := name + "\n"

	for _, v := range theme.colors["normal"] {
		footerText += fmt.Sprintf("[%s]████", v)
	}

	footerText += "\n"

	for _, v := range theme.colors["bright"] {
		footerText += fmt.Sprintf("[%s]████", v)
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

	themes := readColorSchemes("")

	themesNames := functional.ExtractKeys(themes)
	sort.Slice(themesNames, func(i, j int) bool {
		return themesNames[i] < themesNames[j]
	})

	cfg, err := getCurrentConfig()
	if err != nil {
		panic(err.Error())
	}

	currentTheme := parseColorScheme(cfg["colors"].(map[string]interface{}))

	list := getListItem(themesNames, app)

	input := getInputField(app)

	selected_theme, _ := list.GetItemText(list.GetCurrentItem())

	footer_left := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("Current: " + generateFooterColorBars(currentTheme, "current"))

	footer_right := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("Selected: " + generateFooterColorBars(themes[selected_theme], selected_theme))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			list.Blur()
			input.Focus(nil)

			return nil
		} else if event.Key() == tcell.KeyEnter {

			t, _ := list.GetItemText(list.GetCurrentItem())
			err := setTheme(themes[t], cfg)
			if err != nil {
				panic(err)
			}

			footer_left.SetText("Current: " + generateFooterColorBars(themes[t], t))
		}

		return event
	})

	list.SetChangedFunc(func(index int, text, secText string, shortcut rune) {
		footer_right.SetText("Selected: " + generateFooterColorBars(themes[text], text))
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
