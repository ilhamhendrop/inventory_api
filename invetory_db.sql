-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: db:3306
-- Waktu pembuatan: 19 Apr 2026 pada 00.08
-- Versi server: 8.4.8
-- Versi PHP: 8.3.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Basis data: `invetory_db`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `categories`
--

CREATE TABLE `categories` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `categories`
--

INSERT INTO `categories` (`id`, `name`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES
('5e85a123-942c-49d8-a2c3-7a043abe3771', 'PC', 'Komputers', '2026-04-18 00:30:59', '2026-04-18 00:30:31', '2026-04-18 00:30:59'),
('8ef43335-848c-47d3-a3b4-90cca75b2b27', 'Laptop', 'Laptop', '2026-04-18 00:27:43', NULL, NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `finances`
--

CREATE TABLE `finances` (
  `id` varchar(255) NOT NULL,
  `maintenance_id` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `price` int NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `finances`
--

INSERT INTO `finances` (`id`, `maintenance_id`, `description`, `price`, `user_id`, `created_at`, `updated_at`, `deleted_at`) VALUES
('60ff8f85-37de-4f1f-a6b4-943a6bf5b9a9', '46ffd051-495b-492a-a0f9-77d7ae45ea1e', 'Kurakan di cpu M', 10000, 'c1f1c02b-bb4b-4451-9ed1-d1b1c0db92ac', '2026-04-18 16:53:41', '2026-04-19 00:03:52', NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `maintenances`
--

CREATE TABLE `maintenances` (
  `id` varchar(255) NOT NULL,
  `product_id` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `status` varchar(255) NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date DEFAULT NULL,
  `stock` int NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `maintenances`
--

INSERT INTO `maintenances` (`id`, `product_id`, `description`, `status`, `start_date`, `end_date`, `stock`, `user_id`, `created_at`, `updated_at`, `deleted_at`) VALUES
('46ffd051-495b-492a-a0f9-77d7ae45ea1e', '6afeacb5-763e-439e-b0f5-a4e8ffc45638', 'Pc tidak bisa hidup', 'Rusak', '2025-05-09', '2025-05-11', 2, 'c1f1c02b-bb4b-4451-9ed1-d1b1c0db92ac', '2026-04-18 04:39:35', NULL, NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `products`
--

CREATE TABLE `products` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `categorie_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `merek` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `products`
--

INSERT INTO `products` (`id`, `name`, `categorie_id`, `merek`, `created_at`, `updated_at`, `deleted_at`) VALUES
('6afeacb5-763e-439e-b0f5-a4e8ffc45638', 'Laptop A21', '8ef43335-848c-47d3-a3b4-90cca75b2b27', 'Asus', '2026-04-18 07:01:33', NULL, NULL),
('e086bac1-dbc0-461a-9b87-0aea5cebe060', 'Laptops A20', '8ef43335-848c-47d3-a3b4-90cca75b2b27', 'Asus', '2026-04-18 05:10:15', '2026-04-18 05:09:41', '2026-04-18 05:10:15'),
('e165d0fd-8237-4121-8160-5f1d44757062', 'Laptop A20', '8ef43335-848c-47d3-a3b4-90cca75b2b27', 'Asus', '2026-04-18 05:07:23', NULL, NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `role` varchar(100) NOT NULL,
  `status` varchar(150) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `username`, `name`, `role`, `status`, `password`, `created_at`, `updated_at`, `deleted_at`) VALUES
('666559fe-29e1-4780-bf2d-479ac8b1a8e1', 'Admin2', 'Admins2', 'Admin', 'Aktif', '$2a$10$eYwzvjgsB.cZHaKfLWDLxORv8ut.qDokjt2lc4oCUvfgPR41JCRoy', '2026-04-17 22:53:10', '2026-04-17 22:52:36', '2026-04-17 22:53:10'),
('c1f1c02b-bb4b-4451-9ed1-d1b1c0db92ac', 'Admin', 'Admin', 'Admin', 'Aktif', '$2a$12$jbqbBVnYMtlVBQ593NluYOC7033Q0TFLn37k98k5IaEKypBi0J3V6', '2026-04-17 07:32:26', NULL, NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `warehouses`
--

CREATE TABLE `warehouses` (
  `id` varchar(255) NOT NULL,
  `product_id` varchar(255) NOT NULL,
  `stock` int NOT NULL,
  `status` varchar(150) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `warehouses`
--

INSERT INTO `warehouses` (`id`, `product_id`, `stock`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES
('4a59544f-4281-4661-9cf3-36883c73bd24', 'e165d0fd-8237-4121-8160-5f1d44757062', 11, 'Ready', '2026-04-18 07:06:22', '2026-04-18 07:05:42', '2026-04-18 07:06:22'),
('b92df79c-bd16-448e-846a-d61b85ef4c8b', '6afeacb5-763e-439e-b0f5-a4e8ffc45638', 8, 'Ready', '2026-04-18 04:24:49', NULL, NULL);

--
-- Indeks untuk tabel yang dibuang
--

--
-- Indeks untuk tabel `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `finances`
--
ALTER TABLE `finances`
  ADD PRIMARY KEY (`id`),
  ADD KEY `finance_maintenance_constrained` (`maintenance_id`),
  ADD KEY `finance_user_constrained` (`user_id`);

--
-- Indeks untuk tabel `maintenances`
--
ALTER TABLE `maintenances`
  ADD PRIMARY KEY (`id`),
  ADD KEY `maintenance_product_constrainer` (`product_id`),
  ADD KEY `maintenance_user_constrained` (`user_id`);

--
-- Indeks untuk tabel `products`
--
ALTER TABLE `products`
  ADD PRIMARY KEY (`id`),
  ADD KEY `product_categorie_constrained` (`categorie_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `warehouses`
--
ALTER TABLE `warehouses`
  ADD PRIMARY KEY (`id`),
  ADD KEY `warehouse_product_constrained` (`product_id`);

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `finances`
--
ALTER TABLE `finances`
  ADD CONSTRAINT `finance_maintenance_constrained` FOREIGN KEY (`maintenance_id`) REFERENCES `maintenances` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `finance_user_constrained` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ketidakleluasaan untuk tabel `maintenances`
--
ALTER TABLE `maintenances`
  ADD CONSTRAINT `maintenance_product_constrainer` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `maintenance_user_constrained` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ketidakleluasaan untuk tabel `products`
--
ALTER TABLE `products`
  ADD CONSTRAINT `product_categorie_constrained` FOREIGN KEY (`categorie_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ketidakleluasaan untuk tabel `warehouses`
--
ALTER TABLE `warehouses`
  ADD CONSTRAINT `warehouse_product_constrained` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
