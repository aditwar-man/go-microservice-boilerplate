````markdown
# 🚀 Go Microservice Boilerplate - Simple Multi Wallet Currency  

Contoh **REST API dengan Clean Architecture** di Golang.  
Didesain modular, scalable, dan siap dipake buat production.  

---

## 🔧 Tech Stack  

- [echo](https://github.com/labstack/echo) – Web framework  
- [sqlx](https://github.com/jmoiron/sqlx) – DB extensions  
- [pgx](https://github.com/jackc/pgx) – PostgreSQL driver  
- [viper](https://github.com/spf13/viper) – Config manager  
- [go-redis](https://github.com/go-redis/redis) – Redis client  
- [zap](https://github.com/uber-go/zap) – Logger  
- [validator](https://github.com/go-playground/validator) – Request validation  
- [jwt-go](https://github.com/dgrijalva/jwt-go) – JWT auth  
- [uuid](https://github.com/google/uuid) – UUID utils  
- [migrate](https://github.com/golang-migrate/migrate) – DB migration  
- [minio-go](https://github.com/minio/minio-go) – S3 client  
- [bluemonday](https://github.com/microcosm-cc/bluemonday) – HTML sanitizer  
- [swag](https://github.com/swaggo/swag) – Swagger docs  
- [testify](https://github.com/stretchr/testify) – Testing toolkit  
- [gomock](https://github.com/golang/mock) – Mocking  
- [CompileDaemon](https://github.com/githubnemo/CompileDaemon) – Live reload  
- [Docker](https://www.docker.com/) – Containerization  

---

## ⚡ Development  

### Local Dev (recommended)  
```bash
make local   # run dependencies (Postgres, Redis, etc)
make run     # run service (with debugger support)
````

### Docker Dev

```bash
make docker  # run everything in docker
```

---

## 📦 Docker Compose

* `docker-compose.local.yml` → PostgreSQL, Redis, MinIO, Prometheus, Grafana
* `docker-compose.dev.yml` → Development environment
* `docker-compose.delve.yml` → Development + Delve debugger

---

## 📊 Monitoring & Docs

* Swagger → [https://localhost:5000/swagger/index.html](https://localhost:5000/swagger/index.html)
* Jaeger → [http://localhost:16686](http://localhost:16686)
* Prometheus → [http://localhost:9090](http://localhost:9090)
* Grafana → [http://localhost:3000](http://localhost:3000)

---

## 🛠️ Migration

```bash
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down
```

---

## ✅ Testing

```bash
make test
```

---

## 🎯 Features

* Clean Architecture ready
* Fully containerized with Docker
* CI/CD friendly structure
* Observability (Prometheus, Grafana, Jaeger)
* API docs with Swagger
