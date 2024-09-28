package twitch

import "strings"

func replaceMultiple(str string, replacements map[string]string) string {
	for old, new := range replacements {
		str = strings.ReplaceAll(str, old, new)
	}
	return str
}

func getPreferredFormat(formats []string) string {
	for _, format := range formats {
		if format == "animated" {
			return "animated"
		}
	}
	return "static"
}

func getPreferredThemeMode(modes []string) string {
	for _, mode := range modes {
		if mode == "dark" {
			return "dark"
		}
	}
	return "light"
}

func getPreferredScale(scales []string) string {
	for _, scale := range scales {
		if scale == "3.0" {
			return "3.0"
		}
		if scale == "2.0" {
			return "2.0"
		}
	}
	return "1.0"
}
