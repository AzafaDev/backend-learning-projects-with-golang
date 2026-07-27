# Number Guessing Game

Game tebak angka berbasis CLI. Program memilih angka rahasia (1-100) dan pemain menebak sampai benar atau kehabisan kesempatan, dengan hint otomatis saat kesempatan tersisa tinggal separuh.

Project dari [roadmap.sh - Number Guessing Game](https://roadmap.sh/projects/number-guessing-game).

## Menjalankan

```bash
go run .
```

## Cara main

1. Pilih tingkat kesulitan: `easy`, `medium`, atau `hard` (menentukan jumlah maksimum percobaan).
2. Masukkan tebakan angka antara 1-100.
3. Program memberi tahu apakah tebakan terlalu rendah, terlalu tinggi, atau benar.
4. Setelah percobaan tersisa tinggal setengah dari maksimum, program menampilkan hint.
5. Setelah menang atau kehabisan percobaan, program menampilkan durasi permainan dan high score per tingkat kesulitan, lalu menawarkan main lagi.

## Struktur

- `game/` — logika inti game: parsing tingkat kesulitan, state permainan, evaluasi tebakan, hint.
- `score/` — pencatatan dan penyimpanan high score.
- `main.go` — I/O CLI dan alur permainan.

## Testing

```bash
go test ./...
```
