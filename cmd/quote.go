package cmd

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetMotivationalQuote will retrieve a random motivation quote.
func GetMotivationalQuote() (string, error) {
	doc, err := goquery.NewDocument("http://inspirationalshit.com/quotes")
	if err != nil {
		return "", err
	}

	container := doc.Find("#scores blockquote")
	quote := container.Find("p.text-uppercase").Text()
	quote = strings.TrimPrefix(quote, "\n")
	if quote == "" {
		return "", fmt.Errorf("failed to find quote in response body")
	}

	author := container.Find("cite").Text()
	if author == "" {
		author = "Unknown"
	}

	quote = fmt.Sprintf("\n\"%s\" - %s\n", quote, author)
	return quote, nil
}
