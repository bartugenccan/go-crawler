package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	res, err := http.Get("https://books.toscrape.com")
	if err != nil {
		fmt.Println("failed:", err)
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't read the body.", err)
		return
	}

	doc, err := html.Parse(bytes.NewReader(data))

	if err != nil {
		fmt.Println("Parse hatası", err)
		return
	}

	result := find(doc)

	if result != nil {
		targetNode := getAttr(result, "class")

		if targetNode != "" {
			fmt.Println("TargetNode:", targetNode)
			return
		} else {
			fmt.Println("Target node bulunamadı")

		}

		return
	} else {
		fmt.Println("Bulunamadı <nil>")
	}

}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func find(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "ol" {
			for _, attr := range c.Attr {
				if attr.Key == "class" && attr.Val == "row" {
					return c
				}
			}

		}

		if res := find(c); res != nil {
			return res
		}
	}

	return nil
}
