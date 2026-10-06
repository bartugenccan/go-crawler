package main

import (
	"fmt"
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
		fmt.Println("Couldn't read the body.")
		return
	}

	fmt.Println(string(data))

}
