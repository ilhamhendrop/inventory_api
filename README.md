# Inventory API

Inventory API adalah sistem manajemen inventaris berbasis REST API yang digunakan untuk mengelola data produk, kategori, warehouse, maintenance, finance, dan user.

---

# Cara Penggunaan

## 1. Clone Repository

Buka terminal terlebih dahulu, kemudian jalankan perintah berikut untuk clone repository:

```
git clone <repository-url>
```

---

## 2. Masuk ke Folder Project

Masih melalui terminal, masuk ke folder project dengan perintah:

```
cd Inventory_API
```

---

## 3. Buka Menggunakan VS Code

Jalankan project menggunakan Visual Studio Code:

```
code .
```

Atau buka secara manual melalui VS Code.

---

## 4. Konfigurasi Environment

Cari file:

```
.env.production
```

kemudian:

* copy / rename menjadi

```
.env
```

* isi seluruh konfigurasi yang masih kosong sesuai kebutuhan project

Contohnya seperti:

* database configuration
* redis configuration
* phpMyAdmin access
* port application

---

## 5. Jalankan Docker

Buka terminal pada project dan pastikan Docker sudah terinstall serta Docker Desktop sudah berjalan, lalu jalankan:

```
docker compose up -d
```

Perintah ini akan menjalankan:

* MySQL
* Redis
* phpMyAdmin
* Service pendukung lainnya

---

## 6. Import Database

Buka phpMyAdmin menggunakan link.

```
http://localhost:9473
```

Selanjutnya:

* masuk ke database

```
invetory_db
```

* lakukan import file database

```
invetory_db.sql
```

agar seluruh tabel dan data awal tersedia.

---

## 7. Import Postman Collection

Buka aplikasi Postman, lalu import file:

```
inventory.postman_collection.json
```

File ini digunakan untuk mencoba seluruh endpoint API yang tersedia pada sistem.

---

## 8. Jalankan dan Akses API

Setelah semua proses selesai, sistem dapat diakses melalui:

```
http://localhost:8000
```

---

# Catatan

Pastikan:

* Docker Desktop berjalan
* Terminal digunakan untuk menjalankan perintah `git clone` dan `docker compose up -d`
* Port yang digunakan tidak bentrok
* File `.env` sudah terisi dengan benar
* Database berhasil diimport sebelum testing API dilakukan

Jika salah satu konfigurasi belum sesuai, sistem dapat gagal dijalankan.
