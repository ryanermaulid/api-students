# REST API Mahasiswa - Database & Repository Pattern

Aplikasi backend REST API Mahasiswa yang dibangun menggunakan Go Fiber, PostgreSQL via driver `pgx/v5`, dan arsitektur Repository Pattern.

## Skema Database
Skema tabel `students` didefinisikan pada file `migrations/001_create_students.sql`:
- `id` (SERIAL, PRIMARY KEY)
- `nim` (VARCHAR(20), NOT NULL, UNIQUE CASE-INSENSITIVE)
- `name` (VARCHAR(100), NOT NULL, INDEXED LOWER)
- `grade` (DOUBLE PRECISION, NOT NULL, DEFAULT 0.0)
- `is_active` (BOOLEAN, NOT NULL, DEFAULT TRUE)
- `created_at` (TIMESTAMPTZ, NOT NULL, DEFAULT NOW())

## Menyiapkan Basis Data dari Nol
1. Pastikan PostgreSQL telah terpasang dan aktif di komputer lokal.
2. Buat database baru:
   ```bash
   psql -U postgres -c "CREATE DATABASE praktikum_backend;"