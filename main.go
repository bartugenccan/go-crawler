package main

import (
	"fmt"
	"go-crawler/internal/parser"
	"io"
	"net/http"
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

	books, skippedBooksCount, err := parser.ParseBooks(data)

	if err != nil {
		fmt.Println("Kitaplar bulunamadı:", err)
		return
	}

	for index, book := range books {
		fmt.Println("Index:", index+1, "Book Name:", book.Name, "Is Shortened:", book.IsShortened)
	}
	fmt.Println("Skipped Books Count:", skippedBooksCount)
	fmt.Println("Found Books Count:", len(books))

}
