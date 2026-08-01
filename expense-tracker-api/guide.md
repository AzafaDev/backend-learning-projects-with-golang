# Expense Tracker API — Build Guide

Hasil grilling session. Ini panduan desain, bukan kode — implementasinya kamu yang nulis.

Referensi: https://roadmap.sh/projects/expense-tracker-api

---

## 1. Tech Stack

Sama persis kayak `todo-list-api`, gak ganti-ganti:

- Go stdlib `net/http` + `http.NewServeMux()` (Go 1.22+ method+pattern routing, `"POST /expenses"`, `"GET /expenses/{id}"`) — no chi/gin/gorilla
- Postgres, `github.com/jackc/pgx/v5` (raw SQL, no ORM)
- `github.com/golang-jwt/jwt/v4` — JWT auth
- `golang.org/x/crypto/bcrypt` — password hashing
- `github.com/go-playground/validator/v10` — request validation
- `github.com/google/uuid` — primary keys
- `github.com/joho/godotenv` — env loading
- `golang-migrate` — DB migrations (CLI, naming: `000001_description.up.sql` / `.down.sql`)
- Testing: `stretchr/testify` (mock service-layer dependencies)

---

## 2. Project Structure

Domain-per-package (pola `blogging-platform-api`, bukan pola layered `todo-list-api`):

```
expense-tracker-api/
  cmd/api/main.go
  internal/
    config/           # env loading (PORT, DATABASE_URL, JWT_SECRET_KEY)
    database/         # pgx pool connection
    httpserver/        # WriteJSON / WriteError response envelope
    middleware/        # JWT auth middleware
    user/              # model.go, repository.go, service.go, handler.go
    expense/           # model.go, repository.go, service.go, handler.go
  migrations/
```

Tiap domain folder isinya satu file per layer, tapi digroup per fitur — bukan digroup per layer di top level.

---

## 3. Data Model

### `users` (tulis ulang dari nol, jangan copy dari `todo-list-api`)
- `id UUID PK`
- `email TEXT UNIQUE NOT NULL`
- `password_hash TEXT NOT NULL`
- `created_at`, `updated_at TIMESTAMP`

### `expenses`
| Kolom | Tipe | Catatan |
|---|---|---|
| `id` | `UUID PK` | |
| `user_id` | `UUID FK -> users.id` | |
| `amount` | `BIGINT` | **integer, satuan rupiah penuh — jangan pakai float/NUMERIC.** IDR itu zero-decimal currency, jadi integer = representasi yang benar, bukan simplifikasi. |
| `category` | `TEXT` | `CHECK (category IN ('Groceries','Leisure','Electronics','Utilities','Clothing','Health','Others'))` |
| `description` | `TEXT` | nullable, optional — gak ada `validate:"required"` |
| `date` | `DATE` | **event time** — tanggal pengeluaran terjadi (user input), terpisah dari `created_at` |
| `created_at`, `updated_at` | `TIMESTAMP` | record time — system-generated, jangan dicampur sama `date` |

Validasi kategori dobel layer: `validator` tag di request struct **dan** `CHECK` constraint di DB (defense in depth — jangan andalkan app-layer doang, ada jalur insert lain kayak migration/seed yang bisa skip validasi Go).

---

## 4. Endpoints

| Method | Path | Auth | Catatan |
|---|---|---|---|
| POST | `/register` | public | |
| POST | `/login` | public | |
| POST | `/expenses` | JWT | |
| GET | `/expenses` | JWT | filter + pagination, lihat §5 |
| GET | `/expenses/{id}` | JWT | |
| PUT | `/expenses/{id}` | JWT | full replace — semua field wajib dikirim ulang |
| DELETE | `/expenses/{id}` | JWT | hard delete |

### ⚠️ Ownership — baca ini sebelum nulis repository layer

Bug yang kejadian di `todo-list-api` (lihat commit `aa5b88a`, `fa8790a`): endpoint update/get-by-id awalnya query `WHERE id = $1` doang, ownership dicek belakangan di service layer. Itu IDOR bug — user A bisa akses expense user B kalau tau ID-nya.

**Aturan buat proyek ini: setiap repository method yang operate di single row (get-by-id, update, delete) WAJIB nerima `userID` sebagai parameter, dan `user_id` masuk ke SQL `WHERE` clause:**

```sql
SELECT ... FROM expenses WHERE id = $1 AND user_id = $2
```

Bukan `WHERE id = $1` lalu dicek belakangan di Go. Kalau row bukan milik user itu → treat sebagai "not found" (404), jangan "forbidden" (403) — jangan bocorin exist/gak-nya resource ke user yang gak punya akses.

---

## 5. Filtering — `GET /expenses`

Query params:

```
?period=week|month|3months|custom
&start_date=YYYY-MM-DD&end_date=YYYY-MM-DD   (wajib kalau period=custom)
&category=Groceries
&page=1&limit=20
```

- `period` opsional. Gak dikirim → gak ada filter tanggal.
- `period=custom` → `start_date` dan `end_date` wajib ada, validasi `start_date <= end_date`.
- Perhitungan `week`/`month`/`3months` relatif ke `NOW()` **UTC** di server (jangan overengineer timezone-per-user).
- Filter tanggal pakai kolom `date`, **bukan** `created_at`.
- `category` independen, bisa digabung sama `period` apapun (AND).
- Repository: **satu** method `GetExpenses(userID, filter, pagination)` yang build `WHERE` clause secara dinamis (slice of conditions + slice of args, join `AND`, append `LIMIT`/`OFFSET` di akhir). Jangan bikin method terpisah per kombinasi filter.

### Pagination defaults
- `page < 1` → `1`
- `limit < 1` → `20`
- `limit > 100` → clamp ke `100` (celah yang gak ada di `todo-list-api` lama — di situ cuma clamp bawah, gak ada clamp atas, jadi `?limit=999999` bisa nyedot seluruh tabel)

---

## 6. Response Format

Reuse pola `blogging-platform-api/internal/httpserver`:

```go
type APIErrorRespose struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any)
func WriteError(w http.ResponseWriter, status int, message string, err error)
```

---

## 7. Testing

Scope minimal: **unit test service layer** pakai mock repository (`testify/mock`), fokus utama:

- Fungsi hitung `period → (start_date, end_date)` — ini paling gampang salah off-by-one, dan paling gampang ditest karena pure function (input `time.Time`, output `time.Time`), gak butuh DB sama sekali.
- Ownership logic di service (kalau ada logic tambahan di luar SQL WHERE).

Repository-layer test (query SQL beneran, pakai `pgxmock` atau testcontainers) itu opsional/nice-to-have, bukan wajib buat versi ini.

---

## 8. Suggested Build Order

1. `config` + `database` (copy pola, bukan logic, dari proyek lama)
2. Migration: `users` + `expenses` table (termasuk `CHECK` constraint kategori)
3. `user` domain: register, login, bcrypt hash, JWT issue — tulis ulang dari nol
4. `middleware`: JWT auth middleware
5. `expense` domain: model → repository (CRUD + ownership di WHERE) → service → handler
6. Tambahin filtering (period calculation) + pagination ke `GetExpenses`
7. Unit test fungsi period calculation
8. Manual test end-to-end tiap endpoint (curl/Postman), termasuk coba akses expense user lain buat verifikasi ownership beneran ke-block

---

Kalau nanti nulis kode dan mau di-review atau stuck, paste ke chat aja.
