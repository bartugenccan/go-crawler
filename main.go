package main

import (
	"fmt"
	"go-crawler/internal/crawler"
	"os"
)

func main() {

	res, err := crawler.Crawl("https://books.toscrape.com/")

	if err != nil {
		fmt.Println("sayfa sayısı öğrenilemedi:", err)
		os.Exit(1)
	}

	for index, book := range res.Books {
		fmt.Println("Index:", index+1, "Book Name:", book.Name, "Is Shortened:", book.IsShortened)
	}
	fmt.Println("Bulunamayan Kitap Sayısı:", res.SkippedBooksCount)
	fmt.Println("Bulunan Kitap Sayısı:", len(res.Books))
	fmt.Println("Başarısız Sayfa Sayısı:", len(res.FailedPageAdresses))

}
