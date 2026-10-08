package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Book struct {
	Name        string
	IsShortened bool
}

func ParseBooks(data []byte) ([]Book, int, error) {

	doc, err := html.Parse(bytes.NewReader(data))

	if err != nil {
		return nil, 0, fmt.Errorf("html çözümlenemedi: %w", err)
	}

	result := find(doc, "ol", "row")

	if result == nil {
		return nil, 0, errors.New("html ol.row elementi bulunamadı")
	}

	var books []Book
	skippedBooksCount := 0

	for c := result.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "li" {
			article := find(c, "article", "product_pod")

			if article != nil {
				title, isShortened := extractTitle(article)
				if title == "" {
					skippedBooksCount++
					continue
				}
				books = append(books, Book{Name: title, IsShortened: isShortened})
			}

		}
	}

	return books, skippedBooksCount, nil
}

func getAttr(n *html.Node, key string) string {
	if n != nil {
		for _, attr := range n.Attr {
			if attr.Key == key {
				return attr.Val
			}
		}
	}

	return ""
}

func find(n *html.Node, tag, class string) *html.Node {
	if n == nil {
		return nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag && (class == "" || getAttr(c, "class") == class) {
			return c
		}

		if res := find(c, tag, class); res != nil {
			return res
		}
	}

	return nil
}

// a'nın title'ı → img'nin alt'ı → a'nın metni.

func extractTitle(article *html.Node) (string, bool) {
	res := find(article, "h3", "")
	aTag := find(res, "a", "")
	resImg := find(article, "img", "thumbnail")

	title := getAttr(aTag, "title")

	trimmedTitle := strings.TrimSpace(title)

	if trimmedTitle != "" {
		return hasSuffixCheck(trimmedTitle, ". . .", "...", "…")
	}

	title2 := getAttr(resImg, "alt")
	trimmedTitle2 := strings.TrimSpace(title2)

	if trimmedTitle2 != "" {
		return hasSuffixCheck(trimmedTitle2, ". . .", "...", "…")
	}

	if aTag != nil && aTag.FirstChild != nil && aTag.FirstChild.Type == html.TextNode {
		title := aTag.FirstChild.Data
		trimmedTitle := strings.TrimSpace(title)
		if trimmedTitle != "" {
			return hasSuffixCheck(trimmedTitle, ". . .", "...", "…")
		}
	}

	return "", false

}

func hasSuffixCheck(title, suffix, suffix2, suffix3 string) (string, bool) {
	if strings.HasSuffix(title, suffix) || strings.HasSuffix(title, suffix2) || strings.HasSuffix(title, suffix3) {
		return title, true
	} else {
		return title, false
	}
}
