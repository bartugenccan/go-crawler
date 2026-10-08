# Mimari

## Amaç

books.toscrape.com'daki bütün sayfaları gezip kitap adlarını toplayan bir crawler. Go'yu, sistem ve güvenlik konularıyla birlikte öğrenmek için yazıldı.

## Akış

```
main
 └─► crawler.Crawl(baseUrl) ─────────────► Result{Books, SkippedBooksCount, FailedPageAdresses}
       │  1) sayfa 1..3: sayfa sayısını öğren
       │  2) kalan sayfalar: 2 sn bekle, sırayla çek
       │
       ├─► fetcher.FetchPage(url) ───────► []byte (ham HTML) | error
       └─► parser.ParseBooks(data) ──────► PageResult{Booklist, SkippedBookCount, TotalPageCount} | error
```

## Paketler

| Paket | Sorumluluğu |
|---|---|
| `main` | Taramayı başlatır, sonucu yazdırır, hata olursa `1` koduyla çıkar. |
| `internal/crawler` | Sayfadan sayfaya gezer, sonuçları toplar, başarısız sayfaları kaydeder. |
| `internal/fetcher` | Tek bir adresi HTTP ile güvenli şekilde getirir (timeout, tekrar deneme, boyut sınırı). |
| `internal/parser` | Bir sayfanın HTML'inden kitapları ve toplam sayfa sayısını çıkarır. |

Her paket sadece kendi işini bilir: `main` HTML'den, `parser` HTTP'den, `fetcher` içerikten habersizdir.

## Çalıştırma

```
go run .
```

Tam tarama yaklaşık 100 saniye sürer (50 sayfa, aralarda 2 sn bekleme). Beklenen sonuç: 1000 kitap.
