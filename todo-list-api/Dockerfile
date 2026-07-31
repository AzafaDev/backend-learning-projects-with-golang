# ---- Stage 1: build ----
# Pakai image Go lengkap cuma buat compile, gak akan ikut ke image final.
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy go.mod & go.sum dulu (terpisah dari source code) supaya layer ini
# di-cache oleh Docker selama dependency belum berubah -- artinya rebuild
# jadi jauh lebih cepat kalau kamu cuma ngedit file .go, bukan go.mod.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Stage 2: run ----
# Image final cuma berisi binary hasil compile, bukan seluruh toolchain Go.
# Hasilnya: image jauh lebih kecil (~15MB vs ~350MB) dan lebih aman
# (gak ada source code atau compiler ikut ke production).
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]