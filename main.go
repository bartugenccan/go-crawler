package main

import (
	"fmt"
	"go-crawler/internal/fetcher"
	"go-crawler/internal/parser"
	"os"
)

func main() {

	data, err := fetcher.FetchPage("https://books.toscrape.com/")

	if err != nil {
		fmt.Println("sayfa çekilemedi:", err)
		os.Exit(1)
	}

	books, skippedBooksCount, err := parser.ParseBooks(data)

	if err != nil {
		fmt.Println("kitaplar bulunamadı:", err)
		os.Exit(1)
	}

	for index, book := range books {
		fmt.Println("Index:", index+1, "Book Name:", book.Name, "Is Shortened:", book.IsShortened)
	}
	fmt.Println("Skipped Books Count:", skippedBooksCount)
	fmt.Println("Found Books Count:", len(books))

}
