# Petunjuk Pengunaan

Aplikasi manajemen produk (Northwind-style) dengan backend Go (Echo + MySQL) dan frontend JavaScript + Tailwind CSS (Vite).

## Struktur Project

```
est-pemrograman-go/
├── backend/          # API Go (Echo)
│   ├── api/
│   │   └── product/  # API produk, kategori, supplier
│   ├── database/     # Koneksi MySQL
│   ├── report/      # Laporan penjualan
│   ├── main.go
│   ├── go.mod
│   └── .env         # Konfigurasi DB (jangan commit)
├── frontend/        # Web UI (Vite + Tailwind)
│   ├── src/
│   │   ├── main.js
│   │   ├── api.js
│   │   └── style.css
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
└── README.md
```


## Persyaratan

- **Backend:** Go 1.21+, MySQL
- **Frontend:** Node.js 18+, npm

## Setup

### 1. Backend

```bash
cd backend
```

Buat file `.env` (copy dari `.env.example` jika ada) dengan isi:

```env
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=test-pemrograman-go
```

Install dependency dan jalankan:

```bash
go get name_pkg
go run main.go
```

Server API berjalan di **http://localhost:8083**.

### 2. Frontend

```bash
cd frontend
npm install
npm run dev
```

Dev server berjalan di **http://localhost:5173**. Request ke `/api` dan `/report` di-proxy ke backend (port 8083).

## API Backend

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/products` | Daftar produk (query: page, limit, category_id, supplier_id, min_price, max_price, keyword, sort) |
| GET | `/api/products/:id` | Detail produk |
| POST | `/api/products` | Buat produk baru |
| GET | `/api/categories` | Daftar kategori |
| GET | `/api/suppliers` | Daftar supplier |
| GET | `/report/sales/top-customers` | Laporan top customers (text) |

## Fitur Frontend

- **Product List** – Tabel produk dengan filter (kategori, supplier, harga, pencarian, sort) dan pagination.
- **Product Detail** – Modal detail produk (View Details).
- **Remind Me** (jika ada modul reminder) – Daftar reminder dari backend, tambah/tandai selesai/hapus.

## Script NPM (Frontend)

- `npm run dev` – Dev server dengan hot reload
- `npm run build` – Build production ke `dist/`
- `npm run preview` – Preview hasil build

## Lisensi

Private / pendidikan.
