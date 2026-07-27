# Unit Converter

Aplikasi web sederhana untuk mengonversi nilai antar satuan panjang, berat, dan suhu. Form di-submit ke server, server menghitung konversi dan me-render hasilnya kembali di halaman yang sama.

Project dari [roadmap.sh - Unit Converter](https://roadmap.sh/projects/unit-converter).

## Menjalankan

```bash
go run .
```

Lalu buka `http://localhost:8080` di browser.

## Kategori & satuan yang didukung

- **Length**: millimeter, centimeter, meter, kilometer, inch, foot, yard, mile
- **Weight**: milligram, gram, kilogram, ounce, pound
- **Temperature**: celsius, fahrenheit, kelvin

## Cara pakai

1. Pilih tab kategori (Length / Weight / Temperature) di bagian atas halaman.
2. Isi nilai yang ingin dikonversi, lalu pilih satuan asal (`From`) dan satuan tujuan (`To`).
3. Klik **Convert** — hasil atau pesan error akan ditampilkan di halaman yang sama.

## Struktur

- `converter/` — logika konversi murni, terpisah per kategori (`length.go`, `weight.go`, `temperature.go`), masing-masing juga expose fungsi `XxxUnits()` untuk daftar satuan yang didukung.
- `templates/index.html` — template HTML (form, tab, hasil/error) dirender via `html/template`.
- `main.go` — HTTP server dan handler (GET menampilkan form kosong, POST memproses dan menampilkan hasil).

Tidak menggunakan database — semua perhitungan dilakukan langsung saat request.
