# GitHub User Activity CLI

CLI untuk menampilkan aktivitas terbaru seorang user GitHub, diambil langsung dari GitHub public Events API (`/users/:username/events`), tanpa autentikasi.

Project dari [roadmap.sh - GitHub User Activity](https://roadmap.sh/projects/github-user-activity-cli).

## Menjalankan

```bash
go run . <username>
```

Contoh:
```bash
go run . octocat
```

## Perilaku

- Mengambil event terbaru dari `https://api.github.com/users/<username>/events`.
- Menampilkan ringkasan aktivitas per baris, dengan format berbeda tergantung tipe event:
  - `PushEvent` — push ke branch tertentu
  - `CreateEvent` — pembuatan repo/branch/tag
  - `WatchEvent` — starring repo
  - `IssuesEvent` — aksi pada issue (opened/closed/dll)
  - `PullRequestEvent` — aksi pada pull request
  - Tipe event lain ditampilkan dengan format generik
- Jika username tidak ditemukan (404), CLI menampilkan pesan error dan keluar dengan status non-zero.
