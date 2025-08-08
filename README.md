# Golang REST API with PostgreSQL, Redis, and RabbitMQ

## 📦 Stack

- Go
- PostgreSQL (with init schema)
- Redis
- RabbitMQ
- Docker + Docker Compose

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/alviansyahexza/mt_test.git
cd mt_test
```

### 2. Make .sh file executable

```bash
chmod +x db/create-mt-test.sh
```

### 3. Build and start all services

```bash
docker-compose up --build
```

## 🌍 Access Points

API → http://localhost:3000
PostgreSQL → localhost:5432 (user: postgres, pass: mysecretpassword)
Redis → localhost:6379
RabbitMQ UI → http://localhost:15672 (user: user, pass: password)
