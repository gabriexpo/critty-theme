package main

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
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
		} else if normal && strings.Contains(l, `"#`) {
			colors = append(colors, strings.Split(l, `"`)[1])
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
