# 🚀 Todo API (Go + Echo)

![Go](https://img.shields.io/badge/Go-1.22-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-DB-blue)
![JWT](https://img.shields.io/badge/Auth-JWT-green)
![Swagger](https://img.shields.io/badge/API-Swagger-orange)

Simple and clean REST API for managing todos with JWT authentication.

---

## ✨ Features

- User registration
- User login (JWT)
- Authentication via Bearer Token
- Create todo
- Get all todos
- Get todo by ID
- Update todo
- Delete todo
- Swagger API documentation
- Structured logging (logrus)

---

## 🛠 Tech Stack

- **Go**
- **Echo (web framework)**
- **PostgreSQL**
- **pgx**
- **JWT (golang-jwt)**
- **Swagger (swaggo)**
- **Logrus (logging)**

---

## 🏗 Project Structure

```text
todo-app/
│
├── cmd/app                # entry point (main.go)
├── internal/
│   ├── config            # config loader
│   ├── db                # database connection
│   ├── repository        # DB layer
│   ├── service           # business logic
│   ├── http/
│   │   ├── handler       # handlers
│   │   ├── middleware    # JWT, logging
│   │   └── routes.go
│   ├── auth              # JWT logic
│   ├── apperror          # custom errors
│   ├── dto               # request/response DTO
│   └── logger            # logging setup
│
├── database/migrations   # goose migrations
├── docs/                 # swagger docs
└── go.mod

## ⚙️ Setup & Run

### 1. Clone repository

git clone https://github.com/alizhomart/todo-app.git
cd todo-app

### 2. Create .env

APP_PORT=8888
JWT_SECRET=supersecretkey

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_db

### 3. Run PostgreSQL

Make sure PostgreSQL is running locally.

### 4. Run migrations and seeds

goose -dir database/migrations postgres "postgres://postgres:postgres@localhost:5432/todo_db?sslmode=disable" up
psql "postgres://postgres:postgres@localhost:5432/todo_app?sslmode=disable" -f database/seeds/seeds.sql

### 5. Start application

go run cmd/app/main.go

## 🌐 Swagger

Open in browser:  
http://localhost:8888/swagger/index.html

---

## 🔐 Authentication

Use Bearer token:

Authorization: Bearer <your_token>

---

## 📊 Logging

Structured logging using **logrus**:

- Request logs (method, path, status)
- Error logs
- Business actions (create / update / delete)

---

## 📌 Notes

- UserID is always taken from JWT (secure)
- Passwords are hashed using bcrypt
- Clean architecture:


## 🚀 Future Improvements

- Refresh tokens
- Pagination
- Docker support
- Unit tests
- CI/CD

---

## 👨‍💻 Author

**Alizhomart Shukayev**  
GitHub: https://github.com/alizhomarts




