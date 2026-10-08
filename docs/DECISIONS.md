# Kararlar

Her kararın nedeni ve kabul edilen bedeli. Yeni bir karar verdiğinde buraya ekle.

## Ayrıştırma

**HTML parser (`golang.org/x/net/html`), metin arama değil.**
Neden: metin araması (`strings`, regex) esnek ve güvenilir değil. Parser HTML'i ağaç olarak okur.
Bedel: ilk dış bağımlılık. `go.sum` ile sürümü ve içeriği doğrulanıyor.

**Ağaçta sabit yol değil, recursive arama (`find(n, tag, class)`).**
Neden: sabit yol (`html > body > div > ...`) araya bir `div` eklenince bozulur. Arama, sarmalayıcılardan bağımsız çalışır.
Ayrıntı: `find` tek bir genel fonksiyon. `class` boşsa class'a bakılmaz.

**Başlık sırası: `a[title]` → `img[alt]` → `a` metni.**
Neden: ilk ikisi tam adı içerir, link metni kısaltılmış olabilir (`A Light in the ...`).
Kısaltma işareti: ad `...`, `. . .` ya da `…` ile bitiyorsa `IsShortened = true`.
Bedel: adı gerçekten `...` ile biten bir kitap yanlışlıkla "kısaltılmış" işaretlenir. Bu kabul edildi.

**Başlık bulunamazsa kitap atlanır ve sayılır.**
Neden: uydurma bir değer ("Unknown book") yazmak yerine eksik veri. Sayaç kaybı görünür kılar.

## Paket yapısı

**Ayrı paketler (`internal/fetcher`, `internal/parser`, `internal/crawler`).**
Neden: Go'nun paket sınırlarını baştan deneyimlemek. Dışarıya sadece gereken açılır, yardımcılar küçük harfle gizli kalır.

**Parser'ın API'si yüksek seviyeli: `ParseBooks([]byte) (PageResult, error)`.**
Neden: HTML yapısı (`ol.row`, `li`, `product_pod`) sadece parser'da kalır. Site değişirse sadece parser değişir.

**Hatalar döndürülür, `main`'de ele alınır.**
Neden: paketler ne yapılacağına karar vermez (dur, atla, tekrar dene). Kararı çağıran verir. `main` hatada `os.Exit(1)` ile çıkar.

**Bulunan kitap sayısı ayrı bir sayaçta değil, `len(Books)`'tan.**
Neden: aynı bilgiyi iki yerde tutmamak (tek doğruluk kaynağı).

## Ağ ve güvenlik (fetcher)

**Paylaşılan tek `http.Client`, istek başına 20 sn timeout.**
Neden: cevap vermeyen bir sunucu programı sonsuza kadar bekletmesin. Tek istemci bağlantıları tekrar kullanır.

**Sıra: hata kontrolü → `defer Close` → status kontrolü → okuma.**
Neden: `Get` hatasız döndüyse body açıktır, status `404` olsa bile. Önce kapatma garantisi, sonra erken çıkış.

**Status `200` değilse hata.**
Neden: hata sayfasını kitap listesi gibi işlememek.

**Tekrar deneme: 3 deneme, aralarda 5 ve 10 sn.**
Sadece geçici durumlarda: ağ hatası, `429`, `500`, `502`, `503`, `504`. Diğerleri (örn. `404`) hemen döner.
Neden: geçici sorunları atlatmak, ama zor durumdaki sunucuya yüklenmemek (artan bekleme).
Ayrıntı: her deneme ayrı bir fonksiyonda (`fetchOnce`), çünkü döngü içindeki `defer` başarısız denemelerin body'lerini açık bırakırdı.

**Body sınırı 500 KB. Aşılırsa hata, kesilmiş veri değil.**
Neden: sınırsız okuma bellek tüketir. Kesilmiş HTML ise sessiz kitap kaybına ve parser'da yanlış yerde hata aramaya yol açar.
İki katman: `Content-Length` ile erken red, `io.LimitReader` (sınır+1) ile okuma sırasında kesin kontrol. Header yalan söyleyebilir.

## Tarama (crawler)

**Adresler kalıptan üretilir (`catalogue/page-N.html`), "next" linki takip edilmez.**
Neden: bu sitede her sayfa aynı statik yapıda.
Bedel: siteye özel. Kalıp değişirse bozulur.

**Sayfa sayısı sitenin kendisinden okunur (`Page 1 of 50`).**
Neden: site büyürse kod değişmeden uyum sağlar. Sayfa sayısı sabit yazılmaz.
Yedek: ilk sayfa çekilemezse en fazla 3 sayfa denenir. Hiçbirinden öğrenilemezse tarama hata ile biter.

**Başarısız sayfa atlanır ve adresi kaydedilir.**
Neden: tek sayfa yüzünden bütün tarama kaybolmasın, ama kayıp görünür olsun.

**Sayfalar arasında 2 sn bekleme.**
Neden: başkasının sunucusuna yük bindirmemek. Bedel: tam tarama yaklaşık 100 sn sürer.

**Tarama döngüsü `crawler` paketinde, `main`'de değil.**
Neden: `main` sadece başlatır ve gösterir. Dolaşma işi ayrı bir sorumluluk ve büyüyecek (eşzamanlılık).
