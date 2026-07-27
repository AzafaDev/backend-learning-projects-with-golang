# Task Tracker CLI

CLI untuk mengelola daftar tugas (task), disimpan di file JSON lokal (`tasks.json`).

Project dari [roadmap.sh - Task Tracker](https://roadmap.sh/projects/task-tracker).

## Menjalankan

```bash
go run . <command> [arguments]
```

## Perintah

```bash
# Tambah task baru
go run . add "Buy groceries"

# List semua task (bisa difilter status: todo / in-progress / done)
go run . list [status]

# Update deskripsi task
go run . update <id> "New description"

# Hapus task
go run . delete <id>

# Tandai status task
go run . mark-in-progress <id>
go run . mark-done <id>
```

## Struktur data

Setiap task punya field: `id`, `description`, `status`, `createdAt`, `updatedAt`.
Data disimpan di `tasks.json` pada direktori kerja saat command dijalankan; file akan dibuat otomatis jika belum ada.
