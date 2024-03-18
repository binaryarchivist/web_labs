package utils

import (
	"fmt"
	"regexp"
)

func ParseHTML(html string) string {
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	return tagRegex.ReplaceAllString(html, "")
}

func ParseJSON(json string) string {
	fmt.Printf("HTML content: %s", json)

	return "ParseJSON Function"
}
