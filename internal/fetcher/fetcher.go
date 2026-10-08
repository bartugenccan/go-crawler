package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

func FetchPage(url string) ([]byte, error) {
	res, err := httpClient.Get(url)

	if err != nil {
		return nil, fmt.Errorf("%s isteği başarısız oldu: %w", url, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s isteği %s statusu ile başarısız", url, res.Status)
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("%s verisi okunamadı: %w", url, err)
	}

	return data, nil
}
