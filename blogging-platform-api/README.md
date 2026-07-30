# Blogging Platform API

RESTful API sederhana untuk platform blogging pribadi, dibangun sebagai bagian dari [roadmap.sh Blogging Platform API project](https://roadmap.sh/projects/blogging-platform-api).

## Tech Stack

- **Go 1.26** — `net/http` standard library, tanpa framework routing tambahan
- **PostgreSQL 16**
- **Docker & Docker Compose**
- **pgx/v5** — database driver & connection pool
- **go-playground/validator** — request body validation
- **golang-migrate** — database migration
- **pgxmock** + **testify** — unit testing repository layer

## Struktur Project

```
├── cmd/api                 # entry point aplikasi
├── internal
│   ├── config              # load environment variables
│   ├── database             # koneksi database (pgxpool)
│   ├── httpserver            # response helper (JSON, error)
│   └── post                   # domain post: handler, service, repository
├── migrations                 # SQL migration files
├── docker-compose.yml
└── Dockerfile
```

Project ini mengikuti **layered architecture**:

```
Handler  →  Service  →  Repository  →  Database
```

- **Handler** — parsing HTTP request/response, tidak ada business logic
- **Service** — validasi request, orchestration
- **Repository** — query ke database, translate error database ke error domain

## Menjalankan Project

### Opsi 1 — Full Docker (paling mudah)

```bash
docker compose up --build
```

Ini akan menjalankan database, migrasi otomatis, dan API sekaligus. API berjalan di `http://localhost:8080`.

### Opsi 2 — Development lokal (API di host, DB di Docker)

```bash
# jalankan db & migrate saja
docker compose up -d db migrate

# copy .env.example ke .env, sesuaikan jika perlu
cp .env.example .env

# jalankan API di lokal
go run ./cmd/api
```

## Environment Variables

| Variable | Deskripsi | Default |
|---|---|---|
| `PORT` | Port API berjalan | `8080` |
| `DATABASE_URL` | Connection string PostgreSQL | wajib diisi |

Contoh `.env` untuk development lokal:

```env
PORT=8080
DATABASE_URL=postgres://bloguser:blogpass@localhost:5432/blogdb?sslmode=disable
```

## API Endpoints

| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` | `/posts` | Membuat post baru |
| `GET` | `/posts` | Mengambil semua post |
| `GET` | `/posts?term=xxx` | Mencari post berdasarkan title/content/category |
| `GET` | `/posts/{id}` | Mengambil satu post berdasarkan id |
| `PUT` | `/posts/{id}` | Mengupdate post berdasarkan id |
| `DELETE` | `/posts/{id}` | Menghapus post berdasarkan id |

### Contoh Request

**Create Post**

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post.",
    "category": "Technology",
    "tags": ["Tech", "Programming"]
  }'
```

Response `201 Created`:

```json
{
  "id": "a5f33539-5c2c-4aa3-9a41-c89fab41557d",
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post.",
  "category": "Technology",
  "tags": ["Tech", "Programming"],
  "created_at": "2026-07-30T15:40:13.683122Z",
  "updated_at": "2026-07-30T15:40:13.683122Z"
}
```

**Get All Posts**

```bash
curl http://localhost:8080/posts
```

**Search Posts**

```bash
curl "http://localhost:8080/posts?term=tech"
```

**Get Single Post**

```bash
curl http://localhost:8080/posts/{id}
```

**Update Post**

```bash
curl -X PUT http://localhost:8080/posts/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content.",
    "category": "Technology",
    "tags": ["Tech"]
  }'
```

**Delete Post**

```bash
curl -i -X DELETE http://localhost:8080/posts/{id}
```

Response `204 No Content` (tanpa body).

### Error Response

```json
{
  "success": false,
  "message": "post not found",
  "error": "post not found"
}
```

| Status | Kapan terjadi |
|---|---|
| `400 Bad Request` | Request body tidak valid / gagal validasi / id bukan UUID valid |
| `404 Not Found` | Post dengan id tersebut tidak ditemukan |
| `500 Internal Server Error` | Kesalahan tak terduga (misal database error) |

## Menjalankan Test

```bash
go test ./... -v
```

Unit test repository menggunakan `pgxmock`, sehingga tidak memerlukan koneksi database sungguhan.

## Database Schema

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## Catatan

- `id` menggunakan UUID (`gen_random_uuid()`), bukan integer auto-increment.
- Tidak ada implementasi pagination maupun authentication/authorization — sesuai scope project ini.
- Pencarian (`?term=`) melakukan wildcard case-insensitive search (`ILIKE`) pada kolom `title`, `content`, dan `category`.