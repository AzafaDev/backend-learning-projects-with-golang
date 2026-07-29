# Personal Blog

Blog sederhana dengan halaman publik untuk membaca artikel dan panel admin untuk mengelola artikel (create, edit, delete). Server-side rendering pakai `html/template`, tanpa JavaScript, tanpa framework — cuma standard library Go.

Project dari [roadmap.sh - Personal Blog](https://roadmap.sh/projects/personal-blog).

## Menjalankan

```bash
go run .
```

Server jalan di `http://localhost:8080`.

## Rute

**Publik**
```
GET  /                          daftar artikel
GET  /articles/{slug}           detail artikel
```

**Admin** (butuh HTTP Basic Auth)
```
GET  /admin                          dashboard, daftar artikel + aksi edit/delete
GET  /admin/articles/new             form buat artikel baru
POST /admin/articles/new             submit artikel baru
GET  /admin/articles/{slug}/edit     form edit artikel
POST /admin/articles/{slug}/edit     submit perubahan artikel
POST /admin/articles/{slug}/delete   hapus artikel
```

Kredensial admin (hardcoded, sesuai spesifikasi project):
```
username: Akmal
password: AkmalGantengBanget
```

## Struktur data

Setiap artikel disimpan sebagai satu file JSON di `data/`, nama file = slug (contoh: `data/belajar-golang.json`). Field: `title`, `content`, `slug`, `published_at`.

Slug dibuat otomatis dari title saat artikel dibuat, dan bersifat permanen — tidak berubah walau title diedit.

## Validasi

- `title` dan `content` wajib diisi (setelah di-trim whitespace).
- `published_at` harus format `YYYY-MM-DD`.
- Slug yang sudah dipakai tidak boleh dibuat ulang (409 Conflict).
