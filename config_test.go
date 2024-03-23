package main

import (
	"reflect"
	"testing"
)

func Test_parseTOMLColorSchemes(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]ColorScheme
	}{
		{
			"First",
			args{[]byte(`
[schemes.adventure_time.bright]
black = "#4e7cbf"
blue = "#1997c6"
cyan = "#c8faf4"
green = "#9eff6e"
magenta = "#9b5953"
red = "#fc5f5a"
white = "#f6f5fb"
yellow = "#efc11a"
                `)},
			map[string]ColorScheme{
				"adventure_time": {
					map[string]map[string]string{
						"bright": {
							"black":   "#4e7cbf",
							"blue":    "#1997c6",
							"cyan":    "#c8faf4",
							"green":   "#9eff6e",
							"magenta": "#9b5953",
							"red":     "#fc5f5a",
							"white":   "#f6f5fb",
							"yellow":  "#efc11a",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTOMLColorSchemes(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTOMLColorSchemes() = %v, want %v", got, tt.want)
			}
		})
	}
}
