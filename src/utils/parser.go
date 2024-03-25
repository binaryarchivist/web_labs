package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func ParseHTML(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic("Failed to parse HTML: " + err.Error())
	}

	parsedHTML := walk(doc)
	return strings.TrimSpace(parsedHTML)
}

func walk(n *html.Node) string {
	content := ""

	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		return ""
	}
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			content += text + " " + "\n"
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		content += walk(c)
	}

	return content
}

func ParseJSON(jsonContent string) string {
	return jsonContent
}
