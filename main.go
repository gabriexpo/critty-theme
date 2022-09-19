package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func readConfig() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadFile(homedir + "/.config/alacritty/alacritty.yml")
	if err != nil {
		panic(err)
	}

	return string(content)
}

// Read the configuration file, find the current theme name and return it
func getCurrentThemeName() string {
	lines := strings.Split(readConfig(), "\n")

	for _, l := range lines {
		if strings.HasPrefix(l, "colors:") {
			return l[9:]
		}
	}

	panic("color config not found...")
}

// Read the config file and find the given theme colors
func getThemeColors(theme string) []string {
	lines := strings.Split(readConfig(), "\n")

	inTheme := false
	normal := false
	colors := []string{}

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if inTheme && strings.Contains(l, "&") {
			return colors
		}

		if strings.Contains(l, theme) {
			inTheme = true
		}

		if inTheme && strings.Contains(l, "normal:") { // normal: row, start listening
			normal = true
		}

		if normal && strings.Contains(l, "'#") { // color row
			colors = append(colors, strings.Split(l, "'")[1])
		}
	}

	return colors
}

func getThemesList() []string {
	lines := strings.Split(readConfig(), "\n")

	schemes := false
	list := []string{}

	for _, l := range lines {
		if strings.Contains(l, "schemes:") {
			schemes = true
		}

		if schemes && strings.Contains(l, ": &") {
			list = append(list, strings.TrimSpace(strings.Split(l, ":")[0]))
		}

		if schemes && strings.Contains(l, "colors: *") {
			sort.Slice(list, func(i, j int) bool {
				return list[i] < list[j]
			})
			return list
		}
	}

	return list
}

// Change Alacritty theme to newTheme, return true if change goes well
func changeTheme(newTheme string) bool {

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	lines := strings.Split(readConfig(), "\n")
	for i, l := range lines {
		if strings.Contains(l, "colors: *") {
			lines[i] = strings.Split(l, "*")[0] + "*" + newTheme
			break
		}
	}

	newContent := []byte(strings.Join(lines, "\n"))

	err = ioutil.WriteFile(homedir+"/.config/alacritty/alacritty.yml", newContent, 0644)
	if err != nil {
		return false
	}

	return true
}

func footerColorBars() string {
	currentTheme := getCurrentThemeName()
	colors := getThemeColors(currentTheme)

	footerText := currentTheme + "\n"
	for i := 0; i < 16; i++ {
		footerText += fmt.Sprintf("[%s]█████", colors[i])
		if i == 7 {
			footerText += "\n"
		}
	}

	return footerText
}

func setSearchedItem(s string, list *tview.List) {
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

	list := tview.NewList().
		SetDoneFunc(func() {
			app.Stop()
			return
		})
	for i, t := range themes {
		list.AddItem(t, "", rune(i+int('a')), nil)
	}

	footer := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(footerColorBars())

	input := tview.NewInputField().
		SetLabel("Search theme: ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEsc {
				app.Stop()
			}
			return
		})

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
			footer.SetText(footerColorBars())
		}

		return event
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
		SetColumns(50, 0, 50).
		SetBorders(true).
		AddItem(title, 1, 1, 1, 1, 0, 0, false).
		AddItem(list, 2, 1, 1, 1, 0, 0, false).
		AddItem(input, 3, 1, 1, 1, 0, 0, true).
		AddItem(footer, 4, 1, 1, 1, 0, 0, false)

	app.SetRoot(outer_grid, true).EnableMouse(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
