# TechStore – Microservices E-Commerce Platform

## 🧠 Project Overview
**TechStore** is a scalable e-commerce platform built with Go and gRPC, designed using Clean Architecture principles. It demonstrates distributed system communication via gRPC, message queues (NATS), caching (Redis), and container orchestration (Docker Compose).

### 🎯 Key Features
- User registration, authentication, and profile with JWT
- Product CRUD operations and listing
- Order placement and tracking
- Email notifications on user registration
- Redis caching for performance
- NATS for event-driven communication
- Unit and Integration tests
- Docker Compose setup for microservices

---

## ⚙️ Technologies Used
- **Go (Golang)**
- **gRPC + Protocol Buffers**
- **PostgreSQL** – data persistence
- **Redis** – caching
- **NATS** – event messaging
- **Gin** – HTTP gateway
- **Docker & Docker Compose**
- **bcrypt** – password hashing
- **JWT** – user authentication
- **Gomail** – sending email
- **Testify** – testing

---

## 🚀 How to Run Locally

> Requirements: Docker Desktop installed

### 1. Clone the repo
```bash
git clone https://github.com/<your-username>/TechStore.git
cd TechStore
```

### 2. Build & run all services
```bash
docker-compose up --build
```

- API Gateway: http://localhost:8081
- PostgreSQL: localhost:5432 (user: `postgres`, password: `0000`)
- Redis: localhost:6379
- gRPC Ports:
    - InventoryService: 50051
    - OrderService: 50052
    - UserService: 50053

---

## 🧪 How to Run Tests

### Unit Tests
```bash
go test ./user_service/internal/service
go test ./inventory_service/internal/service
go test ./order_service/internal/service
```

### Integration Tests
```bash
go test ./inventory_service/internal/integration
```

---

## 📡 gRPC Endpoints

### UserService (Port: 50053)
| Method            | Description                  |
|------------------|------------------------------|
| `RegisterUser`   | Register new user            |
| `AuthenticateUser` | Login and get JWT token     |
| `GetUserProfile` | Get user profile via token   |

### InventoryService (Port: 50051)
| Method            | Description                  |
|------------------|------------------------------|
| `CreateProduct`   | Add new product              |
| `GetProduct`      | Get product by ID            |
| `UpdateProduct`   | Edit product info            |
| `DeleteProduct`   | Remove product               |
| `ListProducts`    | List all products            |

### OrderService (Port: 50052)
| Method            | Description                  |
|------------------|------------------------------|
| `CreateOrder`     | Place a new order            |
| `GetOrder`        | Retrieve order by ID         |
| `DeleteOrder`     | Cancel order                 |
| `ListOrders`      | View all orders              |

---

## ✅ List of Implemented Features

- [x] Clean Architecture
- [x] 12+ gRPC endpoints
- [x] PostgreSQL integration with migrations
- [x] Redis caching
- [x] NATS message queue for product creation event
- [x] JWT authentication
- [x] Sending welcome emails
- [x] Docker Compose setup
- [x] Unit and integration tests

---

## 📌 Notes

- `docker-compose.yml` contains all services and shared network
- `Makefile` or `.sh` scripts can be added for automation
- Grafana + Prometheus can be integrated for observability (bonus)

---

## 💡 Author
Adel Kenesova (Astana IT University)