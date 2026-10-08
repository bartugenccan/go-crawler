package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const retryCount = 3

var delays = []time.Duration{5 * time.Second, 10 * time.Second}

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

func FetchPage(url string) ([]byte, error) {
	var lastErr error

	for i := 1; i <= retryCount; i++ {
		data, status, err := fetchOnce(url)

		if err != nil {
			lastErr = err
			if status == 0 || shouldRetry(status) {
				if i < retryCount {
					fmt.Printf("Deneme: %d, %v bekleniliyor. \n", i, delays[i-1])
					time.Sleep(delays[i-1])
				}
				continue
			} else {
				return nil, lastErr
			}
		}

		return data, nil
	}

	return nil, fmt.Errorf("%v deneme yapıldı, hata: %w", retryCount, lastErr)
}

func shouldRetry(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests:
		return true
	case http.StatusInternalServerError:
		return true
	case http.StatusBadGateway:
		return true
	case http.StatusServiceUnavailable:
		return true
	case http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func fetchOnce(url string) ([]byte, int, error) {
	res, err := httpClient.Get(url)

	if err != nil {
		return nil, 0, fmt.Errorf("%s isteği başarısız oldu: %w", url, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, fmt.Errorf("%s isteği %s statusu ile başarısız", url, res.Status)
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("%s verisi okunamadı: %w", url, err)
	}

	return data, res.StatusCode, nil
}
