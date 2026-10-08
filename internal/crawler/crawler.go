package crawler

import (
	"fmt"
	"go-crawler/internal/fetcher"
	"go-crawler/internal/parser"
	"time"
)

const delayBetweenPages = 2 * time.Second
const maxProbePages = 3

type Result struct {
	Books              []parser.Book
	SkippedBooksCount  int
	FailedPageAdresses []string
}

func Crawl(baseUrl string) (Result, error) {
	var result Result
	var pageCount int
	var currentPage int

	for currentPage = 1; currentPage <= maxProbePages; currentPage++ {

		if currentPage != 1 {
			time.Sleep(delayBetweenPages)
		}

		currentPageAddress := pageAddress(baseUrl, currentPage)
		res, err := fetchAndParse(currentPageAddress)

		if err != nil {
			result.FailedPageAdresses = append(result.FailedPageAdresses, currentPageAddress)
			continue
		}

		result.Books = append(result.Books, res.Booklist...)
		result.SkippedBooksCount += res.SkippedBookCount
		fmt.Printf("Sayfa %d çekildi.\n", currentPage)
		fmt.Printf("%d kitap çekildi (toplam %d).\n", len(res.Booklist), len(result.Books))

		if res.TotalPageCount > 0 {
			pageCount = res.TotalPageCount
			break
		}

	}

	if pageCount == 0 {
		return result, fmt.Errorf("sayfa sayısı öğrenilemedi")
	}

	for i := currentPage + 1; i <= pageCount; i++ {
		time.Sleep(delayBetweenPages)

		pageAdress := pageAddress(baseUrl, i)
		res, err := fetchAndParse(pageAdress)

		if err != nil {
			result.FailedPageAdresses = append(result.FailedPageAdresses, pageAdress)
			continue
		}

		result.Books = append(result.Books, res.Booklist...)
		result.SkippedBooksCount += res.SkippedBookCount
		fmt.Printf("Sayfa %d/%d çekildi.\n", i, pageCount)
		fmt.Printf("%d kitap çekildi (toplam %d).\n", len(res.Booklist), len(result.Books))

	}
	return result, nil
}

// pageAddress builds the catalogue URL for the given page number.
func pageAddress(baseUrl string, page int) string {
	return fmt.Sprintf("%scatalogue/page-%d.html", baseUrl, page)
}

// fetchAndParse fetches a single page and parses its books.
func fetchAndParse(pageAddress string) (parser.PageResult, error) {
	data, err := fetcher.FetchPage(pageAddress)

	if err != nil {
		return parser.PageResult{}, err
	}

	return parser.ParseBooks(data)
}
