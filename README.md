# Soal Hub API 🚀

Aplikasi backend REST API berkinerja tinggi yang dibangun menggunakan **Go (Golang)** dan **Fiber v3**, menerapkan pola arsitektur *Domain-Driven / Layered Architecture* untuk kemudahan pengujian dan pemeliharaan kode.

---

## 🛠️ Tech Stack & Tooling

* **Bahasa:** [Go (Golang)](https://go.dev/) (v1.21+)
* **Framework Web:** [Go Fiber v3](https://gofiber.io/)
* **Database:** [PostgreSQL](https://www.postgresql.org/)
* **Containerization:** [Docker](https://www.docker.com/) & Docker Compose
* **Live Reloading:** [Air](https://github.com/air-verse/air)
* **Architecture:** Domain-Driven / Layered Architecture (Handler, Service, Repository)

---

## 📁 Struktur Proyek

```text
soalhub/
├── cmd/
│   └── api/
│       └── main.go           # Entry point utama aplikasi
├── config/                   # Konfigurasi environment & database
├── internal/                 # Kode privat aplikasi
│   ├── database/             # Koneksi & migrasi database
│   ├── middleware/           # Middleware HTTP (Auth, CORS, Logger)
│   └── modules/              # Domain/Modul aplikasi (User, Soal, dll.)
│       └── [module_name]/
│           ├── dto.go        # Request/Response struct (Binding Fiber)
│           ├── handler.go    # HTTP Handler (Controller)
│           ├── model.go      # Struct DB / Entity
│           ├── repository.go # Query langsung ke Database
│           └── service.go    # Logika bisnis utama
├── pkg/                      # Utilities & helpers umum
├── .air.toml                 # Konfigurasi Hot Reload (Air)
├── .env.example              # Template environment variables
├── docker-compose.yml        # Orchestration layanan pendukung (Postgres)
├── go.mod
└── go.sum