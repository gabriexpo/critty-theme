package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	fp "path/filepath"

	. "github.com/gabriexpo/go-functional"
	toml "github.com/pelletier/go-toml/v2"
)

type ColorScheme struct {
	colors map[string]map[string]string
}

type data struct {
	schemes interface{}
}

func readColorSchemes(file string) map[string]ColorScheme {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if file == "" {
		file = fp.Join(homedir, ".config", "alacritty", "color_schemes.toml")
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return parseTOMLColorSchemes(content)
}

func parseTOMLColorSchemes(content []byte) map[string]ColorScheme {

	var s map[string]interface{}
	toml.Unmarshal([]byte(content), &s)

	return parseUnmarshalledMap(s)
}

func parseUnmarshalledMap(unmarshalled map[string]interface{}) map[string]ColorScheme {

	m := MapMap(func(x interface{}) map[string]interface{} {
		return x.(map[string]interface{})
	}, unmarshalled["schemes"].(map[string]interface{}))

	colorSchemes := make(map[string]ColorScheme)
	for k, v := range m { // colorschemes
		cs := parseColorScheme(v)

		colorSchemes[k] = cs
	}

	return colorSchemes

}

func parseColorScheme(m map[string]interface{}) ColorScheme {

	cs := ColorScheme{colors: make(map[string]map[string]string)}
	// fmt.Println(k)

	for ki, vi := range m { // bright, normal, cursor, primary
		c := vi.(map[string]interface{})
		cs.colors[ki] = make(map[string]string)

		for kii, vii := range c { // black, blue, cyan, green, etc...
			switch m := vii.(type) {
			case string:
				cs.colors[ki][kii] = m
			default:
				log.Fatalf("Error during decoding of: %v", m)
			}
		}
	}

	return cs
}

func getCurrentConfig() (map[string]interface{}, error) {
	cfgFile := alacrittyConfigFile()

	content, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	var cfg interface{}
	err = toml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	cfgMap := cfg.(map[string]interface{})

	return cfgMap, nil
}

func setTheme(cs ColorScheme, cfg map[string]interface{}) error {

	cfgFile := alacrittyConfigFile()

	cfg["colors"] = cs.colors

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cfgFile, data, fs.FileMode(os.O_WRONLY))

	return nil
}

func alacrittyConfigFile() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return fp.Join(homedir, ".config", "alacritty", "alacritty.toml")
}
