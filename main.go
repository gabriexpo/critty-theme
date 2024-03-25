package main

import (
	"fmt"
	"log"
	"os"
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

func main() {

	themes := readColorSchemes("")

	cfg, err := getCurrentConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(os.Args) < 2 {
		fmt.Println("No flags...")
		fmt.Print(helpText)
		return
	}

	switch os.Args[1] {
	case "-s", "--set":
		if len(os.Args) < 3 {
			log.Fatal("To few arguments after -s/--set flag")
		}

		selectedTheme, ok := themes[os.Args[2]]
		if !ok {
			log.Fatalf("Error: colorscheme %s not present\n", os.Args[2])
		}

		setTheme(selectedTheme, cfg)
	case "-h", "--help":
		fmt.Print(helpText)
	case "-i", "--interactive":
		tui(themes, cfg)
	default:
		fmt.Println("No flags recognized...")
		fmt.Print(helpText)
	}

}
