# Expense Tracker CLI

CLI sederhana untuk mencatat dan mengelola pengeluaran, disimpan di file JSON lokal (`expenses.json`).

Project dari [roadmap.sh - Expense Tracker](https://roadmap.sh/projects/expense-tracker).

## Menjalankan

```bash
go run . <command> [flags]
```

## Perintah

```bash
# Tambah pengeluaran
go run . add --description "Lunch" --amount 20 [--category Food]

# Update pengeluaran
go run . update --id 1 --description "New" --amount 15 [--category Food]

# Hapus pengeluaran
go run . delete --id 1

# List semua pengeluaran (bisa difilter kategori)
go run . list [--category Food]

# Ringkasan total pengeluaran (bisa difilter bulan & kategori)
go run . summary [--month 8] [--category Food]

# Export ke CSV
go run . export [--file expenses.csv]
```

## Struktur data

Setiap pengeluaran punya field: `id`, `date`, `description`, `amount`, `category`.
Data disimpan di `expenses.json` pada direktori kerja saat command dijalankan.
