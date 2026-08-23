# API Mahasiswa (Student API) - REST API & HTTP Deep Dive

Repositori ini merupakan implementasi Tugas Mandiri Modul 2 Praktikum Pemrograman Backend Lanjut.

## Kontrak API

| Metode | Endpoint | Parameter | Contoh Body Permintaan | Status yang Dikembalikan | Contoh Respons |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/students` | **Query:**<br>`page` (int, def: 1)<br>`limit` (int, def: 10, max: 50)<br>`search` (string)<br>`sort` (id, nim, name, grade)<br>`order` (asc, desc)<br>`is_active` (bool)<br>`min_grade` (int) | *(Tidak ada)* | `200 OK` | `{"success": true, "message": "daftar mahasiswa berhasil diambil", "data": [...], "meta": {"page": 1, "limit": 10, "total": 1, "total_pages": 1}}` |
| `GET` | `/api/v1/students/:id` | **Path:**<br>`id` (int, positif) | *(Tidak ada)* | `200 OK`<br>`400 Bad Request`<br>`404 Not Found` | `{"success": true, "message": "data mahasiswa ditemukan", "data": {"id": 1, "nim": "5025211001", "name": "Ryan Ermaulid", "grade": 88, "is_active": true}}` |
| `POST` | `/api/v1/students` | *(Tidak ada)* | `{"nim": "5025211001", "name": "Ryan Ermaulid", "grade": 88, "is_active": true}` | `201 Created`<br>`400 Bad Request`<br>`409 Conflict`<br>`415 Unsupported Media Type`<br>`422 Unprocessable Entity` | `{"success": true, "message": "mahasiswa berhasil ditambah", "data": {"id": 1, "nim": "5025211001", "name": "Ryan Ermaulid", "grade": 88, "is_active": true}}` |
| `PUT` | `/api/v1/students/:id` | **Path:**<br>`id` (int, positif) | `{"nim": "5025211001", "name": "Ryan Updated", "grade": 90, "is_active": false}` | `200 OK`<br>`400 Bad Request`<br>`404 Not Found`<br>`409 Conflict`<br>`415 Unsupported Media Type`<br>`422 Unprocessable Entity` | `{"success": true, "message": "data mahasiswa berhasil diganti seluruhnya", "data": {"id": 1, "nim": "5025211001", "name": "Ryan Updated", "grade": 90, "is_active": false}}` |
| `PATCH` | `/api/v1/students/:id` | **Path:**<br>`id` (int, positif) | `{"is_active": true}` | `200 OK`<br>`400 Bad Request`<br>`404 Not Found`<br>`409 Conflict`<br>`415 Unsupported Media Type`<br>`422 Unprocessable Entity` | `{"success": true, "message": "data mahasiswa berhasil diperbarui sebagian", "data": {"id": 1, "nim": "5025211001", "name": "Ryan Updated", "grade": 90, "is_active": true}}` |
| `DELETE` | `/api/v1/students/:id` | **Path:**<br>`id` (int, positif) | *(Tidak ada)* | `204 No Content`<br>`400 Bad Request`<br>`404 Not Found` | *(Tanpa Body)* |