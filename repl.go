package main

import (
	"strings"
)

func cleanInput(text string) []string {
	lcText := strings.ToLower(text)
	return strings.Fields(lcText)
}
