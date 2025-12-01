# Go Social API

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)
![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=flat&logo=swagger)

A production-ready, layered REST API for a social media platform. Built with **Go**, **Gin**, and **PostgreSQL**, featuring clean architecture, JWT authentication, and automated seeding.

## ⚡ Features

- **🔐 Auth:** Secure JWT authentication (Signup, Login, Password Reset).
- **📝 Content:** CRUD for Posts, Comments, and Nested Replies.
- **❤️ Social:** Like system, Follow/Unfollow users, and Tagging.
- **📰 Feed:** Personalized feeds based on follower graph.
- **🔍 Search:** User search and Tag filtering.
- **🐳 DevOps:** Fully Dockerized with Make automation and Migrations.

## 🚀 Quick Start

### 1. Environment Setup

Create a `.env` file in the root directory:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=social_db
JWT_SECRET=supersecretkey
TEST_DB_URL=postgres://postgres:password@localhost:5432/social_db_test
```

### 2. Run with Docker

The easiest way to start the API and Database:

```bash
make up
```

_API will be available at `http://localhost:3000`_

## 📚 Documentation

Interactive API documentation is available via Swagger UI:

👉 **[http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)**

## 🛠 Make Commands

| Command           | Description                               |
| :---------------- | :---------------------------------------- |
| `make up`         | Start API and DB containers in background |
| `make down`       | Stop and remove containers                |
| `make logs`       | View live container logs                  |
| `make run`        | Run the application locally (requires DB) |
| `make test`       | Run unit tests                            |
| `make swagger`    | Regenerate Swagger documentation          |
| `make migrate-up` | Apply database migrations                 |

## 📂 Architecture

The project follows a **Clean Architecture** pattern to ensure separation of concerns and scalability:

- **`cmd/`**: Application entry points (API server, Seeder).
- **`internal/controller`**: **Gin** handlers responsible for request validation and response formatting.
- **`internal/services`**: Pure business logic layer, decoupled from HTTP and Database details.
- **`internal/database`**: Data access layer using raw SQL for **PostgreSQL** interactions.
- **`internal/model`**: Core domain entities and database schemas.
- **`internal/dto`**: Data Transfer Objects for strict API contract definition.
- **`internal/middleware`**: Cross-cutting concerns like **JWT Auth**, Logging, and Rate Limiting.

#### Note: ( Readme generated via AI )
